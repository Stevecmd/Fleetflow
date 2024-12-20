Commands to use:
Register a new user:
```sh
curl -X POST http://localhost:8000/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{
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
    "profile_image_url": "https://example.com/profile.jpg"
}'

```
Login:
```sh
curl -X POST http://localhost:8000/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{
    "username": "testuser",
    "password": "Test@123"
}'

```

Refresh token:
```sh
curl -X POST http://localhost:8000/api/v1/auth/refresh \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM1MTM4NTcsImlkIjo2LCJ1c2VybmFtZSI6InRlc3R1c2VyIn0.gIGOcZ1wsG9z_ui9h7XI34m7_wvUXXmaOJl6vcFruBI" \
-d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM1MTM4NTcsImlkIjo2LCJ1c2VybmFtZSI6InRlc3R1c2VyIn0.gIGOcZ1wsG9z_ui9h7XI34m7_wvUXXmaOJl6vcFruBI"
}'

```
Logout:
```sh
curl -X POST http://localhost:8000/api/v1/auth/logout \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <Access_Token>"

```
Get user profile:
```sh
curl -X GET http://localhost:8000/api/v1/users/profile \
-H "Authorization: Bearer <Access_Token>"

```

Update user profile:
```sh
curl -X PUT http://localhost:8000/api/v1/users/profile \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <Access_Token>" \
-d '{
    "first_name": "Updated",
    "last_name": "User",
    "phone": "0987654321"
}'

```

List drivers:
```sh
curl -X GET http://localhost:8000/api/v1/drivers \
-H "Authorization: Bearer <Access_Token>"

```

Get Driver:
```sh
curl -X GET http://localhost:8000/api/v1/drivers/1 \
-H "Authorization: Bearer <Access_Token>"

```

Update driver profile:
```sh
curl -X PUT http://localhost:8000/api/v1/drivers/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <Access_Token>" \
-d '{
    "license_number": "DL654321",
    "license_type": "CDL-B",
    "expiry_date": "2026-12-31"
}'

```

# Vehicle endpoints
List vehicles:
```sh
curl -X GET http://localhost:8000/api/v1/vehicles \
-H "Authorization: Bearer <Access_Token>"

```

Get Vehicle: **
```sh
curl -X GET http://localhost:8000/api/v1/vehicles/1 \
-H "Authorization: Bearer <Access_Token>"

```

Create Vehicle:
```sh
curl -X POST http://localhost:8000/api/v1/vehicles \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <Access_Token>" \
-d '{
    "plate_number": "KDD123D",
    "type": "Truck",
    "make": "Ford",
    "model": "F-150",
    "year": 2021,
    "capacity": 5000,
    "fuel_type": "Diesel",
    "status_id": 1
}'

```
Update vehicle:
```sh
curl -X PUT http://localhost:8000/api/v1/vehicles/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <Access_Token>" \
-d '{
    "make": "Chevrolet",
    "model": "Silverado",
    "year": 2022
}'

```

Delete vehicle:
```sh
curl -X DELETE http://localhost:8000/api/v1/vehicles/1 \
-H "Authorization: Bearer <Access_Token>"

```

Interact with my docker container:
```sh
/backend$ docker exec -it postgres_db psql -U postgres_user -d fleetflow bash
```
Load schema into docker container:
```sh
/backend$ docker cp ./db/init/schema.sql postgres_db:/schema.sql
```

Run the schema:
```sh
/backend$ docker exec -it postgres_db psql -U postgres_user -d fleetflow -f /schema.sql
```
Drop database:
```sh
/backend$ docker exec -it postgres_db psql -U postgres_user -d postgres -c "DROP DATABASE IF EXISTS fleetflow;"
```
