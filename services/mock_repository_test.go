package services

import (
	"codeWithUmam/models"
	"errors"
)

// MockCategoryRepository implements repositories.CategoryRepository for testing
type MockCategoryRepository struct {
	GetAllFunc  func() ([]models.Category, error)
	CreateFunc  func(category *models.Category) error
	GetByIDFunc func(id int) (*models.Category, error)
	UpdateFunc  func(category *models.Category) error
	DeleteFunc  func(id int) error
}

func (m *MockCategoryRepository) GetAll() ([]models.Category, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockCategoryRepository) Create(category *models.Category) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(category)
	}
	return nil
}

func (m *MockCategoryRepository) GetByID(id int) (*models.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, errors.New("not found")
}

func (m *MockCategoryRepository) Update(category *models.Category) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(category)
	}
	return nil
}

func (m *MockCategoryRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
