package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/stevecmd/Fleetflow/backend/models"
)

type FleetStats struct {
	TotalVehicles   int     `json:"total_vehicles"`
	TotalDrivers    int     `json:"total_drivers"`
	AvgDriverRating float64 `json:"avg_driver_rating"`
	TotalDeliveries int     `json:"total_deliveries"`
}

type FleetVehicle struct {
	models.Vehicle
	DriverID      *int    `json:"driver_id,omitempty"`
	LicenseNumber *string `json:"license_number,omitempty"`
}

type MaintenanceMetrics struct {
	PendingMaintenance int                   `json:"pending_maintenance"`
	CompletedLastMonth int                   `json:"completed_last_month"`
	AverageCost        float64               `json:"average_cost"`
	UpcomingServices   []MaintenanceSchedule `json:"upcoming_services"`
}

type MaintenanceSchedule struct {
	VehicleID     int       `json:"vehicle_id"`
	PlateNumber   string    `json:"plate_number"`
	NextService   time.Time `json:"next_service"`
	ServiceType   string    `json:"service_type"`
	EstimatedCost float64   `json:"estimated_cost"`
}

type FleetEfficiencyMetrics struct {
	FuelEfficiency    float64 `json:"fuel_efficiency"`
	CarbonEmissions   float64 `json:"carbon_emissions"`
	OperatingCosts    float64 `json:"operating_costs"`
	IdleTime          int     `json:"idle_time"`
	RouteOptimization float64 `json:"route_optimization"`
}

