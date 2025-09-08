.PHONY: all build test clean run lint docker-build docker-run migrate-up migrate-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=saas-billing
MAIN_PATH=cmd/api/main.go

all: test build

build:
	cd cmd/api && $(GOBUILD) -o ../../$(BINARY_NAME)

test:
	$(GOTEST) -v ./... -cover

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	cd cmd/api && $(GOCMD) run .

lint:
	golangci-lint run

docker-build:
	docker build -t saas-billing .

docker-run:
	docker-compose up -d

migrate-up:
	migrate -path internal/db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path internal/db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

setup: docker-run migrate-up

# Security checks
.PHONY: security-check
security-check:
	gosec ./...
	nancy go.sum
