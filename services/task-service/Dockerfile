FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY services/task-service/go.mod services/task-service/go.sum ./services/task-service/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/task-service
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/task-service ./services/task-service/

# Build the application
WORKDIR /app/services/task-service
RUN CGO_ENABLED=0 GOOS=linux go build -o task-service ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/task-service/task-service .

# Copy any necessary configuration files
COPY configs/task-service ./configs/task-service/

# Expose the port the service runs on
EXPOSE 8082
EXPOSE 9082

# Run the application
CMD ["./task-service"]