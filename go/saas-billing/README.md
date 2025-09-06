# SaaS Billing System

A production-grade SaaS subscription billing system built with Go.

## Features (MVP)

- User authentication (JWT)
- Organization management
- Subscription plans
- Secure payment processing (Stripe)
- Webhook handling
- Redis caching
- Background job processing

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL
- Redis

### Local Development Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/saas-billing.git
cd saas-billing
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Start the development environment:
```bash
docker-compose up -d
```

4. Run the migrations:
```bash
go run cmd/migrate/main.go up
```

5. Start the server:
```bash
go run cmd/api/main.go
```

The server will start on http://localhost:8080

## Project Structure

```
.
├── cmd/
│   ├── api/        # Main application
│   └── migrate/    # Database migration tool
├── internal/
│   ├── auth/       # Authentication logic
│   ├── billing/    # Billing logic
│   ├── db/         # Database migrations
│   ├── cache/      # Redis caching
│   └── webhooks/   # Webhook handlers
└── pkg/            # Shared packages
```

## API Documentation

Coming soon...

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
