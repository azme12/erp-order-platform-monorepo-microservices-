#!/bin/bash

# Generate Swagger documentation for all services
# This script ensures Swagger docs are always up-to-date

set -e

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "🔧 Generating Swagger documentation for all services..."
echo "Working directory: $(pwd)"
echo ""

SERVICES=("auth" "contact" "inventory" "purchase" "sales")

for service in "${SERVICES[@]}"; do
    echo "📝 Generating Swagger docs for $service service..."
    
    SERVICE_DIR="services/$service"
    
    if [ ! -d "$SERVICE_DIR" ]; then
        echo "⚠️  Warning: Directory $SERVICE_DIR not found"
        continue
    fi
    
    cd "$SERVICE_DIR"
    
    if [ -f "cmd/main.go" ]; then
        echo "  Running: swag init -g cmd/main.go --parseDependency --parseInternal"
        swag init -g cmd/main.go --parseDependency --parseInternal --output ./docs
        echo "✅ Swagger docs generated for $service"
    else
        echo "⚠️  Warning: cmd/main.go not found for $service service"
    fi
    
    cd "$SCRIPT_DIR"
    echo ""
done

echo "✅ All Swagger documentation generated successfully!"
echo ""
echo "📚 Swagger UI available at:"
echo "  - Auth: http://localhost:8002/swagger/index.html"
echo "  - Contact: http://localhost:8001/swagger/index.html"
echo "  - Inventory: http://localhost:8003/swagger/index.html"
echo "  - Sales: http://localhost:8004/swagger/index.html"
echo "  - Purchase: http://localhost:8005/swagger/index.html"

