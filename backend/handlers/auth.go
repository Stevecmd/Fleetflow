package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stevecmd/Fleetflow/backend/models"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey []byte
var DatabaseURL string
var tokenBlacklist *models.TokenBlacklist
var refreshTokens = make(map[string]string)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}
	if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}

func generateTokenPair(userID int, username string, roleName string) (*TokenPair, error) {
	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        userID,
		"username":  username,
		"role_name": roleName,
		"exp":       time.Now().Add(time.Minute * 15).Unix(), // Access token expires in 15 minutes
	})

	accessTokenString, err := accessToken.SignedString(JwtKey)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	})

	refreshTokenString, err := refreshToken.SignedString(JwtKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	refreshTokens[refreshTokenString] = username

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// LoginHandler godoc
// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "User credentials"
// @Success 200 {object} models.TokenPair
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", DatabaseURL) // Use DatabaseURL from handlers package
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user models.User
	var roleName string
	err = db.QueryRow(`
		SELECT u.id, u.password, u.role_id, r.name as role_name, u.username 
		FROM users u 
		JOIN roles r ON u.role_id = r.id 
		WHERE u.username = $1`,
		creds.Username).Scan(&user.ID, &user.Password, &user.RoleID, &roleName, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenPair, err := generateTokenPair(user.ID, user.Username, roleName)
	if err != nil {
		http.Error(w, "Token generation error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Successfully logged in",
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"user_id":       user.ID,
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newUser.Username == "" || newUser.Password == "" || newUser.Email == "" ||
		newUser.FirstName == "" || newUser.LastName == "" || newUser.Phone == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate password
	err = validatePassword(newUser.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	// Set timestamps
	now := time.Now()
	newUser.CreatedAt = now
	newUser.UpdatedAt = now

	// Open database connection
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Transaction start error: %v", err)
		http.Error(w, "Database transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // This will be a no-op if the transaction has been committed

	// Get role name first
	var roleName string
	err = tx.QueryRow("SELECT name FROM roles WHERE id = $1", newUser.RoleID).Scan(&roleName)
	if err != nil {
		log.Printf("Role query error: %v", err)
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var userID int
	query := `
		INSERT INTO users (
			username, password, email, role_id, first_name, last_name, phone,
			date_of_birth, gender, nationality, preferred_language, profile_image_url,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 
			NULLIF($8, '')::date, NULLIF($9, ''), NULLIF($10, ''), NULLIF($11, ''), NULLIF($12, ''),
			$13, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id`

	err = tx.QueryRow(
		query,
		newUser.Username, newUser.Password, newUser.Email, newUser.RoleID,
		newUser.FirstName, newUser.LastName, newUser.Phone,
		newUser.DateOfBirth, newUser.Gender, newUser.Nationality,
		newUser.PreferredLanguage, newUser.ProfileImageURL,
		"active").Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		log.Printf("User insertion error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Insert addresses
	for _, addr := range newUser.Addresses {
		_, err = tx.Exec(`
			INSERT INTO addresses (
				user_id, street1, street2, city, state, zip, country,
				address_type, is_default, latitude, longitude,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
			userID, addr.Street1, addr.Street2, addr.City, addr.State, addr.Zip, addr.Country,
			addr.Type, addr.IsDefault, addr.Latitude, addr.Longitude)
		if err != nil {
			log.Printf("Address insertion error: %v", err)
			http.Error(w, "Error saving address information", http.StatusInternalServerError)
			return
		}
	}

	// Insert emergency contacts
	for _, contact := range newUser.EmergencyContacts {
		_, err = tx.Exec(`
			INSERT INTO emergency_contacts (
				user_id, name, relationship, phone, email,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
			userID, contact.Name, contact.Relationship, contact.Phone, contact.Email)
		if err != nil {
			log.Printf("Emergency contact insertion error: %v", err)
			http.Error(w, "Error saving emergency contact information", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Transaction commit error: %v", err)
		http.Error(w, "Error completing registration", http.StatusInternalServerError)
		return
	}

	// Generate tokens
	tokenPair, err := generateTokenPair(userID, newUser.Username, roleName)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		http.Error(w, "Token generation error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokenPair)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	// Add token to blacklist
	tokenBlacklist.Add(tokenString)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse and validate refresh token
	token, err := jwt.Parse(refreshRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtKey, nil
	})

	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Get user information
	userID := int(claims["id"].(float64))
	username := claims["username"].(string)

	// Get user's role from database
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var roleName string
	err = db.QueryRow(`
		SELECT r.name 
		FROM users u 
		JOIN roles r ON u.role_id = r.id 
		WHERE u.id = $1`,
		userID).Scan(&roleName)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Error retrieving user role", http.StatusInternalServerError)
		return
	}

	// Generate new token pair
	tokenPair, err := generateTokenPair(userID, username, roleName)
	if err != nil {
		http.Error(w, "Token generation error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}

func init() {
	tokenBlacklist = models.NewTokenBlacklist()
}
