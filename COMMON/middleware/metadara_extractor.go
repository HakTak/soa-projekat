package middleware

import (
	"PROJEKAT/COMMON/utils"
	"context"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Ovaj interceptor samo prepakuje Metadata -> Context
// Ne radi validaciju tokena (to je vec uradio Gateway)
func MetadataExtractorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// Nema metadata (mozda interni poziv), pustamo dalje
		return handler(ctx, req)
	}

	// Pravimo claims mapu od podataka iz headera
	claims := jwt.MapClaims{}

	if vals := md.Get("user-id"); len(vals) > 0 {
		claims["id"] = vals[0]
	}
	if vals := md.Get("user-role"); len(vals) > 0 {
		claims["role"] = vals[0]
	}
	if vals := md.Get("user-username"); len(vals) > 0 {
		claims["username"] = vals[0]
	}
	if vals := md.Get("user-email"); len(vals) > 0 {
		claims["email"] = vals[0]
	}

	// Ako smo nasli podatke, ubacujemo ih u kontekst
	if len(claims) > 0 {
		ctx = utils.ContextWithClaims(ctx, claims)
	}

	return handler(ctx, req)
}
