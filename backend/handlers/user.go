package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/stevecmd/Fleetflow/backend/models"
	"github.com/stevecmd/Fleetflow/backend/pkg/constants"
)

// type contextKey string

// const (
// 	userIDKey contextKey = "user_id"
// 	roleKey   contextKey = "role"
// )

// GetUserProfile handles GET /users/:id requests
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	authUserID, ok := r.Context().Value(constants.UserIDKey).(int)
	if !ok {
		log.Printf("Failed to get userID from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract user ID from URL path or use authenticated user's ID
	var userID int
	parts := strings.Split(r.URL.Path, "/")
	if parts[len(parts)-1] == "profile" || parts[len(parts)-1] == "" {
		userID = authUserID
	} else {
		var err error
		userID, err = strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
	}

	// Only allow users to view their own profile unless they are admins
	userRole, ok := r.Context().Value(constants.RoleKey).(string)
	if !ok {
		log.Printf("Failed to get role from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if userRole != constants.AdminRole && authUserID != userID {
		log.Printf("Access denied: User %d (role: %s) attempted to access profile %d", authUserID, userRole, userID)
		http.Error(w, "Forbidden - Cannot access other user's profile", http.StatusForbidden)
		return
	}

	// Connect to database
	log.Printf("Connecting to database with URL: %s", DatabaseURL)
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query user information
	log.Printf("Querying user information for ID: %d", userID)
	var user models.User
	err = db.QueryRow(`
        SELECT 
            u.id, u.username, u.email, u.role_id, r.name as role_name,
            u.first_name, u.last_name, u.phone, u.profile_image_url,
            u.date_of_birth, u.gender, u.nationality, u.preferred_language,
            u.status, u.last_login_at, u.created_at, u.updated_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.id = $1`,
		userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.RoleID, &user.RoleName,
		&user.FirstName, &user.LastName, &user.Phone, &user.ProfileImageURL,
		&user.DateOfBirth, &user.Gender, &user.Nationality, &user.PreferredLanguage,
		&user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Database error querying user information: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Query addresses
	log.Printf("Querying addresses for user ID: %d", userID)
	rows, err := db.Query(`
		SELECT id, street1, street2, city, state, zip, country,
			   address_type, is_default, latitude, longitude,
			   created_at, updated_at
		FROM addresses
		WHERE user_id = $1`,
		userID)
	if err != nil {
		log.Printf("Database error querying addresses: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	user.Addresses = []models.Address{}
	for rows.Next() {
		var addr models.Address
		err = rows.Scan(
			&addr.ID, &addr.Street1, &addr.Street2, &addr.City, &addr.State,
			&addr.Zip, &addr.Country, &addr.Type, &addr.IsDefault,
			&addr.Latitude, &addr.Longitude, &addr.CreatedAt, &addr.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning address: %v", err)
			continue
		}
		user.Addresses = append(user.Addresses, addr)
	}

	// Query emergency contacts
	log.Printf("Querying emergency contacts for user ID: %d", userID)
	rows, err = db.Query(`
		SELECT id, name, relationship, phone, email,
			   created_at, updated_at
		FROM emergency_contacts
		WHERE user_id = $1`,
		userID)
	if err != nil {
		log.Printf("Database error querying emergency contacts: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	user.EmergencyContacts = []models.EmergencyContact{}
	for rows.Next() {
		var contact models.EmergencyContact
		err = rows.Scan(
			&contact.ID, &contact.Name, &contact.Relationship,
			&contact.Phone, &contact.Email, &contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning emergency contact: %v", err)
			continue
		}
		user.EmergencyContacts = append(user.EmergencyContacts, contact)
	}

	// Hide sensitive information
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// ListUsers handles GET /api/v1/users requests to list all users
func ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user's role from context
	userRole, ok := r.Context().Value(constants.RoleKey).(string)
	if !ok {
		log.Printf("Failed to get role from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only allow admins to list all users
	if userRole != constants.AdminRole {
		log.Printf("Access denied: Non-admin user attempted to list all users")
		http.Error(w, "Forbidden - Admin access required", http.StatusForbidden)
		return
	}

	// Connect to database
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query all users
	rows, err := db.Query(`
        SELECT 
            u.id, u.username, u.email, u.role_id, r.name as role_name,
            u.first_name, u.last_name, u.phone, u.profile_image_url,
            u.date_of_birth, u.gender, u.nationality, u.preferred_language,
            u.status, u.last_login_at, u.created_at, u.updated_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        ORDER BY u.id`)
	if err != nil {
		log.Printf("Database error querying users: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.RoleID, &user.RoleName,
			&user.FirstName, &user.LastName, &user.Phone, &user.ProfileImageURL,
			&user.DateOfBirth, &user.Gender, &user.Nationality, &user.PreferredLanguage,
			&user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}
		// Hide sensitive information
		user.Password = ""
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// UpdateUserProfile handles PUT /users/:id requests
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("Requested user ID: %d", userID)

	// Get authenticated user ID from context
	authUserID, ok := r.Context().Value(constants.UserIDKey).(int)
	if !ok {
		log.Printf("Failed to get userID from context. Context values: %+v", r.Context())
		http.Error(w, "Unauthorized - No user ID in context", http.StatusUnauthorized)
		return
	}
	log.Printf("Authenticated user ID: %d", authUserID)

	// Only allow users to update their own profile unless they are admins
	userRole, ok := r.Context().Value(constants.RoleKey).(string)
	if !ok {
		log.Printf("Failed to get role from context. Context values: %+v", r.Context())
		http.Error(w, "Unauthorized - No role in context", http.StatusUnauthorized)
		return
	}
	log.Printf("User role: %s", userRole)

	if userRole != constants.AdminRole && authUserID != userID {
		log.Printf("Access denied: User %d (role: %s) attempted to access profile %d", authUserID, userRole, userID)
		http.Error(w, "Forbidden - Cannot access other user's profile", http.StatusForbidden)
		return
	}

	// Parse request body
	var updates struct {
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Phone           string `json:"phone"`
		ProfileImageURL string `json:"profile_image_url"`
	}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Connect to database
	log.Printf("Connecting to database with URL: %s", DatabaseURL)
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Update user information
	log.Printf("Updating user information for ID: %d", userID)
	_, err = db.Exec(`
		UPDATE users 
		SET first_name = $1, last_name = $2, phone = $3, profile_image_url = $4, updated_at = NOW()
		WHERE id = $5`,
		updates.FirstName, updates.LastName, updates.Phone, updates.ProfileImageURL, userID)

	if err != nil {
		log.Printf("Error updating user profile: %v", err)
		http.Error(w, "Error updating user profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
