package auth

import (
	"finance-tracker/pkg/jwt"
	"finance-tracker/pkg/request"
	"finance-tracker/pkg/response"
	"net/http"
)

type AuthController struct {
	AuthService AuthServiceInterface
}

func NewAuthController(router *http.ServeMux, service AuthServiceInterface) {
	handler := &AuthController{
		AuthService: service,
	}

	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/login", handler.Login())
}

func (controller AuthController) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dto, err := request.HandleBody[RegisterRequestDto](w, r)
		if err != nil {
			return
		}

		res, _ := controller.AuthService.Register(dto)

		response.Json(w, res, http.StatusCreated)
	}
}

func (controller AuthController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto, err := request.HandleBody[LoginRequestDto](w, r)
		if err != nil {
			return
		}

		email, err := controller.AuthService.Login(dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.GenerateJWT(jwt.JWTData{
			Email: email,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := &LoginResponseDto{
			Access_token: token,
		}

		response.Json(w, res, http.StatusCreated)
	}
}
