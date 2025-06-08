#!/bin/bash

# CloudGate Development Database Setup
# This script sets up PostgreSQL for local development
# Production uses Neon DB (neon.tech) - a serverless PostgreSQL platform

set -e

echo "🚀 Setting up CloudGate Development Database (PostgreSQL)"

# Configuration
DB_NAME="cloudgate_dev"
DB_USER="cloudgate"
DB_PASSWORD="cloudgate_dev_password"
DB_HOST="localhost"
DB_PORT="5432"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ PostgreSQL is not installed. Please install PostgreSQL first.${NC}"
    echo -e "${YELLOW}Ubuntu/Debian: sudo apt-get install postgresql postgresql-contrib${NC}"
    echo -e "${YELLOW}macOS: brew install postgresql${NC}"
    echo -e "${YELLOW}CentOS/RHEL: sudo yum install postgresql-server postgresql-contrib${NC}"
    exit 1
fi

# Check if PostgreSQL service is running
if ! pg_isready -h $DB_HOST -p $DB_PORT &> /dev/null; then
    echo -e "${RED}❌ PostgreSQL service is not running. Please start PostgreSQL service.${NC}"
    echo -e "${YELLOW}Ubuntu/Debian: sudo systemctl start postgresql${NC}"
    echo -e "${YELLOW}macOS: brew services start postgresql${NC}"
    echo -e "${YELLOW}CentOS/RHEL: sudo systemctl start postgresql${NC}"
    exit 1
fi

echo -e "${GREEN}✅ PostgreSQL is installed and running${NC}"

# Create database user if it doesn't exist
echo "📝 Creating database user: $DB_USER"
sudo -u postgres psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || echo -e "${YELLOW}⚠️  User $DB_USER already exists${NC}"

# Grant privileges to user
echo "🔐 Granting privileges to user: $DB_USER"
sudo -u postgres psql -c "ALTER USER $DB_USER CREATEDB;"

# Create database if it doesn't exist
echo "🗄️  Creating database: $DB_NAME"
sudo -u postgres psql -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;" 2>/dev/null || echo -e "${YELLOW}⚠️  Database $DB_NAME already exists${NC}"

# Grant all privileges on database
echo "🔑 Granting all privileges on database: $DB_NAME"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

# Enable UUID extension (needed for our models)
echo "🔧 Enabling UUID extension"
sudo -u postgres psql -d $DB_NAME -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

# Test connection
echo "🧪 Testing database connection"
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version();" &> /dev/null; then
    echo -e "${GREEN}✅ Database connection successful!${NC}"
else
    echo -e "${RED}❌ Database connection failed!${NC}"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f "../.env" ]; then
    echo "📄 Creating .env file from template"
    cp ../env.example ../.env
    echo -e "${YELLOW}⚠️  Please update the .env file with your specific configuration${NC}"
else
    echo -e "${YELLOW}⚠️  .env file already exists. Please ensure database configuration matches:${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Development database setup complete!${NC}"
echo ""
echo "Database Configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Password: $DB_PASSWORD"
echo ""
echo "Add these to your .env file:"
echo "  DB_TYPE=postgres"
echo "  DB_HOST=$DB_HOST"
echo "  DB_PORT=$DB_PORT"
echo "  DB_USER=$DB_USER"
echo "  DB_PASSWORD=$DB_PASSWORD"
echo "  DB_NAME=$DB_NAME"
echo "  DB_SSLMODE=disable"
echo ""
echo -e "${GREEN}For production, use Neon DB (neon.tech) with:${NC}"
echo "  NEON_DATABASE_URL=postgresql://username:password@ep-xxx-xxx.us-east-2.aws.neon.tech/dbname?sslmode=require"
echo ""
echo -e "${GREEN}✅ Ready to run: go run . or ./cloudgate-backend${NC}" 