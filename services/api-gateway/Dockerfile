FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/api-gateway
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/api-gateway ./services/api-gateway/

# Build the application
WORKDIR /app/services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/api-gateway/api-gateway .

# Copy any necessary configuration files
COPY configs/api-gateway.yml ./configs/api-gateway.yml

# Expose the port the service runs on
EXPOSE 8082

# Run the application
CMD ["./api-gateway"]