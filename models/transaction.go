package models

import "time"

// Transaction merepresentasikan header transaksi belanja.
// Struct ini mencerminkan tabel `transactions` di database.
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"` // Relasi: Satu transaksi punya banyak detail (One-to-Many)
}

// TransactionDetail merepresentasikan detail item dalam satu transaksi.
// Ini untuk mencatat barang apa saja yang dibeli dalam satu struk.
type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"` // omitempty: Field ini tidak akan muncul di JSON jika string-nya kosong ""
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"` // Harga satuan * Quantity saat transaksi terjadi
}

// CheckoutItem adalah input dari User/Frontend untuk request checkout.
// Kita pisahkan struct ini karena User hanya perlu kirim ProductID dan Qty, sisanya (Harga, Nama) kita ambil dari DB.
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CheckoutRequest adalah format body JSON untuk endpoint POST /api/checkout.
// Contoh JSON: { "items": [ { "product_id": 1, "quantity": 2 } ] }
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// ProductSales merepresentasikan data penjualan produk (untuk report).
// Digunakan untuk menampilkan produk terlaris di laporan harian.
type ProductSales struct {
	Name       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

// SalesSummary adalah response untuk endpoint report harian.
// Menggabungkan total omset, jumlah transaksi, dan produk best seller dalam satu response JSON.
type SalesSummary struct {
	TotalRevenue   int          `json:"total_revenue"`
	TotalTransaksi int          `json:"total_transaksi"`
	ProdukTerlaris ProductSales `json:"produk_terlaris"`
}
