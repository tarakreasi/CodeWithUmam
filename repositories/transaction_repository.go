package repositories

import (
	"codeWithUmam/models"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction memproses pembelian barang.
// Menggunakan Database Transaction (Begin -> Commit/Rollback) untuk menjaga integritas data.
// Konsep Transaction (ACID):
// - Atomicity: Semua berhasil atau gagal semua. Tidak boleh ada stok berkurang tapi transaksi gagal dicatat.
// - Consistency: Data harus valid sebelum dan sesudah transaksi.
// - Isolation: Transaksi ini tidak boleh terganggu transaksi lain yang berjalan bersamaan.
// - Durability: Setelah commit, data tersimpan permanen.
// CreateTransaction memproses pembelian barang.
func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem, paidAmount int, change int, paymentMethod string) (*models.Transaction, error) {
	// 1. Mulai Database Transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// 2. Loop setiap item yang dibeli
	for _, item := range items {
		var productPrice, stock int
		var productName string

		// Ambil data produk terbaru
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = ?", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Validasi Stok
		if stock < item.Quantity {
			return nil, fmt.Errorf("stok tidak cukup untuk produk %s (sisa: %d)", productName, stock)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Kurangi stok produk
		_, err = tx.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName, // Optional: simpan nama produk history
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Validasi ulang Paid Amount
	// Note: change argument is ignored, we calculate it here to be safe from client manipulation
	if paidAmount < totalAmount {
		return nil, fmt.Errorf("uang pembayaran kurang (Total: %d, Paid: %d)", totalAmount, paidAmount)
	}

	realChange := paidAmount - totalAmount

	// 3. Insert ke tabel transaction header
	var transactionID int64
	// SQLite tidak support RETURNING id secara native di semua versi/driver dengan mudah, jadi pakai LastInsertId
	res, err := tx.Exec("INSERT INTO transactions (total_amount, paid_amount, change, payment_method) VALUES (?, ?, ?, ?)", totalAmount, paidAmount, realChange, paymentMethod)
	if err != nil {
		return nil, err
	}
	transactionID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 4. Insert ke tabel transaction details
	for i := range details {
		details[i].TransactionID = int(transactionID)
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES (?, ?, ?, ?)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	// 5. Commit Transaksi (Simpan permanen)
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:            int(transactionID),
		TotalAmount:   totalAmount,
		PaidAmount:    paidAmount,
		Change:        realChange,
		PaymentMethod: paymentMethod,
		Details:       details,
	}, nil
}

// GetDailySalesSummary mengambil laporan penjualan hari ini.
// Menggunakan fungsi agregasi SQL (SUM, COUNT, MAX) dan JOIN tabel.
func (repo *TransactionRepository) GetDailySalesSummary() (*models.SalesSummary, error) {
	summary := &models.SalesSummary{}

	// Query 1: Total Revenue hari ini
	// COALESCE digunakan agar jika hasilnya NULL (tidak ada penjualan), diganti jadi 0.
	err := repo.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE date(created_at) = date('now')").Scan(&summary.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("gagal hitung revenue: %v", err)
	}

	// Query 2: Total Transaksi hari ini
	// Menghitung berapa baris transaksi yang terjadi hari ini.
	err = repo.db.QueryRow("SELECT COUNT(id) FROM transactions WHERE date(created_at) = date('now')").Scan(&summary.TotalTransaksi)
	if err != nil {
		return nil, fmt.Errorf("gagal hitung transaksi: %v", err)
	}

	// Query 3: Produk Terlaris hari ini
	// Ini query agak kompleks (Intermediate SQL):
	// 1. JOIN 3 tabel: transaction_details -> transactions -> products
	// 2. Filter hanya transaksi hari ini
	// 3. GROUP BY nama produk (kelompokkan penjualan per produk)
	// 4. SUM quantity (hitung total qty terjual per produk)
	// 5. ORDER BY qty DESC (urutkan dari yang paling banyak terjual)
	// 6. LIMIT 1 (ambil juara 1 nya saja)
	queryBestSeller := `
		SELECT p.name, SUM(td.quantity) as qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE date(t.created_at) = date('now')
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`
	err = repo.db.QueryRow(queryBestSeller).Scan(&summary.ProdukTerlaris.Name, &summary.ProdukTerlaris.QtyTerjual)

	if err == sql.ErrNoRows {
		// Belum ada penjualan hari ini, set default strip (-) dan 0
		summary.ProdukTerlaris = models.ProductSales{Name: "-", QtyTerjual: 0}
	} else if err != nil {
		return nil, fmt.Errorf("gagal cari produk terlaris: %v", err)
	}

	return summary, nil
}

// FindAll mengambil semua data transaksi, opsional dengan filter tanggal.
// filter start/end format: YYYY-MM-DD
func (repo *TransactionRepository) FindAll(start, end string) ([]models.Transaction, error) {
	query := "SELECT id, total_amount, created_at FROM transactions"
	args := []interface{}{}

	if start != "" && end != "" {
		// Filter by date range (inclusive)
		// SQLite date function: date(created_at)
		query += " WHERE date(created_at) BETWEEN ? AND ?"
		args = append(args, start, end)
	}

	query += " ORDER BY created_at DESC"

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// FindByID mengambil detail transaksi beserta item-nya (JOIN).
func (repo *TransactionRepository) FindByID(id int) (*models.Transaction, error) {
	// 1. Ambil Header Transaksi
	var t models.Transaction
	err := repo.db.QueryRow("SELECT id, total_amount, created_at FROM transactions WHERE id = ?", id).Scan(&t.ID, &t.TotalAmount, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}

	// 2. Ambil Details (Items)
	rows, err := repo.db.Query("SELECT id, product_id, quantity, subtotal FROM transaction_details WHERE transaction_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []models.TransactionDetail
	for rows.Next() {
		var d models.TransactionDetail
		if err := rows.Scan(&d.ID, &d.ProductID, &d.Quantity, &d.Subtotal); err != nil {
			return nil, err
		}
		// Optional: Ambil nama produk jika perlu, tapi itu butuh JOIN lagi.
		// Untuk efisiensi, kita bisa JOIN di query pertama atau query terpisah.
		// Mari kita ambil nama produk sekalian biar lengkap.
		var productName string
		repo.db.QueryRow("SELECT name FROM products WHERE id = ?", d.ProductID).Scan(&productName)
		d.ProductName = productName

		details = append(details, d)
	}

	t.Details = details
	return &t, nil
}
