#!/bin/bash

# Test script to simulate Cloud Run environment locally
echo "🧪 Testing CloudGate Backend in Cloud Run simulation mode"

# Set Cloud Run environment variables
export PORT=8080
export GIN_MODE=release
export SKIP_DEMO_USER=true
export RUN_MIGRATIONS=false

# Test database connection (use SQLite for quick testing)
export DB_TYPE=sqlite
export DB_NAME=test_cloudgate.db

# Basic Keycloak configuration (you'll need to set these for your deployment)
export KEYCLOAK_URL=${KEYCLOAK_URL:-"http://localhost:8080"}
export KEYCLOAK_REALM=${KEYCLOAK_REALM:-"cloudgate"}
export KEYCLOAK_CLIENT_ID=${KEYCLOAK_CLIENT_ID:-"cloudgate-frontend"}

# CORS configuration
export ALLOWED_ORIGINS="*"

echo "🔧 Environment variables set:"
echo "   PORT=$PORT"
echo "   GIN_MODE=$GIN_MODE"
echo "   DB_TYPE=$DB_TYPE"
echo "   KEYCLOAK_URL=$KEYCLOAK_URL"

echo ""
echo "🚀 Starting backend..."
echo "   This simulates how the backend will run in Cloud Run"
echo "   Press Ctrl+C to stop"
echo ""

# Run the backend
go run ./main.go 