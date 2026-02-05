package handlers

import (
	"bytes"
	"codeWithUmam/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockTransactionService untuk testing handler
// Mock ini adalah tiruan dari Service asli.
// Gunanya: Agar kita bisa test Handler TANPA harus konek ke Database beneran.
// Kita bisa "mengatur" agar mock ini me-return sukses atau error sesuai keinginan kita.
type MockTransactionService struct {
	CheckoutFunc       func(items []models.CheckoutItem) (*models.Transaction, error)
	GetDailyReportFunc func() (*models.SalesSummary, error)
}

func (m *MockTransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	if m.CheckoutFunc != nil {
		// Jika fungsi tiruan diedefinisikan di test, panggil
		return m.CheckoutFunc(items)
	}
	// Default return
	return nil, errors.New("not implemented")
}

func (m *MockTransactionService) GetDailyReport() (*models.SalesSummary, error) {
	if m.GetDailyReportFunc != nil {
		return m.GetDailyReportFunc()
	}
	return nil, nil
}

func TestTransactionHandler_HandleCheckout_Success(t *testing.T) {
	// 1. Setup Mock
	// Kita pura-pura service akan sukses memproses checkout
	mockService := &MockTransactionService{
		CheckoutFunc: func(items []models.CheckoutItem) (*models.Transaction, error) {
			return &models.Transaction{
				ID:          1,
				TotalAmount: 50000,
				Details: []models.TransactionDetail{
					{ProductID: 1, Quantity: 2, Subtotal: 50000},
				},
			}, nil
		},
	}
	handler := NewTransactionHandler(mockService)

	// 2. Create Request (Simulasi Panggilan HTTP)
	// Kita buat body JSON request palsu
	payload := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	body, _ := json.Marshal(payload)
	// http.NewRequest membuat objek request tanpa melakukan network call beneran
	req, err := http.NewRequest("POST", "/api/checkout", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// 3. Response Recorder (Perekam Respon)
	// Ini berfungsi sebagai papan tulis untuk menangkap apa yang ditulis oleh handler (w http.ResponseWriter)
	rr := httptest.NewRecorder()

	// 4. Execute Handler (Panggil fungsi asli)
	handler.HandleCheckout(rr, req)

	// 5. Verifikasi Hasil (Assert)
	// Cek status codenya 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Cek isi Body response
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Pastikan data total_amount sesuai dengan yang kita mock (50000)
	data := response["data"].(map[string]interface{})
	if data["total_amount"].(float64) != 50000 {
		t.Errorf("expected total_amount 50000, got %v", data["total_amount"])
	}
}

func TestTransactionHandler_HandleDailyReport_Success(t *testing.T) {
	mockService := &MockTransactionService{
		GetDailyReportFunc: func() (*models.SalesSummary, error) {
			return &models.SalesSummary{
				TotalRevenue:   100000,
				TotalTransaksi: 5,
				ProdukTerlaris: models.ProductSales{Name: "Kopi", QtyTerjual: 10},
			}, nil
		},
	}
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest("GET", "/api/report/hari-ini", nil)
	rr := httptest.NewRecorder()

	handler.HandleDailyReport(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
