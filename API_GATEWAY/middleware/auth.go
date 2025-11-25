package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		// --- DODAJ OVO ZA DEBUG ---
		fmt.Printf(">>> DEBUG TOKEN CLAIMS: %+v\n", claims)
		// ---------

		// 4. PAKOVANJE PODATAKA U METADATA
		// Ovde uzimamo podatke iz tokena i pakujemo ih za gRPC servise
		// Koristimo 'sub' za ID, 'role', 'email', 'username' (sta god imas u tokenu)

		userID, _ := claims["id"].(string) // Prilagodi kljuc tvom tokenu
		role, _ := claims["role"].(string)
		username, _ := claims["username"].(string)
		email, _ := claims["email"].(string)

		// 5. SLANJE PODATAKA PUTEM HEADER-a (STANDARD ZA GRPC-GATEWAY)
		// Biblioteka ce ovo automatski pretvoriti u Metadata: "user-id", "user-role"...
		r.Header.Set("Grpc-Metadata-User-Id", userID)
		r.Header.Set("Grpc-Metadata-User-Role", role)
		r.Header.Set("Grpc-Metadata-User-Username", username)
		r.Header.Set("Grpc-Metadata-User-Email", email)

		// Takodje prosledjujemo originalni Authorization header
		r.Header.Set("Grpc-Metadata-Authorization", authHeader)

		next.ServeHTTP(w, r)
	})
}
