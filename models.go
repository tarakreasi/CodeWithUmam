package main

// Category adalah struct untuk menyimpan data kategori
// Struct ini simple aja, cuma ada 3 field:
// - ID: nomor unik kategori (otomatis dari sistem)
// - Name: nama kategorinya
// - Description: penjelasan tentang kategori ini
type Category struct {
	ID          int    `json:"id"`          // ID otomatis dari sistem
	Name        string `json:"name"`        // Nama kategori (wajib diisi)
	Description string `json:"description"` // Deskripsi kategori (opsional)
}
