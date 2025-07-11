FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.work go.work.sum ./
COPY services/worker/go.mod services/worker/go.sum ./services/worker/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/worker
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/worker ./services/worker/

# Build the application
WORKDIR /app/services/worker
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/worker/worker .

# Copy any necessary configuration files
COPY configs/worker ./configs/worker/

# Run the application
CMD ["./worker"]