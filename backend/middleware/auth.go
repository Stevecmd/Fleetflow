package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stevecmd/Fleetflow/backend/pkg/constants"
)

var tokenBlacklist = make(map[string]bool)
var JwtKey []byte

// AuthMiddleware is an HTTP middleware that validates an authorization token.
// If the token is valid and not revoked, it adds the user's ID, username, and role name to the request context.
// The middleware returns http.StatusUnauthorized if the token is invalid, expired, or revoked.
// The middleware returns http.StatusForbidden if the user is not authorized to access the endpoint.
// The middleware logs token parsing and validation errors.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			http.Error(w, "Invalid authorization token format", http.StatusUnauthorized)
			return
		}

		if tokenBlacklist[tokenString] {
			http.Error(w, "Token revoked", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtKey, nil
		})

		if err != nil {
			log.Printf("Token parsing error: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["id"].(float64))
			username := claims["username"].(string)
			roleName := claims["role_name"].(string)

			ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
			ctx = context.WithValue(ctx, constants.UsernameKey, username)
			ctx = context.WithValue(ctx, constants.RoleKey, roleName)

			if claims["role_name"] == "driver" {
				// Allow access to driver-specific endpoints
				next.ServeHTTP(w, r.WithContext(ctx))
			} else if claims["role_name"] == "fleet_manager" || claims["role_name"] == "admin" || claims["role_name"] == "customer" {
				// Allow access to all endpoints
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		} else {
			log.Printf("Token validation failed")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	}
}
