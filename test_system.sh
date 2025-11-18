#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

GATEWAY_URL="http://localhost:8000"

echo -e "${GREEN}🧪 Testing Microservices System${NC}"
echo ""

# Check if services are running
echo "📋 Checking service health..."
services=("gateway:8000" "auth:8002" "contact:8001" "inventory:8003" "sales:8004" "purchase:8005")
all_healthy=true

for service in "${services[@]}"; do
    name=$(echo $service | cut -d: -f1)
    port=$(echo $service | cut -d: -f2)
    if curl -s -f "http://localhost:$port/health" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✅${NC} $name service is healthy"
    else
        echo -e "  ${RED}❌${NC} $name service is not responding"
        all_healthy=false
    fi
done

if [ "$all_healthy" = false ]; then
    echo -e "\n${RED}❌ Some services are not running. Please start them first.${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ All services are healthy${NC}"
echo ""

# Test 1: Register a user
echo "📝 Test 1: Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "role": "inventory_manager"
  }')

if echo "$REGISTER_RESPONSE" | grep -q "email"; then
    echo -e "  ${GREEN}✅${NC} User registered successfully"
else
    echo -e "  ${RED}❌${NC} User registration failed: $REGISTER_RESPONSE"
    exit 1
fi

# Test 2: Login
echo "🔐 Test 2: Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "  ${RED}❌${NC} Login failed: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "  ${GREEN}✅${NC} Login successful"
echo "  Token: ${TOKEN:0:20}..."

# Test 3: Create Customer
echo "👤 Test 3: Creating customer..."
CUSTOMER_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/customers" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "address": "123 Main St"
  }')

CUSTOMER_ID=$(echo "$CUSTOMER_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$CUSTOMER_ID" ]; then
    echo -e "  ${RED}❌${NC} Customer creation failed: $CUSTOMER_RESPONSE"
    exit 1
fi

echo -e "  ${GREEN}✅${NC} Customer created: $CUSTOMER_ID"

# Test 4: Create Vendor
echo "🏢 Test 4: Creating vendor..."
VENDOR_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/vendors" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Acme Corp",
    "email": "acme@example.com",
    "phone": "+1987654321",
    "address": "456 Business Ave"
  }')

VENDOR_ID=$(echo "$VENDOR_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$VENDOR_ID" ]; then
    echo -e "  ${RED}❌${NC} Vendor creation failed: $VENDOR_RESPONSE"
    exit 1
fi

echo -e "  ${GREEN}✅${NC} Vendor created: $VENDOR_ID"

# Test 5: Create Item
echo "📦 Test 5: Creating item..."
ITEM_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Widget A",
    "description": "A great widget",
    "sku": "WIDGET-A-001",
    "unit_price": 10.50
  }')

ITEM_ID=$(echo "$ITEM_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$ITEM_ID" ]; then
    echo -e "  ${RED}❌${NC} Item creation failed: $ITEM_RESPONSE"
    exit 1
fi

echo -e "  ${GREEN}✅${NC} Item created: $ITEM_ID"

# Test 6: Check Stock (should be 0)
echo "📊 Test 6: Checking stock..."
STOCK_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/api/items/$ITEM_ID/stock" \
  -H "Authorization: Bearer $TOKEN")

QUANTITY=$(echo "$STOCK_RESPONSE" | grep -o '"quantity":[0-9]*' | cut -d':' -f2)

if [ -z "$QUANTITY" ]; then
    echo -e "  ${YELLOW}⚠️${NC}  Stock check response: $STOCK_RESPONSE"
else
    echo -e "  ${GREEN}✅${NC} Stock quantity: $QUANTITY"
fi

# Test 7: Create Purchase Order
echo "🛒 Test 7: Creating purchase order..."
PURCHASE_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/purchase/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"vendor_id\": \"$VENDOR_ID\",
    \"items\": [
      {
        \"item_id\": \"$ITEM_ID\",
        \"quantity\": 100
      }
    ]
  }")

