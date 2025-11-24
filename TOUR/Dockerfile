# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# 1. Copy COMMON (shared code, generated proto, utils, etc.)
COPY COMMON ./COMMON

# 2. Copy go.mod and go.sum from TOUR folder
COPY TOUR/go.mod TOUR/go.sum ./TOUR/

# 3. Enter the TOUR directory
WORKDIR /app/TOUR

# Download dependencies
RUN go mod download

# 4. Copy the rest of the TOUR service code
COPY TOUR .

# Build the binary
RUN go build -o tour cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /app/TOUR/tour .

# Service port (MUST match main.go -> :8083)
EXPOSE 8083

# Run the service
CMD ["./tour"]
