#!/bin/bash

# CloudGate Backend Start Script

echo "🚀 Starting CloudGate Backend Server..."

# Build the application
echo "📦 Building application..."
go build -o cloudgate-backend .

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"

# Set default environment variables if not set
export PORT=${PORT:-8081}
export KEYCLOAK_URL=${KEYCLOAK_URL:-http://localhost:8080}
export KEYCLOAK_REALM=${KEYCLOAK_REALM:-cloudgate}
export KEYCLOAK_CLIENT_ID=${KEYCLOAK_CLIENT_ID:-cloudgate-frontend}
export ALLOWED_ORIGINS=${ALLOWED_ORIGINS:-http://localhost:3000}

echo "🔧 Configuration:"
echo "   Port: $PORT"
echo "   Keycloak URL: $KEYCLOAK_URL"
echo "   Keycloak Realm: $KEYCLOAK_REALM"
echo "   Client ID: $KEYCLOAK_CLIENT_ID"
echo "   Allowed Origins: $ALLOWED_ORIGINS"
echo ""

# Start the server
echo "🌟 Starting server on port $PORT..."
./cloudgate-backend 