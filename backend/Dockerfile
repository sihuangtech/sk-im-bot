# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY config/config.yaml ./config/config.yaml

EXPOSE 8888

CMD ["./main"]
