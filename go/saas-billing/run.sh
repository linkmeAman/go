#!/bin/bash

echo "🚀 Starting SaaS Billing System..."

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
docker exec -i $(docker-compose ps -q postgres) psql -U postgres -c "CREATE DATABASE saas_billing;" > /dev/null 2>&1 || true
docker exec -i $(docker-compose ps -q postgres) psql -U postgres -d saas_billing -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;" > /dev/null 2>&1
docker exec -i $(docker-compose ps -q postgres) psql -U postgres -d saas_billing -f migrations/001_initial_schema.sql > /dev/null 2>&1

echo "🏗️ Building the application..."
cd cmd/api && go build -o ../../saas-billing && cd ../..

echo "🌟 Starting the server..."
./saas-billing
