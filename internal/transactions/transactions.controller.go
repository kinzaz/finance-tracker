package transactions

import (
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/pkg/request"
	"finance-tracker/pkg/response"
	"net/http"
)

type TransactionsController struct {
	TransactionService TransactionsServiceInterface
}

func NewTransactionsController(router *http.ServeMux, transactionsService TransactionsServiceInterface) {
	handler := &TransactionsController{
		TransactionService: transactionsService,
	}

	router.HandleFunc("POST /transaction", handler.Create())
	router.HandleFunc("DELETE /transaction/{id}", handler.Delete())
}

func (controller *TransactionsController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto, err := request.HandleBody[TransactionRequestDto](w, r)
		if err != nil {
			return
		}

		res, err := controller.TransactionService.CreateTransaction(dto)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidTransactionType) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else if errors.Is(err, errs.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		response.Json(w, res, http.StatusCreated)
	}
}

func (controller *TransactionsController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := request.GetParam[uint](r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = controller.TransactionService.DeleteTransaction(id)
		if err != nil {
			if errors.Is(err, errs.ErrTransactionNotFound) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response.Json(w, nil, http.StatusOK)
	}
}
