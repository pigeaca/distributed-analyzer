#!/bin/bash

# Exit on error
set -e

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Array of services that need Swagger documentation
SERVICES=(
    "api-gateway"
)

# Generate Swagger specs for each service
for SERVICE in "${SERVICES[@]}"; do
    SERVICE_PATH="./services/${SERVICE}"
    
    # Check if the service directory exists
    if [ -d "$SERVICE_PATH" ]; then
        echo "Generating Swagger specs for ${SERVICE}..."
        
        # Check if the service has a cmd directory (where main.go with Swagger annotations should be)
        if [ -d "${SERVICE_PATH}/cmd" ]; then
            # Create docs directory if it doesn't exist
            mkdir -p "${SERVICE_PATH}/docs"
            
            # Generate Swagger specs
            (cd "$SERVICE_PATH" && swag init -g cmd/main.go -o docs)
            
            echo "Swagger specs generated for ${SERVICE}"
        else
            echo "Skipping ${SERVICE} - no cmd directory found"
        fi
    else
        echo "Skipping ${SERVICE} - directory not found"
    fi
done

echo "Swagger spec generation complete!"