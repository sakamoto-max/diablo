package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/diablo/internal/pkg/token"
)

var userIdKey ctxKey = "userId"

type Auth struct{}

func (a *Auth) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get("Authorization")
		if tkn == "" {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "token is missing",
			})
			return
		}

		claims, err := token.ValidateToken(tkn)
		if err != nil {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "token is invalid",
			})
			return
		}

		newCtx := context.WithValue(r.Context(), userIdKey, claims.UserId)

		next.ServeHTTP(w, r.WithContext(newCtx))

	})
}

func GetUserId(ctx context.Context) string {
	id, ok := ctx.Value(userIdKey).(string)
	if !ok {
		return ""
	}
	return id
}
