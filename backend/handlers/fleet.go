package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetFleetAnalytics(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Get vehicle statistics
		vehicleStats := struct {
			Total           int     `json:"total"`
			Active          int     `json:"active"`
			Maintenance     int     `json:"maintenance"`
			Idle            int     `json:"idle"`
			UtilizationRate float64 `json:"utilizationRate"`
		}{}

		err := db.QueryRow(`
            SELECT 
                COALESCE(COUNT(*), 0) as total,
                COALESCE(COUNT(CASE WHEN vs.name = 'available' THEN 1 END), 0) as active,
                COALESCE(COUNT(CASE WHEN vs.name = 'in_maintenance' THEN 1 END), 0) as maintenance,
                COALESCE(COUNT(CASE WHEN vs.name = 'out_of_service' THEN 1 END), 0) as idle
            FROM vehicles v
            LEFT JOIN vehicle_statuses vs ON v.status_id = vs.id
        `).Scan(&vehicleStats.Total, &vehicleStats.Active, &vehicleStats.Maintenance, &vehicleStats.Idle)

		if err != nil {
			log.Printf("Error fetching vehicle stats: %v", err)
			http.Error(w, "Error fetching vehicle stats", http.StatusInternalServerError)
			return
		}

		// Calculate utilization rate
		if vehicleStats.Total > 0 {
			vehicleStats.UtilizationRate = float64(vehicleStats.Active) / float64(vehicleStats.Total) * 100
		}

		// Get maintenance statistics with COALESCE
		maintenanceStats := struct {
			Pending   int `json:"pending"`
			Completed int `json:"completed"`
			Overdue   int `json:"overdue"`
			Upcoming  int `json:"upcoming"`
		}{}

		err = db.QueryRow(`
            SELECT 
                COALESCE(COUNT(CASE WHEN next_service_date > NOW() AND service_date IS NULL THEN 1 END), 0) as pending,
                COALESCE(COUNT(CASE WHEN service_date IS NOT NULL THEN 1 END), 0) as completed,
                COALESCE(COUNT(CASE WHEN next_service_date < NOW() AND service_date IS NULL THEN 1 END), 0) as overdue,
                COALESCE(COUNT(CASE WHEN next_service_date > NOW() AND service_date IS NULL THEN 1 END), 0) as upcoming
            FROM maintenance_records
            WHERE service_date >= NOW() - INTERVAL '30 days' OR next_service_date >= NOW()
        `).Scan(&maintenanceStats.Pending, &maintenanceStats.Completed,
			&maintenanceStats.Overdue, &maintenanceStats.Upcoming)

		if err != nil {
			log.Printf("Error fetching maintenance stats: %v", err)
			http.Error(w, "Error fetching maintenance stats", http.StatusInternalServerError)
			return
		}

		// Get driver performance statistics with COALESCE
		driverStats := struct {
			TotalDrivers  int     `json:"totalDrivers"`
			AvgRating     float64 `json:"avgRating"`
			TopPerformers int     `json:"topPerformers"`
			OnTimeRate    float64 `json:"on_time"`
		}{}

		err = db.QueryRow(`
            SELECT 
                COALESCE(COUNT(DISTINCT user_id), 0) as total_drivers,
                COALESCE(AVG(customer_rating_avg), 0) as avg_rating,
                COALESCE(COUNT(CASE WHEN safety_score >= 90 THEN 1 END), 0) as top_performers,
                COALESCE(AVG(on_time_delivery_rate), 0) as on_time_rate
            FROM driver_performance_metrics
            WHERE metric_date >= NOW() - INTERVAL '30 days'
        `).Scan(&driverStats.TotalDrivers, &driverStats.AvgRating,
			&driverStats.TopPerformers, &driverStats.OnTimeRate)

		if err != nil {
			log.Printf("Error fetching driver stats: %v", err)
			http.Error(w, "Error fetching driver stats", http.StatusInternalServerError)
			return
		}

		// Combine all statistics
		response := struct {
			VehicleUtilization  interface{} `json:"vehicleUtilization"`
			MaintenanceSchedule interface{} `json:"maintenanceSchedule"`
			DriverPerformance   interface{} `json:"driverPerformance"`
			FleetStatus         interface{} `json:"fleetStatus"`
		}{
			VehicleUtilization:  vehicleStats,
			MaintenanceSchedule: maintenanceStats,
			DriverPerformance:   driverStats,
			FleetStatus: map[string]int{
				"operational":      vehicleStats.Active,
				"underMaintenance": vehicleStats.Maintenance,
				"outOfService":     vehicleStats.Idle,
				"total":            vehicleStats.Total,
			},
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
