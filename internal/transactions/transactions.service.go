package transactions

import (
	"finance-tracker/internal/dto"
	"finance-tracker/internal/models"
	"finance-tracker/internal/user"
)

type TransactionsServiceInterface interface {
	GetUserTransactions(userId uint, filters TransactionsFilter, pagination dto.PaginationRequestDto) (dto.PaginationResponseDto[models.Transaction], error)
	GetUserTransaction(userId, transactionId uint) (*models.Transaction, error)
	CreateTransaction(userId uint, dto *TransactionRequestDto) (*TransactionResponseDto, error)
	DeleteTransaction(id uint) error
	UpdateTransaction(id uint, dto *TransactionUpdateRequestDto) (*models.Transaction, error)
}

type TransactionsService struct {
	TransactionsRepository TransactionsRepositoryInterface
	UserRepository         user.UserRepositoryInterface
}

func NewTransactionsService(transactionsRepository TransactionsRepositoryInterface, userRepository user.UserRepositoryInterface) *TransactionsService {
	return &TransactionsService{
		TransactionsRepository: transactionsRepository,
		UserRepository:         userRepository,
	}
}

func (service *TransactionsService) CreateTransaction(userId uint, dto *TransactionRequestDto) (*TransactionResponseDto, error) {
	_, err := service.UserRepository.FindById(userId)
	if err != nil {
		return nil, err
	}

	transactionEntity := &models.Transaction{
		UserID:      userId,
		Type:        dto.Type,
		Amount:      dto.Amount,
		Description: dto.Description,
		Date:        dto.Date,
	}

	_, err = service.TransactionsRepository.Create(transactionEntity)
	if err != nil {
		return nil, err
	}

	response := &TransactionResponseDto{
		ID:          transactionEntity.ID,
		UserID:      transactionEntity.UserID,
		Type:        transactionEntity.Type,
		Amount:      transactionEntity.Amount,
		Description: transactionEntity.Description,
		Date:        transactionEntity.Date,
	}

	return response, nil
}

func (service *TransactionsService) DeleteTransaction(id uint) error {
	if err := service.TransactionsRepository.Delete(id); err != nil {
		return err
	}
	return nil
}
func (service *TransactionsService) UpdateTransaction(id uint, dto *TransactionUpdateRequestDto) (*models.Transaction, error) {
	updatedTransaction, err := service.TransactionsRepository.Update(id, dto)
	if err != nil {
		return nil, err
	}
	return updatedTransaction, nil
}

func (service *TransactionsService) GetUserTransactions(userId uint, filters TransactionsFilter, pagination dto.PaginationRequestDto) (dto.PaginationResponseDto[models.Transaction], error) {
	transactions, totalCount, err := service.TransactionsRepository.GetTransactionsByUserId(userId, filters, pagination)
	if err != nil {
		return dto.PaginationResponseDto[models.Transaction]{}, err
	}

	response := dto.PaginationResponseDto[models.Transaction]{
		Items:      transactions,
		TotalCount: totalCount,
		Limit:      pagination.Limit,
		Offset:     pagination.Offset,
	}

	return response, nil
}

func (service *TransactionsService) GetUserTransaction(userId, transactionId uint) (*models.Transaction, error) {
	result, err := service.TransactionsRepository.GetTransactionById(userId, transactionId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
