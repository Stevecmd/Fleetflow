package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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
