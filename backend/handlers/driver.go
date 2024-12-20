package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/stevecmd/Fleetflow/backend/models"
)

// CreateDriverProfile creates a new driver profile
func CreateDriverProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var profile models.DriverProfile
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO driver_profiles (
				user_id, license_number, license_type, license_expiry,
				vehicle_type, years_experience, certification,
				status_id, current_vehicle_id, rating, total_trips
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, created_at, updated_at`

		err := db.QueryRow(
			query,
			profile.UserID,
			profile.LicenseNumber,
			profile.LicenseType,
			profile.LicenseExpiry,
			profile.VehicleType,
			profile.YearsExperience,
			pq.Array(profile.Certification),
			profile.StatusID,
			profile.CurrentVehicleID,
			profile.Rating,
			profile.TotalTrips,
		).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating driver profile: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

// GetDriverProfile retrieves a driver profile by ID
func GetDriverProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var profile models.DriverProfile
		query := `
			SELECT 
				id, user_id, license_number, license_type, vehicle_type, 
				years_experience, certification, status, license_expiry,
				status_id, current_vehicle_id, rating, total_trips, created_at, updated_at
			FROM driver_profiles
			LEFT JOIN driver_statuses ON driver_profiles.status_id = driver_statuses.id
			WHERE driver_profiles.id = $1`

		var licenseExpiry sql.NullTime
		var statusID, currentVehicleID sql.NullInt64
		var rating sql.NullFloat64
		var totalTrips sql.NullInt64
		var yearsExperience string

		err = db.QueryRow(query, id).Scan(
			&profile.ID,
			&profile.UserID,
			&profile.LicenseNumber,
			&profile.LicenseType,
			&profile.VehicleType,
			&yearsExperience,
			pq.Array(&profile.Certification),
			&profile.Status,
			&licenseExpiry,
			&statusID,
			&currentVehicleID,
			&rating,
			&totalTrips,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			http.Error(w, "Driver profile not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Handle nullable fields
		if licenseExpiry.Valid {
			profile.LicenseExpiry = &licenseExpiry.Time
		}
		if statusID.Valid {
			intVal := int(statusID.Int64)
			profile.StatusID = &intVal
		}
		if currentVehicleID.Valid {
			intVal := int(currentVehicleID.Int64)
			profile.CurrentVehicleID = &intVal
		}
		if rating.Valid {
			floatVal := rating.Float64
			profile.Rating = &floatVal
		}
		if totalTrips.Valid {
			intVal := int(totalTrips.Int64)
			profile.TotalTrips = &intVal
		}

		// Convert years_experience back to string
		profile.YearsExperience = yearsExperience

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

// UpdateDriverProfile updates a driver profile
func UpdateDriverProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		var profile models.DriverProfile
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query := `
			UPDATE driver_profiles SET
				user_id = $1,
				license_number = $2,
				license_type = $3,
				license_expiry = $4,
				vehicle_type = $5,
				years_experience = $6,
				certification = $7,
				status_id = $8,
				current_vehicle_id = $9,
				rating = $10,
				total_trips = $11,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $12
			RETURNING id, created_at, updated_at`

		err = db.QueryRow(
			query,
			profile.UserID,
			profile.LicenseNumber,
			profile.LicenseType,
			profile.LicenseExpiry,
			profile.VehicleType,
			profile.YearsExperience,
			pq.Array(profile.Certification),
			profile.StatusID,
			profile.CurrentVehicleID,
			profile.Rating,
			profile.TotalTrips,
			id,
		).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Profile not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error updating driver profile: %v", err), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

// ListDriverProfiles lists all driver profiles
func ListDriverProfiles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT 
				dp.id, dp.user_id, dp.license_number,
				COALESCE(dp.license_type, '') as license_type,
				COALESCE(dp.vehicle_type, '') as vehicle_type,
				COALESCE(dp.years_experience::text, '') as years_experience,
				COALESCE(dp.certification, ARRAY[]::TEXT[]) as certification,
				COALESCE(ds.name, '') as status,
				dp.license_expiry,
				dp.status_id,
				dp.current_vehicle_id,
				dp.rating,
				dp.total_trips,
				dp.created_at, dp.updated_at
			FROM driver_profiles dp
			LEFT JOIN driver_statuses ds ON dp.status_id = ds.id`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database query error: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var profiles []models.DriverProfile
		for rows.Next() {
			var profile models.DriverProfile
			var licenseExpiry sql.NullTime
			var statusID, currentVehicleID sql.NullInt64
			var rating sql.NullFloat64
			var totalTrips sql.NullInt64
			var yearsExperience string

			err := rows.Scan(
				&profile.ID,
				&profile.UserID,
				&profile.LicenseNumber,
				&profile.LicenseType,
				&profile.VehicleType,
				&yearsExperience,
				pq.Array(&profile.Certification),
				&profile.Status,
				&licenseExpiry,
				&statusID,
				&currentVehicleID,
				&rating,
				&totalTrips,
				&profile.CreatedAt,
				&profile.UpdatedAt,
			)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
				return
			}

			// Handle nullable fields
			if licenseExpiry.Valid {
				profile.LicenseExpiry = &licenseExpiry.Time
			}
			if statusID.Valid {
				intVal := int(statusID.Int64)
				profile.StatusID = &intVal
			}
			if currentVehicleID.Valid {
				intVal := int(currentVehicleID.Int64)
				profile.CurrentVehicleID = &intVal
			}
			if rating.Valid {
				floatVal := rating.Float64
				profile.Rating = &floatVal
			}
			if totalTrips.Valid {
				intVal := int(totalTrips.Int64)
				profile.TotalTrips = &intVal
			}

			// Convert years_experience back to string
			profile.YearsExperience = yearsExperience

			profiles = append(profiles, profile)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Error iterating rows: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(profiles); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// DeleteDriverProfile deletes a driver profile
func DeleteDriverProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "DELETE FROM driver_profiles WHERE id = $1"
		result, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Driver profile not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// getDriverVehicle retrieves the vehicle associated with a driver
func GetDriverVehicle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"] // This should be the actual user ID

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Printf("Error converting userID to int: %v", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		log.Printf("Received request for vehicle with userID: %s", userID)
		log.Printf("Querying vehicle for userID: %d", userIDInt)

		// Query to get the vehicle associated with the user
		query := `
            SELECT v.id, v.plate_number, v.type, v.make, v.model, v.year, v.capacity, 
            v.fuel_type, v.status_id, v.gps_unit_id, v.last_maintenance, v.next_maintenance, 
            v.mileage, v.insurance_expiry
            FROM vehicles v
            JOIN driver_profiles dp ON v.id = dp.current_vehicle_id
            WHERE dp.user_id = $1
        `
		var vehicle models.Vehicle
		err = db.QueryRow(query, userIDInt).Scan(
			&vehicle.ID, &vehicle.PlateNumber, &vehicle.Type, &vehicle.Make,
			&vehicle.Model, &vehicle.Year, &vehicle.Capacity, &vehicle.FuelType,
			&vehicle.StatusID, &vehicle.GPSUnitID, &vehicle.LastMaintenance,
			&vehicle.NextMaintenance, &vehicle.Mileage, &vehicle.InsuranceExpiry,
		)
		if err == sql.ErrNoRows {
			http.Error(w, "No vehicle found for this driver", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving vehicle: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vehicle)
	}
}

// GetDriverOrders retrieves orders for a specific driver
func GetDriverOrders(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Printf("Error converting userID to int: %v", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		log.Printf("Received request for orders with userID: %s", userID)

		log.Printf("Querying orders for userID: %d", userIDInt)

		// Query to get orders for the specific driver using user_id
		query := `
			SELECT id, from_location, to_location, weight, status, 
			created_at, updated_at
			FROM orders
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT 10
		`

		rows, err := db.Query(query, userIDInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orders []models.Delivery
		for rows.Next() {
			var order models.Delivery
			err := rows.Scan(
				&order.ID, &order.FromLocation, &order.ToLocation,
				&order.Weight, &order.Status, &order.CreatedAt, &order.UpdatedAt,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			orders = append(orders, order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}

// GetDriverPerformance retrieves performance metrics for a specific driver
func GetDriverPerformance(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Query performance metrics
		query := `
            SELECT 
                total_deliveries,
                on_time_deliveries,
                avg_rating,
                total_distance
            FROM driver_performance_metrics
            WHERE user_id = $1
            ORDER BY created_at DESC
            LIMIT 1
        `

		var performance struct {
			TotalDeliveries  int     `json:"total_deliveries"`
			OnTimeDeliveries int     `json:"on_time_deliveries"`
			AvgRating        float64 `json:"avg_rating"`
			TotalDistance    float64 `json:"total_distance"`
		}

		err = db.QueryRow(query, userIDInt).Scan(
			&performance.TotalDeliveries,
			&performance.OnTimeDeliveries,
			&performance.AvgRating,
			&performance.TotalDistance,
		)

		if err == sql.ErrNoRows {
			http.Error(w, "Performance data not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(performance)
	}
}
