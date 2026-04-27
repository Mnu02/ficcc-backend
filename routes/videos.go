package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"ficcc-backend/youtube"

	"github.com/gorilla/mux"
)

// SetupVideoRoutes registers all video-related routes.
func SetupVideoRoutes(router *mux.Router) {
	router.HandleFunc("/videos", getVideosHandler).Methods(http.MethodGet)
}

func getVideosHandler(w http.ResponseWriter, r *http.Request) {
	maxResults := 0
	if value := r.URL.Query().Get("maxResults"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, "maxResults must be a number", http.StatusBadRequest)
			return
		}
		maxResults = parsed
	}

	videos, err := youtube.GetVideos(r.Context(), maxResults)
	if err != nil {
		log.Printf("Error fetching YouTube videos: %v", err)
		http.Error(w, "Failed to fetch videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(videos); err != nil {
		log.Printf("Error encoding videos response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
