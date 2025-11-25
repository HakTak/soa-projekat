package utils

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Definisemo kljuc za kontekst. Mora biti ne-eksportovan tip da ne bi bilo kolizije.
type claimsKey struct{}

var ClaimsContextKey = claimsKey{}

// 1. ContextWithClaims: Ubacuje podatke (claims) u kontekst
func ContextWithClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, claims)
}

// 2. ClaimsFromContext: Vadi podatke iz konteksta
func ClaimsFromContext(ctx context.Context) jwt.MapClaims {
	if claims, ok := ctx.Value(ClaimsContextKey).(jwt.MapClaims); ok {
		return claims
	}
	return nil
}

// 3. Authorize: Univerzalna funkcija za proveru uloge
// Ovu funkciju ces moci da zoves iz BILO KOG servisa (Auth, Stakeholders, Blog...)
func Authorize(ctx context.Context, requiredRoles ...string) error {
	claims := ClaimsFromContext(ctx)
	if claims == nil {
		return status.Error(codes.Unauthenticated, "No user claims found in context")
	}

	// Token claims su obicno map[string]interface{}, pa moramo da kastujemo
	userRole, ok := claims["role"].(string)
	if !ok {
		return status.Error(codes.Unauthenticated, "User role not found in token")
	}

	// Provera da li se uloga poklapa sa nekom od trazenih
	for _, role := range requiredRoles {
		if userRole == role {
			return nil // OK!
		}
	}

	return status.Errorf(codes.PermissionDenied, "Access denied. User role '%s' is not allowed. Required: %v", userRole, requiredRoles)
}
