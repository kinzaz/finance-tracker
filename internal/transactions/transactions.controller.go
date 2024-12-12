package transactions

import (
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/internal/user"
	"finance-tracker/pkg/request"
	"finance-tracker/pkg/response"
	"net/http"
)

type TransactionsController struct {
	TransactionService TransactionsServiceInterface
	UserRepository     user.UserRepositoryInterface
}

func NewTransactionsController(router *http.ServeMux, transactionsService TransactionsServiceInterface, userRepository user.UserRepositoryInterface) {
	handler := &TransactionsController{
		TransactionService: transactionsService,
		UserRepository:     userRepository,
	}

	router.HandleFunc("POST /transaction", handler.Create())
}

func (controller *TransactionsController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto, err := request.HandleBody[TransactionRequestDto](w, r)
		if err != nil {
			return
		}

		_, err = controller.UserRepository.FindById(dto.UserID)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		res, err := controller.TransactionService.CreateTransaction(dto)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidTransactionType) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response.Json(w, res, http.StatusCreated)
	}
}