// GetFleetVehicles retrieves all the vehicles in the fleet, along with their current
// driver and license number (if applicable).
//
// The response is a JSON array of FleetVehicle objects, with the following fields:
//
// - id
// - plate_number
// - type
// - make
// - model
// - year
// - capacity
// - fuel_type
// - status_id
// - gps_unit_id
// - last_maintenance
// - next_maintenance
// - mileage
// - insurance_expiry
// - current_location_latitude
// - current_location_longitude
// - current_location_updated_at
// - fuel_efficiency_rating
// - total_fuel_consumption
// - total_maintenance_cost
// - vehicle_images
// - registration_document_image_url
// - insurance_document_image_url
// - created_at
// - updated_at
// - driver_id (optional)
// - license_number (optional)
func GetFleetVehicles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
            SELECT 
                v.id, v.plate_number, v.type, v.make, v.model, v.year,
                v.capacity, v.fuel_type, v.status_id, v.gps_unit_id,
                v.last_maintenance, v.next_maintenance, v.mileage,
                v.insurance_expiry, v.current_location_latitude,
                v.current_location_longitude, v.current_location_updated_at,
                v.fuel_efficiency_rating, v.total_fuel_consumption,
                v.total_maintenance_cost, v.vehicle_images,
                v.registration_document_image_url, v.insurance_document_image_url,
                v.created_at, v.updated_at,
                d.user_id as driver_id, d.license_number
            FROM vehicles v
            LEFT JOIN driver_profiles d ON v.id = d.current_vehicle_id
            ORDER BY v.id
        `

		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error querying fleet vehicles: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var vehicles []FleetVehicle
		for rows.Next() {
			var v FleetVehicle
			var lastMaintenance, nextMaintenance sql.NullTime

			err := rows.Scan(
				&v.ID, &v.PlateNumber, &v.Type, &v.Make, &v.Model, &v.Year,
				&v.Capacity, &v.FuelType, &v.StatusID, &v.GPSUnitID,
				&lastMaintenance, &nextMaintenance, &v.Mileage,
				&v.InsuranceExpiry, &v.CurrentLocationLat, &v.CurrentLocationLong,
				&v.CurrentLocationUpdated, &v.FuelEfficiencyRating,
				&v.TotalFuelConsumption, &v.TotalMaintenanceCost,
				&v.VehicleImages, &v.RegistrationDocURL, &v.InsuranceDocURL,
				&v.CreatedAt, &v.UpdatedAt,
				&v.DriverID, &v.LicenseNumber,
			)
			if err != nil {
				log.Printf("Error scanning vehicle row: %v", err)
				continue
			}

			if lastMaintenance.Valid {
				v.LastMaintenance = &lastMaintenance.Time
			}
			if nextMaintenance.Valid {
				v.NextMaintenance = &nextMaintenance.Time
			}

			vehicles = append(vehicles, v)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vehicles)
	}
}

// GetFleetPerformance retrieves aggregate performance metrics for the fleet,
// including the total number of vehicles, total number of drivers, average
// driver rating, and total number of deliveries. It executes a SQL query to
// calculate these metrics and returns the data as a JSON response. If an error
// occurs during the query, it logs the error and sends a 500 Internal Server
// Error response.

func GetFleetPerformance(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Starting GetFleetPerformance")

		query := `
            SELECT 
                COUNT(DISTINCT v.id) as total_vehicles,
                COUNT(DISTINCT d.id) as total_drivers,
                COALESCE(AVG(d.rating), 0) as avg_driver_rating,
                COUNT(DISTINCT del.id) as total_deliveries
            FROM vehicles v
            LEFT JOIN driver_profiles d ON v.id = d.current_vehicle_id
            LEFT JOIN deliveries del ON d.id = del.driver_id
        `

		var stats FleetStats
		err := db.QueryRow(query).Scan(
			&stats.TotalVehicles,
			&stats.TotalDrivers,
			&stats.AvgDriverRating,
			&stats.TotalDeliveries,
		)

		if err != nil {
			log.Printf("Error querying fleet performance: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

// GetFleetDrivers retrieves a list of all drivers in the fleet, along with
// their profile information and the vehicle they are currently assigned to.
// It executes a SQL query to retrieve the data and returns the list of
// drivers as a JSON response. If an error occurs during the query, it logs
// the error and sends a 500 Internal Server Error response.
func GetFleetDrivers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
            SELECT 
                d.id, d.user_id, d.license_number, d.license_type,
                d.license_expiry, d.vehicle_type, d.years_experience,
                d.certification, d.status, d.status_id,
                d.current_vehicle_id, d.rating, d.total_trips,
                d.created_at, d.updated_at,
                u.first_name, u.last_name, u.email,
                v.plate_number as assigned_vehicle
            FROM driver_profiles d
            LEFT JOIN users u ON d.user_id = u.id
            LEFT JOIN vehicles v ON d.current_vehicle_id = v.id
            ORDER BY d.id
        `

		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error querying fleet drivers: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var drivers []models.DriverProfile
		for rows.Next() {
			var d models.DriverProfile
			var licenseExpiry sql.NullTime
			var certification []string

			err := rows.Scan(
				&d.ID, &d.UserID, &d.LicenseNumber, &d.LicenseType,
				&licenseExpiry, &d.VehicleType, &d.YearsExperience,
				pq.Array(&certification), &d.Status, &d.StatusID,
				&d.CurrentVehicleID, &d.Rating, &d.TotalTrips,
				&d.CreatedAt, &d.UpdatedAt,
				&d.FirstName, &d.LastName, &d.Email,
				&d.AssignedVehicle,
			)
			if err != nil {
				log.Printf("Error scanning driver row: %v", err)
				continue
			}

			if licenseExpiry.Valid {
				d.LicenseExpiry = &licenseExpiry.Time
			}
			d.Certification = certification

			drivers = append(drivers, d)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(drivers)
	}
}

