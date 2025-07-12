#!/bin/bash

# Exit on error
set -e

# Directory containing proto files
PROTO_DIR="./api/proto"
# Directory for generated Go code
OUTPUT_DIR="../libs/proto"

# Create output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Install protoc plugins if not already installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate code for each proto file
echo "Generating Go code from proto files..."

# Task service
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    --go-grpc_out=. --go-grpc_opt=module=distributed-analyzer \
    $PROTO_DIR/task/task.proto

# Scheduler service
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    --go-grpc_out=. --go-grpc_opt=module=distributed-analyzer \
    $PROTO_DIR/scheduler/scheduler.proto

# Worker service
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    --go-grpc_out=. --go-grpc_opt=module=distributed-analyzer \
    $PROTO_DIR/worker/worker.proto

# Result service
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    --go-grpc_out=. --go-grpc_opt=module=distributed-analyzer \
    $PROTO_DIR/result/result.proto

# Billing service
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    --go-grpc_out=. --go-grpc_opt=module=distributed-analyzer \
    $PROTO_DIR/billing/billing.proto

# Audit service
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    --go-grpc_out=. --go-grpc_opt=module=distributed-analyzer \
    $PROTO_DIR/audit/audit.proto

# User service
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    --go-grpc_out=. --go-grpc_opt=module=distributed-analyzer \
    $PROTO_DIR/user/user.proto

# Kafka messages
protoc -I . \
    --go_out=. --go_opt=module=distributed-analyzer \
    $PROTO_DIR/kafka/messages.proto

echo "Code generation complete!"
