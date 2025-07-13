FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.work go.work.sum ./
COPY services/result-service/go.mod services/result-service/go.sum ./services/result-service/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/result-service
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/result-service ./services/result-service/

# Build the application
WORKDIR /app/services/result-service
RUN CGO_ENABLED=0 GOOS=linux go build -o result-service ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/result-service/result-service .

# Copy any necessary configuration files
COPY configs/result-service ./configs/result-service/

# Expose the port the service runs on
EXPOSE 8084
EXPOSE 9084

# Run the application
CMD ["./result-service"]