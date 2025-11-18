package http

import (
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get(("Authorization"))
		if h == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		var subject string
		var roles []string
		switch token {
		case "token-alice":
			subject = "123e4567-e89b-12d3-a456-426614174001"
			roles = []string{"guide"}
		case "token-bob":
			subject = "123e4567-e89b-12d3-a456-426614174002"
			roles = []string{"tourist"}
		case "token-carol":
			subject = "123e4567-e89b-12d3-a456-426614174003"
			roles = []string{"admin"}
		default:
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "subject", subject)
		ctx = context.WithValue(ctx, "roles", roles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
