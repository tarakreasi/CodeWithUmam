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
func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	// 1. Mulai Database Transaction
	// `tx` adalah objek khusus (mirip `db`) tapi semua operasinya tertahan (pending) sampai kita panggil Commit().
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	// Defer Rollback: Mekanisme safety net.
	// Jika fungsi return error di tengah jalan (sebelum Commit), Rollback akan membatalkan semua perubahan.
	// Jika Commit berhasil dijalankan, Rollback tidak akan melakukan apa-apa (no-op).
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// 2. Loop setiap item yang dibeli user
	for _, item := range items {
		var productPrice, stock int
		var productName string

		// Ambil data produk terbaru (Harga & Stok) dari DB.
		// PENTING: Jangan percaya harga dari frontend/user input, fatal! Selalu ambil dari DB.
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = ?", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Validasi Stok: Cek apakah stok cukup sebelum dijual.
		if stock < item.Quantity {
			return nil, fmt.Errorf("stok tidak cukup untuk produk %s (sisa: %d)", productName, stock)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Kurangi stok produk
		// Kita jalankan query UPDATE update products stok.
		// Karena ini dalam block `tx`, perubahan ini belum permanen sampai `tx.Commit()`.
		_, err = tx.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Simpan detail ke memory dulu (slice details), nanti diinsert sekaligus.
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName, // Optional: simpan nama produk history
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// 3. Insert ke tabel transaction header
	// Kita simpan total belanjaan dan dapatkan Transaction ID yang baru digenerate.
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
	// Kita simpan rincian barang apa saja yang dibeli dengan TransactionID yang baru kita dapat.
	for i := range details {
		details[i].TransactionID = int(transactionID)
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES (?, ?, ?, ?)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	// 5. Commit Transaksi (Simpan permanen)
	// Ini adalah titik penentuan. Jika baris ini sukses, semua perubahan (stok berkurang, insert transaksi) jadi permanen.
	// Jika gagal, semua dibatalkan oleh defer Rollback() di atas.
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
