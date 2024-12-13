package middleware

import (
	"context"
	"finance-tracker/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
	ContextIdKey    key = "ContextIdKey"
)

func writeUnauthed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func extractTokenFromHeader(authHeader string) string {
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractTokenFromHeader(r.Header.Get("Authorization"))

		isValid, claims := jwt.ParseJWT(token)

		if !isValid {
			writeUnauthed(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, claims.Email)
		ctx = context.WithValue(ctx, ContextIdKey, claims.ID)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
