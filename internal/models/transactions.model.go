package models

import (
	"finance-tracker/internal/types"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	UserID      uint                  `json:"user_id"`
	Type        types.TransactionType `json:"type"`
	Amount      float64               `json:"amount"`
	Description string                `json:"description"`
	Date        time.Time             `json:"date"`
}
