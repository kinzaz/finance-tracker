package transactions

import (
	"encoding/json"
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/internal/user"
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
		/* Вынести в пакет */
		var dto TransactionRequestDto
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&dto); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := controller.UserRepository.FindById(dto.UserID)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		response, err := controller.TransactionService.CreateTransaction(dto)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidTransactionType) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		/* Вынести в пакет */
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
