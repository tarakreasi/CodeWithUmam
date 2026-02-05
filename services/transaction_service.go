package services

import (
	"codeWithUmam/models"
	"codeWithUmam/repositories"
)

// TransactionServiceImpl adalah implementasi dari interface TransactionService.
// Struct ini menjembatani antara Handler (HTTP) dan Repository (Database).
type TransactionServiceImpl struct {
	repo *repositories.TransactionRepository
}

// NewTransactionService adalah Constructor.
// Menerima dependency Repository (Dependency Injection).
func NewTransactionService(repo *repositories.TransactionRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{repo: repo}
}

// Checkout menangani logika pembelian.
// Saat ini hanya wrapper ke repository, tapi di sinilah tempatnya jika kita mau tambah logic bisnis.
// Contoh Logic Bisnis: Cek member, hitung diskon, validasi jam operasional, kirim email notifikasi, dll.
func (s *TransactionServiceImpl) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	// Di masa depan, logic seperti perhitungan diskon atau validasi member bisa ditaruh di sini.
	return s.repo.CreateTransaction(items)
}

// GetDailyReport mengambil rekap laporan harian.
func (s *TransactionServiceImpl) GetDailyReport() (*models.SalesSummary, error) {
	return s.repo.GetDailySalesSummary()
}
