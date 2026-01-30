package models

// Category merepresentasikan data kategori di dalam sistem.
type Category struct {
	// ID unik kategori.
	// Tag `json:"id"` berarti saat diubah jadi JSON (API response), field ini akan bernama "id".
	ID int `json:"id"`

	// Nama kategori.
	Name string `json:"name"`

	// Deskripsi singkat kategori.
	Description string `json:"description"`
}
