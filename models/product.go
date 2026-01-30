package models

// Product merepresentasikan data produk di dalam sistem.
type Product struct {
	// ID unik produk.
	// Tag `json:"id"` berarti saat diubah jadi JSON (API response), field ini akan bernama "id".
	ID int `json:"id"`

	// Nama produk.
	Name string `json:"name"`

	// Harga produk dalam integer (Rupiah tidak punya desimal penting).
	Price int `json:"price"`

	// Jumlah stok tersedia.
	Stock int `json:"stock"`

	// Foreign Key: ID dari kategori produk ini.
	CategoryID int `json:"category_id"`

	// Category adalah relasi (join).
	// Pointer (*) berarti field ini bisa bernilai nil (kosong) jika tidak ada datanya.
	// `omitempty`: Field ini tidak akan muncul di JSON jika nil.
	Category *Category `json:"category,omitempty"`
}
