package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stevecmd/Fleetflow/backend/models"
	"github.com/stevecmd/Fleetflow/backend/pkg/constants"
)

func GetCustomerDeliveries(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context using constants
		userID, ok := r.Context().Value(constants.UserIDKey).(int)
		if !ok {
			log.Printf("Failed to get userID from context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get role from context using constants
		role, ok := r.Context().Value(constants.RoleKey).(string)
		if !ok {
			log.Printf("Failed to get role from context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify user role
		if role != "customer" {
			http.Error(w, "Unauthorized role", http.StatusForbidden)
			return
		}

		// Query deliveries using ID from context
		query := `
            SELECT d.id, d.tracking_number, d.status_id, 
                   d.estimated_delivery_time, d.cargo_type, 
                   d.from_location, d.proof_of_delivery_image_url
            FROM deliveries d
            WHERE d.user_id = $1
            ORDER BY d.created_at DESC
        `

		rows, err := db.Query(query, userID)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var deliveries []models.Delivery
		for rows.Next() {
			var d models.Delivery
			err := rows.Scan(
				&d.ID, &d.TrackingNumber, &d.StatusID,
				&d.EstimatedDeliveryTime, &d.CargoType,
				&d.FromLocation, &d.ProofOfDeliveryImageURL,
			)
			if err != nil {
				log.Printf("Row scan error: %v", err)
				continue
			}
			deliveries = append(deliveries, d)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deliveries)
	}
}

func GetCustomerInvoices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		query := `
            SELECT id, delivery_id, amount, status, due_date, payment_method
            FROM invoices
            WHERE user_id = $1
            ORDER BY created_at DESC
        `

		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var invoices []models.Invoice
		for rows.Next() {
			var i models.Invoice
			err := rows.Scan(&i.ID, &i.DeliveryID, &i.Amount,
				&i.Status, &i.DueDate, &i.PaymentMethod)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			invoices = append(invoices, i)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}

func GetCustomerFeedback(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		query := `
            SELECT id, delivery_id, rating, feedback_text, 
                   timeliness_rating, driver_rating, package_condition_rating
            FROM delivery_feedback
            WHERE user_id = $1
            ORDER BY created_at DESC
        `

		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var feedback []models.DeliveryFeedback
		for rows.Next() {
			var f models.DeliveryFeedback
			err := rows.Scan(&f.ID, &f.DeliveryID, &f.Rating,
				&f.FeedbackText, &f.TimelinessRating,
				&f.DriverRating, &f.PackageConditionRating)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			feedback = append(feedback, f)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feedback)
	}
}
