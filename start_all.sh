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

# Step 1.5: Run Database Migrations
echo "🗄️  Step 1.5: Running database migrations..."
if command -v psql > /dev/null; then
    # Run Auth migrations
    echo "  Running Auth migrations..."
    PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST:-localhost} -p ${DB_PORT:-5432} -U ${DB_USER} -d ${DB_NAME} -f migrations/auth/000001_create_users_table.up.sql > /dev/null 2>&1 || echo "  Auth migration may already exist"
    
    # Run Contact migrations
    echo "  Running Contact migrations..."
    PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST:-localhost} -p ${DB_PORT:-5432} -U ${DB_USER} -d ${DB_NAME} -f migrations/contact/000001_create_customers_vendors_tables.up.sql > /dev/null 2>&1 || echo "  Contact migration may already exist"
    
    # Run Inventory migrations
    echo "  Running Inventory migrations..."
    PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST:-localhost} -p ${DB_PORT:-5432} -U ${DB_USER} -d ${DB_NAME} -f migrations/inventory/000001_create_items_stock_tables.up.sql > /dev/null 2>&1 || echo "  Inventory migration may already exist"
    
    # Run Sales migrations
    echo "  Running Sales migrations..."
    PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST:-localhost} -p ${DB_PORT:-5432} -U ${DB_USER} -d ${DB_NAME} -f migrations/sales/000001_create_sales_orders_order_items_tables.up.sql > /dev/null 2>&1 || echo "  Sales migration may already exist"
    
    # Run Purchase migrations
    echo "  Running Purchase migrations..."
    PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST:-localhost} -p ${DB_PORT:-5432} -U ${DB_USER} -d ${DB_NAME} -f migrations/purchase/000001_create_purchase_orders_purchase_order_items_tables.up.sql > /dev/null 2>&1 || echo "  Purchase migration may already exist"
    
    echo "✅ Migrations completed"
else
    echo "⚠️  psql not found, skipping migrations. Please run migrations manually."
fi
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

# Step 4: Start Inventory Service
echo "📦 Step 4: Starting Inventory Service on port ${INVENTORY_SERVICE_PORT:-8003}..."
PORT=${INVENTORY_SERVICE_PORT:-8003} \
DB_HOST=${DB_HOST} \
DB_PORT=${DB_PORT} \
DB_USER=${DB_USER} \
DB_PASSWORD=${DB_PASSWORD} \
DB_NAME=${DB_NAME} \
JWT_SECRET=${JWT_SECRET} \
NATS_URL=${NATS_URL} \
go run services/inventory/cmd/main.go > /tmp/inventory.log 2>&1 &
INVENTORY_PID=$!
echo $INVENTORY_PID > /tmp/inventory.pid
sleep 3

if curl -s http://localhost:${INVENTORY_SERVICE_PORT:-8003}/health > /dev/null; then
    echo "✅ Inventory Service started (PID: $INVENTORY_PID)"
else
    echo "❌ Inventory Service failed to start. Check /tmp/inventory.log"
    exit 1
fi
echo ""

