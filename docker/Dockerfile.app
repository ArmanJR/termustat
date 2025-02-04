# Build stage
FROM golang:1.23.5-alpine AS builder
WORKDIR /app
COPY .. .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o termustat cmd/api/main.go

# Runtime stage
FROM alpine:3.21.2
WORKDIR /app
COPY --from=builder /app/termustat .
COPY ../.env /app/.env

EXPOSE 80
ENTRYPOINT ["./termustat"]