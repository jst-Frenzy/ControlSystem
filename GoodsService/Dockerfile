FROM golang:1.24.1-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o goods-service ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/goods-service .

CMD ["./goods-service"]