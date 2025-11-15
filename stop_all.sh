#!/bin/bash

echo "🛑 Stopping all services..."
echo ""

# Kill services by PID if they exist
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

if [ -f /tmp/gateway.pid ]; then
    GATEWAY_PID=$(cat /tmp/gateway.pid)
    if ps -p $GATEWAY_PID > /dev/null 2>&1; then
        kill $GATEWAY_PID 2>/dev/null || true
        echo "✅ Gateway stopped"
    fi
    rm -f /tmp/gateway.pid
fi

# Kill any remaining processes on the ports
lsof -ti:8000 | xargs kill -9 2>/dev/null || true
lsof -ti:8001 | xargs kill -9 2>/dev/null || true
lsof -ti:8002 | xargs kill -9 2>/dev/null || true

# Stop Docker services
cd /home/admn/Documents/project/SR-BE-interview-1
docker compose stop db nats 2>/dev/null || true

echo ""
echo "✅ All services stopped"
echo ""

