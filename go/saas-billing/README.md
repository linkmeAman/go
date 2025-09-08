# SaaS Billing System

A production-grade SaaS subscription billing system built with Go, featuring secure authentication, organization management, and subscription billing.

## Features

### Core Features
- 🔐 JWT-based Authentication
- 👥 Organization Management
- 💰 Subscription Plans & Billing
- 🔑 Role-Based Access Control (RBAC)
- 📊 Usage Tracking
- 💾 PostgreSQL Database
- ⚡ Redis Caching

### Technical Features
- ✅ Comprehensive Test Coverage
- 📝 OpenAPI/Swagger Documentation
- 🔍 Request/Response Logging
- 🚦 Rate Limiting
- 🔄 Background Jobs
- 📈 Prometheus Metrics
- 📊 Grafana Dashboards

## Quick Start

The easiest way to get started is using our setup script:

```bash
# Clone the repository
git clone https://github.com/yourusername/saas-billing.git
cd saas-billing

# Run the setup script
chmod +x scripts/setup.sh
./scripts/setup.sh
```

This will:
- Check prerequisites
- Set up environment variables
- Start required services
- Run database migrations
- Create initial admin user

## Manual Setup

### Prerequisites

- Go 1.18+
- Docker and Docker Compose
- PostgreSQL 12+
- Redis 6+ (optional)

### Development Setup

1. Clone and setup:
```bash
# Clone repository
git clone https://github.com/linkmeAman/go.git
cd saas-billing

# Copy environment file
cp env.example .env

# Install dependencies
go mod download
```

2. Start services:
```bash
# Start PostgreSQL and Redis
docker-compose up -d

# Verify services are running
docker-compose ps
```

3. Run migrations and start server:
```bash
# Create database
psql -U postgres -c "CREATE DATABASE saas_billing;"

# Run database migrations
psql -U postgres -d saas_billing -f migrations/001_initial_schema.sql

# Start the server
go run cmd/api/main.go
```

The API will be available at http://localhost:8080

## Project Structure

```
saas-billing/
├── api/            # API specifications and Swagger docs
├── cmd/
│   └── api/        # Main application entry point
├── docs/           # Documentation files
├── internal/       # Internal packages
│   ├── auth/       # Authentication & authorization
│   ├── billing/    # Billing & subscription logic
│   ├── cache/      # Redis caching implementation
│   ├── db/         # Database operations & migrations
│   ├── logger/     # Structured logging
│   ├── middleware/ # HTTP middleware
│   ├── orgs/       # Organization management
│   ├── types/      # Shared types and interfaces
│   └── users/      # User management
├── pkg/            # Public packages
├── scripts/        # Setup and maintenance scripts
└── tests/          # Integration tests
```

## API Documentation

Full API documentation is available in multiple formats:

1. OpenAPI/Swagger Documentation:
   - View `api/swagger.yaml`
   - Or visit `/api/docs` when server is running

2. Markdown Documentation:
   - See `docs/API.md` for detailed endpoint documentation

3. Postman Collection:
   - Import `docs/postman_collection.json`

## Development

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test ./internal/auth -v

# Run with coverage
make test-coverage
```

### Development Tools

```bash
# Install development tools
make install-tools

# Run linter
make lint

# Format code
make fmt

# Check for security issues
make security-check
```

### Docker Commands

```bash
# Start development environment
docker-compose up -d

# Start production environment
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Monitoring

### Prometheus Metrics

Available at `/metrics` endpoint, including:
- Request latencies
- Error rates
- Resource usage
- Custom business metrics

### Grafana Dashboards

Access Grafana at `http://localhost:3000` with:
- Default username: `admin`
- Default password: `admin`

Pre-configured dashboards available for:
- API metrics
- System metrics
- Business metrics

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Create a Pull Request

### Commit Guidelines

We follow conventional commits:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Adding or modifying tests
- `chore:` Maintenance tasks

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- Documentation: See `docs/` directory
- Issues: Please use GitHub issues
- Security: Report security issues to security@example.com
