package services

import (
	"codeWithUmam/models"
	"codeWithUmam/repositories"
)

type TransactionServiceImpl struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{repo: repo}
}

func (s *TransactionServiceImpl) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	// Di masa depan, logic seperti perhitungan diskon atau validasi member bisa ditaruh di sini.
	return s.repo.CreateTransaction(items)
}

func (s *TransactionServiceImpl) GetDailyReport() (*models.SalesSummary, error) {
	return s.repo.GetDailySalesSummary()
}
