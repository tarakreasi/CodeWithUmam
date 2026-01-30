package repositories

import (
	"codeWithUmam/models"
	"database/sql"
)

// CategoryRepositoryImpl bertugas melakukan komunikasi langsung ke Database.
// Semua Query SQL (SELECT, INSERT, UPDATE, DELETE) ada di sini.
// Struct ini mengimplementasikan interface CategoryRepository dari package repositories.
type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db: db}
}

// GetAll mengambil semua baris data dari tabel categories.
func (r *CategoryRepositoryImpl) GetAll() ([]models.Category, error) {
	// Menjalankan query SELECT standar.
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	// Pastikan koneksi row ditutup setelah fungsi selesai agar tidak memory leak.
	defer rows.Close()

	var categories []models.Category
	// Loop setiap baris hasil query (Next)
	for rows.Next() {
		var c models.Category
		// Scan: Memindahkan data dari database ke variabel struct Go.
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		// Masukkan ke slice (array dinamis)
		categories = append(categories, c)
	}
	return categories, nil
}

// Create menyimpan data kategori baru ke database.
func (r *CategoryRepositoryImpl) Create(category *models.Category) error {
	// Query INSERT dengan placeholder (?) untuk mencegah SQL Injection.
	query := "INSERT INTO categories (name, description) VALUES (?, ?)"

	// Exec: Menjalankan query yang mengubah data (tidak mengembalikan baris data).
	result, err := r.db.Exec(query, category.Name, category.Description)
	if err != nil {
		return err
	}

	// Ambil ID yang baru saja digenerate oleh database (AUTOINCREMENT).
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Update ID di struct category agar pemanggil fungsi tau ID barunya.
	category.ID = int(id)
	return nil
}

// GetByID mengambil satu kategori berdasarkan ID.
func (r *CategoryRepositoryImpl) GetByID(id int) (*models.Category, error) {
	var c models.Category
	query := "SELECT id, name, description FROM categories WHERE id = ?"

	// QueryRow: Untuk mengambil 1 baris data saja.
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Update mengubah data kategori yang sudah ada.
func (r *CategoryRepositoryImpl) Update(category *models.Category) error {
	query := "UPDATE categories SET name = ?, description = ? WHERE id = ?"
	_, err := r.db.Exec(query, category.Name, category.Description, category.ID)
	return err
}

// Delete menghapus kategori dari database.
func (r *CategoryRepositoryImpl) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
