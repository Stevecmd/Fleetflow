package repository

// DeliveryRepository handles database operations for deliveries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/stevecmd/Fleetflow/backend/models"
)

type DeliveryRepository struct {
	db *sql.DB
}

func NewDeliveryRepository(db *sql.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

// Create inserts a new delivery record
func (r *DeliveryRepository) Create(req *models.CreateDeliveryRequest) (*models.Delivery, error) {
	query := `
		INSERT INTO deliveries (
			tracking_number, customer_id, cargo_type, cargo_weight,
			special_instructions, pickup_warehouse_id, delivery_warehouse_id,
			pickup_latitude, pickup_longitude, delivery_latitude, delivery_longitude,
			estimated_delivery_time, status_id, payment_status,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		) RETURNING id, tracking_number, customer_id, cargo_type, cargo_weight,
			special_instructions, pickup_warehouse_id, delivery_warehouse_id,
			pickup_latitude, pickup_longitude, delivery_latitude, delivery_longitude,
			estimated_delivery_time, status_id, payment_status,
			created_at, updated_at`

	trackingNumber := generateTrackingNumber()
	var result models.Delivery

	err := r.db.QueryRow(
		query,
		trackingNumber,
		req.CustomerID,
		req.CargoType,
		req.CargoWeight,
		req.SpecialInstructions,
		req.PickupWarehouseID,
		req.DeliveryWarehouseID,
		req.PickupLatitude,
		req.PickupLongitude,
		req.DeliveryLatitude,
		req.DeliveryLongitude,
		req.EstimatedDeliveryTime,
		1,         // Initial status (pending)
		"pending", // Initial payment status
	).Scan(
		&result.ID,
		&result.TrackingNumber,
		&result.CustomerID,
		&result.CargoType,
		&result.CargoWeight,
		&result.SpecialInstructions,
		&result.PickupWarehouseID,
		&result.DeliveryWarehouseID,
		&result.PickupLatitude,
		&result.PickupLongitude,
		&result.DeliveryLatitude,
		&result.DeliveryLongitude,
		&result.EstimatedDeliveryTime,
		&result.StatusID,
		&result.PaymentStatus,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating delivery: %v", err)
	}

	return &result, nil
}

// GetByID retrieves a delivery by its ID
func (r *DeliveryRepository) GetByID(id int64) (*models.DeliveryResponse, error) {
	query := `
		SELECT 
			d.id, d.tracking_number, d.customer_id, d.cargo_type, d.cargo_weight,
			d.special_instructions, d.pickup_warehouse_id, d.delivery_warehouse_id,
			d.pickup_latitude, d.pickup_longitude, d.delivery_latitude, d.delivery_longitude,
			d.estimated_delivery_time, d.status_id, d.payment_status, d.created_at, d.updated_at,
			u.id as customer_id, u.username as customer_username, u.email as customer_email,
			dp.id as driver_id, dp.license_number, dp.years_experience,
			v.id as vehicle_id, v.plate_number, v.type as vehicle_type,
			ds.name as status_name, ds.description as status_description
		FROM deliveries d
		LEFT JOIN users u ON d.customer_id = u.id
		LEFT JOIN driver_profiles dp ON d.driver_id = dp.id
		LEFT JOIN vehicles v ON d.vehicle_id = v.id
		LEFT JOIN delivery_statuses ds ON d.status_id = ds.id
		WHERE d.id = $1`

	var resp models.DeliveryResponse
	resp.Delivery = &models.Delivery{}
	resp.Customer = &models.UserBasic{}
	resp.Status = &models.DeliveryStatus{}

	// Temporary variables for nullable fields
	var (
		driverID            sql.NullInt64
		licenseNumber       sql.NullString
		yearsExperience     sql.NullInt64
		vehicleID           sql.NullInt64
		plateNumber         sql.NullString
		vehicleType         sql.NullString
		statusName          sql.NullString
		statusDescription   sql.NullString
		pickupWarehouseID   sql.NullInt64
		deliveryWarehouseID sql.NullInt64
		pickupLatitude      sql.NullFloat64
		pickupLongitude     sql.NullFloat64
		deliveryLatitude    sql.NullFloat64
		deliveryLongitude   sql.NullFloat64
	)

	err := r.db.QueryRow(query, id).Scan(
		&resp.Delivery.ID,
		&resp.Delivery.TrackingNumber,
		&resp.Delivery.CustomerID,
		&resp.Delivery.CargoType,
		&resp.Delivery.CargoWeight,
		&resp.Delivery.SpecialInstructions,
		&pickupWarehouseID,
		&deliveryWarehouseID,
		&pickupLatitude,
		&pickupLongitude,
		&deliveryLatitude,
		&deliveryLongitude,
		&resp.Delivery.EstimatedDeliveryTime,
		&resp.Delivery.StatusID,
		&resp.Delivery.PaymentStatus,
		&resp.Delivery.CreatedAt,
		&resp.Delivery.UpdatedAt,
		&resp.Customer.ID,
		&resp.Customer.Username,
		&resp.Customer.Email,
		&driverID,
		&licenseNumber,
		&yearsExperience,
		&vehicleID,
		&plateNumber,
		&vehicleType,
		&statusName,
		&statusDescription,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting delivery: %v", err)
	}

	// Set driver info if exists
	if driverID.Valid {
		resp.Driver = &models.DriverProfile{
			ID:              int(driverID.Int64),
			LicenseNumber:   licenseNumber.String,
			YearsExperience: fmt.Sprintf("%d", yearsExperience.Int64),
		}
	}

	// Set vehicle info if exists
	if vehicleID.Valid {
		resp.Vehicle = &models.Vehicle{
			ID:          int64(vehicleID.Int64),
			PlateNumber: plateNumber.String,
			Type:        vehicleType.String,
		}
	}

	// Set status info if exists
	if statusName.Valid {
		resp.Status.Name = statusName.String
		resp.Status.Description = statusDescription.String
	}

	// Set warehouse IDs if exists
	if pickupWarehouseID.Valid {
		id := pickupWarehouseID.Int64
		resp.Delivery.PickupWarehouseID = &id
	}
	if deliveryWarehouseID.Valid {
		id := deliveryWarehouseID.Int64
		resp.Delivery.DeliveryWarehouseID = &id
	}

	// Set latitude and longitude if exists
	if pickupLatitude.Valid {
		resp.Delivery.PickupLatitude = pickupLatitude.Float64
	}
	if pickupLongitude.Valid {
		resp.Delivery.PickupLongitude = pickupLongitude.Float64
	}
	if deliveryLatitude.Valid {
		resp.Delivery.DeliveryLatitude = deliveryLatitude.Float64
	}
	if deliveryLongitude.Valid {
		resp.Delivery.DeliveryLongitude = deliveryLongitude.Float64
	}

	return &resp, nil
}

// Update updates a delivery record
func (r *DeliveryRepository) Update(id int64, req *models.UpdateDeliveryRequest) (*models.Delivery, error) {
	setValues := []string{}
	args := []interface{}{id}
	argCount := 2

	if req.StatusID != nil {
		setValues = append(setValues, fmt.Sprintf("status_id = $%d", argCount))
		args = append(args, *req.StatusID)
		argCount++
	}

	if req.DriverID != nil {
		setValues = append(setValues, fmt.Sprintf("driver_id = $%d", argCount))
		args = append(args, *req.DriverID)
		argCount++
	}

	if req.VehicleID != nil {
		setValues = append(setValues, fmt.Sprintf("vehicle_id = $%d", argCount))
		args = append(args, *req.VehicleID)
		argCount++
	}

	if req.PickupTime != nil {
		setValues = append(setValues, fmt.Sprintf("pickup_time = $%d", argCount))
		args = append(args, *req.PickupTime)
		argCount++
	}

	if req.DeliveryTime != nil {
		setValues = append(setValues, fmt.Sprintf("delivery_time = $%d", argCount))
		args = append(args, *req.DeliveryTime)
		argCount++
	}

	if req.ActualDeliveryTime != nil {
		setValues = append(setValues, fmt.Sprintf("actual_delivery_time = $%d", argCount))
		args = append(args, *req.ActualDeliveryTime)
		argCount++
	}

	if req.CustomerSignature != nil {
		setValues = append(setValues, fmt.Sprintf("customer_signature = $%d", argCount))
		args = append(args, *req.CustomerSignature)
		argCount++
	}

	if req.PaymentStatus != nil {
		setValues = append(setValues, fmt.Sprintf("payment_status = $%d", argCount))
		args = append(args, *req.PaymentStatus)
		argCount++
	}

	if req.ActFuelConsumption != nil {
		setValues = append(setValues, fmt.Sprintf("actual_fuel_consumption = $%d", argCount))
		args = append(args, *req.ActFuelConsumption)
		argCount++
	}

	if req.RouteEfficiencyScore != nil {
		setValues = append(setValues, fmt.Sprintf("route_efficiency_score = $%d", argCount))
		args = append(args, *req.RouteEfficiencyScore)
		argCount++
	}

	if len(setValues) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setValues = append(setValues, "updated_at = CURRENT_TIMESTAMP")
	query := fmt.Sprintf(`
		UPDATE deliveries
		SET %s
		WHERE id = $1
		RETURNING id, tracking_number, customer_id, cargo_type, cargo_weight,
			special_instructions, pickup_warehouse_id, delivery_warehouse_id,
			pickup_latitude, pickup_longitude, delivery_latitude, delivery_longitude,
			estimated_delivery_time, status_id, payment_status,
			created_at, updated_at`,
		strings.Join(setValues, ", "))

	var result models.Delivery
	err := r.db.QueryRow(query, args...).Scan(
		&result.ID,
		&result.TrackingNumber,
		&result.CustomerID,
		&result.CargoType,
		&result.CargoWeight,
		&result.SpecialInstructions,
		&result.PickupWarehouseID,
		&result.DeliveryWarehouseID,
		&result.PickupLatitude,
		&result.PickupLongitude,
		&result.DeliveryLatitude,
		&result.DeliveryLongitude,
		&result.EstimatedDeliveryTime,
		&result.StatusID,
		&result.PaymentStatus,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error updating delivery: %v", err)
	}

	return &result, nil
}

// List retrieves deliveries based on filter criteria
func (r *DeliveryRepository) List(filter *models.DeliveryFilter) (*models.DeliveryListResponse, error) {
	where := []string{"1 = 1"}
	args := []interface{}{}
	argCount := 1

	if filter.CustomerID != nil {
		where = append(where, fmt.Sprintf("d.customer_id = $%d", argCount))
		args = append(args, *filter.CustomerID)
		argCount++
	}

	if filter.UserID != nil {
		where = append(where, fmt.Sprintf("d.user_id = $%d", argCount))
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.StatusID != nil {
		where = append(where, fmt.Sprintf("d.status_id = $%d", argCount))
		args = append(args, *filter.StatusID)
		argCount++
	}

	if filter.DateFrom != nil {
		where = append(where, fmt.Sprintf("d.created_at >= $%d", argCount))
		args = append(args, *filter.DateFrom)
		argCount++
	}

	if filter.DateTo != nil {
		where = append(where, fmt.Sprintf("d.created_at <= $%d", argCount))
		args = append(args, *filter.DateTo)
		argCount++
	}

	// Remove the DriverID check since it no longer exists in DeliveryFilter
	// if filter.DriverID != nil {
	// 	where = append(where, fmt.Sprintf("d.driver_id = $%d", argCount))
	// 	args = append(args, *filter.DriverID)
	// 	argCount++
	// }

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM deliveries d WHERE %s`,
		strings.Join(where, " AND "))

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting deliveries: %v", err)
	}

	// Add pagination
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	// Add sorting
	orderBy := "d.created_at DESC"
	if filter.SortBy != "" {
		if filter.SortOrder != "" && (filter.SortOrder == "ASC" || filter.SortOrder == "DESC") {
			orderBy = fmt.Sprintf("d.%s %s", filter.SortBy, filter.SortOrder)
		}
	}

	query := fmt.Sprintf(`
		SELECT 
			d.id, d.tracking_number, d.customer_id, d.cargo_type, d.cargo_weight,
			d.special_instructions, d.pickup_warehouse_id, d.delivery_warehouse_id,
			d.pickup_latitude, d.pickup_longitude, d.delivery_latitude, d.delivery_longitude,
			d.estimated_delivery_time, d.status_id, d.payment_status, d.created_at, d.updated_at,
			u.id, u.username, u.email,
			dp.id, dp.license_number, dp.years_experience,
			v.id, v.plate_number, v.type,
			ds.name, ds.description
		FROM deliveries d
		LEFT JOIN users u ON d.customer_id = u.id
		LEFT JOIN driver_profiles dp ON d.driver_id = dp.id
		LEFT JOIN vehicles v ON d.vehicle_id = v.id
		LEFT JOIN delivery_statuses ds ON d.status_id = ds.id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "),
		orderBy,
		argCount,
		argCount+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying deliveries: %v", err)
	}
	defer rows.Close()

	var deliveries []models.DeliveryResponse
	for rows.Next() {
		var resp models.DeliveryResponse
		resp.Delivery = &models.Delivery{}
		resp.Customer = &models.UserBasic{}
		resp.Status = &models.DeliveryStatus{}
		resp.Driver = &models.DriverProfile{}
		resp.Vehicle = &models.Vehicle{}

		// Temporary variables for nullable fields
		var (
			cargoWeight         sql.NullFloat64
			cargoType           sql.NullString
			specialInstructions sql.NullString
			pickupWarehouseID   sql.NullInt64
			deliveryWarehouseID sql.NullInt64
			pickupLatitude      sql.NullFloat64
			pickupLongitude     sql.NullFloat64
			deliveryLatitude    sql.NullFloat64
			deliveryLongitude   sql.NullFloat64
			paymentStatus       sql.NullString
			driverID            sql.NullInt64
			licenseNumber       sql.NullString
			yearsExperience     sql.NullInt64
			vehicleID           sql.NullInt64
			plateNumber         sql.NullString
			vehicleType         sql.NullString
			statusName          sql.NullString
			statusDescription   sql.NullString
		)

		err := rows.Scan(
			&resp.Delivery.ID,
			&resp.Delivery.TrackingNumber,
			&resp.Delivery.CustomerID,
			&cargoType,
			&cargoWeight,
			&specialInstructions,
			&pickupWarehouseID,
			&deliveryWarehouseID,
			&pickupLatitude,
			&pickupLongitude,
			&deliveryLatitude,
			&deliveryLongitude,
			&resp.Delivery.EstimatedDeliveryTime,
			&resp.Delivery.StatusID,
			&paymentStatus,
			&resp.Delivery.CreatedAt,
			&resp.Delivery.UpdatedAt,
			&resp.Customer.ID,
			&resp.Customer.Username,
			&resp.Customer.Email,
			&driverID,
			&licenseNumber,
			&yearsExperience,
			&vehicleID,
			&plateNumber,
			&vehicleType,
			&statusName,
			&statusDescription,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning delivery row: %v", err)
		}

		// Set nullable fields if they exist
		if cargoType.Valid {
			resp.Delivery.CargoType = cargoType.String
		}
		if cargoWeight.Valid {
			resp.Delivery.CargoWeight = cargoWeight.Float64
		}
		if specialInstructions.Valid {
			resp.Delivery.SpecialInstructions = specialInstructions.String
		}
		if pickupWarehouseID.Valid {
			id := pickupWarehouseID.Int64
			resp.Delivery.PickupWarehouseID = &id
		}
		if deliveryWarehouseID.Valid {
			id := deliveryWarehouseID.Int64
			resp.Delivery.DeliveryWarehouseID = &id
		}
		if pickupLatitude.Valid {
			resp.Delivery.PickupLatitude = pickupLatitude.Float64
		}
		if pickupLongitude.Valid {
			resp.Delivery.PickupLongitude = pickupLongitude.Float64
		}
		if deliveryLatitude.Valid {
			resp.Delivery.DeliveryLatitude = deliveryLatitude.Float64
		}
		if deliveryLongitude.Valid {
			resp.Delivery.DeliveryLongitude = deliveryLongitude.Float64
		}
		if paymentStatus.Valid {
			resp.Delivery.PaymentStatus = paymentStatus.String
		}

		// Set driver info if exists
		if driverID.Valid {
			resp.Driver = &models.DriverProfile{
				ID:              int(driverID.Int64),
				LicenseNumber:   licenseNumber.String,
				YearsExperience: fmt.Sprintf("%d", yearsExperience.Int64),
			}
		} else {
			resp.Driver = nil
		}

		// Set vehicle info if exists
		if vehicleID.Valid {
			resp.Vehicle = &models.Vehicle{
				ID:          int64(vehicleID.Int64),
				PlateNumber: plateNumber.String,
				Type:        vehicleType.String,
			}
		} else {
			resp.Vehicle = nil
		}

		// Set status info if exists
		if statusName.Valid {
			resp.Status.Name = statusName.String
			resp.Status.Description = statusDescription.String
		}

		deliveries = append(deliveries, resp)
	}

	return &models.DeliveryListResponse{
		Deliveries: deliveries,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

// GetDeliveryStatistics retrieves comprehensive delivery statistics
func (r *DeliveryRepository) GetDeliveryStatistics() (*models.DeliveryStatistics, error) {
	query := `
		WITH status_counts AS (
			SELECT 
				status_id, 
				ds.name as status_name, 
				COUNT(*) as count 
			FROM deliveries d
			JOIN delivery_statuses ds ON d.status_id = ds.id
			GROUP BY status_id, ds.name
		)
		SELECT 
			(SELECT COUNT(*) FROM deliveries) as total_deliveries,
			(SELECT COUNT(*) FROM deliveries WHERE created_at >= NOW() - INTERVAL '30 days') as deliveries_last_30_days,
			(SELECT AVG(EXTRACT(EPOCH FROM (updated_at - created_at)) / 3600) FROM deliveries WHERE status_id = 3) as avg_delivery_time_hours,
			(SELECT SUM(cargo_weight) FROM deliveries WHERE status_id = 3) as total_cargo_delivered,
			(SELECT json_agg(json_build_object('status', status_name, 'count', count)) FROM status_counts) as status_distribution
	`

	var stats models.DeliveryStatistics
	var statusDistributionJSON []byte

	err := r.db.QueryRow(query).Scan(
		&stats.TotalDeliveries,
		&stats.DeliveriesLast30Days,
		&stats.AverageDeliveryTimeHours,
		&stats.TotalCargoDelivered,
		&statusDistributionJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("error fetching delivery statistics: %v", err)
	}

	// Unmarshal status distribution
	err = json.Unmarshal(statusDistributionJSON, &stats.StatusDistribution)
	if err != nil {
		return nil, fmt.Errorf("error parsing status distribution: %v", err)
	}

	return &stats, nil
}

func generateTrackingNumber() string {
	return fmt.Sprintf("TRK%d", time.Now().UnixNano())
}
