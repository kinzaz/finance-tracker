package transactions

import (
	"finance-tracker/internal/types"
	"time"
)

type TransactionRequestDto struct {
	Type        types.TransactionType `json:"type" validate:"required,oneof=income expense"`
	Amount      float64               `json:"amount" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Date        time.Time             `json:"date" validate:"required"`
}

type TransactionResponseDto struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	Type        types.TransactionType `json:"type"`
	Amount      float64               `json:"amount"`
	Description string                `json:"description"`
	Date        time.Time             `json:"date"`
}

type TransactionUpdateRequestDto struct {
	Type        *types.TransactionType `json:"type"`
	Amount      *float64               `json:"amount"`
	Description *string                `json:"description"`
	Date        *time.Time             `json:"date"`
}

type TransactionUpdateResponseDto struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	Type        types.TransactionType `json:"type"`
	Amount      float64               `json:"amount"`
	Description string                `json:"description"`
	Date        time.Time             `json:"date"`
}

type TransactionsFilter struct {
	DateTo    *time.Time             `json:"date_to"`
	DateFrom  *time.Time             `json:"date_from"`
	Type      *types.TransactionType `json:"type"`
	MinAmount *float64               `json:"min_amount"`
	MaxAmount *float64               `json:"max_amount"`
	SortBy    *string                `json:"sort_by"`
	SortOrder *string                `json:"sort_order"`
}
