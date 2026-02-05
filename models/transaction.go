package models

import "time"

// Transaction merepresentasikan header transaksi belanja.
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

// TransactionDetail merepresentasikan detail item dalam satu transaksi.
type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"` // omitempty: tidak muncul jika kosong
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

// CheckoutItem adalah input dari User/Frontend untuk request checkout.
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CheckoutRequest adalah format body JSON untuk endpoint POST /api/checkout.
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// ProductSales merepresentasikan data penjualan produk (untuk report).
type ProductSales struct {
	Name       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

// SalesSummary adalah response untuk endpoint report harian.
type SalesSummary struct {
	TotalRevenue   int          `json:"total_revenue"`
	TotalTransaksi int          `json:"total_transaksi"`
	ProdukTerlaris ProductSales `json:"produk_terlaris"`
}
