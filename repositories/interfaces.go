package repositories

import "codeWithUmam/models"

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	Create(category *models.Category) error
	GetByID(id int) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id int) error
}

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Create(product *models.Product) error
	GetByID(id int) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id int) error
}
