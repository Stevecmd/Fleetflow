package models

// Credentials represents the request body for authenticating a user

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
