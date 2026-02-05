FROM golang:1.24.1-alpine AS builder
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations/

CMD ["./auth-service"]