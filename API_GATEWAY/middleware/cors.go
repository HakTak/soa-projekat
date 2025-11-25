package middleware

import "net/http"

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dozvoljavamo pristup sa Angular porta
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		// Dozvoljavamo metode
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Dozvoljavamo headere (ukljucujuci Authorization)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Grpc-Metadata-User-Id, Grpc-Metadata-User-Role")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Ako je "OPTIONS" zahtev (pre-flight check), vratimo OK odmah
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}