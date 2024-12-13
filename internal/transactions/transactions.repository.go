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
	Update(id uint, dto *TransactionUpdateRequestDto) (*models.Transaction, error)
	FindTransactionById(id uint) error
	GetTransactionsByUserId(userId uint) ([]models.Transaction, error)
	GetTransactionById(userId, transactionId uint) (*models.Transaction, error)
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

func (repo *TransactionsRepository) Update(id uint, dto *TransactionUpdateRequestDto) (*models.Transaction, error) {
	var transaction models.Transaction

	if err := repo.Database.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrTransactionNotFound
		}
		return nil, err
	}

	// Пока нет необходимости
	// updates := buildUpdates(dto)

	if err := repo.Database.DB.Model(&transaction).Updates(dto).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
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

func (repo *TransactionsRepository) GetTransactionsByUserId(userId uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := repo.Database.DB.Where("user_id = ?", userId).Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (repo *TransactionsRepository) GetTransactionById(userId, id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := repo.Database.DB.Where("user_id = ? AND id = ?", userId, id).First(&transaction)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: id = %d", errs.ErrTransactionNotFound, id)
		}
	}

	return &transaction, nil
}

// func buildUpdates(dto *TransactionUpdateRequestDto) map[string]interface{} {
// 	updates := make(map[string]interface{})

// 	fields := map[string]interface{}{
// 		"type":        dto.Type,
// 		"amount":      dto.Amount,
// 		"description": dto.Description,
// 		"date":        dto.Date,
// 	}

// 	for key, value := range fields {
// 		if value != nil {
// 			updates[key] = reflect.ValueOf(value).Elem().Interface()
// 		}
// 	}

// 	return updates
// }
