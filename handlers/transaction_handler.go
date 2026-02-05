package handlers

import (
	"encoding/json"
	"net/http"

	"codeWithUmam/models"
	"codeWithUmam/services"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleCheckout menangani request pembelian barang.
// Endpoint: POST /api/checkout
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Menggunakan helper sendJSON yang sudah ada
	sendJSON(w, transaction)
}

// HandleDailyReport menangani request laporan harian.
// Endpoint: GET /api/report/hari-ini
func (h *TransactionHandler) HandleDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	summary, err := h.service.GetDailyReport()
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, summary)
}
