package models

import (
	"time"
)

// DriverProfile represents a driver profile in the system
type DriverProfile struct {
	ID               int        `json:"id"`
	UserID           int        `json:"user_id"`
    FirstName        string     `json:"first_name"`     
    LastName         string     `json:"last_name"`     
    Email            string     `json:"email"`          
    AssignedVehicle  string     `json:"assigned_vehicle"`
	LicenseNumber    string     `json:"license_number"`
	LicenseType      string     `json:"license_type"`
	LicenseExpiry    *time.Time `json:"license_expiry,omitempty"`
	VehicleType      string     `json:"vehicle_type"`
	YearsExperience  string     `json:"years_experience"`
	Certification    []string   `json:"certification"`
	Status           string     `json:"status"`
	StatusID         *int       `json:"status_id,omitempty"`
	CurrentVehicleID *int       `json:"current_vehicle_id,omitempty"`
	Rating           *float64   `json:"rating,omitempty"`
	TotalTrips       *int       `json:"total_trips,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
