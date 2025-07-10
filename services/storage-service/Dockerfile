FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.work go.work.sum ./
COPY services/storage-service/go.mod services/storage-service/go.sum ./services/storage-service/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/storage-service
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/storage-service ./services/storage-service/

# Build the application
WORKDIR /app/services/storage-service
RUN CGO_ENABLED=0 GOOS=linux go build -o storage-service ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/storage-service/storage-service .

# Copy any necessary configuration files
COPY configs/storage-service ./configs/storage-service/

# Expose the port the service runs on
EXPOSE 8080

# Run the application
CMD ["./storage-service"]