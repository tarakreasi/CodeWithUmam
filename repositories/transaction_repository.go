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
func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	// 1. Mulai Database Transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	// Defer Rollback: Jika terjadi error di tengah jalan, batalkan semua perubahan.
	// Jika nanti berhasil Commit, Rollback tidak akan melakukan apa-apa.
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// 2. Loop setiap item yang dibeli
	for _, item := range items {
		var productPrice, stock int
		var productName string

		// Cek harga, stok, dan nama produk
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = ?", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Cek apakah stok mencukupi
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

		// Simpan detail untuk nanti diinsert ke transaction_details
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName, // Optional: simpan nama produk history
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// 3. Insert ke tabel transaction header
	var transactionID int64
	// SQLite tidak support RETURNING id secara native di semua versi/driver dengan mudah, jadi pakai LastInsertId
	res, err := tx.Exec("INSERT INTO transactions (total_amount) VALUES (?)", totalAmount)
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
		ID:          int(transactionID),
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

// GetDailySalesSummary mengambil laporan penjualan hari ini.
func (repo *TransactionRepository) GetDailySalesSummary() (*models.SalesSummary, error) {
	summary := &models.SalesSummary{}

	// Query 1: Total Revenue hari ini (SQLite: date('now', 'localtime'))
	err := repo.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE date(created_at) = date('now')").Scan(&summary.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("gagal hitung revenue: %v", err)
	}

	// Query 2: Total Transaksi hari ini
	err = repo.db.QueryRow("SELECT COUNT(id) FROM transactions WHERE date(created_at) = date('now')").Scan(&summary.TotalTransaksi)
	if err != nil {
		return nil, fmt.Errorf("gagal hitung transaksi: %v", err)
	}

	// Query 3: Produk Terlaris hari ini
	// Join transaction_details, transactions, dan products
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
		// Belum ada penjualan hari ini
		summary.ProdukTerlaris = models.ProductSales{Name: "-", QtyTerjual: 0}
	} else if err != nil {
		return nil, fmt.Errorf("gagal cari produk terlaris: %v", err)
	}

	return summary, nil
}
