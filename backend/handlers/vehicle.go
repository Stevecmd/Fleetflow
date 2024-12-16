package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stevecmd/Fleetflow/backend/models"
)

// CreateVehicle creates a new vehicle
func CreateVehicle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vehicle models.Vehicle
		if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var lastMaintenance, nextMaintenance sql.NullTime
		if vehicle.LastMaintenance != nil {
			lastMaintenance.Time = *vehicle.LastMaintenance
			lastMaintenance.Valid = true
		}
		if vehicle.NextMaintenance != nil {
			nextMaintenance.Time = *vehicle.NextMaintenance
			nextMaintenance.Valid = true
		}

		query := `
			INSERT INTO vehicles (
				plate_number, type, make, model, year, capacity, fuel_type, 
				status_id, gps_unit_id, last_maintenance, next_maintenance, mileage,
				insurance_expiry, current_location_latitude, current_location_longitude,
				current_location_updated_at, fuel_efficiency_rating, total_fuel_consumption,
				total_maintenance_cost, vehicle_images, registration_document_image_url,
				insurance_document_image_url
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
				$16, $17, $18, $19, $20, $21, $22
			) RETURNING id, created_at, updated_at`

		err := db.QueryRow(
			query,
			vehicle.PlateNumber, vehicle.Type, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.Capacity,
			vehicle.FuelType, vehicle.StatusID, vehicle.GPSUnitID, lastMaintenance,
			nextMaintenance, vehicle.Mileage, vehicle.InsuranceExpiry,
			vehicle.CurrentLocationLat, vehicle.CurrentLocationLong, vehicle.CurrentLocationUpdated,
			vehicle.FuelEfficiencyRating, vehicle.TotalFuelConsumption, vehicle.TotalMaintenanceCost,
			vehicle.VehicleImages, vehicle.RegistrationDocURL, vehicle.InsuranceDocURL,
		).Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vehicle)
	}
}

// GetVehicle retrieves a vehicle by ID
func GetVehicle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		log.Printf("Received ID: %s", idStr) // Debugging
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var vehicle models.Vehicle
		query := `
			SELECT id, plate_number, type, make, model, year, capacity, fuel_type, 
			status_id, gps_unit_id, last_maintenance, next_maintenance, mileage,
			insurance_expiry, current_location_latitude, current_location_longitude,
			current_location_updated_at, fuel_efficiency_rating, total_fuel_consumption,
			total_maintenance_cost, vehicle_images, registration_document_image_url,
			insurance_document_image_url, created_at, updated_at
			FROM vehicles
			WHERE id = $1`

		var lastMaintenance, nextMaintenance sql.NullTime
		err = db.QueryRow(query, id).Scan(
			&vehicle.ID, &vehicle.PlateNumber, &vehicle.Type, &vehicle.Make, &vehicle.Model, &vehicle.Year,
			&vehicle.Capacity, &vehicle.FuelType, &vehicle.StatusID, &vehicle.GPSUnitID,
			&lastMaintenance, &nextMaintenance, &vehicle.Mileage,
			&vehicle.InsuranceExpiry, &vehicle.CurrentLocationLat, &vehicle.CurrentLocationLong,
			&vehicle.CurrentLocationUpdated, &vehicle.FuelEfficiencyRating,
			&vehicle.TotalFuelConsumption, &vehicle.TotalMaintenanceCost,
			&vehicle.VehicleImages, &vehicle.RegistrationDocURL, &vehicle.InsuranceDocURL,
			&vehicle.CreatedAt, &vehicle.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if lastMaintenance.Valid {
			vehicle.LastMaintenance = &lastMaintenance.Time
		}
		if nextMaintenance.Valid {
			vehicle.NextMaintenance = &nextMaintenance.Time
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vehicle)
	}
}

// UpdateVehicle updates a vehicle
func UpdateVehicle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var vehicle models.Vehicle
		if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var lastMaintenance, nextMaintenance sql.NullTime
		if vehicle.LastMaintenance != nil {
			lastMaintenance.Time = *vehicle.LastMaintenance
			lastMaintenance.Valid = true
		}
		if vehicle.NextMaintenance != nil {
			nextMaintenance.Time = *vehicle.NextMaintenance
			nextMaintenance.Valid = true
		}

		query := `
			UPDATE vehicles SET
				plate_number = $1, type = $2, make = $3, model = $4, year = $5,
				capacity = $6, fuel_type = $7, status_id = $8, gps_unit_id = $9,
				last_maintenance = $10, next_maintenance = $11, mileage = $12,
				insurance_expiry = $13, current_location_latitude = $14,
				current_location_longitude = $15, current_location_updated_at = $16,
				fuel_efficiency_rating = $17, total_fuel_consumption = $18,
				total_maintenance_cost = $19, vehicle_images = $20,
				registration_document_image_url = $21, insurance_document_image_url = $22,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $23
			RETURNING created_at, updated_at`

		err = db.QueryRow(
			query,
			vehicle.PlateNumber, vehicle.Type, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.Capacity,
			vehicle.FuelType, vehicle.StatusID, vehicle.GPSUnitID, lastMaintenance,
			nextMaintenance, vehicle.Mileage, vehicle.InsuranceExpiry,
			vehicle.CurrentLocationLat, vehicle.CurrentLocationLong, vehicle.CurrentLocationUpdated,
			vehicle.FuelEfficiencyRating, vehicle.TotalFuelConsumption, vehicle.TotalMaintenanceCost,
			vehicle.VehicleImages, vehicle.RegistrationDocURL, vehicle.InsuranceDocURL, id,
		).Scan(&vehicle.CreatedAt, &vehicle.UpdatedAt)

		if err == sql.ErrNoRows {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vehicle)
	}
}

// ListVehicles lists all vehicles
func ListVehicles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT id, plate_number, type, make, model, year, capacity, fuel_type, 
			status_id, gps_unit_id, last_maintenance, next_maintenance, mileage,
			insurance_expiry, current_location_latitude, current_location_longitude,
			current_location_updated_at, fuel_efficiency_rating, total_fuel_consumption,
			total_maintenance_cost, vehicle_images, registration_document_image_url,
			insurance_document_image_url, created_at, updated_at
			FROM vehicles`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var vehicles []models.Vehicle
		for rows.Next() {
			var vehicle models.Vehicle
			var lastMaintenance, nextMaintenance sql.NullTime
			err := rows.Scan(
				&vehicle.ID, &vehicle.PlateNumber, &vehicle.Type, &vehicle.Make, &vehicle.Model, &vehicle.Year,
				&vehicle.Capacity, &vehicle.FuelType, &vehicle.StatusID, &vehicle.GPSUnitID,
				&lastMaintenance, &nextMaintenance, &vehicle.Mileage,
				&vehicle.InsuranceExpiry, &vehicle.CurrentLocationLat, &vehicle.CurrentLocationLong,
				&vehicle.CurrentLocationUpdated, &vehicle.FuelEfficiencyRating,
				&vehicle.TotalFuelConsumption, &vehicle.TotalMaintenanceCost,
				&vehicle.VehicleImages, &vehicle.RegistrationDocURL, &vehicle.InsuranceDocURL,
				&vehicle.CreatedAt, &vehicle.UpdatedAt,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if lastMaintenance.Valid {
				vehicle.LastMaintenance = &lastMaintenance.Time
			}
			if nextMaintenance.Valid {
				vehicle.NextMaintenance = &nextMaintenance.Time
			}

			vehicles = append(vehicles, vehicle)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vehicles)
	}
}

// DeleteVehicle deletes a vehicle
func DeleteVehicle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "DELETE FROM vehicles WHERE id = $1"
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
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
