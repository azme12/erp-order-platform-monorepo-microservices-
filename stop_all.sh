#!/bin/bash

echo "🛑 Stopping all services..."
echo ""


if [ -f /tmp/auth.pid ]; then
    AUTH_PID=$(cat /tmp/auth.pid)
    if ps -p $AUTH_PID > /dev/null 2>&1; then
        kill $AUTH_PID 2>/dev/null || true
        echo "✅ Auth Service stopped"
    fi
    rm -f /tmp/auth.pid
fi

if [ -f /tmp/contact.pid ]; then
    CONTACT_PID=$(cat /tmp/contact.pid)
    if ps -p $CONTACT_PID > /dev/null 2>&1; then
        kill $CONTACT_PID 2>/dev/null || true
        echo "✅ Contact Service stopped"
    fi
    rm -f /tmp/contact.pid
fi

if [ -f /tmp/inventory.pid ]; then
    INVENTORY_PID=$(cat /tmp/inventory.pid)
    if ps -p $INVENTORY_PID > /dev/null 2>&1; then
        kill $INVENTORY_PID 2>/dev/null || true
        echo "✅ Inventory Service stopped"
    fi
    rm -f /tmp/inventory.pid
fi

if [ -f /tmp/sales.pid ]; then
    SALES_PID=$(cat /tmp/sales.pid)
    if ps -p $SALES_PID > /dev/null 2>&1; then
        kill $SALES_PID 2>/dev/null || true
        echo "✅ Sales Service stopped"
    fi
    rm -f /tmp/sales.pid
fi

if [ -f /tmp/purchase.pid ]; then
    PURCHASE_PID=$(cat /tmp/purchase.pid)
    if ps -p $PURCHASE_PID > /dev/null 2>&1; then
        kill $PURCHASE_PID 2>/dev/null || true
        echo "✅ Purchase Service stopped"
    fi
    rm -f /tmp/purchase.pid
fi

if [ -f /tmp/gateway.pid ]; then
    GATEWAY_PID=$(cat /tmp/gateway.pid)
    if ps -p $GATEWAY_PID > /dev/null 2>&1; then
        kill $GATEWAY_PID 2>/dev/null || true
        echo "✅ Gateway stopped"
    fi
    rm -f /tmp/gateway.pid
fi


lsof -ti:8000 | xargs kill -9 2>/dev/null || true
lsof -ti:8001 | xargs kill -9 2>/dev/null || true
lsof -ti:8002 | xargs kill -9 2>/dev/null || true
lsof -ti:8003 | xargs kill -9 2>/dev/null || true
lsof -ti:8004 | xargs kill -9 2>/dev/null || true
lsof -ti:8005 | xargs kill -9 2>/dev/null || true


cd /home/admn/Documents/project/SR-BE-interview-1
docker compose stop db nats 2>/dev/null || true

echo ""
echo "✅ All services stopped"
echo ""

