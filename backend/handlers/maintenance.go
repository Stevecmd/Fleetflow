package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/stevecmd/Fleetflow/backend/models"
)

type MaintenanceError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func validateMaintenanceRecord(record *models.MaintenanceRecord) error {
	if record.VehicleID <= 0 {
		return errors.New("invalid vehicle ID")
	}
	if record.Type == "" {
		return errors.New("maintenance type is required")
	}
	if record.Description == "" {
		return errors.New("description is required")
	}
	if record.ServiceDate.IsZero() {
		return errors.New("service date is required")
	}
	if record.Cost < 0 {
		return errors.New("cost cannot be negative")
	}
	if record.OdometerReading < 0 {
		return errors.New("odometer reading cannot be negative")
	}
	if record.PerformedBy == "" {
		return errors.New("performed by is required")
	}
	return nil
}

func handleMaintenanceError(err error, w http.ResponseWriter) {
	var response MaintenanceError
	switch {
	case errors.Is(err, sql.ErrNoRows):
		response = MaintenanceError{
			Code:    "RECORD_NOT_FOUND",
			Message: "Maintenance record not found",
		}
		w.WriteHeader(http.StatusNotFound)
	case err.Error() == "invalid vehicle ID":
		response = MaintenanceError{
			Code:    "INVALID_VEHICLE_ID",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
	default:
		response = MaintenanceError{
			Code:    "INTERNAL_ERROR",
			Message: "An internal error occurred",
			Details: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func ListMaintenanceRecords(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := models.NewMaintenanceRepository(db)

		// Parse query parameters
		vehicleID := 0
		if vid := r.URL.Query().Get("vehicle_id"); vid != "" {
			var err error
			vehicleID, err = strconv.Atoi(vid)
			if err != nil {
				handleMaintenanceError(errors.New("invalid vehicle ID"), w)
				return
			}
		}

		var startDate, endDate *time.Time
		if sd := r.URL.Query().Get("start_date"); sd != "" {
			t, err := time.Parse(time.RFC3339, sd)
			if err != nil {
				handleMaintenanceError(errors.New("invalid start date format"), w)
				return
			}
			startDate = &t
		}

		if ed := r.URL.Query().Get("end_date"); ed != "" {
			t, err := time.Parse(time.RFC3339, ed)
			if err != nil {
				handleMaintenanceError(errors.New("invalid end date format"), w)
				return
			}
			endDate = &t
		}

		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")

		records, err := repo.List(vehicleID, startDate, endDate, sortBy, sortOrder)
		if err != nil {
			handleMaintenanceError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)
	}
}

func CreateMaintenanceRecord(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var record models.MaintenanceRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			handleMaintenanceError(err, w)
			return
		}

		if err := validateMaintenanceRecord(&record); err != nil {
			handleMaintenanceError(err, w)
			return
		}

		repo := models.NewMaintenanceRepository(db)
		if err := repo.Create(&record); err != nil {
			handleMaintenanceError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(record)
	}
}

func GetMaintenanceRecord(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			handleMaintenanceError(errors.New("invalid maintenance record ID"), w)
			return
		}

		repo := models.NewMaintenanceRepository(db)
		record, err := repo.Get(id)
		if err != nil {
			handleMaintenanceError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(record)
	}
}

func UpdateMaintenanceRecord(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			handleMaintenanceError(errors.New("invalid maintenance record ID"), w)
			return
		}

		var record models.MaintenanceRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			handleMaintenanceError(err, w)
			return
		}

		if err := validateMaintenanceRecord(&record); err != nil {
			handleMaintenanceError(err, w)
			return
		}

		repo := models.NewMaintenanceRepository(db)
		if err := repo.Update(id, &record); err != nil {
			handleMaintenanceError(err, w)
			return
		}

		record.ID = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(record)
	}
}

func DeleteMaintenanceRecord(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			handleMaintenanceError(errors.New("invalid maintenance record ID"), w)
			return
		}

		repo := models.NewMaintenanceRepository(db)
		if err := repo.Delete(id); err != nil {
			handleMaintenanceError(err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
