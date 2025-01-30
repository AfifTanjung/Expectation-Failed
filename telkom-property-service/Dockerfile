# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/main.go

# Production stage
FROM alpine:latest
COPY --from=builder /app/user-service /user-service
EXPOSE 8080
CMD ["main.go"]