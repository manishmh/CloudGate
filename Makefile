.PHONY: help build test clean dev up down logs frontend backend

# Default target
help:
	@echo "CloudGate SSO Portal - Development Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make dev        - Start development environment"
	@echo "  make up         - Start all services with Docker Compose"
	@echo "  make down       - Stop all services"
	@echo "  make logs       - View logs from all services"
	@echo "  make build      - Build all Docker images"
	@echo "  make test       - Run all tests"
	@echo "  make clean      - Clean up containers and volumes"
	@echo "  make frontend   - Start frontend development server"
	@echo "  make backend    - Start backend development server"
	@echo "  make lint       - Run linting for all components"
	@echo "  make typecheck  - Run TypeScript type checking"
	@echo "  make setup      - Initial setup and configuration"

# Development environment
dev: up

# Docker Compose commands
up:
	@echo "Starting CloudGate SSO Portal..."
	docker-compose up -d
	@echo "Services started. Access:"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Keycloak: http://localhost:8080"
	@echo "  Backend:  http://localhost:8081"

down:
	@echo "Stopping all services..."
	docker-compose down

logs:
	docker-compose logs -f

# Build commands
build:
	@echo "Building all Docker images..."
	docker-compose build

build-backend:
	@echo "Building backend..."
	cd backend && docker build -t cloudgate/backend .

build-frontend:
	@echo "Building frontend..."
	cd frontend && docker build -f Dockerfile.dev -t cloudgate/frontend .

# Test commands
test: test-frontend test-backend

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm run lint && npm run typecheck && npm run build

test-backend:
	@echo "Running backend tests..."
	cd backend && CGO_ENABLED=0 /usr/local/go/bin/go vet ./... && CGO_ENABLED=0 /usr/local/go/bin/go test -v ./...

# Lint commands
lint: lint-frontend lint-backend

lint-frontend:
	@echo "Linting frontend..."
	cd frontend && npm run lint && npm run typecheck

lint-backend:
	@echo "Linting backend..."
	cd backend && CGO_ENABLED=0 /usr/local/go/bin/go vet ./... && /usr/local/go/bin/go fmt ./...

# TypeScript type checking
typecheck: typecheck-frontend

typecheck-frontend:
	@echo "Running TypeScript type checking..."
	cd frontend && npm run typecheck

# Development servers
frontend:
	@echo "Starting frontend development server..."
	cd frontend && npm run dev

backend:
	@echo "Starting backend development server..."
	cd backend && CGO_ENABLED=0 /usr/local/go/bin/go run main.go

# Setup and installation
setup:
	@echo "Setting up CloudGate development environment..."
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Installing backend dependencies..."
	cd backend && /usr/local/go/bin/go mod tidy
	@echo "Copying environment template..."
	cp frontend/env.example frontend/.env.local
	@echo "Setup complete! Run 'make up' to start the services."

# Cleanup commands
clean:
	@echo "Cleaning up containers and volumes..."
	docker-compose down -v
	docker system prune -f

clean-all: clean
	@echo "Removing all CloudGate images..."
	docker images | grep cloudgate | awk '{print $$3}' | xargs -r docker rmi -f

# Health checks
health:
	@echo "Checking service health..."
	@echo "Backend health:"
	@curl -s http://localhost:8081/health | jq . || echo "Backend not responding"
	@echo "Keycloak health:"
	@curl -s http://localhost:8080/health/ready || echo "Keycloak not responding"

# Database commands
db-reset:
	@echo "Resetting database..."
	docker-compose down postgres
	docker volume rm cloudgate_postgres_data || true
	docker-compose up -d postgres

# Keycloak management
keycloak-admin:
	@echo "Keycloak Admin Console: http://localhost:8080/admin"
	@echo "Username: admin"
	@echo "Password: admin_password"

# Production commands
prod-build:
	@echo "Building production images..."
	docker build -t cloudgate/backend:latest ./backend
	docker build -t cloudgate/frontend:latest ./frontend

# Monitoring
monitor:
	@echo "Monitoring services..."
	watch -n 2 'docker-compose ps && echo "" && docker stats --no-stream' 