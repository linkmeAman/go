#!/bin/bash

echo "🚀 Starting SaaS Billing System..."

# Check Go version and install Go 1.21 if needed
REQUIRED_GO_VERSION="1.21"
CURRENT_GO_VERSION=$(go version 2>/dev/null | grep -oP "go\K[0-9]+\.[0-9]+")

if [ -z "$CURRENT_GO_VERSION" ] || [ "$CURRENT_GO_VERSION" != "$REQUIRED_GO_VERSION" ]; then
    echo "📥 Installing Go $REQUIRED_GO_VERSION..."
    # Download and install Go 1.21
    wget -q https://go.dev/dl/go${REQUIRED_GO_VERSION}.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${REQUIRED_GO_VERSION}.0.linux-amd64.tar.gz
    rm go${REQUIRED_GO_VERSION}.0.linux-amd64.tar.gz
    
    # Add Go to PATH if not already there
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    fi
    export PATH=$PATH:/usr/local/go/bin
    
    # Verify installation
    NEW_GO_VERSION=$(/usr/local/go/bin/go version | grep -oP "go\K[0-9]+\.[0-9]+")
    if [ "$NEW_GO_VERSION" != "$REQUIRED_GO_VERSION" ]; then
        echo "❌ Failed to install Go $REQUIRED_GO_VERSION"
        exit 1
    fi
    echo "✅ Go $REQUIRED_GO_VERSION installed successfully"
fi

# Check for environment file
if [ ! -f ".env" ]; then
    if [ -f "env.example" ]; then
        echo "📄 Creating .env file from example..."
        cp env.example .env
    else
        echo "❌ No .env file found and no env.example to copy from!"
        exit 1
    fi
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first!"
    exit 1
fi

echo "🔄 Starting database and Redis..."
docker-compose up -d

echo "⏳ Waiting for database to be ready..."
sleep 5

echo "🗄️ Setting up database..."
# Try to get Postgres container ID
PG_CONTAINER=$(docker-compose ps -q postgres)
if [ -z "$PG_CONTAINER" ]; then
    echo "❌ Postgres container not found"
    exit 1
fi

# Create database (ignore error if it exists)
docker exec -i $PG_CONTAINER psql -U postgres -c "CREATE DATABASE saas_billing;" > /dev/null 2>&1 || true

# Enable pgcrypto extension
docker exec -i $PG_CONTAINER psql -U postgres -d saas_billing -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ Failed to create pgcrypto extension"
    exit 1
fi

# Run schema migrations
echo "📋 Running database migrations..."
docker cp migrations/001_initial_schema.sql $PG_CONTAINER:/tmp/
docker exec -i $PG_CONTAINER psql -U postgres -d saas_billing -f /tmp/001_initial_schema.sql
if [ $? -ne 0 ]; then
    echo "❌ Database migration failed"
    exit 1
fi

echo "📦 Downloading dependencies..."
/usr/local/go/bin/go mod download
if [ $? -ne 0 ]; then
    echo "❌ Failed to download dependencies"
    exit 1
fi

echo "🧹 Cleaning up old builds..."
rm -f saas-billing

echo "🏗️ Building the application..."
cd cmd/api && /usr/local/go/bin/go build -o ../../saas-billing
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    cd ../..
    exit 1
fi
cd ../..

if [ ! -f "./saas-billing" ]; then
    echo "❌ Build file not found"
    exit 1
fi

echo "🌟 Starting the server..."
./saas-billing
