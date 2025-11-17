# Stage 1: build
FROM golang:1.25.4-alpine AS builder
WORKDIR /app

# Kopiramo module fajlove i skidamo dependencies
COPY go.mod ./
RUN go mod download

# Kopiramo ostatak koda
COPY . .

# Build komanda za gateway
RUN go build -o gateway cmd/main.go

# Stage 2: runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway .

# Port na kojem servis radi
EXPOSE 8080

# Komanda koja se izvršava kad container startuje
CMD ["./gateway"]
