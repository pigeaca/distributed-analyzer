FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.work go.work.sum ./
COPY services/scheduler-service/go.mod services/scheduler-service/go.sum ./services/scheduler-service/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/scheduler-service
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/scheduler-service ./services/scheduler-service/

# Build the application
WORKDIR /app/services/scheduler-service
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler-service ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/scheduler-service/scheduler-service .

# Copy any necessary configuration files
COPY configs/scheduler-service ./configs/scheduler-service/

# Expose the port the service runs on
EXPOSE 8080

# Run the application
CMD ["./scheduler-service"]