# Step 5: Start Sales Service
echo "💰 Step 5: Starting Sales Service on port ${SALES_SERVICE_PORT:-8004}..."
PORT=${SALES_SERVICE_PORT:-8004} \
DB_HOST=${DB_HOST} \
DB_PORT=${DB_PORT} \
DB_USER=${DB_USER} \
DB_PASSWORD=${DB_PASSWORD} \
DB_NAME=${DB_NAME} \
JWT_SECRET=${JWT_SECRET} \
NATS_URL=${NATS_URL} \
CONTACT_SERVICE_URL=${CONTACT_SERVICE_URL:-http://localhost:8001} \
INVENTORY_SERVICE_URL=${INVENTORY_SERVICE_URL:-http://localhost:8003} \
go run services/sales/cmd/main.go > /tmp/sales.log 2>&1 &
SALES_PID=$!
echo $SALES_PID > /tmp/sales.pid
sleep 3

if curl -s http://localhost:${SALES_SERVICE_PORT:-8004}/health > /dev/null; then
    echo "✅ Sales Service started (PID: $SALES_PID)"
else
    echo "❌ Sales Service failed to start. Check /tmp/sales.log"
    exit 1
fi
echo ""

# Step 6: Start Purchase Service
echo "🛒 Step 6: Starting Purchase Service on port ${PURCHASE_SERVICE_PORT:-8005}..."
PORT=${PURCHASE_SERVICE_PORT:-8005} \
DB_HOST=${DB_HOST} \
DB_PORT=${DB_PORT} \
DB_USER=${DB_USER} \
DB_PASSWORD=${DB_PASSWORD} \
DB_NAME=${DB_NAME} \
JWT_SECRET=${JWT_SECRET} \
NATS_URL=${NATS_URL} \
CONTACT_SERVICE_URL=${CONTACT_SERVICE_URL:-http://localhost:8001} \
INVENTORY_SERVICE_URL=${INVENTORY_SERVICE_URL:-http://localhost:8003} \
go run services/purchase/cmd/main.go > /tmp/purchase.log 2>&1 &
PURCHASE_PID=$!
echo $PURCHASE_PID > /tmp/purchase.pid
sleep 3

if curl -s http://localhost:${PURCHASE_SERVICE_PORT:-8005}/health > /dev/null; then
    echo "✅ Purchase Service started (PID: $PURCHASE_PID)"
else
    echo "❌ Purchase Service failed to start. Check /tmp/purchase.log"
    exit 1
fi
echo ""

# Step 7: Start Gateway
echo "🌐 Step 7: Starting Gateway on port ${GATEWAY_PORT:-8000}..."
PORT=${GATEWAY_PORT:-8000} \
AUTH_SERVICE_URL=${AUTH_SERVICE_URL} \
CONTACT_SERVICE_URL=${CONTACT_SERVICE_URL} \
INVENTORY_SERVICE_URL=${INVENTORY_SERVICE_URL} \
SALES_SERVICE_URL=${SALES_SERVICE_URL:-http://localhost:8004} \
PURCHASE_SERVICE_URL=${PURCHASE_SERVICE_URL:-http://localhost:8005} \
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
echo "  📦 Inventory:   http://localhost:${INVENTORY_SERVICE_PORT:-8003}"
echo "  💰 Sales:       http://localhost:${SALES_SERVICE_PORT:-8004}"
echo "  🛒 Purchase:    http://localhost:${PURCHASE_SERVICE_PORT:-8005}"
echo "  📡 NATS:        http://localhost:8222 (monitoring)"
echo ""
echo "Health Checks:"
echo "  Gateway:   $(curl -s http://localhost:${GATEWAY_PORT:-8000}/health)"
echo "  Auth:      $(curl -s http://localhost:${AUTH_SERVICE_PORT:-8002}/health | jq -r '.status // "OK"' 2>/dev/null || echo 'OK')"
echo "  Contact:   $(curl -s http://localhost:${CONTACT_SERVICE_PORT:-8001}/health | jq -r '.status // "OK"' 2>/dev/null || echo 'OK')"
echo "  Inventory: $(curl -s http://localhost:${INVENTORY_SERVICE_PORT:-8003}/health | jq -r '.status // "OK"' 2>/dev/null || echo 'OK')"
echo "  Sales:     $(curl -s http://localhost:${SALES_SERVICE_PORT:-8004}/health | jq -r '.status // "OK"' 2>/dev/null || echo 'OK')"
echo "  Purchase:  $(curl -s http://localhost:${PURCHASE_SERVICE_PORT:-8005}/health | jq -r '.status // "OK"' 2>/dev/null || echo 'OK')"
echo ""
echo "Swagger UI:"
echo "  Auth:      http://localhost:${AUTH_SERVICE_PORT:-8002}/swagger/index.html"
echo "  Contact:   http://localhost:${CONTACT_SERVICE_PORT:-8001}/swagger/index.html"
echo "  Inventory: http://localhost:${INVENTORY_SERVICE_PORT:-8003}/swagger/index.html"
echo "  Sales:     http://localhost:${SALES_SERVICE_PORT:-8004}/swagger/index.html"
echo "  Purchase:  http://localhost:${PURCHASE_SERVICE_PORT:-8005}/swagger/index.html"
echo ""
echo "Logs:"
echo "  Auth:      tail -f /tmp/auth.log"
echo "  Contact:   tail -f /tmp/contact.log"
echo "  Inventory: tail -f /tmp/inventory.log"
echo "  Sales:     tail -f /tmp/sales.log"
echo "  Purchase:  tail -f /tmp/purchase.log"
echo "  Gateway:   tail -f /tmp/gateway.log"
echo ""
echo "To stop all services:"
echo "  bash stop_all.sh"
echo ""

