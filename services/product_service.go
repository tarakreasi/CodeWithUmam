package services

import (
	"codeWithUmam/models"
	"codeWithUmam/repositories"
)

// ProductServiceImpl berisi Bisnis Logic aplikasi.
// Di sinilah tempat validasi data, kalkulasi, dll terjadi SEBELUM disimpan ke database.
// Saat ini isinya masih "pass-through" (langsung panggil repo), tapi nanti logic komplek ada di sini.
type ProductServiceImpl struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{repo: repo}
}

func (s *ProductServiceImpl) GetAll() ([]models.Product, error) {
	// Di sini bisa ditambah logic, misalnya: Filter data yang aktif saja, atau sort.
	return s.repo.GetAll()
}

func (s *ProductServiceImpl) Create(product *models.Product) error {
	// Contoh Bisnis Logic yang bisa ditambahkan:
	// if product.Price < 0 { return error("Harga tidak boleh minus") }
	return s.repo.Create(product)
}

func (s *ProductServiceImpl) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductServiceImpl) Update(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}
