#!/bin/bash
set -e

# Start PostgreSQL server in the background
docker-entrypoint.sh postgres &

# Wait for PostgreSQL to be ready
until pg_isready -h localhost -p 5432 -U "$POSTGRES_USER"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Run the custom initialization script
/docker-entrypoint-initdb.d/always-init.sh

# Bring PostgreSQL server to the foreground
wait
