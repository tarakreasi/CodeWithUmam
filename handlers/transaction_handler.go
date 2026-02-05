package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"codeWithUmam/models"
	"codeWithUmam/services"
)

// TransactionHandler menangani HTTP request terkait transaksi dan report.
// Seperti kasir yang melayani pelanggan langsung (menerima order, memberikan struk).
type TransactionHandler struct {
	// Kita bergantung pada Interface, bukan struct konkret.
	// Ini membuat code "Loosely Coupled" dan mudah di-test (Mocking).
	service services.TransactionService
}

// NewTransactionHandler adalah Constructor.
func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleCheckout menangani request pembelian barang.
// Endpoint: POST /api/checkout
// Body JSON: { "items": [ { "product_id": 1, "quantity": 2 } ] }
// @Summary      Checkout Transaction
// @Description  Create a new transaction with items and payment info
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body models.CheckoutRequest true "Checkout Request"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /checkout [post]
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	// Validasi Method: Hanya boleh POST
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parsing JSON Body: Mengubah JSON mentah dari request body menjadi struct Go.
	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil Service untuk proses checkout
	transaction, err := h.service.Checkout(req)
	if err != nil {
		// Jika ada error (misal stok habis), kirim response 500
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sukses: Kirim balik detail transaksi dalam format JSON
	sendJSON(w, transaction)
}

// HandleDailyReport menangani request laporan harian.
// Endpoint: GET /api/report/hari-ini
// Digunakan oleh Owner/Manajer untuk melihat omset hari ini.
// @Summary      Get Daily Report
// @Description  Get sales summary for today
// @Tags         transactions
// @Produce      json
// @Success      200  {object}  models.SalesSummary
// @Router       /report/hari-ini [get]
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

// HandleHistory menangani request daftar transaksi.
// Endpoint: GET /api/transactions
// Params: ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD (Optional)
// @Summary      Get Transaction History
// @Description  Get list of transactions with optional date filter
// @Tags         transactions
// @Produce      json
// @Param        start_date query string false "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string false "End Date (YYYY-MM-DD)"
// @Success      200  {array}   models.Transaction
// @Router       /transactions [get]
func (h *TransactionHandler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := r.URL.Query().Get("start_date")
	end := r.URL.Query().Get("end_date")

	transactions, err := h.service.GetHistory(start, end)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, transactions)
}

// HandleDetail menangani request detail satu transaksi.
// Endpoint: GET /api/transactions/{id}
// @Summary      Get Transaction Detail
// @Description  Get detailed transaction by ID
// @Tags         transactions
// @Produce      json
// @Param        id   path      int  true  "Transaction ID"
// @Success      200  {object}  models.Transaction
// @Failure      404  {object}  map[string]string
// @Router       /transactions/{id} [get]
func (h *TransactionHandler) HandleDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil ID dari URL Path manual (karena kita pakai net/http standard)
	// Asumsi URL: /api/transactions/123
	path := r.URL.Path
	// Split by slash
	// Contoh: ["", "api", "transactions", "123"]
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	// Validasi jika idStr kosong atau bukan angka
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid Transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.GetDetail(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if transaction == nil {
		sendError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	sendJSON(w, transaction)
}
