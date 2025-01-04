package models

// User represents a user in the system

import (
	"time"
)

type Address struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	Street1   string    `json:"street1"`
	Street2   string    `json:"street2,omitempty"`
	City      string    `json:"city"`
	State     string    `json:"state,omitempty"`
	Zip       string    `json:"zip,omitempty"`
	Country   string    `json:"country"`
	Type      string    `json:"address_type,omitempty"`
	IsDefault bool      `json:"is_default,omitempty"`
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type EmergencyContact struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type User struct {
	ID                int                `json:"id"`
	Username          string             `json:"username"`
	Password          string             `json:"password,omitempty"`
	Email             string             `json:"email"`
	RoleID            int                `json:"role_id"`
	RoleName          string             `json:"role_name"`
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	Phone             string             `json:"phone"`
	DateOfBirth       string             `json:"date_of_birth,omitempty"`
	Gender            string             `json:"gender,omitempty"`
	Nationality       string             `json:"nationality,omitempty"`
	PreferredLanguage string             `json:"preferred_language,omitempty"`
	ProfileImageURL   string             `json:"profile_image_url,omitempty"`
	Addresses         []Address          `json:"addresses,omitempty"`
	Status            string             `json:"status,omitempty"`
	EmergencyContacts []EmergencyContact `json:"emergency_contacts,omitempty"`
	LastLoginAt       *time.Time         `json:"last_login_at,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}
