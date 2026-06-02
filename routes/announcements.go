package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"ficcc-backend/sheets"

	"github.com/gorilla/mux"
)

// SetupAnnouncementRoutes registers all announcement-related routes.
func SetupAnnouncementRoutes(router *mux.Router) {
	router.HandleFunc("/announcements", getAnnouncementsHandler).Methods(http.MethodGet)
}

func getAnnouncementsHandler(w http.ResponseWriter, r *http.Request) {
	announcements, err := sheets.GetAnnouncements(r.Context())
	if err != nil {
		log.Printf("Error fetching announcements from Google Sheets: %v", err)
		http.Error(w, "Failed to fetch announcements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(announcements); err != nil {
		log.Printf("Error encoding announcements response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
