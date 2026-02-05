package database

import (
	"database/sql"
	"log"

	// Import driver sqlite3
	// Tanda underscore (_) berarti kita meng-import package ini hanya untuk side-effect nya (init function),
	// tapi tidak menggunakan function apapun dari package tersebut secara langsung di kode kita.
	// Driver ini diperlukan agar `database/sql` tau cara bicara dengan SQLite.
	_ "github.com/mattn/go-sqlite3"
)

// InitDB membuka koneksi ke database dan memastikan tabel-tabel yang dibutuhkan sudah ada.
func InitDB(connStr string) (*sql.DB, error) {
	// Membuka koneksi ke database.
	// `sqlite3` adalah nama driver yang didaftarkan oleh `github.com/mattn/go-sqlite3`.
	// `connStr` adalah path ke file database (misal: ./data.db).
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}

	// Cek apakah database benar-benar bisa diakses (Ping).
	// Open() terkadang hanya memvalidasi argumen tanpa benar-benar connect.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Buat tabel jika belum ada (Migration sederhana).
	createTables(db)

	log.Println("✅ Database berhasil terkoneksi")
	return db, nil
}

// createTables membuat tabel 'categories' dan 'products' jika belum ada.
func createTables(db *sql.DB) {
	// Query untuk membuat tabel categories
	// AUTOINCREMENT: ID akan bertambah otomatis (1, 2, 3...)
	queryCategories := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT
	);`

	// Query untuk membuat tabel products
	// FOREIGN KEY: Menandakan bahwa category_id merujuk ke id di tabel categories.
	queryProducts := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price INTEGER,
		stock INTEGER,
		category_id INTEGER,
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);`

	// Eksekusi query pembuatan tabel categories
	if _, err := db.Exec(queryCategories); err != nil {
		log.Fatal("Gagal membuat tabel categories:", err)
	}

	// Eksekusi query pembuatan tabel products
	if _, err := db.Exec(queryProducts); err != nil {
		log.Fatal("Gagal membuat tabel products:", err)
	}

	// ==========================================
	// Bootcamp Session 3: Transaction Tables
	// ==========================================

	// Query untuk membuat tabel transactions
	// Menyimpan header transaksi (total belanja, waktu transaksi)
	queryTransactions := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_amount INTEGER NOT NULL,
		paid_amount INTEGER,
		change INTEGER,
		payment_method TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Query untuk membuat tabel transaction_details
	// Menyimpan detail barang yang dibeli per transaksi
	queryTransactionDetails := `
	CREATE TABLE IF NOT EXISTS transaction_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_id INTEGER,
		product_id INTEGER,
		quantity INTEGER NOT NULL,
		subtotal INTEGER NOT NULL,
		FOREIGN KEY(transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
		FOREIGN KEY(product_id) REFERENCES products(id)
	);`

	if _, err := db.Exec(queryTransactions); err != nil {
		log.Fatal("Gagal membuat tabel transactions:", err)
	}

	if _, err := db.Exec(queryTransactionDetails); err != nil {
		log.Fatal("Gagal membuat tabel transaction_details:", err)
	}
}
