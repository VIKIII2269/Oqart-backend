.PHONY: help run build test migrate-up migrate-down migrate-create seed docker-up docker-down clean lint

# Variables
APP_NAME=oqart-api
MAIN_PATH=./cmd/api
BUILD_DIR=./bin
MIGRATIONS_PATH=./migrations
DATABASE_URL?=postgresql://oqart:oqart@localhost:5432/oqart?sslmode=disable

## help: Display this help message
help:
	@echo "Available commands:"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application binary"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback database migrations"
	@echo "  make migrate-create NAME=<name> - Create new migration"
	@echo "  make seed          - Seed the database with test data"
	@echo "  make docker-up     - Start services with docker-compose"
	@echo "  make docker-down   - Stop docker-compose services"
	@echo "  make docker-build  - Build docker image"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make lint          - Run linters"
	@echo "  make fmt           - Format code"
	@echo "  make deps          - Install dependencies"
	@echo "  make tidy          - Tidy go modules"

## run: Run the application
run:
	@echo "Starting application..."
	go run $(MAIN_PATH)/main.go

## build: Build the application binary
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)/main.go
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## migrate-up: Run database migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up
	@echo "Migrations complete"

## migrate-down: Rollback database migrations
migrate-down:
	@echo "Rolling back migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1
	@echo "Rollback complete"

## migrate-down-all: Rollback all migrations
migrate-down-all:
	@echo "Rolling back all migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down
	@echo "All migrations rolled back"

## migrate-create: Create a new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
	@echo "Migration created"

## migrate-force: Force migration version (use with caution)
migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@echo "Forcing migration version to $(VERSION)..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(VERSION)

## migrate-version: Display current migration version
migrate-version:
	@echo "Current migration version:"
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

## seed: Seed the database
seed:
	@echo "Seeding database..."
	go run $(MAIN_PATH)/main.go seed
	@echo "Database seeded"

## docker-up: Start docker-compose services
docker-up:
	@echo "Starting docker services..."
	docker-compose up -d
	@echo "Services started"

## docker-down: Stop docker-compose services
docker-down:
	@echo "Stopping docker services..."
	docker-compose down
	@echo "Services stopped"

## docker-build: Build docker image
docker-build:
	@echo "Building docker image..."
	docker build -t $(APP_NAME):latest .
	@echo "Docker image built"

## docker-logs: View docker logs
docker-logs:
	docker-compose logs -f

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

## lint: Run linters
lint:
	@echo "Running linters..."
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

## deps: Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	@echo "Dependencies installed"

## tidy: Tidy go modules
tidy:
	@echo "Tidying modules..."
	go mod tidy
	@echo "Modules tidied"

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installed"

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g $(MAIN_PATH)/main.go -o ./docs/swagger
	@echo "Swagger documentation generated"

## dev: Run in development mode with hot reload
dev:
	@echo "Starting in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to standard run..."; \
		make run; \
	fi

## db-reset: Reset database (drop, create, migrate)
db-reset:
	@echo "Resetting database..."
	make migrate-down-all
	make migrate-up
	make seed
	@echo "Database reset complete"

.DEFAULT_GOAL := help
