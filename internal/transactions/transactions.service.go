package transactions

import "finance-tracker/internal/models"

type TransactionsServiceInterface interface {
	CreateTransaction(dto *TransactionRequestDto) (*TransactionResponseDto, error)
}

type TransactionsService struct {
	TransactionsRepository TransactionsRepositoryInterface
}

func NewTransactionsService(transactionsRepository TransactionsRepositoryInterface) *TransactionsService {
	return &TransactionsService{
		TransactionsRepository: transactionsRepository,
	}
}

func (service *TransactionsService) CreateTransaction(dto *TransactionRequestDto) (*TransactionResponseDto, error) {
	if err := dto.Type.Validate(); err != nil {
		return nil, err
	}

	transactionEntity := &models.Transaction{
		UserID:      dto.UserID,
		Type:        dto.Type,
		Amount:      dto.Amount,
		Description: dto.Description,
		Date:        dto.Date,
	}

	_, err := service.TransactionsRepository.Create(transactionEntity)
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
