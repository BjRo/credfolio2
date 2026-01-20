#!/bin/bash
set -e

# Create test database if it doesn't exist
# This script runs during postgres container initialization

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    SELECT 'CREATE DATABASE ${POSTGRES_DB_TEST}'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${POSTGRES_DB_TEST}')\gexec
EOSQL

echo "Test database '${POSTGRES_DB_TEST}' created or already exists"
