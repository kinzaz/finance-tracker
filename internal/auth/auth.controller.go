package auth

import (
	"encoding/json"
	"finance-tracker/pkg/jwt"
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
		/* Вынести в пакет */
		var dto RegisterRequestDto
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&dto)

		response, _ := controller.AuthService.Register(dto)

		/* Вынести в пакет */
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func (controller AuthController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* Вынести в пакет */
		var dto LoginRequestDto
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&dto)

		// response, err := controller.AuthRepository.login(dto)
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

		response := &LoginResponseDto{
			Access_token: token,
		}

		/* Вынести в пакет */
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
