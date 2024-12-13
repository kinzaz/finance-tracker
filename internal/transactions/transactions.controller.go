package transactions

import (
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/pkg/middleware"
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

	router.Handle("GET /transactions", middleware.IsAuthed(handler.GetUserTransactions()))

	router.Handle("GET /transaction/{id}", middleware.IsAuthed(handler.GetUserTransaction()))
	router.Handle("POST /transaction", middleware.IsAuthed(handler.Create()))
	router.Handle("DELETE /transaction/{id}", middleware.IsAuthed(handler.Delete()))
	router.Handle("PATCH /transaction/{id}", middleware.IsAuthed(handler.Update()))
}

func (controller *TransactionsController) GetUserTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := r.Context().Value(middleware.ContextIdKey).(uint)
		transactionId, err := request.GetParam[uint](r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := controller.TransactionService.GetUserTransaction(userId, transactionId)
		if err != nil {
			if errors.Is(err, errs.ErrTransactionNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response.Json(w, res, http.StatusOK)
	}
}

func (controller *TransactionsController) GetUserTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := r.Context().Value(middleware.ContextIdKey).(uint)

		res, err := controller.TransactionService.GetUserTransactions(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, res, http.StatusOK)
	}
}

func (controller *TransactionsController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := r.Context().Value(middleware.ContextIdKey).(uint)
		dto, err := request.HandleBody[TransactionRequestDto](w, r)
		if err != nil {
			return
		}

		res, err := controller.TransactionService.CreateTransaction(id, dto)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidTransactionType) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else if errors.Is(err, errs.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
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
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response.Json(w, nil, http.StatusOK)
	}
}

func (controller *TransactionsController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := request.GetParam[uint](r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dto, err := request.HandleBody[TransactionUpdateRequestDto](w, r)
		if err != nil {
			return
		}

		updatedTransaction, err := controller.TransactionService.UpdateTransaction(id, dto)
		if err != nil {
			if errors.Is(err, errs.ErrTransactionNotFound) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response.Json(w, TransactionResponseDto{
			ID:          updatedTransaction.ID,
			UserID:      updatedTransaction.UserID,
			Type:        updatedTransaction.Type,
			Description: updatedTransaction.Description,
			Amount:      updatedTransaction.Amount,
			Date:        updatedTransaction.Date,
		}, http.StatusOK)
	}
}
