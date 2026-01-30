package repositories

import (
	"codeWithUmam/models"
	"database/sql"
)

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepositoryImpl) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES (?, ?)"
	result, err := r.db.Exec(query, category.Name, category.Description)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = int(id)
	return nil
}

func (r *CategoryRepositoryImpl) GetByID(id int) (*models.Category, error) {
	var c models.Category
	query := "SELECT id, name, description FROM categories WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepositoryImpl) Update(category *models.Category) error {
	query := "UPDATE categories SET name = ?, description = ? WHERE id = ?"
	_, err := r.db.Exec(query, category.Name, category.Description, category.ID)
	return err
}

func (r *CategoryRepositoryImpl) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
