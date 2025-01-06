package services

import (
	"squad-checkout/internal/models"
	"squad-checkout/internal/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) StoreTransaction(transaction models.Transaction) error {
	return s.repo.Save(transaction)
}

func (s *TransactionService) RetrieveTransaction(id string) (*models.Transaction, error) {
	return s.repo.FindByID(id)
}

func (s *TransactionService) RetrieveAllTransactions() ([]models.Transaction, error) {
	return s.repo.FindAll()
}
