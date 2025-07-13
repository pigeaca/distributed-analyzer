FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY services/worker-manager/go.mod services/worker-manager/go.sum ./services/worker-manager/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/worker-manager
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/worker-manager ./services/worker-manager/

# Build the application
WORKDIR /app/services/worker-manager
RUN CGO_ENABLED=0 GOOS=linux go build -o worker-manager ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/worker-manager/worker-manager .

# Copy any necessary configuration files
COPY configs/worker-manager.yml ./configs/worker-manager.yml

# Expose the port the service runs on
EXPOSE 8086
EXPOSE 9086

# Run the application
CMD ["./worker-manager"]