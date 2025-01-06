// Package main FleetFlow API
//
// @title FleetFlow API
// @version 1.0
// @description Fleet management system API
// @host localhost:8000
// @BasePath /api/v1
package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
	"github.com/rs/cors"
	_ "github.com/stevecmd/Fleetflow/backend/docs" // Import generated docs
	"github.com/stevecmd/Fleetflow/backend/handlers"
	"github.com/stevecmd/Fleetflow/backend/middleware"
	"github.com/stevecmd/Fleetflow/backend/repository"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title FleetFlow API
// @version 1.0
// @description This is a sample server for FleetFlow.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1
func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get JWT secret and decode from base64
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	jwtKey, err := base64.StdEncoding.DecodeString(jwtSecret)
	if err != nil {
		log.Fatal("Error decoding JWT secret:", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	log.Printf("Connecting to database with URL: %s", databaseURL)

	// Initialize database connection
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
	log.Println("Successfully connected to database")

	defer db.Close()

	handlers.DatabaseURL = databaseURL // Pass the database URL to the handlers package
	handlers.JwtKey = jwtKey
	middleware.JwtKey = jwtKey

	// Create a new router
	router := mux.NewRouter()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend URL only
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,  // Maximum value not ignored by any of major browsers
		Debug:            true, // Enable debugging for troubleshooting
	})

	// Create handler chain with CORS and rate limiting
	handler := c.Handler(router)
	handler = middleware.RateLimitMiddleware(handler)

	// Add middleware to all routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add CORS headers to all responses
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Auth routes
	router.HandleFunc("/api/v1/auth/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/refresh", handlers.RefreshTokenHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/logout", handlers.LogoutHandler).Methods("POST")
	router.HandleFunc("/api/v1/protected", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Protected route accessed")
	})).Methods("GET")

	// User profile routes
	router.HandleFunc("/api/v1/users/profile", middleware.AuthMiddleware(handlers.GetUserProfile)).Methods("GET")
	router.HandleFunc("/api/v1/users", middleware.AuthMiddleware(handlers.ListUsers)).Methods("GET")
	router.HandleFunc("/api/v1/users/{userID}", middleware.AuthMiddleware(handlers.GetUserProfile)).Methods("GET")
	router.HandleFunc("/api/v1/users/{userID}", middleware.AuthMiddleware(handlers.UpdateUserProfile)).Methods("PUT")
	router.HandleFunc("/api/v1/users/{id}", middleware.AuthMiddleware(handlers.DeleteUser(db))).Methods("DELETE")

	// Customer routes
	router.HandleFunc("/api/v1/users/{userID}/deliveries", middleware.AuthMiddleware(handlers.GetCustomerDeliveries(db))).Methods("GET")
	router.HandleFunc("/api/v1/users/{userID}/invoices", middleware.AuthMiddleware(handlers.GetCustomerInvoices(db))).Methods("GET")
	router.HandleFunc("/api/v1/users/{userID}/feedback", middleware.AuthMiddleware(handlers.GetCustomerFeedback(db))).Methods("GET")

	// Driver profile routes
	router.HandleFunc("/api/v1/drivers", middleware.AuthMiddleware(handlers.ListDriverProfiles(db))).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{userID}", middleware.AuthMiddleware(handlers.GetDriverProfile(db))).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{userID}", middleware.AuthMiddleware(handlers.UpdateDriverProfile(db))).Methods("PUT")
	router.HandleFunc("/api/v1/drivers/{userID}", middleware.AuthMiddleware(handlers.DeleteDriverProfile(db))).Methods("DELETE")
	router.HandleFunc("/api/v1/drivers/{userID}/vehicle", middleware.AuthMiddleware(handlers.GetDriverVehicle(db))).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{userID}/performance", middleware.AuthMiddleware(handlers.GetDriverPerformance(db))).Methods("GET")

	// Driver routes
	router.HandleFunc("/api/v1/drivers/{userID}/orders", middleware.AuthMiddleware(handlers.GetDriverOrders(db))).Methods("GET")

	// Fleet Manager routes
	router.HandleFunc("/api/v1/fleet-manager/vehicles", middleware.AuthMiddleware(handlers.GetFleetVehicles(db))).Methods("GET")
	router.HandleFunc("/api/v1/fleet-manager/performance", middleware.AuthMiddleware(handlers.GetFleetPerformance(db))).Methods("GET")
	router.HandleFunc("/api/v1/fleet-manager/drivers", middleware.AuthMiddleware(handlers.GetFleetDrivers(db))).Methods("GET")

	// Vehicle routes
	router.HandleFunc("/api/v1/vehicles", handlers.ListVehicles(db)).Methods("GET")
	router.HandleFunc("/api/v1/vehicles", middleware.AuthMiddleware(handlers.CreateVehicle(db))).Methods("POST")
	router.HandleFunc("/api/v1/vehicles/{userID}", handlers.GetVehicle(db)).Methods("GET")
	router.HandleFunc("/api/v1/vehicles/{userID}", middleware.AuthMiddleware(handlers.UpdateVehicle(db))).Methods("PUT")
	router.HandleFunc("/api/v1/vehicles/{userID}", middleware.AuthMiddleware(handlers.DeleteVehicle(db))).Methods("DELETE")

	// Maintenance routes
	router.HandleFunc("/api/v1/maintenance", handlers.ListMaintenanceRecords(db)).Methods("GET")
	router.HandleFunc("/api/v1/maintenance", middleware.AuthMiddleware(handlers.CreateMaintenanceRecord(db))).Methods("POST")
	router.HandleFunc("/api/v1/maintenance/{userID}", handlers.GetMaintenanceRecord(db)).Methods("GET")
	router.HandleFunc("/api/v1/maintenance/{userID}", middleware.AuthMiddleware(handlers.UpdateMaintenanceRecord(db))).Methods("PUT")
	router.HandleFunc("/api/v1/maintenance/{userID}", middleware.AuthMiddleware(handlers.DeleteMaintenanceRecord(db))).Methods("DELETE")

	// Delivery routes
	deliveryRepo := repository.NewDeliveryRepository(db)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryRepo)

	router.HandleFunc("/api/v1/deliveries/stats", middleware.AuthMiddleware(deliveryHandler.GetDeliveryStatistics)).Methods("GET")
	router.HandleFunc("/api/v1/deliveries", middleware.AuthMiddleware(deliveryHandler.CreateDelivery)).Methods("POST")
	router.HandleFunc("/api/v1/deliveries", middleware.AuthMiddleware(deliveryHandler.ListDeliveries)).Methods("GET")
	router.HandleFunc("/api/v1/deliveries/{userID}", middleware.AuthMiddleware(deliveryHandler.GetDelivery)).Methods("GET")
	router.HandleFunc("/api/v1/deliveries/{userID}", middleware.AuthMiddleware(deliveryHandler.UpdateDelivery)).Methods("PUT")

	// Warehouse routes
	warehouseHandler := handlers.NewWarehouseHandler(db)
	router.HandleFunc("/api/v1/warehouses", warehouseHandler.ListWarehouses).Methods("GET")
	router.HandleFunc("/api/v1/warehouses", middleware.AuthMiddleware(warehouseHandler.CreateWarehouse)).Methods("POST")
	router.HandleFunc("/api/v1/warehouses/{userID}", warehouseHandler.GetWarehouse).Methods("GET")
	router.HandleFunc("/api/v1/warehouses/{userID}", middleware.AuthMiddleware(warehouseHandler.UpdateWarehouse)).Methods("PUT")
	router.HandleFunc("/api/v1/warehouses/{userID}", middleware.AuthMiddleware(warehouseHandler.DeleteWarehouse)).Methods("DELETE")

	// Fleet analytics route
	router.HandleFunc("/api/v1/fleet-analytics", middleware.AuthMiddleware(handlers.GetFleetAnalytics(db))).Methods("GET")

	// Swagger documentation route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("Server starting on :8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal(err)
	}
}
