# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# 1. Kopiramo COMMON
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum u STAKEHOLDERS folder
COPY STAKEHOLDERS/go.mod STAKEHOLDERS/go.sum ./STAKEHOLDERS/

# 3. Ulazimo u STAKEHOLDERS folder
WORKDIR /app/STAKEHOLDERS

# Download
RUN go mod download

# 4. Kopiramo ostatak koda
COPY STAKEHOLDERS .

# Build
RUN go build -o stakeholders cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/STAKEHOLDERS/stakeholders .

EXPOSE 8081
CMD ["./stakeholders"]