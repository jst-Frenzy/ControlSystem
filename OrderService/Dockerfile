FROM golang:1.26.0-alpine as builder
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/order-service .
COPY --from=builder /app/migrations ./migrations/

CMD ["./order-service"]