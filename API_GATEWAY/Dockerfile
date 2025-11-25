# Stage 1: build
FROM golang:1.25.4-alpine AS builder
WORKDIR /app

# 1. Kopiramo COMMON folder
COPY COMMON ./COMMON

# --- DODAJ OVO (DEBUG) ---
# Ova komanda ce srusiti build ako Docker ne vidi "admin/users" u fajlu.
# Takodje ce ispisati liniju u terminalu ako je nadje.
RUN echo "PROVERAVAM FAJL U DOCKERU..." && \
    grep "admin/users" COMMON/auth/proto/auth.pb.gw.go && \
    echo "FAJL JE DOBAR!"
# -------------------------

# 2. Kopiramo go.mod i go.sum
COPY API_GATEWAY/go.mod API_GATEWAY/go.sum ./API_GATEWAY/

# ... (ostatak fajla je isti) ...
WORKDIR /app/API_GATEWAY
RUN go mod download
COPY API_GATEWAY/ .
RUN go build -o gateway cmd/main.go

# Stage 2: runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/API_GATEWAY/gateway .
EXPOSE 8080
CMD ["./gateway"]