# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# 1. Kopiramo COMMON
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum u STAKEHOLDERS folder
COPY FOLLOWER/go.mod FOLLOWER/go.sum ./FOLLOWER/

# 3. Ulazimo u STAKEHOLDERS folder
WORKDIR /app/FOLLOWER

# Download
RUN go mod download

# 4. Kopiramo ostatak koda
COPY FOLLOWER .

# Build
RUN go build -o follower cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/FOLLOWER/follower .

RUN chmod +x follower

EXPOSE 8084
EXPOSE 9090
CMD ["./follower"]