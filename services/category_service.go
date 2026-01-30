package services

import (
	"codeWithUmam/models"
	"codeWithUmam/repositories"
)

type CategoryServiceImpl struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{repo: repo}
}

func (s *CategoryServiceImpl) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryServiceImpl) Create(category *models.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryServiceImpl) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryServiceImpl) Update(category *models.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}
