# To-do

## Maintenance Routes
- [ ] **GET `/api/v1/maintenance`**: List all maintenance records
- [ ] **POST `/api/v1/maintenance`**: Create a new maintenance record
- [ ] **GET `/api/v1/maintenance/{id}`**: Get details of a specific maintenance record
- [ ] **PUT `/api/v1/maintenance/{id}`**: Update an existing maintenance record
- [ ] **DELETE `/api/v1/maintenance/{id}`**: Delete a maintenance record

## Route History Routes
- [ ] **GET `/api/v1/routes`**: List all route histories
- [ ] **POST `/api/v1/routes`**: Create a new route history record
- [ ] **GET `/api/v1/routes/{id}`**: Get details of a specific route history
- [ ] **PUT `/api/v1/routes/{id}`**: Update an existing route history record
- [ ] **DELETE `/api/v1/routes/{id}`**: Delete a route history record

## Warehouse Routes
- [ ] **GET `/api/v1/warehouses`**: List all warehouses
- [ ] **POST `/api/v1/warehouses`**: Create a new warehouse
- [ ] **GET `/api/v1/warehouses/{id}`**: Get details of a specific warehouse
- [ ] **PUT `/api/v1/warehouses/{id}`**: Update an existing warehouse
- [ ] **DELETE `/api/v1/warehouses/{id}`**: Delete a warehouse

## Inventory Routes
- [ ] **GET `/api/v1/inventory`**: List all inventory items
- [ ] **POST `/api/v1/inventory`**: Add a new inventory item
- [ ] **GET `/api/v1/inventory/{id}`**: Get details of a specific inventory item
- [ ] **PUT `/api/v1/inventory/{id}`**: Update an existing inventory item
- [ ] **DELETE `/api/v1/inventory/{id}`**: Delete an inventory item

## Notification Routes
- [ ] **GET `/api/v1/notifications`**: List all notifications
- [ ] **POST `/api/v1/notifications`**: Create a new notification
- [ ] **GET `/api/v1/notifications/{id}`**: Get details of a specific notification
- [ ] **PUT `/api/v1/notifications/{id}`**: Update an existing notification
- [ ] **DELETE `/api/v1/notifications/{id}`**: Delete a notification

## Driver Schedule Routes
- [ ] **GET `/api/v1/driver-schedules`**: List all driver schedules
- [ ] **POST `/api/v1/driver-schedules`**: Create a new driver schedule
- [ ] **GET `/api/v1/driver-schedules/{id}`**: Get details of a specific driver schedule
- [ ] **PUT `/api/v1/driver-schedules/{id}`**: Update an existing driver schedule
- [ ] **DELETE `/api/v1/driver-schedules/{id}`**: Delete a driver schedule

## Fleet Analytics Routes
- [ ] **GET `/api/v1/fleet-analytics`**: List all fleet analytics
- [ ] **POST `/api/v1/fleet-analytics`**: Create a new fleet analytics record
- [ ] **GET `/api/v1/fleet-analytics/{id}`**: Get details of a specific fleet analytics record
- [ ] **PUT `/api/v1/fleet-analytics/{id}`**: Update an existing fleet analytics record
- [ ] **DELETE `/api/v1/fleet-analytics/{id}`**: Delete a fleet analytics record

## Route Optimization Routes
- [ ] **GET `/api/v1/route-optimization`**: List all route optimizations
- [ ] **POST `/api/v1/route-optimization`**: Create a new route optimization record
- [ ] **GET `/api/v1/route-optimization/{id}`**: Get details of a specific route optimization record
- [ ] **PUT `/api/v1/route-optimization/{id}`**: Update an existing route optimization record
- [ ] **DELETE `/api/v1/route-optimization/{id}`**: Delete a route optimization record

## Delivery Feedback Routes
- [ ] **GET `/api/v1/delivery-feedback`**: List all delivery feedback
- [ ] **POST `/api/v1/delivery-feedback`**: Create a new delivery feedback record
- [ ] **GET `/api/v1/delivery-feedback/{id}`**: Get details of a specific delivery feedback record
- [ ] **PUT `/api/v1/delivery-feedback/{id}`**: Update an existing delivery feedback record
- [ ] **DELETE `/api/v1/delivery-feedback/{id}`**: Delete a delivery feedback record



My routes journey:
Login
```sh
curl -X POST http://localhost:8000/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"customer1","passwo
rd":"password123"}'
{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM1NTIwMTgsImlkIjoxLCJyb2xlX25hbWUiOiJjdXN0b21lciIsInVzZXJuYW1lIjoiY3VzdG9tZXIxIn0.TXivif7iDMsSlXN5WcPAEw1LckwxNwJPFYHtcODNnZM","message":"Successfully logged in","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQxNTU5MTgsImlkIjoxLCJ1c2VybmFtZSI6ImN1c3RvbWVyMSJ9.xUFk1rFfHr1rTkLvCTJo_ksngtK3110AtNSbe6vGZfM"}
```
Logout
```sh
curl -X POST http://localhost:8000/api/v1/auth/logout -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
{"message":"Successfully logged out"}
```
```sh
curl -X POST http://localhost:8000/api/v1/auth/logout -H Content-Type: application/json -H Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM1NTIwMTgsImlkIjoxLCJyb2xlX25hbWUiOiJjdXN0b21lciIsInVzZXJuYW1lIjoiY3VzdG9tZXIxIn0.TXivif7iDMsSlXN5WcPAEw1LckwxNwJPFYHtcODNnZM
{"message":"Successfully logged out"}

```
List all vehicles
```sh
curl -X GET http://localhost:8000/api/v1/vehicles \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI5NzU3ODEsImlkIjo5LCJyb2xlX25hbWUiOiJhZG1pbiIsInVzZXJuYW1lIjoidGVzdHVzZXIyIn0.koZN5yeNV2NzsCiCA-GmdkFXAMkBQ7BOagLRqCcH-qo"
[{"id":1,"plate_number":"KAA123A","type":"Truck","capacity":10000,"status":"1","created_at":"2024-11-29T22:01:07.836131Z","updated_at":"2024-11-29T22:01:07.836131Z"},{"id":2,"plate_number":"KBB456B","type":"Van","capacity":2000,"status":"2","created_at":"2024-11-29T22:01:07.836131Z","updated_at":"2024-11-29T22:01:07.836131Z"},{"id":3,"plate_number":"KCC789C","type":"Pickup","capacity":1000,"status":"3","created_at":"2024-11-29T22:01:07.836131Z","updated_at":"2024-11-29T22:01:07.836131Z"}]
```

