#!/bin/bash

# Exit on error
set -e

echo "🚀 Setting up SaaS Billing System..."

# Check for required tools
echo "📝 Checking prerequisites..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.18 or later."
    exit 1
fi

# Check for PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed. Please install PostgreSQL 12 or later."
    exit 1
fi

# Check for environment file
if [ ! -f .env ]; then
    echo "📄 Creating .env file from example..."
    cp .env.example .env
    echo "⚠️ Please edit .env file with your configuration"
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Create database
echo "🗄️ Setting up database..."
DB_NAME=$(grep DB_NAME .env | cut -d '=' -f2)
DB_USER=$(grep DB_USER .env | cut -d '=' -f2)
DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d '=' -f2)

# Create database if it doesn't exist
psql -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;" || true

# Run migrations
echo "🔄 Running database migrations..."
# Assuming migrations are stored in SQL files
for migration in internal/db/migrations/*.sql; do
    psql -U "$DB_USER" -d "$DB_NAME" -f "$migration"
done

# Create initial admin user
echo "👤 Creating initial admin user..."
go run cmd/api/main.go seed

echo "✅ Setup complete!"
echo ""
echo "To start the server, run:"
echo "go run cmd/api/main.go"
echo ""
echo "The API will be available at http://localhost:8080"
