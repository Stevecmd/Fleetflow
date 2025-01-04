package models

import (
	"database/sql"
	"fmt"
	"time"
)

// MaintenanceRecord represents a vehicle maintenance record

type MaintenanceRecord struct {
	ID              int       `json:"id"`
	VehicleID       int       `json:"vehicle_id"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
	ServiceDate     time.Time `json:"service_date"`
	Cost            float64   `json:"cost"`
	OdometerReading float64   `json:"odometer_reading"`
	PerformedBy     string    `json:"performed_by"`
	NextServiceDate time.Time `json:"next_service_date"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MaintenanceRepository handles database operations for maintenance records
type MaintenanceRepository struct {
	db *sql.DB
}

// NewMaintenanceRepository creates a new maintenance repository
func NewMaintenanceRepository(db *sql.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

// Create inserts a new maintenance record
func (r *MaintenanceRepository) Create(record *MaintenanceRecord) error {
	query := `
		INSERT INTO maintenance_records (
			vehicle_id, type, description, service_date,
			cost, odometer_reading, performed_by, next_service_date,
			notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		record.VehicleID,
		record.Type,
		record.Description,
		record.ServiceDate,
		record.Cost,
		record.OdometerReading,
		record.PerformedBy,
		record.NextServiceDate,
		record.Notes,
	).Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt)
}

// Get retrieves a maintenance record by ID
func (r *MaintenanceRepository) Get(id int) (*MaintenanceRecord, error) {
	record := &MaintenanceRecord{}
	query := `
		SELECT id, vehicle_id, type, description, service_date,
			   cost, odometer_reading, performed_by, next_service_date,
			   notes, created_at, updated_at
		FROM maintenance_records 
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&record.ID,
		&record.VehicleID,
		&record.Type,
		&record.Description,
		&record.ServiceDate,
		&record.Cost,
		&record.OdometerReading,
		&record.PerformedBy,
		&record.NextServiceDate,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// Update modifies an existing maintenance record
func (r *MaintenanceRepository) Update(id int, record *MaintenanceRecord) error {
	query := `
		UPDATE maintenance_records 
		SET vehicle_id = $1, type = $2, description = $3,
			service_date = $4, cost = $5, odometer_reading = $6,
			performed_by = $7, next_service_date = $8, notes = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		RETURNING created_at, updated_at`

	return r.db.QueryRow(
		query,
		record.VehicleID,
		record.Type,
		record.Description,
		record.ServiceDate,
		record.Cost,
		record.OdometerReading,
		record.PerformedBy,
		record.NextServiceDate,
		record.Notes,
		id,
	).Scan(&record.CreatedAt, &record.UpdatedAt)
}

// Delete removes a maintenance record
func (r *MaintenanceRepository) Delete(id int) error {
	query := `DELETE FROM maintenance_records WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// List retrieves maintenance records with optional filters
func (r *MaintenanceRepository) List(vehicleID int, startDate, endDate *time.Time, sortBy string, sortOrder string) ([]MaintenanceRecord, error) {
	query := "SELECT * FROM maintenance_records WHERE 1=1"
	params := []interface{}{}
	paramCount := 0

	if vehicleID > 0 {
		paramCount++
		query += " AND vehicle_id = $" + fmt.Sprint(paramCount)
		params = append(params, vehicleID)
	}

	if startDate != nil {
		paramCount++
		query += " AND service_date >= $" + string(paramCount)
		params = append(params, startDate)
	}

	if endDate != nil {
		paramCount++
		query += " AND service_date <= $" + string(paramCount)
		params = append(params, endDate)
	}

	if sortBy != "" {
		query += " ORDER BY " + sortBy
		if sortOrder == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	} else {
		query += " ORDER BY service_date DESC"
	}

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []MaintenanceRecord
	for rows.Next() {
		var record MaintenanceRecord
		err := rows.Scan(
			&record.ID,
			&record.VehicleID,
			&record.Type,
			&record.Description,
			&record.ServiceDate,
			&record.Cost,
			&record.OdometerReading,
			&record.PerformedBy,
			&record.NextServiceDate,
			&record.Notes,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
