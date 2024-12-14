package user

import (
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/pkg/middleware"
	"finance-tracker/pkg/response"
	"net/http"
)

type UserController struct {
	UserService UserServiceInterface
}

func NewUserController(router *http.ServeMux, userService UserServiceInterface) {
	handler := &UserController{
		UserService: userService,
	}

	router.Handle("GET /user/profile", middleware.IsAuthed(handler.GetUserProfile()))
}

func (controller *UserController) GetUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := r.Context().Value(middleware.ContextIdKey).(uint)

		userProfile, err := controller.UserService.GetUserProfile(id)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		response.Json(w, &UserProfileResponseDto{
			Email:   userProfile.Email,
			Name:    userProfile.Name,
			Balance: userProfile.Balance,
		}, http.StatusOK)
	}
}
