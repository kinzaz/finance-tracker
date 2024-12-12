package types

import "finance-tracker/internal/errs"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

func (t TransactionType) Validate() error {
	switch t {
	case Income, Expense:
		return nil
	default:
		return errs.ErrInvalidTransactionType
	}
}
