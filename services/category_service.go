package services

import (
	"codeWithUmam/models"
	"codeWithUmam/repositories"
)

// CategoryServiceImpl berisi Bisnis Logic aplikasi untuk kategori.
// Di sinilah tempat validasi data, kalkulasi, dll terjadi SEBELUM disimpan ke database.
// Saat ini isinya masih "pass-through" (langsung panggil repo), tapi nanti logic komplek ada di sini.
type CategoryServiceImpl struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{repo: repo}
}

func (s *CategoryServiceImpl) GetAll() ([]models.Category, error) {
	// Di sini bisa ditambah logic, misalnya: Filter data yang aktif saja, atau sort.
	return s.repo.GetAll()
}

func (s *CategoryServiceImpl) Create(category *models.Category) error {
	// Contoh Bisnis Logic yang bisa ditambahkan:
	// if category.Name == "" { return error("Nama kategori tidak boleh kosong") }
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
