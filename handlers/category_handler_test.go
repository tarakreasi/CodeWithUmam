package handlers

import (
	"codeWithUmam/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockCategoryService untuk testing handler
type MockCategoryService struct {
	GetAllFunc  func() ([]models.Category, error)
	CreateFunc  func(category *models.Category) error
	GetByIDFunc func(id int) (*models.Category, error)
	UpdateFunc  func(category *models.Category) error
	DeleteFunc  func(id int) error
}

func (m *MockCategoryService) GetAll() ([]models.Category, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockCategoryService) Create(category *models.Category) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(category)
	}
	return nil
}

func (m *MockCategoryService) GetByID(id int) (*models.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockCategoryService) Update(category *models.Category) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(category)
	}
	return nil
}

func (m *MockCategoryService) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestCategoryHandler_GetAll(t *testing.T) {
	// Setup Mock
	mockService := &MockCategoryService{
		GetAllFunc: func() ([]models.Category, error) {
			return []models.Category{
				{ID: 1, Name: "Makanan"},
			}, nil
		},
	}
	handler := NewCategoryHandler(mockService)

	// Create Request
	req, err := http.NewRequest("GET", "/api/v1/categories", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create Response Recorder
	rr := httptest.NewRecorder()

	// Call Handler
	handler.HandleCategories(rr, req)

	// Check Status Code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check Response Body
	// Trim newline from recorder body if any
	// Simple check (JSON string matching can be brittle, but ok for simple test)
	// Better to decode and check content, but let's check non-empty for now
	var response map[string][]models.Category
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(response["data"]) != 1 {
		t.Errorf("expected 1 category, got %d", len(response["data"]))
	}
	if response["data"][0].Name != "Makanan" {
		t.Errorf("expected name Makanan, got %s", response["data"][0].Name)
	}
}
