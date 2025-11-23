package middleware

import (
	"context"
	"strings"

	"auth/internal/services"

	// IMPORTUJEMO TVOJ NOVI UTILS PAKET
	"PROJEKAT/COMMON/utils"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcAuthMiddleware struct {
	jwtService *services.JWTService
}

func NewGrpcAuthMiddleware(jwtService *services.JWTService) *GrpcAuthMiddleware {
	return &GrpcAuthMiddleware{jwtService}
}

func (m *GrpcAuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// 1. Javne metode (Login/Register) preskacemo
		publicMethods := map[string]bool{
			"/auth.AuthService/Login":    true,
			"/auth.AuthService/Register": true,
		}
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// 2. Vadimo token
		tokenStr, err := m.extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Missing token")
		}

		// 3. Validacija tokena
		token, err := m.jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "Invalid token")
		}

		// 4. Izvlacenje podataka (Claims)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Internal, "Invalid claims")
		}

		// --- OVDE JE BILA GRESKA, SADA JE POPRAVLJENO ---
		// Koristimo funkciju iz COMMON/utils
		newCtx := utils.ContextWithClaims(ctx, claims)

		// Prosledjujemo novi kontekst dalje
		return handler(newCtx, req)
	}
}

// extractToken ostaje isti...
func (m *GrpcAuthMiddleware) extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "No metadata")
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "No auth header")
	}
	return strings.TrimPrefix(values[0], "Bearer "), nil
}
