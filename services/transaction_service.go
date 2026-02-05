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
func (s *TransactionServiceImpl) Checkout(req models.CheckoutRequest) (*models.Transaction, error) {
	// Logic Hitung Total (Idealnya logic ini ada sebelum create transaction, tapi untuk simple approach kita hitung estimasi dulu atau biarkan repo yang handle total real)
	// Masalah: Kita butuh Total Amount Real untuk hitung Change.
	// Solusi: Kita minta Repo hitung dulu (dry run) atau kita lakukan kalkulasi di service (ini butuh akses harga produk).
	// Alternatif Cepat: Kita percaya bahwa Repo akan melakukan validasi akhir. Tapi Service harus bisa hitung Change.
	// AGAR AMAN & SESUAI PLAN:
	// Repository CreateTransaction akan kita ubah sedikit flow-nya: Hitung total -> Validasi Bayar -> Baru Insert.
	// Di sini Service pass control ke Repo.

	// Tapi tunggu, repo CreateTransaction sekarang terima (paidAmount, change). Repo tidak tahu totalnya berapa SEBELUM query harga produk.
	// Jadi Service Idealnya:
	// 1. Get Product Prices (Repo.GetPrices?) --> Ribet
	// 2. Repo.CreateTransaction yang menghitung kembalian.

	// Mari kita serahkan logika "Hitung Kembalian" ke Repository saja, karena Repository yang pegang Data Harga (Single Source of Truth).
	// Kita kirim paidAmount ke Repo. Repo yang akan return error jika kurang, dan repo yang akan hitung change.

	// Namuun, signature Repo tadi: CreateTransaction(items, paidAmount, change, method)
	// Berarti Repo dikasih tau "Change-nya sekian". Ini aneh kalau Repo belum tau Totalnya.

	// KOREKSI APPROACH:
	// Biarkan Repository yang menghitung Change. Kita pass 0 sebagai dummy change, nanti repo update.
	// ATAU Service hitung Change? Service gak tau harga.
	// OK, Repo yang hitung Change. Signature Repo tadi sudah terlanjur `change int`.
	// Mari kita sesuaikan Repo implementation dikit nanti atau kita pass 0 di sini, dan Repo yang override.

	// Mari kita update Repository Implementation LAGI nanti biar dia yang hitung Change (Total - Paid).
	// Untuk sekarang, kita panggil Repo.

	return s.repo.CreateTransaction(req.Items, req.PaidAmount, 0, req.PaymentMethod)
	// Note: Repo akan kita update lagi untuk menghitung Change otomatis.
}

// GetDailyReport mengambil rekap laporan harian.
func (s *TransactionServiceImpl) GetDailyReport() (*models.SalesSummary, error) {
	return s.repo.GetDailySalesSummary()
}

func (s *TransactionServiceImpl) GetHistory(start, end string) ([]models.Transaction, error) {
	return s.repo.FindAll(start, end)
}

func (s *TransactionServiceImpl) GetDetail(id int) (*models.Transaction, error) {
	return s.repo.FindByID(id)
}
