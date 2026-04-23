package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"ficcc-backend/sheets"

	"github.com/gorilla/mux"
)

// SetupEventRoutes registers all event-related routes.
func SetupEventRoutes(router *mux.Router) {
	router.HandleFunc("/events", getEventsHandler).Methods(http.MethodGet)
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := sheets.GetEvents(r.Context())
	if err != nil {
		log.Printf("Error fetching events from Google Sheets: %v", err)
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Printf("Error encoding events response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
