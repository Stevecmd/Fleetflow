package constants

import (
	"os"
	"time"
)

// Database connection string
var DatabaseURL = os.Getenv("DATABASE_URL")

// Context keys
type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	UsernameKey ContextKey = "username"
	RoleKey     ContextKey = "role"
)

// Token settings
const (
	AccessTokenDuration  = 24 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// Role names
const (
	AdminRole        = "admin"
	CustomerRole     = "customer"
	DriverRole       = "driver"
	FleetManagerRole = "fleet_manager"
	LoaderRole       = "loader"
)
