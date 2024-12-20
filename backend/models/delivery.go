package models

import (
	"time"
)

// Delivery represents a delivery record in the database
type Delivery struct {
	ID                      int64      `json:"id"`
	FromLocation            string     `json:"from_location"`
	ToLocation              string     `json:"to_location"`
	ProofOfDeliveryImageURL string     `json:"proof_of_delivery_image_url"`
	Weight                  float64    `json:"weight"`
	Status                  string     `json:"status"`
	TrackingNumber          string     `json:"tracking_number"`
	UserID                  *int64     `json:"user_id,omitempty"`
	VehicleID               *int64     `json:"vehicle_id,omitempty"`
	StatusID                int64      `json:"status_id"`
	PickupTime              *time.Time `json:"pickup_time,omitempty"`
	DeliveryTime            *time.Time `json:"delivery_time,omitempty"`
	EstimatedDeliveryTime   *time.Time `json:"estimated_delivery_time"`
	ActualDeliveryTime      *time.Time `json:"actual_delivery_time,omitempty"`
	CargoType               string     `json:"cargo_type"`
	CargoWeight             float64    `json:"cargo_weight"`
	SpecialInstructions     string     `json:"special_instructions"`
	CustomerID              int64      `json:"customer_id"`
	CustomerSignature       string     `json:"customer_signature"`
	PickupWarehouseID       *int64     `json:"pickup_warehouse_id,omitempty"`
	DeliveryWarehouseID     *int64     `json:"delivery_warehouse_id,omitempty"`
	PickupLatitude          float64    `json:"pickup_latitude"`
	PickupLongitude         float64    `json:"pickup_longitude"`
	DeliveryLatitude        float64    `json:"delivery_latitude"`
	DeliveryLongitude       float64    `json:"delivery_longitude"`
	PaymentStatus           string     `json:"payment_status"`
	EstFuelConsumption      float64    `json:"estimated_fuel_consumption"`
	ActFuelConsumption      float64    `json:"actual_fuel_consumption"`
	RouteEfficiencyScore    float64    `json:"route_efficiency_score"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// CreateDeliveryRequest represents the request body for creating a new delivery
type CreateDeliveryRequest struct {
	CustomerID            int64     `json:"customer_id" validate:"required"`
	CargoType             string    `json:"cargo_type" validate:"required"`
	CargoWeight           float64   `json:"cargo_weight" validate:"required,gt=0"`
	SpecialInstructions   string    `json:"special_instructions"`
	PickupWarehouseID     *int64    `json:"pickup_warehouse_id"`
	DeliveryWarehouseID   *int64    `json:"delivery_warehouse_id"`
	PickupLatitude        float64   `json:"pickup_latitude" validate:"required"`
	PickupLongitude       float64   `json:"pickup_longitude" validate:"required"`
	DeliveryLatitude      float64   `json:"delivery_latitude" validate:"required"`
	DeliveryLongitude     float64   `json:"delivery_longitude" validate:"required"`
	EstimatedDeliveryTime time.Time `json:"estimated_delivery_time" validate:"required"`
}

// UpdateDeliveryRequest represents the request body for updating a delivery
type UpdateDeliveryRequest struct {
	StatusID             *int64     `json:"status_id"`
	DriverID             *int64     `json:"driver_id"`
	VehicleID            *int64     `json:"vehicle_id"`
	PickupTime           *time.Time `json:"pickup_time"`
	DeliveryTime         *time.Time `json:"delivery_time"`
	ActualDeliveryTime   *time.Time `json:"actual_delivery_time"`
	CustomerSignature    *string    `json:"customer_signature"`
	PaymentStatus        *string    `json:"payment_status"`
	ActFuelConsumption   *float64   `json:"actual_fuel_consumption"`
	RouteEfficiencyScore *float64   `json:"route_efficiency_score"`
}

// UserBasic represents basic user information
type UserBasic struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// DeliveryStatus represents the status of a delivery
type DeliveryStatus struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DeliveryResponse represents the response body for delivery endpoints
type DeliveryResponse struct {
	*Delivery
	Customer *UserBasic      `json:"customer,omitempty"`
	Driver   *DriverProfile  `json:"driver,omitempty"`
	Vehicle  *Vehicle        `json:"vehicle,omitempty"`
	Status   *DeliveryStatus `json:"status,omitempty"`
}

// DeliveryListResponse represents the paginated response for listing deliveries
type DeliveryListResponse struct {
	Deliveries []DeliveryResponse `json:"deliveries"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
}

// DeliveryFilter represents the filter options for listing deliveries
type DeliveryFilter struct {
	CustomerID *int64 `json:"customer_id"`
	// DriverID   *int64     `json:"driver_id"`
	StatusID  *int64     `json:"status_id"`
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
	SortBy    string     `json:"sort_by"`
	SortOrder string     `json:"sort_order"`
	UserID    *int64     `json:"user_id"`
}

type DeliveryStatistics struct {
	TotalDeliveries          int64                 `json:"total_deliveries"`
	DeliveriesLast30Days     int64                 `json:"deliveries_last_30_days"`
	AverageDeliveryTimeHours float64               `json:"average_delivery_time_hours"`
	TotalCargoDelivered      float64               `json:"total_cargo_delivered"`
	StatusDistribution       []DeliveryStatusCount `json:"status_distribution"`
}

type DeliveryStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}
