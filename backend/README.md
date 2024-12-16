# FleetFlow Backend API Documentation

## Overview
FleetFlow is a comprehensive fleet management system built with Go. This document details the available API endpoints, their expected inputs, and responses.

### Backend contains:
1. Authentication System:
    - JWT-based authentication with access and refresh tokens
    - Access tokens expire in 15 minutes
    - Refresh tokens expire in 7 days
    - Secure password validation with multiple requirements
    - Token blacklisting for logout

2. API Endpoints:
Auth: Login, Register, Logout, Refresh Token
    - Users: Profile management
    - Vehicles: CRUD operations
    - Drivers: Management and statistics
    - Deliveries: Tracking and management
    - Warehouses: Inventory management
    - Maintenance: Vehicle maintenance records

3. Security Features:
    - CORS configured for localhost:3000
    - Rate limiting middleware
    - Secure password hashing with bcrypt
    - JWT-based authorization middleware

## Authentication
All routes except `/auth/login` and `/auth/register` require a valid JWT token in the Authorization header:
```
Authorization: Bearer your_access_token_here
```

## API Endpoints

### Authentication Routes

#### 1. Register User
- **Endpoint**: `POST /auth/register`
- **Body**:
```json
{
    "username": "testuser",
    "password": "Test@123",
    "email": "test@example.com",
    "role_id": 5,
    "first_name": "Test",
    "last_name": "User",
    "phone": "1234567890",
    "date_of_birth": "1990-01-01",
    "gender": "male",
    "nationality": "US",
    "preferred_language": "en",
    "profile_image_url": "https://example.com/profile.jpg",
    "addresses": [
        {
            "street1": "123 Main St",
            "street2": "Apt 4B",
            "city": "New York",
            "state": "NY",
            "zip": "10001",
            "country": "USA",
            "address_type": "primary",
            "is_default": true,
            "latitude": 40.7128,
            "longitude": -74.0060
        }
    ],
    "emergency_contacts": [
        {
            "name": "John Doe",
            "relationship": "Father",
            "phone": "1234567890",
            "email": "john.doe@example.com"
        }
    ]
}
```
- **Success Response** (201 Created):
```json
{
    "message": "Registration successful",
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "user": {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "role_id": 5,
        "role_name": "customer",
        "first_name": "Test",
        "last_name": "User",
        "phone": "1234567890",
        "profile_image_url": "https://example.com/profile.jpg",
        "date_of_birth": "1990-01-01",
        "gender": "male",
        "nationality": "US",
        "preferred_language": "en",
        "addresses": [...],
        "emergency_contacts": [...],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 2. Login
- **Endpoint**: `POST /auth/login`
```sh
curl -X POST http://localhost:8000/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"testuser","password":"Test@123"}
```

```sh
curl -X POST http://localhost:8000/api/v1/auth/register -H Content-Type: application/json -d {"username":"testuser","password":"Test@123","email":"testuser@example.com","first_name":"Test","last_name":"User","phone":"+1234567890","role_id":1}

