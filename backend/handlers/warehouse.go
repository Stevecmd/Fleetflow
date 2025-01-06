package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stevecmd/Fleetflow/backend/models"
	"github.com/stevecmd/Fleetflow/backend/repository"
)

type WarehouseHandler struct {
	repo *repository.WarehouseRepository
}

// NewWarehouseHandler initializes and returns a new instance of WarehouseHandler.
// It requires a database connection to create a new WarehouseRepository.

func NewWarehouseHandler(db *sql.DB) *WarehouseHandler {
	return &WarehouseHandler{
		repo: repository.NewWarehouseRepository(db),
	}
}

// CreateWarehouse handles POST /api/v1/warehouses
// It creates a new warehouse entry in the database based on the provided
// request body, and returns the newly created warehouse entry in the
// response body.
func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var req models.CreateWarehouseRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	warehouse, err := h.repo.CreateWarehouse(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(warehouse)
}

// GetWarehouse handles GET /api/v1/warehouses/{id}
// It retrieves a single warehouse entry by ID from the database and
// returns it in the response body.
func (h *WarehouseHandler) GetWarehouse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	warehouse, err := h.repo.GetWarehouseByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouse)
}

// ListWarehouses handles GET /api/v1/warehouses
// It retrieves a list of warehouse entries from the database with pagination support.
// The response is a JSON array of warehouses. Query parameters "limit" and "offset"
// are used to control pagination. If "limit" is not provided, a default of 10 is used.

func (h *WarehouseHandler) ListWarehouses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	if limit == 0 {
		limit = 10 // Default limit
	}

	warehouses, err := h.repo.ListWarehouses(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouses)
}

// UpdateWarehouse handles PUT /api/v1/warehouses/{id}
// It updates a single warehouse entry by ID in the database with the
// provided request body and returns the updated warehouse entry in the
// response body.
func (h *WarehouseHandler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateWarehouseRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	warehouse, err := h.repo.UpdateWarehouse(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouse)
}

// DeleteWarehouse handles DELETE /api/v1/warehouses/{id}
// It deletes a single warehouse entry by ID from the database.
// Returns a 204 No Content status on success, or an error status
// if the warehouse ID is invalid or the deletion fails.

func (h *WarehouseHandler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteWarehouse(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
