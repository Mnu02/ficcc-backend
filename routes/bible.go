package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const bibleBaseURL = "https://rest.api.bible/v1"

// SetupBibleRoutes registers all Bible-related routes
func SetupBibleRoutes(router *mux.Router) {
	router.HandleFunc("/bible/versions", getBibleVersionsHandler).Methods("GET")
	router.HandleFunc("/bible/books", getBibleBooksHandler).Methods("GET")
	router.HandleFunc("/bible/chapters", getBibleChaptersHandler).Methods("GET")
	router.HandleFunc("/bible/verses", getBibleVersesHandler).Methods("GET")
}

// proxyBibleRequest forwards a request to api.bible and writes the response back
func proxyBibleRequest(w http.ResponseWriter, apiURL string) {
	apiKey := os.Getenv("BIBLE_API_KEY")
	if apiKey == "" {
		http.Error(w, "Bible API key not configured", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("Error creating Bible API request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error calling Bible API: %v", err)
		http.Error(w, "Failed to reach Bible API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error writing Bible API response: %v", err)
	}
}

// getBibleVersionsHandler returns all available English Bible versions
func getBibleVersionsHandler(w http.ResponseWriter, r *http.Request) {
	proxyBibleRequest(w, fmt.Sprintf("%s/bibles?language=eng", bibleBaseURL))
}

// getBibleBooksHandler returns all books for a given Bible version
// Query param: versionId
func getBibleBooksHandler(w http.ResponseWriter, r *http.Request) {
	versionID := r.URL.Query().Get("versionId")
	if versionID == "" {
		http.Error(w, "versionId is required", http.StatusBadRequest)
		return
	}
	proxyBibleRequest(w, fmt.Sprintf("%s/bibles/%s/books", bibleBaseURL, versionID))
}

// getBibleChaptersHandler returns all chapters for a book in a given Bible version
// Query params: versionId, bookId
func getBibleChaptersHandler(w http.ResponseWriter, r *http.Request) {
	versionID := r.URL.Query().Get("versionId")
	bookID := r.URL.Query().Get("bookId")
	if versionID == "" || bookID == "" {
		http.Error(w, "versionId and bookId are required", http.StatusBadRequest)
		return
	}
	proxyBibleRequest(w, fmt.Sprintf("%s/bibles/%s/books/%s/chapters", bibleBaseURL, versionID, bookID))
}

// getBibleVersesHandler returns all verses for a chapter in a given Bible version
// Query params: versionId, chapterId
func getBibleVersesHandler(w http.ResponseWriter, r *http.Request) {
	versionID := r.URL.Query().Get("versionId")
	chapterID := r.URL.Query().Get("chapterId")
	if versionID == "" || chapterID == "" {
		http.Error(w, "versionId and chapterId are required", http.StatusBadRequest)
		return
	}
	proxyBibleRequest(w, fmt.Sprintf("%s/bibles/%s/chapters/%s", bibleBaseURL, versionID, chapterID))
}
