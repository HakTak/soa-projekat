# Stage 1: build
FROM golang:1.25.4-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o stakeholders cmd/main.go

# Stage 2: run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/stakeholders .

EXPOSE 8081
CMD ["./stakeholders"]
