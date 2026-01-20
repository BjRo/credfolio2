#!/bin/bash
set -e

# Create dev and test databases if they don't exist
# This script runs during postgres container initialization

# Create dev database (credfolio_dev)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    SELECT 'CREATE DATABASE credfolio_dev'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'credfolio_dev')\gexec
EOSQL
echo "Database 'credfolio_dev' created or already exists"

# Create test database (credfolio_test)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    SELECT 'CREATE DATABASE credfolio_test'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'credfolio_test')\gexec
EOSQL
echo "Database 'credfolio_test' created or already exists"
