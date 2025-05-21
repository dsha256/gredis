FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/air-verse/air@latest

# Build the application
RUN go build -o gredis ./cmd/api

# Development stage with hot reload
FROM builder AS development
WORKDIR /app
COPY . .
COPY .air.toml .
EXPOSE 8090
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:3.21 AS production
WORKDIR /app
COPY --from=builder /app/gredis /usr/local/bin/gredis
COPY config.yaml /app/config.yaml
RUN apk add --no-cache bash curl
EXPOSE 8090
CMD ["gredis"]
