package jwt

import (
	"finance-tracker/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email string `json:"email"`
}

func GenerateJWT(data JWTData) (string, error) {
	conf := config.LoadConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	})

	return token.SignedString([]byte(conf.SECRET))
}
