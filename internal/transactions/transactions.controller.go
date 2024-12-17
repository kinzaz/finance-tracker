package transactions

import (
	"errors"
	"finance-tracker/internal/dto"
	"finance-tracker/internal/errs"
	"finance-tracker/internal/types"
	"finance-tracker/pkg/middleware"
	"finance-tracker/pkg/request"
	"finance-tracker/pkg/response"
	"net/http"
	"strconv"
	"time"
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
		queryParams := r.URL.Query()
		filters := TransactionsFilter{}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 0 {
			offset = 0
		}

		pagination := dto.PaginationRequestDto{
			Limit:  limit,
			Offset: offset,
		}

		if err := parseDateFilter(queryParams.Get("date_from"), &filters.DateFrom); err != nil {
			http.Error(w, "Invalid date_from format", http.StatusBadRequest)
			return
		}

		if err := parseDateFilter(queryParams.Get("date_to"), &filters.DateTo); err != nil {
			http.Error(w, "Invalid date_to format", http.StatusBadRequest)
			return
		}

		if t := queryParams.Get("type"); t != "" {
			tp := types.TransactionType(t)
			filters.Type = &tp
		}

		if err := parseFloatFilter(queryParams.Get("min_amount"), &filters.MinAmount); err != nil {
			http.Error(w, "Invalid min_amount format", http.StatusBadRequest)
			return
		}

		if err := parseFloatFilter(queryParams.Get("max_amount"), &filters.MaxAmount); err != nil {
			http.Error(w, "Invalid max_amount format", http.StatusBadRequest)
			return
		}

		if sortBy := queryParams.Get("sort_by"); sortBy != "" {
			filters.SortBy = &sortBy
		}

		if sortOrder := queryParams.Get("sort_order"); sortOrder != "" {
			filters.SortOrder = &sortOrder
		}

		res, err := controller.TransactionService.GetUserTransactions(userId, filters, pagination)
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
			} else if errors.Is(err, errs.ErrInsufficientBalance) {
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

func parseFloatFilter(amountStr string, target **float64) error {
	if amountStr != "" {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return err
		}
		*target = &amount
	}
	return nil
}

func parseDateFilter(dateStr string, target **time.Time) error {
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
		*target = &parsedDate
	}
	return nil
}
