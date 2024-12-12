package transactions

import (
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/internal/models"
	"finance-tracker/pkg/database"
	"fmt"

	"gorm.io/gorm"
)

type TransactionsRepositoryInterface interface {
	Create(transaction *models.Transaction) (*models.Transaction, error)
	Delete(id uint) error
	FindTransactionById(id uint) error
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

func (repo *TransactionsRepository) Delete(id uint) error {

	result := repo.Database.DB.Delete(&models.Transaction{}, id)

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: id = %d", errs.ErrTransactionNotFound, id)
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *TransactionsRepository) FindTransactionById(id uint) error {
	if err := repo.Database.DB.First(&models.Transaction{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: id = %d", errs.ErrTransactionNotFound, id)
		}
		return err
	}
	return nil
}
