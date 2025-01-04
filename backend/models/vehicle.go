package models

// Vehicle represents a vehicle in the system

import (
	"time"
)

type Vehicle struct {
	ID                     int64      `json:"id" db:"id"`
	PlateNumber            string     `json:"plate_number" db:"plate_number"`
	Type                   string     `json:"type" db:"type"`
	Make                   string     `json:"make" db:"make"`
	Model                  string     `json:"model" db:"model"`
	Year                   int        `json:"year" db:"year"`
	Capacity               float64    `json:"capacity" db:"capacity"`
	FuelType               string     `json:"fuel_type" db:"fuel_type"`
	StatusID               int        `json:"status_id" db:"status_id"`
	GPSUnitID              *string    `json:"gps_unit_id,omitempty" db:"gps_unit_id"`
	LastMaintenance        *time.Time `json:"last_maintenance,omitempty" db:"last_maintenance"`
	NextMaintenance        *time.Time `json:"next_maintenance,omitempty" db:"next_maintenance"`
	Mileage                *float64   `json:"mileage,omitempty" db:"mileage"`
	InsuranceExpiry        *time.Time `json:"insurance_expiry,omitempty" db:"insurance_expiry"`
	CurrentLocationLat     *float64   `json:"current_location_latitude,omitempty" db:"current_location_latitude"`
	CurrentLocationLong    *float64   `json:"current_location_longitude,omitempty" db:"current_location_longitude"`
	CurrentLocationUpdated *time.Time `json:"current_location_updated_at,omitempty" db:"current_location_updated_at"`
	FuelEfficiencyRating   *float64   `json:"fuel_efficiency_rating,omitempty" db:"fuel_efficiency_rating"`
	TotalFuelConsumption   *float64   `json:"total_fuel_consumption,omitempty" db:"total_fuel_consumption"`
	TotalMaintenanceCost   *float64   `json:"total_maintenance_cost,omitempty" db:"total_maintenance_cost"`
	VehicleImages          *string    `json:"vehicle_images,omitempty" db:"vehicle_images"`
	RegistrationDocURL     *string    `json:"registration_document_image_url,omitempty" db:"registration_document_image_url"`
	InsuranceDocURL        *string    `json:"insurance_document_image_url,omitempty" db:"insurance_document_image_url"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
}