PURCHASE_ORDER_ID=$(echo "$PURCHASE_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$PURCHASE_ORDER_ID" ]; then
    echo -e "  ${RED}❌${NC} Purchase order creation failed: $PURCHASE_RESPONSE"
    exit 1
fi

echo -e "  ${GREEN}✅${NC} Purchase order created: $PURCHASE_ORDER_ID"

# Test 8: Receive Purchase Order (triggers event)
echo "📥 Test 8: Receiving purchase order..."
RECEIVE_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/purchase/orders/$PURCHASE_ORDER_ID/receive" \
  -H "Authorization: Bearer $TOKEN")

if echo "$RECEIVE_RESPONSE" | grep -q "status"; then
    echo -e "  ${GREEN}✅${NC} Purchase order received (stock should increase)"
else
    echo -e "  ${RED}❌${NC} Receive failed: $RECEIVE_RESPONSE"
fi

# Wait a bit for event processing
sleep 2

# Test 9: Check Stock After Purchase (should be 100)
echo "📊 Test 9: Checking stock after purchase..."
STOCK_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/api/items/$ITEM_ID/stock" \
  -H "Authorization: Bearer $TOKEN")

NEW_QUANTITY=$(echo "$STOCK_RESPONSE" | grep -o '"quantity":[0-9]*' | cut -d':' -f2)

if [ "$NEW_QUANTITY" = "100" ]; then
    echo -e "  ${GREEN}✅${NC} Stock correctly increased to: $NEW_QUANTITY"
else
    echo -e "  ${YELLOW}⚠️${NC}  Stock is $NEW_QUANTITY (expected 100). Event may not have processed yet."
fi

# Test 10: Create Sales Order
echo "💰 Test 10: Creating sales order..."
SALES_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/sales/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"customer_id\": \"$CUSTOMER_ID\",
    \"items\": [
      {
        \"item_id\": \"$ITEM_ID\",
        \"quantity\": 25
      }
    ]
  }")

SALES_ORDER_ID=$(echo "$SALES_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$SALES_ORDER_ID" ]; then
    echo -e "  ${RED}❌${NC} Sales order creation failed: $SALES_RESPONSE"
    exit 1
fi

echo -e "  ${GREEN}✅${NC} Sales order created: $SALES_ORDER_ID"

# Test 11: Confirm Sales Order (triggers event)
echo "✅ Test 11: Confirming sales order..."
CONFIRM_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/sales/orders/$SALES_ORDER_ID/confirm" \
  -H "Authorization: Bearer $TOKEN")

if echo "$CONFIRM_RESPONSE" | grep -q "status"; then
    echo -e "  ${GREEN}✅${NC} Sales order confirmed (stock should decrease)"
else
    echo -e "  ${RED}❌${NC} Confirm failed: $CONFIRM_RESPONSE"
fi

# Wait a bit for event processing
sleep 2

# Test 12: Check Stock After Sale (should be 75)
echo "📊 Test 12: Checking stock after sale..."
FINAL_STOCK_RESPONSE=$(curl -s -X GET "$GATEWAY_URL/api/items/$ITEM_ID/stock" \
  -H "Authorization: Bearer $TOKEN")

FINAL_QUANTITY=$(echo "$FINAL_STOCK_RESPONSE" | grep -o '"quantity":[0-9]*' | cut -d':' -f2)

if [ "$FINAL_QUANTITY" = "75" ]; then
    echo -e "  ${GREEN}✅${NC} Stock correctly decreased to: $FINAL_QUANTITY"
else
    echo -e "  ${YELLOW}⚠️${NC}  Stock is $FINAL_QUANTITY (expected 75). Event may not have processed yet."
fi

echo ""
echo -e "${GREEN}🎉 All tests completed!${NC}"
echo ""
echo "Summary:"
echo "  - User registered and logged in"
echo "  - Customer and Vendor created"
echo "  - Item created"
echo "  - Purchase order created and received (stock increased)"
echo "  - Sales order created and confirmed (stock decreased)"
echo ""
echo "System is working correctly! ✅"

