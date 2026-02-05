package repositories

import (
	"codeWithUmam/models"
	"database/sql"
)

// ProductRepositoryImpl bertugas melakukan komunikasi langsung ke Database.
// Semua Query SQL (SELECT, INSERT, UPDATE, DELETE) ada di sini.
// Struct ini mengimplementasikan interface ProductRepository dari package repositories.
type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

// GetAll mengambil semua baris data dari tabel products.
// Jika parameter name tidak kosong, akan dilakukan filter search by name.
func (r *ProductRepositoryImpl) GetAll(name string) ([]models.Product, error) {
	query := "SELECT id, name, price, stock, category_id FROM products"
	args := []interface{}{}

	// Jika ada filter nama, tambahkan WHERE clause
	// Kita pakai LIKE untuk pencarian partial (misal: "indom" -> "Indomie")
	if name != "" {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	// Masukkan args... (spread operator) ke dalam Query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err // Kembalikan error jika query gagal
	}
	// Pastikan koneksi row ditutup setelah fungsi selesai agar tidak memory leak.
	defer rows.Close()

	var products []models.Product

	// Loop setiap baris hasil query (Next)
	for rows.Next() {
		var p models.Product
		// Scan: Memindahkan data dari database ke variabel struct Go.
		// Urutan Scan HARUS SAMA dengan urutan SELECT di atas.
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		// Masukkan ke slice (array dinamis)
		products = append(products, p)
	}
	return products, nil
}

// Create menyimpan data produk baru ke database.
func (r *ProductRepositoryImpl) Create(product *models.Product) error {
	// Query INSERT. Tanda tanya (?) adalah placeholder untuk mencegah SQL Injection.
	query := "INSERT INTO products (name, price, stock, category_id) VALUES (?, ?, ?, ?)"

	// Exec: Menjalankan query yang mengubah data (tidak mengembalikan baris data).
	result, err := r.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID)
	if err != nil {
		return err
	}

	// Ambil ID yang baru saja digenerate oleh database (AUTOINCREMENT).
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Update ID di struct product agar pemanggil fungsi tau ID barunya.
	product.ID = int(id)
	return nil
}

// GetByID mengambil satu produk dan DETAIL KATEGORINYA menggunakan JOIN.
func (r *ProductRepositoryImpl) GetByID(id int) (*models.Product, error) {
	// Query JOIN: Menggabungkan tabel products (p) dan categories (c).
	// LEFT JOIN: Ambil produk meskipun kategori-nya tidak ada.
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?`

	var p models.Product
	var c models.Category

	// QueryRow: Untuk mengambil 1 baris data saja.
	// Kita scan kolom produk ke struct p, dan kolom kategori ke struct c.
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
		&c.ID, &c.Name, &c.Description,
	)
	if err != nil {
		return nil, err
	}

	// Masukkan struct category ke dalam struct product (Nested Struct).
	p.Category = &c
	return &p, nil
}

// Update mengubah data produk yang sudah ada.
func (r *ProductRepositoryImpl) Update(product *models.Product) error {
	query := "UPDATE products SET name = ?, price = ?, stock = ?, category_id = ? WHERE id = ?"
	_, err := r.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	return err
}

// Delete menghapus produk dari database.
func (r *ProductRepositoryImpl) Delete(id int) error {
	query := "DELETE FROM products WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
