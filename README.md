# FleetFlow

![FleetFlow Logo](assets/logo.svg)

FleetFlow is a comprehensive transportation and delivery management system built with modern technologies. It provides robust features for fleet management, driver coordination, delivery tracking, and warehouse operations.

## 🚀 Features

### Driver Management
- Driver profile management with license and certification tracking
- Real-time driver status and availability monitoring
- Performance tracking and rating system
- Document management for licenses and certifications

### Vehicle Fleet Management
- Comprehensive vehicle tracking and management
- Maintenance scheduling and history
- Vehicle capacity and status monitoring
- Fleet analytics and reporting

### Delivery Operations
- Real-time delivery tracking
- Efficient route planning and optimization
- Package and cargo management
- Delivery status updates and notifications

### Warehouse Management
- Warehouse capacity monitoring
- Storage location management
- Inventory tracking
- Loading/unloading coordination

### User Roles
- Fleet Managers: Overall system administration
- Drivers: Delivery and vehicle management
- Loaders: Warehouse operations
- Customers: Delivery tracking and management

## 🛠 Technology Stack

### Frontend
- TBA
- Vite
- Tailwind CSS
- Modern UI/UX design principles

### Backend
- Go (Golang)
- PostgreSQL
- JWT Authentication
- RESTful API architecture

## 📋 Prerequisites

- Go 1.19 or higher
- Node.js 16.x or higher
- PostgreSQL 13 or higher
- Docker (optional)

## 🚀 Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yourusername/fleetflow.git
cd fleetflow
```

2. Set up the backend:
```bash
cd backend
# Copy environment file
cp .env.example .env
# Install dependencies
go mod download
# Start the server
go run main.go
```

3. Set up the frontend:
```bash
cd frontend
# Install dependencies
npm install
# Start development server
npm run dev
```

4. Set up the database:
```bash
# Create database and run migrations
cd backend/db/init
psql -U postgres -f 01-init.sql
```

## 🔒 Environment Variables

Create a `.env` file in the backend directory with the following variables:

```env
DATABASE_URL=postgresql://postgres:password@localhost:5432/fleetflow?sslmode=disable
JWT_SECRET=your_jwt_secret_key
PORT=8000
```

## 📚 API Documentation

### Authentication Endpoints
- POST `/auth/v1/login`: User login
- POST `/auth/v1/register`: User registration
- POST `/auth/v1/refresh`: Refresh JWT token
- POST `/auth/v1/logout`: User logout

### Driver Endpoints
- GET `/api/v1/drivers`: List all drivers
- POST `/api/v1/drivers`: Create new driver profile
- GET `/api/v1/drivers/{id}`: Get driver details
- PUT `/api/v1/drivers/{id}`: Update driver profile
- DELETE `/api/v1/drivers/{id}`: Delete driver profile

### Vehicle Endpoints
- GET `/api/v1/vehicles`: List all vehicles
- POST `/api/v1/vehicles`: Add new vehicle
- GET `/api/v1/vehicles/{id}`: Get vehicle details
- PUT `/api/v1/vehicles/{id}`: Update vehicle
- DELETE `/api/v1/vehicles/{id}`: Delete vehicle

### Delivery Endpoints
- GET `/api/v1/deliveries`: List all deliveries
- POST `/api/v1/deliveries`: Create new delivery
- GET `/api/v1/deliveries/{id}`: Get delivery details
- PUT `/api/v1/deliveries/{id}`: Update delivery status

### Warehouse Endpoints
- GET `/api/v1/warehouses`: List all warehouses
- POST `/api/v1/warehouses`: Add new warehouse
- GET `/api/v1/warehouses/{id}`: Get warehouse details
- PUT `/api/v1/warehouses/{id}`: Update warehouse information

## 🔐 Security

- JWT-based authentication
- Role-based access control
- Secure password hashing
- Rate limiting
- CORS protection

## 🧪 Testing

Run backend tests:
```bash
cd backend
go test ./...
```

Run frontend tests:
```bash
cd frontend
npm test
```

## 📈 Future Enhancements

- Real-time GPS tracking
- Mobile application
- Advanced analytics dashboard
- Automated route optimization
- Integration with external mapping services
- Mobile notifications
- Customer mobile app
- Automated dispatch system
- Advanced reporting features

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- SteveCMD - *Initial work* - [SteveCMD](https://github.com/stevecmd)

## 🙏 Acknowledgments

- Thanks to all contributors who have helped shape FleetFlow
- Inspired by modern logistics and transportation needs
- Built with best practices in software development
