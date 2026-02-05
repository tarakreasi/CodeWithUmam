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
type MockTransactionService struct {
	CheckoutFunc       func(items []models.CheckoutItem) (*models.Transaction, error)
	GetDailyReportFunc func() (*models.SalesSummary, error)
}

func (m *MockTransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	if m.CheckoutFunc != nil {
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
	// Setup Mock
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

	// Create Request Payload
	payload := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/api/checkout", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Response Recorder
	rr := httptest.NewRecorder()

	// Call Handler
	handler.HandleCheckout(rr, req)

	// Verify Status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify Body
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

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
