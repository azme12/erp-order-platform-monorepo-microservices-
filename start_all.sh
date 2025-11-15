#!/bin/bash
set -e

cd /home/admn/Documents/project/SR-BE-interview-1

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "📋 Please create .env file:"
    echo "   cp .env.example .env"
    echo "   Then edit .env with your secure passwords and secrets"
    exit 1
fi

# Load environment variables
set -a
source .env
set +a

echo "🚀 Starting All Services..."
echo ""

# Step 1: Start Database and NATS
echo "📦 Step 1: Starting database and NATS..."
docker compose up -d db nats
sleep 5
echo "✅ Database and NATS started"
echo ""

# Step 2: Start Auth Service
echo "🔐 Step 2: Starting Auth Service on port ${AUTH_SERVICE_PORT:-8002}..."
PORT=${AUTH_SERVICE_PORT:-8002} \
DB_HOST=${DB_HOST} \
DB_PORT=${DB_PORT} \
DB_USER=${DB_USER} \
DB_PASSWORD=${DB_PASSWORD} \
DB_NAME=${DB_NAME} \
JWT_SECRET=${JWT_SECRET} \
JWT_USER_EXPIRATION_HOURS=${JWT_USER_EXPIRATION_HOURS:-24} \
NATS_URL=${NATS_URL} \
go run services/auth/cmd/main.go > /tmp/auth.log 2>&1 &
AUTH_PID=$!
echo $AUTH_PID > /tmp/auth.pid
sleep 3

if curl -s http://localhost:${AUTH_SERVICE_PORT:-8002}/health > /dev/null; then
    echo "✅ Auth Service started (PID: $AUTH_PID)"
else
    echo "❌ Auth Service failed to start. Check /tmp/auth.log"
    exit 1
fi
echo ""

# Step 3: Start Contact Service
echo "📇 Step 3: Starting Contact Service on port ${CONTACT_SERVICE_PORT:-8001}..."
PORT=${CONTACT_SERVICE_PORT:-8001} \
DB_HOST=${DB_HOST} \
DB_PORT=${DB_PORT} \
DB_USER=${DB_USER} \
DB_PASSWORD=${DB_PASSWORD} \
DB_NAME=${DB_NAME} \
JWT_SECRET=${JWT_SECRET} \
NATS_URL=${NATS_URL} \
go run services/contact/cmd/main.go > /tmp/contact.log 2>&1 &
CONTACT_PID=$!
echo $CONTACT_PID > /tmp/contact.pid
sleep 3

if curl -s http://localhost:${CONTACT_SERVICE_PORT:-8001}/health > /dev/null; then
    echo "✅ Contact Service started (PID: $CONTACT_PID)"
else
    echo "❌ Contact Service failed to start. Check /tmp/contact.log"
    exit 1
fi
echo ""

# Step 4: Start Gateway
echo "🌐 Step 4: Starting Gateway on port ${GATEWAY_PORT:-8000}..."
PORT=${GATEWAY_PORT:-8000} \
AUTH_SERVICE_URL=${AUTH_SERVICE_URL} \
CONTACT_SERVICE_URL=${CONTACT_SERVICE_URL} \
JWT_SECRET=${JWT_SECRET} \
go run gateway/cmd/main.go > /tmp/gateway.log 2>&1 &
GATEWAY_PID=$!
echo $GATEWAY_PID > /tmp/gateway.pid
sleep 5

if curl -s http://localhost:${GATEWAY_PORT:-8000}/health > /dev/null; then
    echo "✅ Gateway started (PID: $GATEWAY_PID)"
else
    echo "❌ Gateway failed to start. Check /tmp/gateway.log"
    exit 1
fi
echo ""

echo "=========================================="
echo "✅ All services started successfully!"
echo "=========================================="
echo ""
echo "Services:"
echo "  🌐 Gateway:     http://localhost:${GATEWAY_PORT:-8000}"
echo "  🔐 Auth:        http://localhost:${AUTH_SERVICE_PORT:-8002}"
echo "  📇 Contact:     http://localhost:${CONTACT_SERVICE_PORT:-8001}"
echo "  📡 NATS:        http://localhost:8222 (monitoring)"
echo ""
echo "Health Checks:"
echo "  Gateway:  $(curl -s http://localhost:${GATEWAY_PORT:-8000}/health)"
echo "  Auth:     $(curl -s http://localhost:${AUTH_SERVICE_PORT:-8002}/health | jq -r '.status // "OK"' 2>/dev/null || echo 'OK')"
echo "  Contact:  $(curl -s http://localhost:${CONTACT_SERVICE_PORT:-8001}/health | jq -r '.status // "OK"' 2>/dev/null || echo 'OK')"
echo ""
echo "Swagger UI:"
echo "  Auth:     http://localhost:${AUTH_SERVICE_PORT:-8002}/swagger/index.html"
echo "  Contact:  http://localhost:${CONTACT_SERVICE_PORT:-8001}/swagger/index.html"
echo ""
echo "Logs:"
echo "  Auth:     tail -f /tmp/auth.log"
echo "  Contact:  tail -f /tmp/contact.log"
echo "  Gateway:  tail -f /tmp/gateway.log"
echo ""
echo "To stop all services:"
echo "  bash stop_all.sh"
echo ""

