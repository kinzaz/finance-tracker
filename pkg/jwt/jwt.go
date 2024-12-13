package jwt

import (
	"errors"
	"finance-tracker/config"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(data CustomClaims) (string, error) {
	conf := config.LoadConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"id":    data.ID,
	})

	return token.SignedString([]byte(conf.SECRET))
}

func ParseJWT(token string) (bool, *CustomClaims) {
	conf := config.LoadConfig()

	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(conf.SECRET), nil
	})

	if err != nil || !parsedToken.Valid {
		return false, nil
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)

	if !ok {
		return false, nil
	}

	return true, claims
}
