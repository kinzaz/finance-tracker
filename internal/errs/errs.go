package errs

import "errors"

/* user errors */
var ErrUserNotFound = errors.New("user not found")

/* transactions errors */
var ErrInvalidTransactionType = errors.New("invalid transaction type")
