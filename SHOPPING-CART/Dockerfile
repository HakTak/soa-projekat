# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# 1. Kopiramo COMMON
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum u STAKEHOLDERS folder
COPY SHOPPING-CART/go.mod SHOPPING-CART/go.sum ./SHOPPING-CART/

# 3. Ulazimo u STAKEHOLDERS folder
WORKDIR /app/SHOPPING-CART

# Download
RUN go mod download

# 4. Kopiramo ostatak koda
COPY SHOPPING-CART .

# Build
RUN go build -o shopping-cart cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /app/SHOPPING-CART/shopping-cart .


EXPOSE 9092
CMD ["./shopping-cart"]