#!/bin/bash

# CloudGate Full Local Development Setup
# This script sets up PostgreSQL, Keycloak, Backend, and Frontend all locally

set -e

echo "🚀 CloudGate Full Local Development Setup"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker compose &> /dev/null; then
    print_error "Docker Compose is not available. Please install Docker Compose."
    exit 1
fi

# Check if Node.js is available
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js 18+ first."
    exit 1
fi

print_status "Starting local development environment..."

# Stop any existing containers
print_status "Stopping existing containers..."
docker compose down 2>/dev/null || true

# Start backend services (PostgreSQL, Keycloak, Backend API)
print_status "Starting backend services (PostgreSQL, Keycloak, Backend)..."
docker compose up -d postgres keycloak backend

# Wait for services to be ready
print_status "Waiting for services to start..."
echo "This may take 2-3 minutes for first-time setup..."

# Wait for PostgreSQL
print_status "Waiting for PostgreSQL..."
timeout=60
while ! docker exec cloudgate-postgres pg_isready -U keycloak >/dev/null 2>&1; do
    sleep 2
    timeout=$((timeout - 2))
    if [ $timeout -le 0 ]; then
        print_error "PostgreSQL failed to start within 60 seconds"
        docker logs cloudgate-postgres
        exit 1
    fi
done
print_success "PostgreSQL is ready!"

# Wait for Keycloak
print_status "Waiting for Keycloak (this takes longer on first run)..."
timeout=180
while ! curl -f http://localhost:8080/health/ready >/dev/null 2>&1; do
    sleep 5
    timeout=$((timeout - 5))
    if [ $timeout -le 0 ]; then
        print_error "Keycloak failed to start within 3 minutes"
        docker logs cloudgate-keycloak
        exit 1
    fi
    echo -n "."
done
echo ""
print_success "Keycloak is ready!"

# Wait for Backend
print_status "Waiting for Backend API..."
timeout=60
while ! curl -f http://localhost:8081/health >/dev/null 2>&1; do
    sleep 2
    timeout=$((timeout - 2))
    if [ $timeout -le 0 ]; then
        print_error "Backend API failed to start within 60 seconds"
        docker logs cloudgate-backend
        exit 1
    fi
done
print_success "Backend API is ready!"

# Setup frontend dependencies
print_status "Setting up frontend dependencies..."
cd frontend
if [ ! -d "node_modules" ]; then
    print_status "Installing frontend dependencies..."
    npm install
fi

print_success "All services are ready!"
echo ""
echo "🎉 Local Development Environment Started Successfully!"
echo "=================================================="
echo ""
echo "📱 Services:"
echo "   Frontend:  http://localhost:3000 (will start when you run npm run dev)"
echo "   Backend:   http://localhost:8081"
echo "   Keycloak:  http://localhost:8080"
echo "   Database:  localhost:5432"
echo ""
echo "🔐 Keycloak Admin Console:"
echo "   URL:      http://localhost:8080/admin"
echo "   Username: admin"
echo "   Password: admin_password"
echo ""
echo "🚀 To start frontend development:"
echo "   cd frontend"
echo "   npm run dev"
echo ""
echo "🛠️ Development Commands:"
echo "   View logs:     docker compose logs -f [service]"
echo "   Stop services: docker compose down"
echo "   Restart:       ./scripts/local-dev.sh"
echo ""
echo "📋 Next Steps:"
echo "1. Open a new terminal"
echo "2. cd frontend"
echo "3. npm run dev"
echo "4. Open http://localhost:3000 in your browser"
echo ""
print_success "Happy coding! 🎉" 