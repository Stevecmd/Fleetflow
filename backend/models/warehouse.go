package models

// Warehouse represents a warehouse in the system

import (
	"encoding/json"
	"time"
)

type Warehouse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Street1   string    `json:"street1"`
	Street2   *string   `json:"street2,omitempty"`
	City      string    `json:"city"`
	State     *string   `json:"state,omitempty"`
	Zip       string    `json:"zip"`
	Country   string    `json:"country"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Capacity  int       `json:"capacity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (w Warehouse) MarshalJSON() ([]byte, error) {
	type Alias Warehouse
	return json.Marshal(struct {
		Street2 *string `json:"street2,omitempty"`
		State   *string `json:"state,omitempty"`
		*Alias
	}{
		Street2: w.Street2,
		State:   w.State,
		Alias:   (*Alias)(&w),
	})
}

type CreateWarehouseRequest struct {
	Name      string  `json:"name" validate:"required"`
	Street1   string  `json:"street1" validate:"required"`
	Street2   string  `json:"street2,omitempty"`
	City      string  `json:"city" validate:"required"`
	State     string  `json:"state,omitempty"`
	Zip       string  `json:"zip" validate:"required"`
	Country   string  `json:"country" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Capacity  int     `json:"capacity" validate:"required,min=0"`
}

type UpdateWarehouseRequest struct {
	Name      string  `json:"name,omitempty"`
	Street1   string  `json:"street1,omitempty"`
	Street2   string  `json:"street2,omitempty"`
	City      string  `json:"city,omitempty"`
	State     string  `json:"state,omitempty"`
	Zip       string  `json:"zip,omitempty"`
	Country   string  `json:"country,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Capacity  int     `json:"capacity,omitempty" validate:"omitempty,min=0"`
	Status    string  `json:"status,omitempty"`
}
