#!/bin/bash

echo "🚀 Starting SaaS Billing System..."

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
go mod download
if [ $? -ne 0 ]; then
    echo "❌ Failed to download dependencies"
    exit 1
fi

echo "🧹 Cleaning up old builds..."
rm -f saas-billing

echo "🏗️ Building the application..."
cd cmd/api && go build -o ../../saas-billing
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
