package main

import (
	"fmt"
	"log"
	"net/http"
)

// main adalah function utama yang jalan pas program di-run
// Di sini kita cuma setup route dan start server, simple aja
func main() {
	// Daftarin route /categories ke function handleCategories
	// Route ini handle GET semua kategori dan POST kategori baru
	http.HandleFunc("/categories", handleCategories)

	// Route dengan slash di akhir untuk handle /categories/{id}
	// Ini buat GET by ID, PUT update, dan DELETE
	http.HandleFunc("/categories/", handleCategories)

	// Kasih tau user kalo server udah jalan
	fmt.Println("Server running on http://localhost:8080")

	// Start server di port 8080
	// Kalo error, program langsung stop (Fatal)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
