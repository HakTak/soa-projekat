package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

var jwtSecret = []byte("my_super_duper_secret_mega_gg_key_123")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Javne rute
		path := r.URL.Path
		if strings.Contains(path, "/login") || strings.Contains(path, "/register") {
			next.ServeHTTP(w, r)
			return
		}

		// 2. Provera Headera
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// 3. Validacija Tokena
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Prosledjivanje kroz Metadata
		md := metadata.Pairs("authorization", authHeader)
		ctx := metadata.NewOutgoingContext(r.Context(), md)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
