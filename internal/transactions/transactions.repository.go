package transactions

import (
	"finance-tracker/internal/models"
	"finance-tracker/pkg/database"
)

type TransactionsRepositoryInterface interface {
	Create(transaction *models.Transaction) (*models.Transaction, error)
}

type TransactionsRepository struct {
	Database *database.Database
}

func NewTransactionRepository(database *database.Database) *TransactionsRepository {
	return &TransactionsRepository{
		Database: database,
	}
}

func (repo *TransactionsRepository) Create(transaction *models.Transaction) (*models.Transaction, error) {
	result := repo.Database.DB.Create(transaction)
	if result.Error != nil {
		return nil, result.Error
	}

	return transaction, nil
}
