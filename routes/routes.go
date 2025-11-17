package routes

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes and returns the router
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Register routes for each table
	SetupSermonRoutes(router)

	return router
}
