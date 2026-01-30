package repositories

import (
	"codeWithUmam/models"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	query := `
	CREATE TABLE categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestCategoryRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)

	category := &models.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if category.ID == 0 {
		t.Error("ID should be set after create")
	}
}

func TestCategoryRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewCategoryRepository(db)

	repo.Create(&models.Category{Name: "C1", Description: "D1"})
	repo.Create(&models.Category{Name: "C2", Description: "D2"})

	categories, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

func TestCategoryRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewCategoryRepository(db)

	cat := &models.Category{Name: "C1", Description: "D1"}
	repo.Create(cat)

	res, err := repo.GetByID(cat.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if res.Name != "C1" {
		t.Errorf("expected name C1, got %s", res.Name)
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewCategoryRepository(db)

	cat := &models.Category{Name: "C1", Description: "D1"}
	repo.Create(cat)

	cat.Name = "Updated C1"
	err := repo.Update(cat)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	res, _ := repo.GetByID(cat.ID)
	if res.Name != "Updated C1" {
		t.Errorf("expected updated name, got %s", res.Name)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewCategoryRepository(db)

	cat := &models.Category{Name: "C1", Description: "D1"}
	repo.Create(cat)

	err := repo.Delete(cat.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(cat.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}
