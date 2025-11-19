package middleware

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type claimsKey string

const claimsContextKey claimsKey = "jwtClaims"

func ContextWithClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) jwt.MapClaims {
	if claims, ok := ctx.Value(claimsContextKey).(jwt.MapClaims); ok {
		return claims
	}
	return nil
}
