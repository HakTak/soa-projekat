# Stage 1: build
FROM golang:1.25.4-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o tour cmd/main.go

# Stage 2: run
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/tour .

EXPOSE 8093

CMD ["./tour"]
