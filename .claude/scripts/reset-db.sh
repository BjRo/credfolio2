#!/bin/bash
# Reset database: drop, recreate, and run all migrations
#
# Usage: .claude/scripts/reset-db.sh [environment]
# Arguments:
#   environment  Optional. One of: dev (default), test, all
#
# Examples:
#   .claude/scripts/reset-db.sh          # Reset dev database (with confirmation)
#   .claude/scripts/reset-db.sh test     # Reset test database (no confirmation)
#   .claude/scripts/reset-db.sh all      # Reset both databases (with confirmation)

set -e

# Colors for output (matching start-work.sh pattern)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Connection constants
POSTGRES_HOST="credfolio2-postgres"
POSTGRES_PORT="5432"
POSTGRES_USER="credfolio"
POSTGRES_PASSWORD="credfolio_dev"
MIGRATIONS_DIR="/workspace/src/backend/migrations"

# Usage message
usage() {
    echo "Usage: $0 [environment]"
    echo ""
    echo "Arguments:"
    echo "  environment  Optional. One of: dev (default), test, all"
    echo ""
    echo "Examples:"
    echo "  $0          # Reset dev database (with confirmation)"
    echo "  $0 test     # Reset test database (no confirmation)"
    echo "  $0 all      # Reset both databases (with confirmation)"
    exit 1
}

# Confirmation prompt
confirm_reset() {
    local db_description="$1"
    echo -e "${YELLOW}WARNING: Reset ${db_description}? This will DELETE ALL DATA.${NC}"
    read -p "Continue? [y/N]: " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        return 1
    fi
    return 0
}

# Reset a single database
reset_database() {
    local db_name="$1"

    echo -e "${GREEN}Resetting database: ${db_name}${NC}"

    # Set password to avoid prompts
    export PGPASSWORD="$POSTGRES_PASSWORD"

    # Step 1: Drop database
    echo "  [1/4] Dropping database..."
    psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d postgres -c "DROP DATABASE IF EXISTS $db_name" > /dev/null

    # Step 2: Create database
    echo "  [2/4] Creating database..."
    psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE $db_name" > /dev/null

    # Step 3: Run migrations
    echo "  [3/4] Running migrations..."
    local db_url="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$db_name?sslmode=disable"
    migrate -path "$MIGRATIONS_DIR" -database "$db_url" up > /dev/null

    # Step 4: Show migration status
    echo "  [4/4] Migration status:"
    local version=$(migrate -path "$MIGRATIONS_DIR" -database "$db_url" version 2>&1 | tail -n 1)
    echo "         $version"

    echo -e "${GREEN}Successfully reset ${db_name}${NC}\n"
}

# Parse environment argument (default to dev)
ENV="${1:-dev}"

# Validate environment
case "$ENV" in
    dev|test|all)
        # Valid
        ;;
    *)
        echo -e "${RED}Error: Invalid environment '${ENV}'${NC}"
        echo ""
        usage
        ;;
esac

# Execute based on environment
case "$ENV" in
    dev)
        if confirm_reset "credfolio_dev"; then
            reset_database "credfolio_dev"
        else
            echo -e "${YELLOW}Cancelled.${NC}"
            exit 0
        fi
        ;;
    test)
        reset_database "credfolio_test"
        ;;
    all)
        if confirm_reset "BOTH dev and test databases"; then
            reset_database "credfolio_dev"
            reset_database "credfolio_test"
        else
            echo -e "${YELLOW}Cancelled.${NC}"
            exit 0
        fi
        ;;
esac

echo -e "${GREEN}Done!${NC}"
