package handlers

import (
	"codeWithUmam/models"
	"codeWithUmam/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// CategoryHandler bertanggung jawab menangani request HTTP terkait kategori.
// Handler tidak boleh berisi logic bisnis, tugasnya hanya:
// 1. Menerima request (validasi input, baca body/params)
// 2. Memanggil Service (logic bisnis)
// 3. Mengembalikan response (JSON)
type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleCategories adalah "router" sederhana di dalam handler ini.
// Ia menentukan fungsi mana yang dipanggil berdasarkan URL dan Method (GET/POST/dll).
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	// Setup CORS agar API bisa diakses dari browser/frontend yang berbeda domain.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Preflight request: Browser tanya "Boleh gak saya kirim request?".
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Jika URL persis "/api/v1/categories" -> ini untuk GetAll atau Create.
	if r.URL.Path == "/api/v1/categories" {
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

	// Jika URL diawali "/api/v1/categories/" -> berarti ada ID di belakangnya (misal /api/v1/categories/123).
	// Ini untuk GetByID, Update, atau Delete.
	if strings.HasPrefix(r.URL.Path, "/api/v1/categories/") {
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

// @Summary Get all categories
// @Description Get list of all categories
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Category
// @Router /categories [get]
// GetAll mengambil semua data kategori.
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Panggil service untuk ambil data
	categories, err := h.service.GetAll()
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Kirim response sukses (Utility function sendJSON ada di utils.go)
	sendJSON(w, categories)
}

// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param category body models.Category true "Category Data"
// @Success 200 {object} models.Category
// @Router /categories [post]
// Create membuat kategori baru.
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	// Decode JSON dari request body ke struct Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil service untuk simpan data
	if err := h.service.Create(&category); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Kembalikan data yang baru dibuat (lengkap dengan ID baru)
	sendJSON(w, category)
}

// @Summary Get category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Router /categories/{id} [get]
// GetByID mengambil satu kategori berdasarkan ID di URL.
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL (potong prefix path-nya)
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		sendError(w, "Category not found", http.StatusNotFound)
		return
	}
	sendJSON(w, category)
}

// @Summary Update a category
// @Description Update an existing category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Category Data"
// @Success 200 {object} models.Category
// @Router /categories/{id} [put]
// Update mengubah data kategori yang sudah ada.
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Pastikan ID di struct sama dengan ID di URL
	category.ID = id

	if err := h.service.Update(&category); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, category)
}

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Success 200 {boolean} true
// @Router /categories/{id} [delete]
// Delete menghapus kategori berdasarkan ID.
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		sendError(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}
	sendJSON(w, true)
}
