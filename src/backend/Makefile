# Database Migration Makefile
# Uses golang-migrate for database schema migrations

# Database connection settings (loaded from environment or defaults)
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= credfolio
POSTGRES_PASSWORD ?= credfolio_dev
POSTGRES_DB ?= credfolio
POSTGRES_DB_TEST ?= credfolio_test
POSTGRES_SSLMODE ?= disable

# Migration settings
MIGRATIONS_DIR = migrations
MIGRATE_CMD = migrate

# Build database URLs
DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)
DATABASE_URL_TEST ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB_TEST)?sslmode=$(POSTGRES_SSLMODE)

.PHONY: help migration migrate-up migrate-down migrate-up-test migrate-down-test migrate-force migrate-version migrate-status

help: ## Show this help
	@echo "Database Migration Commands:"
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

migrate-up: ## Run all pending migrations on development database
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down: ## Rollback the last migration on development database
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-down-all: ## Rollback all migrations on development database
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down -all

migrate-up-test: ## Run all pending migrations on test database
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL_TEST)" up

migrate-down-test: ## Rollback the last migration on test database
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL_TEST)" down 1

migrate-down-all-test: ## Rollback all migrations on test database
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL_TEST)" down -all

migrate-force: ## Force set migration version (usage: make migrate-force version=1)
ifndef version
	$(error version is required. Usage: make migrate-force version=1)
endif
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

migrate-version: ## Show current migration version
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-status: ## Show migration status (same as version but clearer name)
	@echo "Development database:"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version 2>&1 || echo "  No migrations applied yet"
	@echo ""
	@echo "Test database:"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL_TEST)" version 2>&1 || echo "  No migrations applied yet"
