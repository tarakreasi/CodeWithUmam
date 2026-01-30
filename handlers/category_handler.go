package handlers

import (
	"codeWithUmam/models"
	"codeWithUmam/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&category); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
