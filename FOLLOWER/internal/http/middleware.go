package http

import (
	"context"
	"net/http"
)

// AuthMiddleware expects an X-User-ID header (from Gateway or Postman)
// and puts it into the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Look for the User ID in the header
		// In Prod: The Gateway sends this.
		// In Dev: Postman sends this.
		userID := r.Header.Get("X-User-ID")

		if userID != "" {
			// 2. Put it in the context so ctxSubject() can find it
			ctx := context.WithValue(r.Context(), "subject", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// If no header, just pass through.
			// The handler checks if user is missing and returns 401.
			next.ServeHTTP(w, r)
		}
	})
}
