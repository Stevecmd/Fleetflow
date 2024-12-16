package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/stevecmd/Fleetflow/backend/models"
)

type WarehouseRepository struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) CreateWarehouse(w *models.CreateWarehouseRequest) (*models.Warehouse, error) {
	query := `
        INSERT INTO warehouses (
            name, street1, street2, city, state, zip, country, 
            latitude, longitude, capacity, status, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
        ) RETURNING id, created_at, updated_at
    `

	now := time.Now()
	status := "active"
	if w.Capacity < 0 {
		return nil, errors.New("warehouse capacity cannot be negative")
	}

	var street2Ptr, statePtr *string
	if w.Street2 != "" {
		street2Ptr = &w.Street2
	}
	if w.State != "" {
		statePtr = &w.State
	}

	warehouse := &models.Warehouse{
		Name:      w.Name,
		Street1:   w.Street1,
		Street2:   street2Ptr,
		City:      w.City,
		State:     statePtr,
		Zip:       w.Zip,
		Country:   w.Country,
		Latitude:  w.Latitude,
		Longitude: w.Longitude,
		Capacity:  w.Capacity,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := r.db.QueryRow(query,
		w.Name, w.Street1, street2Ptr, w.City, statePtr, w.Zip, w.Country,
		w.Latitude, w.Longitude, w.Capacity, status, now, now,
	).Scan(&warehouse.ID, &warehouse.CreatedAt, &warehouse.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error creating warehouse: %w", err)
	}

	return warehouse, nil
}

func (r *WarehouseRepository) GetWarehouseByID(id int) (*models.Warehouse, error) {
	query := `
        SELECT id, name, street1, street2, city, state, zip, country, 
               latitude, longitude, capacity, status, created_at, updated_at 
        FROM warehouses 
        WHERE id = $1
    `

	warehouse := &models.Warehouse{}
	var street2, state sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&warehouse.ID, &warehouse.Name, &warehouse.Street1, &street2,
		&warehouse.City, &state, &warehouse.Zip, &warehouse.Country,
		&warehouse.Latitude, &warehouse.Longitude, &warehouse.Capacity,
		&warehouse.Status, &warehouse.CreatedAt, &warehouse.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("warehouse with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving warehouse: %w", err)
	}

	if street2.Valid {
		warehouse.Street2 = &street2.String
	} else {
		warehouse.Street2 = nil
	}
	if state.Valid {
		warehouse.State = &state.String
	} else {
		warehouse.State = nil
	}

	return warehouse, nil
}

func (r *WarehouseRepository) ListWarehouses(limit, offset int) ([]models.Warehouse, error) {
	query := `
        SELECT id, name, street1, street2, city, state, zip, country, 
               latitude, longitude, capacity, status, created_at, updated_at 
        FROM warehouses 
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing warehouses: %w", err)
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var w models.Warehouse
		var street2, state sql.NullString
		err := rows.Scan(
			&w.ID, &w.Name, &w.Street1, &street2,
			&w.City, &state, &w.Zip, &w.Country,
			&w.Latitude, &w.Longitude, &w.Capacity,
			&w.Status, &w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning warehouse row: %w", err)
		}

		if street2.Valid {
			w.Street2 = &street2.String
		} else {
			w.Street2 = nil
		}

		if state.Valid {
			w.State = &state.String
		} else {
			w.State = nil
		}

		warehouses = append(warehouses, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning warehouses: %w", err)
	}

	return warehouses, nil
}

func (r *WarehouseRepository) UpdateWarehouse(id int, update *models.UpdateWarehouseRequest) (*models.Warehouse, error) {
	query := `
        UPDATE warehouses 
        SET 
            name = COALESCE(NULLIF($2, ''), name),
            street1 = COALESCE(NULLIF($3, ''), street1),
            street2 = COALESCE(NULLIF($4, ''), street2),
            city = COALESCE(NULLIF($5, ''), city),
            state = COALESCE(NULLIF($6, ''), state),
            zip = COALESCE(NULLIF($7, ''), zip),
            country = COALESCE(NULLIF($8, ''), country),
            latitude = COALESCE(NULLIF($9, 0), latitude),
            longitude = COALESCE(NULLIF($10, 0), longitude),
            capacity = COALESCE(NULLIF($11, 0), capacity),
            status = COALESCE(NULLIF($12, ''), status),
            updated_at = $13
        WHERE id = $1
        RETURNING id, name, street1, street2, city, state, zip, country, 
                  latitude, longitude, capacity, status, created_at, updated_at
    `

	now := time.Now()
	warehouse := &models.Warehouse{}
	var street2, state sql.NullString

	err := r.db.QueryRow(
		query,
		id,
		update.Name,
		update.Street1,
		update.Street2,
		update.City,
		update.State,
		update.Zip,
		update.Country,
		update.Latitude,
		update.Longitude,
		update.Capacity,
		update.Status,
		now,
	).Scan(
		&warehouse.ID, &warehouse.Name, &warehouse.Street1, &street2,
		&warehouse.City, &state, &warehouse.Zip, &warehouse.Country,
		&warehouse.Latitude, &warehouse.Longitude, &warehouse.Capacity,
		&warehouse.Status, &warehouse.CreatedAt, &warehouse.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("warehouse with ID %d not found", id)
		}
		return nil, fmt.Errorf("error updating warehouse: %w", err)
	}

	if street2.Valid {
		warehouse.Street2 = &street2.String
	} else {
		warehouse.Street2 = nil
	}
	if state.Valid {
		warehouse.State = &state.String
	} else {
		warehouse.State = nil
	}

	return warehouse, nil
}

func (r *WarehouseRepository) DeleteWarehouse(id int) error {
	query := `DELETE FROM warehouses WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting warehouse: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("warehouse with ID %d not found", id)
	}

	return nil
}