func GetMaintenanceMetrics(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
            SELECT 
                COUNT(CASE WHEN next_service_date > NOW() THEN 1 END) as pending,
                COUNT(CASE WHEN service_date >= NOW() - INTERVAL '30 days' THEN 1 END) as completed,
                AVG(cost) as avg_cost
            FROM maintenance_records
        `
		var metrics MaintenanceMetrics
		err := db.QueryRow(query).Scan(&metrics.PendingMaintenance, &metrics.CompletedLastMonth, &metrics.AverageCost)
		if err != nil {
			log.Printf("Error querying maintenance metrics: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Get upcoming services
		upcomingQuery := `
            SELECT v.id, v.plate_number, m.next_service_date, m.type, m.cost
            FROM maintenance_records m
            JOIN vehicles v ON m.vehicle_id = v.id
            WHERE m.next_service_date > NOW()
            ORDER BY m.next_service_date
            LIMIT 5
        `
		rows, err := db.Query(upcomingQuery)
		if err != nil {
			log.Printf("Error querying upcoming services: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var service MaintenanceSchedule
			err := rows.Scan(
				&service.VehicleID,
				&service.PlateNumber,
				&service.NextService,
				&service.ServiceType,
				&service.EstimatedCost,
			)
			if err != nil {
				log.Printf("Error scanning upcoming service row: %v", err)
				continue
			}
			metrics.UpcomingServices = append(metrics.UpcomingServices, service)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	}
}

func GetFleetEfficiency(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
            SELECT 
                COALESCE(AVG(NULLIF(fuel_efficiency_rating, 0)), 0) as fuel_efficiency,
                COALESCE(SUM(carbon_emissions), 0) as total_emissions,
                COALESCE(SUM(operating_cost), 0) as total_costs,
                COALESCE(AVG(NULLIF(idle_time, 0)), 0) as avg_idle_time,
                COALESCE(AVG(NULLIF(efficiency_score, 0)), 0) as route_efficiency
            FROM fleet_analytics
            WHERE analysis_date >= NOW() - INTERVAL '30 days'
        `
		var metrics FleetEfficiencyMetrics
		err := db.QueryRow(query).Scan(
			&metrics.FuelEfficiency,
			&metrics.CarbonEmissions,
			&metrics.OperatingCosts,
			&metrics.IdleTime,
			&metrics.RouteOptimization,
		)
		if err != nil {
			log.Printf("Error querying fleet efficiency: %v", err)
			// Return default values instead of error
			metrics = FleetEfficiencyMetrics{
				FuelEfficiency:    0,
				CarbonEmissions:   0,
				OperatingCosts:    0,
				IdleTime:          0,
				RouteOptimization: 0,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(metrics); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func GetDeliveryAnalytics(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
            SELECT 
                COUNT(*) as total_deliveries,
                COALESCE(COUNT(CASE 
                    WHEN actual_delivery_time IS NOT NULL 
                    AND actual_delivery_time <= estimated_delivery_time 
                    THEN 1 
                END), 0) as on_time,
                COALESCE(AVG(NULLIF(route_efficiency_score, 0)), 0) as efficiency,
                COUNT(DISTINCT CASE WHEN status_id = 2 THEN user_id END) as active_drivers
            FROM deliveries
            WHERE delivery_time >= NOW() - INTERVAL '30 days'
        `
		var analytics struct {
			TotalDeliveries  int     `json:"total_deliveries"`
			OnTimeDeliveries int     `json:"on_time_deliveries"`
			Efficiency       float64 `json:"efficiency"`
			ActiveDrivers    int     `json:"active_drivers"`
		}

		err := db.QueryRow(query).Scan(
			&analytics.TotalDeliveries,
			&analytics.OnTimeDeliveries,
			&analytics.Efficiency,
			&analytics.ActiveDrivers,
		)
		if err != nil {
			log.Printf("Error querying delivery analytics: %v", err)
			// Return default values instead of error
			analytics = struct {
				TotalDeliveries  int     `json:"total_deliveries"`
				OnTimeDeliveries int     `json:"on_time_deliveries"`
				Efficiency       float64 `json:"efficiency"`
				ActiveDrivers    int     `json:"active_drivers"`
			}{
				TotalDeliveries:  0,
				OnTimeDeliveries: 0,
				Efficiency:       0,
				ActiveDrivers:    0,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(analytics); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
