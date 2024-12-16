#!/bin/bash
set -e

# Check if the database exists
if psql -U "$POSTGRES_USER" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'fleetflow'" | grep -q 1; then
  echo "Database fleetflow already exists. Dropping and recreating..."
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    DROP DATABASE IF EXISTS fleetflow;
    CREATE DATABASE fleetflow;
EOSQL
else
  echo "Database fleetflow does not exist. Creating..."
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    CREATE DATABASE fleetflow;
EOSQL
fi

# Run the schema initialization script
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/schema.sql

echo "Database initialization completed successfully!"