get a specific vehicle
```sh
curl -X GET http://localhost:8000/api/v1/vehicles/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI5NzU3ODEsImlkIjo5LCJyb2xlX25hbWUiOiJhZG1pbiIsInVzZXJuYW1lIjoidGVzdHVzZXIyIn0.koZN5yeNV2NzsCiCA-GmdkFXAMkBQ7BOagLRqCcH-qo"
{"id":1,"plate_number":"KAA123A","type":"Truck","capacity":10000,"status":"1","created_at":"2024-11-29T22:01:07.836131Z","updated_at":"2024-11-29T22:01:07.836131Z"}
```
Create a new vehicle
```sh


```

Update an existing vehicle
```sh


```

delete a vehicle
```sh


```



List all maintenance records
```sh
curl -X GET http://localhost:8000/api/v1/maintenance \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI5NzU3ODEsImlkIjo5LCJyb2xlX25hbWUiOiJhZG1pbiIsInVzZXJuYW1lIjoidGVzdHVzZXIyIn0.koZN5yeNV2NzsCiCA-GmdkFXAMkBQ7BOagLRqCcH-qo"
[{"id":8,"vehicle_id":1,"type":"Oil Change","description":"Regular maintenance - oil change","service_date":"2024-01-15T10:00:00Z","cost":50,"odometer_reading":15000,"performed_by":"John Doe","next_service_date":"2024-04-15T10:00:00Z","notes":"Used synthetic oil","created_at":"2024-11-29T22:10:36.281763Z","updated_at":"2024-11-29T22:10:36.281763Z"},{"id":5,"vehicle_id":1,"type":"Routine","description":"Oil change and filter replacement"...

```

Update a maintenance record
```sh
curl -X POST http://localhost:8000/api/v1/maintenance \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI5NzU3ODEsImlkIjo5LCJyb2xlX25hbWUiOiJhZG1pbiIsInVzZXJuYW1lIjoidGVzdHVzZXIyIn0.koZN5yeNV2NzsCiCA-GmdkFXAMkBQ7BOagLRqCcH-qo" \
-d '{
    "vehicle_id": 1,
    "type": "Routine",
    "description": "Oil change and filter replacement",
    "service_date": "2024-01-25T10:00:00Z",
    "cost": 150.00,
    "odometer_reading": 55000,
    "performed_by": "John Auto Service",
}'  "notes": "Regular maintenance completed"0Z",
{"id":9,"vehicle_id":1,"type":"Routine","description":"Oil change and filter replacement","service_date":"2024-01-25T10:00:00Z","cost":150,"odometer_reading":55000,"performed_by":"John Auto Service","next_service_date":"2024-04-25T10:00:00Z","notes":"Regular maintenance completed","created_at":"2024-11-30T14:07:24.670826Z","updated_at":"2024-11-30T14:07:24.670826Z"}

```

Update an existing maintenance record
```sh
curl -X PUT http://localhost:8000/api/v1/maintenance/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI5NzU3ODEsImlkIjo5LCJyb2xlX25hbWUiOiJhZG1pbiIsInVzZXJuYW1lIjoidGVzdHVzZXIyIn0.koZN5yeNV2NzsCiCA-GmdkFXAMkBQ7BOagLRqCcH-qo" \
-d '{
    "vehicle_id": 1,
    "type": "Routine",
    "description": "Oil change, filter replacement, and brake check",
    "service_date": "2024-01-25T10:00:00Z",
    "cost": 200.00,
    "odometer_reading": 55000,
    "performed_by": "John Auto Service",
}'  "notes": "Additional brake inspection performed"
{"id":1,"vehicle_id":1,"type":"Routine","description":"Oil change, filter replacement, and brake check","service_date":"2024-01-25T10:00:00Z","cost":200,"odometer_reading":55000,"performed_by":"John Auto Service","next_service_date":"2024-04-25T10:00:00Z","notes":"Additional brake inspection performed","created_at":"2024-11-29T22:01:08.835279Z","updated_at":"2024-11-30T14:07:54.679412Z"}

```

Delivery Routes
```sh


```

Now you can log in with these credentials:

Username: admin
Password: password123

Drop Database:
```sh
docker exec -it postgres_db psql -U postgres_user -d postgres -c "DROP DATABASE IF EXISTS fleetflow;"
```
