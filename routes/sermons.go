package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"ficcc-backend/db"

	"github.com/gorilla/mux"
)

// SetupSermonRoutes registers all sermon-related routes
func SetupSermonRoutes(router *mux.Router) {
	// Sermons routes
	router.HandleFunc("/sermons", getSermonsHandler).Methods("GET")

}

// getSermonsHandler demonstrates how to use GetSermons() to return JSON
func getSermonsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get sermons from database
	sermons, err := db.GetSermons(ctx)
	if err != nil {
		log.Printf("Error fetching sermons: %v", err)
		http.Error(w, "Failed to fetch sermons", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sermons); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
