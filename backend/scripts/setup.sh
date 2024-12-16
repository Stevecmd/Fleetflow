#!/bin/bash

# Setup script for FleetFlow backend development environment

echo "Setting up FleetFlow development environment..."

# Check for required tools
command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }
command -v psql >/dev/null 2>&1 || { echo "PostgreSQL client is required but not installed. Aborting." >&2; exit 1; }

# Wait for database to be ready
sleep 5

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod verify

echo "Setup complete! You can now run 'go run main.go' to start the server."
