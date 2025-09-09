# Project Dependencies

This document lists all the external dependencies required for the SaaS Billing System.

## Core Dependencies

```bash
# Web Framework
go get -u github.com/gin-gonic/gin@v1.8.1

# Environment Variables
go get -u github.com/joho/godotenv@v1.4.0

# PostgreSQL Driver
go get -u github.com/lib/pq@v1.10.7

# JWT Authentication
go get -u github.com/golang-jwt/jwt/v4@v4.5.2

# Redis Cache (Optional)
go get -u github.com/go-redis/redis/v8@v8.11.5

# Password Hashing
go get -u golang.org/x/crypto@v0.0.0-20211215153901-e495a2d5b3d3
```

## Installing Dependencies

You can install all dependencies by running:

```bash
go mod download
```

Or simply:

```bash
go mod tidy
```

## Version Requirements

- Go: 1.18 or later
- PostgreSQL: 12 or later
- Redis: 6 or later (optional)

## Development Dependencies

For development, these additional tools are recommended:

```bash
# Hot Reloading
go install github.com/cosmtrek/air@latest

# Database Migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Testing
go get -u github.com/stretchr/testify
```

## Managing Dependencies

1. Adding a new dependency:
   ```bash
   go get github.com/example/package@latest
   ```

2. Upgrading dependencies:
   ```bash
   go get -u ./...
   ```

3. Cleaning up unused dependencies:
   ```bash
   go mod tidy
   ```

4. Verifying dependencies:
   ```bash
   go mod verify
   ```

## Dependency Tree

Here's a breakdown of why each dependency is needed:

- **gin-gonic/gin**: Web framework for building the REST API
  - Fast and lightweight
  - Built-in middleware support
  - Good for building RESTful APIs

- **joho/godotenv**: Loading environment variables
  - Loads .env files
  - Development configuration management

- **lib/pq**: PostgreSQL driver
  - Database connectivity
  - Required for PostgreSQL operations

- **golang-jwt/jwt**: JWT token handling
  - User authentication
  - Secure token generation and validation

- **go-redis/redis**: Caching layer (optional)
  - Performance optimization
  - Session management
  - Rate limiting

- **golang.org/x/crypto**: Security utilities
  - Password hashing
  - Secure cryptographic operations

## Manual Installation

If you need to manually install the project on a new machine:

1. Install Go 1.18 or later:
   ```bash
   wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/linkmeAman/saas-billing.git
   cd saas-billing
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Set up environment:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. Run the setup script:
   ```bash
   ./scripts/setup.sh
   ```

## Troubleshooting Dependencies

If you encounter issues:

1. Clean Go's cache:
   ```bash
   go clean -modcache
   ```

2. Reset the module:
   ```bash
   rm go.sum
   go mod tidy
   ```

3. Verify dependencies:
   ```bash
   go mod verify
   ```

4. Check for updates:
   ```bash
   go list -u -m all
   ```