```

```sh
curl -X POST http://localhost:8000/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"testuser2","password":"Test123!"}'
```
- **Body**:
```json
{
    "username": "testuser",
    "password": "Test@123"
}
```
- **Success Response** (200 OK):
```json
{
    "message": "Successfully logged in",
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG..."
}
```

#### 3. Refresh Token
- **Endpoint**: `POST /auth/refresh`
- **Body**:
```json
{
    "refresh_token": "eyJhbG..."
}
```
- **Success Response** (200 OK):
```json
{
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG..."
}
```

#### 4. Logout
- **Endpoint**: `POST /auth/logout`
- **Headers**: Requires Authorization
- **Success Response** (200 OK):
```json
{
    "message": "Successfully logged out"
}
```

### User Routes

#### 1. Get User Profile
- **Endpoint**: `GET /users/{user_id}`
- **Headers**: Requires Authorization
- **Success Response** (200 OK):
```json
{
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "first_name": "Test",
    "last_name": "User",
    "phone": "1234567890",
    "role_id": 5,
    "role_name": "customer",
    "profile_image_url": "https://example.com/profile.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```
User profile routes with the access token:
```sh
curl -X GET http://localhost:8000/api/v1/users/1 -H "Authorization: Bearer <access_token>"
```

#### 2. Update User Profile
- **Endpoint**: `PUT /users/{user_id}`
- **Headers**: Requires Authorization
- **Body**:
```json
{
    "first_name": "Updated",
    "last_name": "Name",
    "phone": "9876543210",
    "profile_image_url": "https://example.com/profile.jpg"
}
```
- **Success Response** (200 OK):
```json
{
    "message": "Profile updated successfully"
}
```

### Driver Routes

#### 1. List Drivers
- **Endpoint**: `GET /api/drivers`
- **Headers**: Requires Authorization
- **Success Response** (200 OK):
```json
{
    "drivers": [
        {
            "id": 1,
            "user_id": 2,
            "license_number": "DL123456",
            "license_type": "CDL-A",
            "expiry_date": "2025-12-31",
            "vehicle_type": "Truck",
            "years_experience": 5,
            "certifications": ["Hazmat", "Tanker"],
            "medical_cert_expiry": "2024-12-31"
        }
    ]
}
```

#### 2. Create Driver Profile
- **Endpoint**: `POST /api/drivers`
- **Headers**: Requires Authorization
- **Body**:
```json
{
    "license_number": "DL999999",
    "license_type": "CDL-A",
    "expiry_date": "2025-12-31",
    "vehicle_type": "Truck",
    "years_experience": 5,
    "certifications": ["Hazmat", "Tanker"],
    "medical_cert_expiry": "2024-12-31"
}
```
- **Success Response** (201 Created):
```json
{
    "message": "Driver profile created successfully",
    "driver_id": 1
}
```

#### 3. Get Driver Profile
- **Endpoint**: `GET /api/drivers/{driver_id}`
- **Headers**: Requires Authorization
- **Success Response** (200 OK):
```json
{
    "id": 1,
    "user_id": 2,
    "license_number": "DL123456",
    "license_type": "CDL-A",
    "expiry_date": "2025-12-31",
    "vehicle_type": "Truck",
    "years_experience": 5,
    "certifications": ["Hazmat", "Tanker"],
    "medical_cert_expiry": "2024-12-31"
}
```

#### 4. Update Driver Profile
- **Endpoint**: `PUT /api/drivers/{driver_id}`
- **Headers**: Requires Authorization
- **Body**:
```json
{
    "license_type": "CDL-B",
    "years_experience": 6,
    "certifications": ["Hazmat", "Passenger"]
}
```
- **Success Response** (200 OK):
```json
{
    "message": "Driver profile updated successfully"
}
```

#### 5. Delete Driver Profile
- **Endpoint**: `DELETE /api/drivers/{driver_id}`
- **Headers**: Requires Authorization
- **Success Response** (200 OK):
```json
{
    "message": "Driver profile deleted successfully"
}
```

### Vehicle Routes

#### 1. List Vehicles
- **Endpoint**: `GET /api/v1/vehicles`
- **Headers**: Requires Authorization
- **Success Response** (200 OK):
```json
{
    "vehicles": [
        {
            "id": 1,
            "plate_number": "ABC123",
            "type": "Truck",
            "make": "Volvo",
            "model": "VNL",
            "year": 2022,
            "capacity": 40000,
            "status": "available"
        }
    ]
}
```

#### 2. Create Vehicle
- **Endpoint**: `POST /api/v1/vehicles`
- **Headers**: Requires Authorization
- **Body**:
```json
{
    "plate_number": "ABC123",
    "type": "Truck",
    "make": "Volvo",
    "model": "VNL",
    "year": 2022,
    "capacity": 40000,
    "status": "available"
}
```
- **Success Response** (201 Created):
```json
{
    "message": "Vehicle created successfully",
    "vehicle_id": 1
}
```

#### 2. Maintenance routes
- **Endpoint**: `POST /api/v1/maintenance`
- **Headers**: Requires Authorization
- **Body**:

Create a new Maintenance record using `PUT` and update a record using `POST`.
```json
curl -X PUT http://localhost:8000/api/v1/maintenance/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-d '{
    "vehicle_id": 1,
    "type": "Oil Change",
    "description": "Updated description - Oil and filter change",
    "service_date": "2024-01-15T10:00:00Z",
    "cost": 80.00,
    "odometer_reading": 51000,
    "performed_by": "John Smith",
    "next_service_date": "2024-04-15T10:00:00Z",
    "notes": "Used synthetic oil - Updated"
}'
```
- **Success Response** (201 Created):
```json
{"id":1,"vehicle_id":1,"type":"Oil Change","description":"Regular maintenance - Oil and filter change","service_date":"2023-12-15T10:00:00Z","cost":75,"odometer_reading":50000,"performed_by":"John Smith","next_service_date":"2024-03-15T10:00:00Z","notes":"Used synthetic oil","created_at":"2024-11-29T22:01:08.835279Z","updated_at":"2024-11-29T22:01:08.835279Z"}
```
### Get Details of a Specific Maintenance Record
```sh
curl -X GET http://localhost:8000/api/v1/maintenance/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN"

```

### Update an Existing Maintenance Record
```sh
curl -X PUT http://localhost:8000/api/v1/maintenance/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-d '{
    "vehicle_id": 1,
    "type": "Oil Change",
    "description": "Updated description - Oil and filter change",
    "service_date": "2024-01-15T10:00:00Z",
    "cost": 80.00,
    "odometer_reading": 51000,
    "performed_by": "John Smith",
    "next_service_date": "2024-04-15T10:00:00Z",
    "notes": "Used synthetic oil - Updated"
}'

```

### Delete a maintenance record
```sh
curl -X DELETE http://localhost:8000/api/v1/maintenance/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN"

```

## Reinitialize Database
Drop the database:
```sh
docker exec -it postgres_db psql -U postgres_user -d postgres -c "DROP DATABASE IF EXISTS fleetflow;"
```
Create the database:
```sh
docker exec -it postgres_db psql -U postgres_user -d postgres -c "CREATE DATABASE fleetflow;"
```
To copy the schema into the docker container use:
```bash
Fleetflow/backend$ docker cp ./db/init/schema.sql postgres_db:/schema.sql
```
You can use the following command to reinitialize the database:
```sh
docker exec -i postgres_db psql -d postgres -f /schema.sql
```
```sh
docker exec -it postgres_db psql -U postgres_user -d fleetflow -f /schema.sql
```
or 

```bash
psql -U postgres -d fleetflow -f backend/db/init/schema.sql
```
This command will drop and create the tables, and insert sample data into the tables.



To initialize the database use:
```bash
Fleetflow/backend$ docker exec -it postgres_db psql -U postgres_user -d fleetflow -f /schema.sql
```

List all tables in \database:
```sh
docker exec postgres_db psql -U postgres_user -d fleetflow -c "\dt"
```

Interactively check on users table:
```bash
docker exec postgres_db psql -U postgres_user -d fleetflow -c "SELECT * FROM users;
```

Interactively check on deliveries table:
```sh
docker exec postgres_db psql -U postgres_user -d fleetflow -c "\d deliveries;"
```

Interactively check on vehicles table:
```sh
docker exec postgres_db psql -U postgres_user -d fleetflow -c "SELECT * FROM vehicles;"
```

## Error Responses
All endpoints may return the following error responses:

- **400 Bad Request**: Invalid input data
- **401 Unauthorized**: Missing or invalid authentication token
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists
- **500 Internal Server Error**: Server-side error

## Role IDs
- 1: admin
- 2: fleet_manager
- 3: driver
- 4: loader
- 5: customer

## Testing
You can use curl commands to test these endpoints. Remember to:
1. Start with registration
2. Login to get access token
3. Use the access token in the Authorization header for subsequent requests
4. Use the refresh token endpoint if the access token expires
