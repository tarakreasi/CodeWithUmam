package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"codeWithUmam/database"
	"codeWithUmam/handlers"
	"codeWithUmam/repositories"
	"codeWithUmam/services"

	"github.com/spf13/viper"
)

// Config adalah struct untuk menyimpan konfigurasi aplikasi.
// Tag `mapstructure` digunakan oleh Viper untuk memetakan key dari .env atau environment variable.
type Config struct {
	Port   string `mapstructure:"PORT"`    // Port dimana server akan berjalan
	DBConn string `mapstructure:"DB_CONN"` // String koneksi database (untuk SQLite path filenya)
}

func main() {
	// ==========================================
	// 1. Initialisasi Konfigurasi (Viper)
	// ==========================================
	// Viper membantu kita membaca konfigurasi dari file .env atau environment variable.
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Cek apakah file .env ada? Jika ada, kita baca isinya.
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Error membaca file config:", err)
		}
	}

	// Masukkan nilai config ke struct agar mudah diakses
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// ==========================================
	// 2. Setup Database
	// ==========================================
	// Kita inisialisasi koneksi ke database SQLite.
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Gagal menginisialisasi database:", err)
	}
	// Pastikan koneksi database ditutup ketika aplikasi berhenti.
	defer db.Close()

	// ==========================================
	// 3. Setup Dependencies (Layered Architecture)
	// ==========================================
	// Di sini kita merakit aplikasi kita seperti tumpukan lego (Dependency Injection).
	// Urutannya: Repository (Data) -> Service (Logic) -> Handler (HTTP)

	// Setup Category
	categoryRepo := repositories.NewCategoryRepository(db)          // Layer Data: butuh koneksi DB
	categoryService := services.NewCategoryService(categoryRepo)    // Layer Logic: butuh Repository
	categoryHandler := handlers.NewCategoryHandler(categoryService) // Layer HTTP: butuh Service

	// Setup Product
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// ==========================================
	// 4. Setup Routes
	// ==========================================
	// Kita daftarkan alamat URL (endpoint) ke handler yang sesuai.
	// Prefix /api/v1 digunakan untuk versioning API (praktek yang baik).

	// Routes untuk Categories
	http.HandleFunc("/api/v1/categories", categoryHandler.HandleCategories)  // Match persis
	http.HandleFunc("/api/v1/categories/", categoryHandler.HandleCategories) // Match dengan suffix ID

	// Routes untuk Products
	http.HandleFunc("/api/v1/products", productHandler.HandleProducts)
	http.HandleFunc("/api/v1/products/", productHandler.HandleProducts)

	// Health Check - Endpoint sederhana untuk mengecek aplikasi hidup atau mati
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json := `{"status":"OK","message":"API Running"}`
		w.Write([]byte(json))
	})

	// ==========================================
	// 5. Start Server
	// ==========================================
	addr := ":" + config.Port
	fmt.Println("Server berjalan di http://localhost" + addr)

	// ListenAndServe akan menjalankan web server.
	// Jika terjadi error fatal (misal port sudah terpakai), aplikasi akan berhenti.
	log.Fatal(http.ListenAndServe(addr, nil))
}
