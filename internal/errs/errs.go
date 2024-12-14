package errs

import "errors"

/* user errors */
var ErrUserNotFound = errors.New("user not found")
var ErrRegisterUser = errors.New("user with your email has already registered")
var ErrInsufficientBalance = errors.New("insufficient balance")

/* transactions errors */
var ErrInvalidTransactionType = errors.New("invalid transaction type")

var ErrTransactionNotFound = errors.New("transaction not found")
