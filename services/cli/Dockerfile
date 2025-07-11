FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.work go.work.sum ./
COPY services/cli/go.mod services/cli/go.sum ./services/cli/
COPY libs ./libs/

# Download dependencies
WORKDIR /app/services/cli
RUN go mod download

# Copy the source code
WORKDIR /app
COPY services/cli ./services/cli/

# Build the application
WORKDIR /app/services/cli
RUN CGO_ENABLED=0 GOOS=linux go build -o cli ./cmd/main.go

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/cli/cli .

# Copy any necessary configuration files
COPY configs/cli ./configs/cli/

# Run the application
ENTRYPOINT ["./cli"]