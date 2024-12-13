package transactions

import (
	"finance-tracker/internal/types"
	"time"
)

type TransactionRequestDto struct {
	Type        types.TransactionType `json:"type"`
	Amount      float64               `json:"amount"`
	Description string                `json:"description"`
	Date        time.Time             `json:"date"`
}

type TransactionResponseDto struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	Type        types.TransactionType `json:"type"`
	Amount      float64               `json:"amount"`
	Description string                `json:"description"`
	Date        time.Time             `json:"date"`
}
