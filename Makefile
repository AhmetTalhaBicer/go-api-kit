APP_NAME   := go-api-kit
BUILD_DIR  := ./build
CMD_DIR    := ./cmd/api
BINARY     := $(BUILD_DIR)/$(APP_NAME)

GOCMD   := go
GOBUILD := $(GOCMD) build
GOTEST  := $(GOCMD) test
GOVET   := $(GOCMD) vet
GOMOD   := $(GOCMD) mod

.PHONY: all build run test test-integration coverage vet lint tidy clean \
        sqlc \
        migrate-up migrate-down migrate-status migrate-create \
        docker-build docker-up docker-down help

## all: Download dependencies and build the binary
all: tidy build

## build: Compile the application binary
build:
	@echo ">> Building binary..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="-s -w" -o $(BINARY) $(CMD_DIR)/...

## run: Run the application directly with go run
run:
	$(GOCMD) run $(CMD_DIR)/main.go

## watch: Run live-reloading dev server with Air (go install github.com/air-verse/air@latest)
watch:
	@which air > /dev/null 2>&1 || (echo "air not found: go install github.com/air-verse/air@latest" && exit 1)
	air

## test: Run all tests with race detection and coverage
test:
	@echo ">> Running tests..."
	$(GOTEST) -v -race -cover ./...

## test-integration: Run integration tests only
test-integration:
	$(GOTEST) -v ./test/...

## coverage: Generate an HTML coverage report
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo ">> coverage.html generated"

## vet: Run go vet
vet:
	$(GOVET) ./...

## lint: Run golangci-lint (must be installed separately)
lint:
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not found: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

## tidy: Run go mod tidy
tidy:
	$(GOMOD) tidy

## clean: Remove build artifacts and coverage files
clean:
	@rm -rf $(BUILD_DIR) coverage.out coverage.html

# ── sqlc ─────────────────────────────────────────────────────────────────────

## sqlc: Generate type-safe Go code from SQL queries in sqlc/queries/
##       Install: https://docs.sqlc.dev/en/latest/overview/install.html
sqlc:
	@which sqlc > /dev/null 2>&1 || (echo "sqlc not found: https://docs.sqlc.dev/en/latest/overview/install.html" && exit 1)
	sqlc generate
	@echo ">> sqlc/gen/ regenerated"

# ── Database migrations (goose) ───────────────────────────────────────────────

## migrate-up: Apply all pending database migrations
##             Install goose: go install github.com/pressly/goose/v3/cmd/goose@latest
migrate-up:
	@bash scripts/migrate.sh up

## migrate-down: Roll back the last database migration
migrate-down:
	@bash scripts/migrate.sh down

## migrate-status: Show current migration status
migrate-status:
	@bash scripts/migrate.sh status

## migrate-create: Create a new SQL migration file. Usage: make migrate-create NAME=add_users
migrate-create:
	@bash scripts/migrate.sh create $(NAME)

# ── Docker ───────────────────────────────────────────────────────────────────

## docker-build: Build the Docker image
docker-build:
	docker build -t $(APP_NAME):latest -f deployments/Dockerfile .

## docker-up: Start services with Docker Compose
docker-up:
	docker compose -f deployments/docker-compose.yml up -d

## docker-down: Stop Docker Compose services
docker-down:
	docker compose -f deployments/docker-compose.yml down

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
