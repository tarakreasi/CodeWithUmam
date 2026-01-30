package services

import (
	"codeWithUmam/models"
	"errors"
	"testing"
)

func TestCategoryService_GetAll(t *testing.T) {
	mockRepo := &MockCategoryRepository{
		GetAllFunc: func() ([]models.Category, error) {
			return []models.Category{
				{ID: 1, Name: "C1"},
				{ID: 2, Name: "C2"},
			}, nil
		},
	}
	service := NewCategoryService(mockRepo)

	cats, err := service.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 2 {
		t.Errorf("expected 2 cats, got %d", len(cats))
	}
}

func TestCategoryService_Create(t *testing.T) {
	mockRepo := &MockCategoryRepository{
		CreateFunc: func(c *models.Category) error {
			if c.Name == "" {
				return errors.New("empty name")
			}
			c.ID = 1
			return nil
		},
	}
	service := NewCategoryService(mockRepo)

	cat := &models.Category{Name: "New"}
	err := service.Create(cat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.ID != 1 {
		t.Error("expected ID to be set")
	}
}
