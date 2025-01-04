# FleetFlow Technical Documentation

## Architecture Overview

FleetFlow is built using a modern, scalable architecture:

### Backend (Go)
- **Framework**: Native Go HTTP server
- **Database**: PostgreSQL
- **Authentication**: JWT-based authentication
- **API Style**: RESTful API

### Directory Structure
```
backend/
├── cmd/           # Command line tools
├── db/            # Database migrations and schemas
├── handlers/      # HTTP request handlers
├── middleware/    # HTTP middleware components
├── models/        # Data models and business logic
├── pkg/           # Shared packages and utilities
└── main.go        # Application entry point
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/refresh`

### Fleet Management
- `GET /api/v1/vehicles`
- `POST /api/v1/vehicles`
- `PUT /api/v1/vehicles/{id}`
- `DELETE /api/v1/vehicles/{id}`

### Maintenance
- `GET /api/v1/maintenance`
- `POST /api/v1/maintenance`
- `PUT /api/v1/maintenance/{id}`

## Security

### Authentication
- JWT-based authentication
- Token refresh mechanism
- Role-based access control (RBAC)

### Data Protection
- TLS encryption for all API endpoints
- Password hashing using bcrypt
- Input validation and sanitization

## Database Schema
The schema is defined in `schema.sql`

## Error Handling

### Error Codes
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

### Error Response Format
```json
{
    "error": {
        "code": "ERROR_CODE",
        "message": "Human readable message",
        "details": {}
    }
}
```

## Rate Limiting
- Rate limiting implemented using token bucket algorithm
- Default: 100 requests per minute per IP
- Configurable through environment variables

## Monitoring and Logging
- Structured logging using zerolog, consolelog

## Performance Considerations
- Connection pooling for database
- Caching layer for frequently accessed data
- Pagination for large data sets

## Development Setup

### Prerequisites
- Go 1.19 or higher
- PostgreSQL 13 or higher
- Make

### Environment Variables
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fleetflow
DB_USER=postgres
DB_PASSWORD=your_password
JWT_SECRET=your_secret
RATE_LIMIT=100
```

### Building and Running
```bash
docker compose build
docker compose up
```

## Routes
Authenticated routes:
- `POST /api/v1/vehicles` - Create a new vehicle
- `PUT /api/v1/vehicles/{id}` - Update an existing vehicle
- `DELETE /api/v1/vehicles/{id}` - Delete a vehicle
- `GET /api/v1/vehicles` - List all vehicles

- `GET /api/v1/maintenance` - List all maintenance records
- `POST /api/v1/maintenance` - Create a new maintenance record
- `PUT /api/v1/maintenance/{id}` - Update an existing maintenance record

Unauthenticated routes:
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/register`

- `GET /api/v1/vehicles` - list all vehicles
- `GET /api/v1/vehicles/{id}` - Get a specific vehicle
