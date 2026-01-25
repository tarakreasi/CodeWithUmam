package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Variabel untuk nyimpen data kategori di memory (RAM)
// Data bakal ilang pas server restart, tapi gpp buat belajar
var categories []Category

// nextID adalah counter untuk ID kategori berikutnya
// Setiap bikin kategori baru, ID ini auto naik
var nextID = 1

// sendJSON adalah helper function buat kirim response JSON
// Parameter data bisa apa aja (interface{}), nanti di-convert ke JSON
func sendJSON(w http.ResponseWriter, data interface{}) {
	// Set header biar browser tau ini JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode data jadi JSON dan kirim ke client
	// Format: {"data": isi_data_kamu}
	json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
}

// sendError adalah helper function buat kirim pesan error
// Bikin lebih gampang kirim response error dengan format konsisten
func sendError(w http.ResponseWriter, message string, code int) {
	// Set header biar browser tau ini JSON
	w.Header().Set("Content-Type", "application/json")

	// Set status code (misal 404, 400, dll)
	w.WriteHeader(code)

	// Kirim pesan error dalam format JSON
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// handleCategories adalah router utama untuk semua endpoint kategori
// Function ini yang nge-handle semua request ke /categories
func handleCategories(w http.ResponseWriter, r *http.Request) {
	// Cek kalo path nya persis "/categories" (tanpa ID)
	if r.URL.Path == "/categories" {
		// Cek method HTTP apa yang dipake
		switch r.Method {
		case "GET":
			// GET /categories -> ambil semua data
			getCategories(w, r)
		case "POST":
			// POST /categories -> bikin kategori baru
			createCategory(w, r)
		default:
			// Method lain ditolak
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Kalo path nya "/categories/{id}" (ada ID nya)
	if strings.HasPrefix(r.URL.Path, "/categories/") {
		switch r.Method {
		case "GET":
			// GET /categories/1 -> ambil satu data by ID
			getCategoryByID(w, r)
		case "PUT":
			// PUT /categories/1 -> update data by ID
			updateCategory(w, r)
		case "DELETE":
			// DELETE /categories/1 -> hapus data by ID
			deleteCategory(w, r)
		default:
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Kalo gak match route mana pun, kirim 404
	sendError(w, "Not found", http.StatusNotFound)
}

// getCategories ngambil semua data kategori
// Simple aja, cuma balikin array categories
func getCategories(w http.ResponseWriter, r *http.Request) {
	// Kirim semua isi array categories
	sendJSON(w, categories)
}

// createCategory bikin kategori baru dari data JSON yang dikirim user
// Data dari request body diparsing jadi struct Category
func createCategory(w http.ResponseWriter, r *http.Request) {
	// Bikin variable kosong buat nampung input dari user
	var input Category

	// Parse JSON dari request body ke variable input
	// Kalo gagal (misal JSON nya salah format), kirim error
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set ID otomatis (user gak perlu kirim ID)
	input.ID = nextID

	// Naikin counter buat kategori berikutnya
	nextID++

	// Tambahin kategori baru ke array
	categories = append(categories, input)

	// Kirim balik data yang baru dibuat
	sendJSON(w, input)
}

// getCategoryByID ngambil satu kategori berdasarkan ID
// ID diambil dari URL, misal /categories/1
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL path
	// Misal path "/categories/1" -> idStr jadi "1"
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	// Convert string ke integer
	// Kalo gagal (misal "/categories/abc"), kirim error
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Loop semua kategori, cari yang ID nya cocok
	for _, category := range categories {
		if category.ID == id {
			// Ketemu! Kirim data kategori ini
			sendJSON(w, category)
			return
		}
	}

	// Kalo gak ketemu, kirim 404
	sendError(w, "Category not found", http.StatusNotFound)
}

// updateCategory update data kategori yang udah ada
// Butuh ID di URL dan data baru di request body
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL (sama kayak getCategoryByID)
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Parse data baru dari request body
	var input Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cari kategori yang mau di-update
	for i, category := range categories {
		if category.ID == id {
			// Ketemu! Update Name dan Description nya
			// ID nya gak di-update, tetep sama
			categories[i].Name = input.Name
			categories[i].Description = input.Description

			// Kirim balik data yang udah di-update
			sendJSON(w, categories[i])
			return
		}
	}

	// Kalo gak ketemu kategori yang mau di-update
	sendError(w, "Category not found", http.StatusNotFound)
}

// deleteCategory hapus kategori berdasarkan ID
// Kategori yang dihapus bakal ilang dari array
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Cari kategori yang mau dihapus
	for i, category := range categories {
		if category.ID == id {
			// Ketemu! Hapus dari array dengan cara:
			// Gabungin bagian sebelum index i dan setelah index i
			// categories[:i] = semua sebelum i
			// categories[i+1:] = semua setelah i
			categories = append(categories[:i], categories[i+1:]...)

			// Kirim response sukses
			sendJSON(w, true)
			return
		}
	}

	// Kalo gak ketemu kategori yang mau dihapus
	sendError(w, "Category not found", http.StatusNotFound)
}
