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

## Quick Start 🚀

### Prerequisites
- Docker Desktop (Download from https://www.docker.com/products/docker-desktop)
- Git (Download from https://git-scm.com/downloads)
- Go (Download from https://golang.org/dl/)

### Super Simple Setup ✨

1. Clone the project:
```bash
git clone https://github.com/linkmeAman/go.git
cd saas-billing
```

2. Copy the environment file:
```bash
cp env.example .env
```

3. Run everything with ONE command:
```bash
./run.sh
```

That's it! The script will:
- 🐳 Start Docker services
- 📦 Set up the database
- 🏗️ Build the application
- 🚀 Start the server

When you see `"Server starting on :8080"`, the application is ready!

### Quick Test 🧪

1. Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "yourpassword123"}'
```

2. Login to get your token:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "yourpassword123"}'
```

3. Create an organization (replace TOKEN with your token from login):
```bash
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "My Company"}'
```

The API will be available at http://localhost:8080

Note: The `go run` command might not work directly from the root directory due to module dependencies. Using the build command is more reliable for production use.

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
