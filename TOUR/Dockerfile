# Stage 1: Build
FROM golang:1.25.4-alpine AS builder
WORKDIR /app

# Install git (required for go modules fetching)
RUN apk update && apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o tour ./cmd/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app

# Install CA certificates for PostgreSQL TLS
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/tour .

# Expose correct port
EXPOSE 8094

CMD ["./tour"]
