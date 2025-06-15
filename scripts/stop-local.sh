#!/bin/bash

# Stop CloudGate Local Development Services

echo "🛑 Stopping CloudGate Local Development Services"
echo "==============================================="

# Stop Docker services
echo "Stopping Docker services..."
docker compose down

# Stop frontend if running
echo "Stopping frontend development server..."
pkill -f "npm run dev" 2>/dev/null || true
pkill -f "next dev" 2>/dev/null || true

echo "✅ All services stopped!"
echo ""
echo "To restart, run: ./scripts/local-dev.sh" 