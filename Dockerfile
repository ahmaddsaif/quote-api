# Build stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go mod tidy
RUN go build -o quote-api ./cmd/api

# Run Stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/quote-api .

EXPOSE 8080

CMD ["./quote-api"]
