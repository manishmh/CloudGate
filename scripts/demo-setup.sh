#!/bin/bash

# CloudGate Demo Setup Script
# This script sets up the demo environment and tests the new SaaS application functionality

set -e

echo "🚀 CloudGate Demo Setup"
echo "======================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

print_status "Docker is running"

# Navigate to project root
cd "$(dirname "$0")/.."

# Start services with Docker Compose
print_info "Starting CloudGate services..."
docker-compose up -d

# Wait for services to be healthy
print_info "Waiting for services to be ready..."

# Wait for Keycloak
print_info "Waiting for Keycloak to be ready..."
timeout=300
counter=0
while ! curl -s http://localhost:8080/health/ready > /dev/null 2>&1; do
    if [ $counter -ge $timeout ]; then
        print_error "Keycloak failed to start within $timeout seconds"
        exit 1
    fi
    sleep 2
    counter=$((counter + 2))
    echo -n "."
done
echo ""
print_status "Keycloak is ready"

# Wait for Backend
print_info "Waiting for Backend API to be ready..."
timeout=60
counter=0
while ! curl -s http://localhost:8081/health > /dev/null 2>&1; do
    if [ $counter -ge $timeout ]; then
        print_error "Backend API failed to start within $timeout seconds"
        exit 1
    fi
    sleep 2
    counter=$((counter + 2))
    echo -n "."
done
echo ""
print_status "Backend API is ready"

# Wait for Frontend
print_info "Waiting for Frontend to be ready..."
timeout=60
counter=0
while ! curl -s http://localhost:3000 > /dev/null 2>&1; do
    if [ $counter -ge $timeout ]; then
        print_error "Frontend failed to start within $timeout seconds"
        exit 1
    fi
    sleep 2
    counter=$((counter + 2))
    echo -n "."
done
echo ""
print_status "Frontend is ready"

# Test API endpoints
print_info "Testing API endpoints..."

# Test health endpoint
if curl -s http://localhost:8081/health | grep -q "healthy"; then
    print_status "Health endpoint working"
else
    print_error "Health endpoint failed"
fi

# Test SaaS apps endpoint
if curl -s -H "Authorization: Bearer demo-token" http://localhost:8081/apps | grep -q "apps"; then
    print_status "SaaS apps endpoint working"
else
    print_error "SaaS apps endpoint failed"
fi

# Test API info endpoint
if curl -s http://localhost:8081/api/info | grep -q "CloudGate"; then
    print_status "API info endpoint working"
else
    print_error "API info endpoint failed"
fi

echo ""
echo "🎉 CloudGate Demo Environment Ready!"
echo "===================================="
echo ""
echo "📋 Service URLs:"
echo "   Frontend:        http://localhost:3000"
echo "   Backend API:     http://localhost:8081"
echo "   Keycloak Admin:  http://localhost:8080/admin"
echo ""
echo "🔑 Demo Credentials:"
echo "   Keycloak Admin:  admin / admin_password"
echo "   Test User:       testuser / password123"
echo ""
echo "🧪 API Testing:"
echo "   Health Check:    curl http://localhost:8081/health"
echo "   SaaS Apps:       curl -H 'Authorization: Bearer demo-token' http://localhost:8081/apps"
echo "   API Info:        curl http://localhost:8081/api/info"
echo ""
echo "🚀 New Features Implemented:"
echo "   ✅ Real SaaS application management"
echo "   ✅ Functional Connect/Launch buttons"
echo "   ✅ Backend API for app management"
echo "   ✅ OAuth flow simulation"
echo "   ✅ Dynamic app status updates"
echo ""
echo "📖 Next Steps:"
echo "   1. Open http://localhost:3000 in your browser"
echo "   2. Click 'Sign in with SSO'"
echo "   3. Use testuser / password123 to login"
echo "   4. Test the SaaS application Connect/Launch functionality"
echo ""
print_info "To stop the demo: docker-compose down"
print_info "To view logs: docker-compose logs -f" 