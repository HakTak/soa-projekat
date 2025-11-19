package middleware

import (
	"net/http"
	"strings"

	"auth/internal/services"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtService *services.JWTService
}

func NewAuthMiddleware(jwtService *services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService}
}

// Ekstrakcija tokena iz Authorization header-a
func (m *AuthMiddleware) extractToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", http.ErrNoCookie
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return "", http.ErrNoCookie
	}

	return header[7:], nil // skip "Bearer "
}

// Middleware: validira token i vraca claims ako je validan
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := m.extractToken(r)
		if err != nil {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token, err := m.jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// token is valid → send claims forward
		claims := token.Claims.(jwt.MapClaims)

		// attach claims to request context
		ctx := r.Context()
		ctx = ContextWithClaims(ctx, claims)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// Middleware: dozvoli samo ADMIN korisnike
func (m *AuthMiddleware) AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return m.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		claims := ClaimsFromContext(r.Context())
		role := claims["role"].(string)

		if role != "ADMIN" {
			http.Error(w, "Access restricted — admin only", http.StatusForbidden)
			return
		}

		next(w, r)
	})
}
