package models

type Invoice struct {
	ID int `json:"id"`

	DeliveryID int `json:"delivery_id"`

	Amount float64 `json:"amount"`

	Status string `json:"status"`

	DueDate string `json:"due_date"`

	PaymentMethod string `json:"payment_method"`
}
