# Database Migration Makefile
# Uses golang-migrate for database schema migrations

# Environment selection (dev or test)
# Usage: CREDFOLIO_ENV=test make migrate-up
CREDFOLIO_ENV ?= dev

# Database connection settings (loaded from environment or defaults)
# Default host is the docker container name when running in devcontainer
POSTGRES_HOST ?= credfolio2-postgres
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= credfolio
POSTGRES_PASSWORD ?= credfolio_dev
POSTGRES_SSLMODE ?= disable

# Database name is constructed from base name + environment
POSTGRES_DB = credfolio_$(CREDFOLIO_ENV)

# Migration settings
MIGRATIONS_DIR = migrations
MIGRATE_CMD = migrate

# Build database URL
DATABASE_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

.PHONY: help migration migrate-up migrate-down migrate-down-all migrate-force migrate-version migrate-status

help: ## Show this help
	@echo "Database Migration Commands:"
	@echo ""
	@echo "Environment: CREDFOLIO_ENV=$(CREDFOLIO_ENV) (database: $(POSTGRES_DB))"
	@echo "Set CREDFOLIO_ENV=test to run against test database"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-20s %s\n", $$1, $$2}'

migration: ## Create a new migration (usage: make migration name=create_users)
ifndef name
	$(error name is required. Usage: make migration name=create_users)
endif
	@timestamp=$$(date +%Y%m%d%H%M%S) && \
	touch $(MIGRATIONS_DIR)/$${timestamp}_$(name).up.sql && \
	touch $(MIGRATIONS_DIR)/$${timestamp}_$(name).down.sql && \
	echo "Created migrations:" && \
	echo "  $(MIGRATIONS_DIR)/$${timestamp}_$(name).up.sql" && \
	echo "  $(MIGRATIONS_DIR)/$${timestamp}_$(name).down.sql"

migrate-up: ## Run all pending migrations
	@echo "Running migrations on $(POSTGRES_DB)..."
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down: ## Rollback the last migration
	@echo "Rolling back last migration on $(POSTGRES_DB)..."
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-down-all: ## Rollback all migrations
	@echo "Rolling back all migrations on $(POSTGRES_DB)..."
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down -all

migrate-force: ## Force set migration version (usage: make migrate-force version=1)
ifndef version
	$(error version is required. Usage: make migrate-force version=1)
endif
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

migrate-version: ## Show current migration version
	@echo "Migration version for $(POSTGRES_DB):"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version 2>&1 || echo "  No migrations applied yet"

migrate-status: ## Show migration status for all environments
	@echo "Development database (credfolio_dev):"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/credfolio_dev?sslmode=$(POSTGRES_SSLMODE)" version 2>&1 || echo "  No migrations applied yet"
	@echo ""
	@echo "Test database (credfolio_test):"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/credfolio_test?sslmode=$(POSTGRES_SSLMODE)" version 2>&1 || echo "  No migrations applied yet"
