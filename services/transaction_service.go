package services

import (
	"codeWithUmam/models"
	"codeWithUmam/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	// Di masa depan, logic seperti perhitungan diskon atau validasi member bisa ditaruh di sini.
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetDailyReport() (*models.SalesSummary, error) {
	return s.repo.GetDailySalesSummary()
}
