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

// validateMaintenanceRecord checks the validity of a MaintenanceRecord.
// It ensures that all required fields are provided and meet specified constraints.
// Returns an error if any validation check fails, otherwise returns nil.

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

// handleMaintenanceError writes a JSON response to the given http.ResponseWriter
// with a MaintenanceError object in the body. The object's Code and Message
// fields are set according to the type of error provided. The response status
// code is also set accordingly. If the error is not recognized, the error is
// returned as a string in the MaintenanceError object's Details field, and
// the response status code is set to http.StatusInternalServerError.
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

// ListMaintenanceRecords handles HTTP requests to list maintenance records.
// It supports filtering by vehicle ID, start date, and end date, as well as
// sorting by specified fields and order. The function extracts query parameters
// from the request to apply these filters and sorting options. If any parameter
// is invalid or an error occurs during retrieval, an appropriate error response
// is generated. The response includes a JSON encoded list of maintenance records.

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

// CreateMaintenanceRecord handles HTTP requests to create a new maintenance record.
// The function expects a JSON encoded maintenance record in the request body
// and validates the record before inserting it into the database. If any error
// occurs during validation or insertion, an appropriate error response is
// generated. Otherwise, the created record is returned as JSON encoded response
// with a 201 status code.
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

// GetMaintenanceRecord handles HTTP requests to retrieve a maintenance record by ID.
// The function extracts the record ID from the request URL, validates it, and
// retrieves the corresponding record from the database. If the ID is invalid or
// the record cannot be found, an appropriate error response is generated.
// Otherwise, the retrieved record is returned as a JSON encoded response.

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

// UpdateMaintenanceRecord handles HTTP requests to update a maintenance record by ID.
// The function expects a JSON encoded maintenance record in the request body and
// validates the record before updating it in the database. If any error occurs
// during validation or update, an appropriate error response is generated.
// Otherwise, the updated record is returned as a JSON encoded response.
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

// DeleteMaintenanceRecord handles HTTP requests to delete a maintenance record by ID.
// The function extracts the record ID from the request URL, validates it, and
// deletes the corresponding record from the database. If the ID is invalid or
// the record cannot be found, an appropriate error response is generated.
// Otherwise, the response status code is set to 204 (No Content).
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
