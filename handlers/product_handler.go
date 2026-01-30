package handlers

import (
	"codeWithUmam/models"
	"codeWithUmam/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// ProductHandler bertanggung jawab menangani request HTTP terkait produk.
// Handler tidak boleh berisi logic bisnis, tugasnya hanya:
// 1. Menerima request (validasi input, baca body/params)
// 2. Memanggil Service (logic bisnis)
// 3. Mengembalikan response (JSON)
type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts adalah "router" sederhana di dalam handler ini.
// Ia menentukan fungsi mana yang dipanggil berdasarkan URL dan Method (GET/POST/dll).
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	// Setup CORS agar API bisa diakses dari browser/frontend yang berbeda domain.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Preflight request: Browser tanya "Boleh gak saya kirim request?".
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Jika URL persis "/api/v1/products" -> ini untuk GetAll atau Create.
	if r.URL.Path == "/api/v1/products" {
		switch r.Method {
		case "GET":
			h.GetAll(w, r)
		case "POST":
			h.Create(w, r)
		default:
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Jika URL diawali "/api/v1/products/" -> berarti ada ID di belakangnya (misal /api/v1/products/123).
	// Ini untuk GetByID, Update, atau Delete.
	if strings.HasPrefix(r.URL.Path, "/api/v1/products/") {
		switch r.Method {
		case "GET":
			h.GetByID(w, r)
		case "PUT":
			h.Update(w, r)
		case "DELETE":
			h.Delete(w, r)
		default:
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	sendError(w, "Not found", http.StatusNotFound)
}

// GetAll mengambil semua data produk.
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Panggil service untuk ambil data
	products, err := h.service.GetAll()
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Kirim response sukses (Utility function sendJSON ada di utils.go)
	sendJSON(w, products)
}

// Create membuat produk baru.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	// Decode JSON dari request body ke struct Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil service untuk simpan data
	if err := h.service.Create(&product); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Kembalikan data yang baru dibuat (lengkap dengan ID baru)
	sendJSON(w, product)
}

// GetByID mengambil satu produk berdasarkan ID di URL.
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL (potong prefix path-nya)
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		sendError(w, "Product not found", http.StatusNotFound)
		return
	}
	sendJSON(w, product)
}

// Update mengubah data produk yang sudah ada.
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Pastikan ID di struct sama dengan ID di URL
	product.ID = id

	if err := h.service.Update(&product); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, product)
}

// Delete menghapus produk berdasarkan ID.
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		sendError(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	sendJSON(w, true)
}
