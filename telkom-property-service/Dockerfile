# Build stage
FROM --platform=linux/amd64 golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o telkom-property-service ./cmd/main.go

# Production stage
FROM alpine:latest
WORKDIR /app

# Tambahkan ini untuk debugging
RUN apk add --no-cache ca-certificates tzdata curl

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Pastikan permissions benar
COPY --from=builder --chown=appuser:appgroup /app/telkom-property-service /app/telkom-property-service

# Explicitly set PORT environment variable
ENV PORT=8080

# Explicitly expose port
EXPOSE 8080

USER appuser

# Ubah CMD menjadi lebih explicit
CMD ["/app/telkom-property-service"]