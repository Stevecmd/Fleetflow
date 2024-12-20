package models

// Other model definitions...

type DeliveryFeedback struct {
	ID int `json:"id"`

	DeliveryID int `json:"delivery_id"`

	Rating int `json:"rating"`

	FeedbackText string `json:"feedback_text"`

	TimelinessRating int `json:"timeliness_rating"`

	DriverRating int `json:"driver_rating"`

	PackageConditionRating int `json:"package_condition_rating"`
}
