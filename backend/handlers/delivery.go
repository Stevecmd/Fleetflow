package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stevecmd/Fleetflow/backend/models"
	"github.com/stevecmd/Fleetflow/backend/repository"
)

type DeliveryHandler struct {
	repo *repository.DeliveryRepository
}

func NewDeliveryHandler(repo *repository.DeliveryRepository) *DeliveryHandler {
	return &DeliveryHandler{repo: repo}
}

// ListDeliveries handles GET /api/v1/deliveries
func (h *DeliveryHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	// Log the entire request URL and parameters for debugging
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Query Parameters: %v", r.URL.Query())

	// Parse query parameters
	filter := &models.DeliveryFilter{
		Page:     1,
		PageSize: 10,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			filter.Page = pageNum
		}
	}

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil && size > 0 {
			filter.PageSize = size
		}
	}

	if customerID := r.URL.Query().Get("customer_id"); customerID != "" {
		if id, err := strconv.ParseInt(customerID, 10, 64); err == nil {
			filter.CustomerID = &id
		}
	}

	if UserID := r.URL.Query().Get("user_id"); UserID != "" {
		if id, err := strconv.ParseInt(UserID, 10, 64); err == nil {
			filter.UserID = &id
			// Log the retrieved user_id immediately after parsing
			log.Printf("Retrieved user_id from query: %s", UserID)
		} else {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Ensure user_id is not empty before proceeding
	if filter.UserID == nil {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Log the user ID for debugging
	log.Printf("User ID retrieved: %d", *filter.UserID)

	if status := r.URL.Query().Get("status_id"); status != "" {
		if id, err := strconv.ParseInt(status, 10, 64); err == nil {
			filter.StatusID = &id
		}
	}

	filter.SortBy = r.URL.Query().Get("sort_by")
	filter.SortOrder = r.URL.Query().Get("sort_order")

	// Get deliveries from repository
	deliveries, err := h.repo.List(filter)
	if err != nil {
		log.Printf("Error listing deliveries: %v", err)
		http.Error(w, "Error retrieving deliveries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

// GetDelivery handles GET /api/v1/deliveries/{id}
func (h *DeliveryHandler) GetDelivery(w http.ResponseWriter, r *http.Request) {
	// Log the entire request URL and parameters for debugging
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Query Parameters: %v", r.URL.Query())

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}

	delivery, err := h.repo.GetByID(id)
	if err != nil {
		log.Printf("Error getting delivery: %v", err)
		http.Error(w, "Error retrieving delivery", http.StatusInternalServerError)
		return
	}

	if delivery == nil {
		http.Error(w, "Delivery not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delivery)
}

// CreateDelivery handles POST /api/v1/deliveries
func (h *DeliveryHandler) CreateDelivery(w http.ResponseWriter, r *http.Request) {
	// Log the entire request URL and parameters for debugging
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Query Parameters: %v", r.URL.Query())

	var req models.CreateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Add validation for required fields

	delivery, err := h.repo.Create(&req)
	if err != nil {
		log.Printf("Error creating delivery: %v", err)
		http.Error(w, "Error creating delivery", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(delivery)
}

// UpdateDelivery handles PUT /api/v1/deliveries/{id}
func (h *DeliveryHandler) UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	// Log the entire request URL and parameters for debugging
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Query Parameters: %v", r.URL.Query())

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	delivery, err := h.repo.Update(id, &req)
	if err != nil {
		log.Printf("Error updating delivery: %v", err)
		http.Error(w, "Error updating delivery", http.StatusInternalServerError)
		return
	}

	if delivery == nil {
		http.Error(w, "Delivery not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delivery)
}

// GetDeliveryStatistics handles GET /api/v1/deliveries/stats
func (h *DeliveryHandler) GetDeliveryStatistics(w http.ResponseWriter, r *http.Request) {
	// Log the entire request URL and parameters for debugging
	log.Printf("Request URL: %s", r.URL.String())
	log.Printf("Query Parameters: %v", r.URL.Query())

	stats, err := h.repo.GetDeliveryStatistics()
	if err != nil {
		log.Printf("Error getting delivery statistics: %v", err)
		http.Error(w, "Error retrieving delivery statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
