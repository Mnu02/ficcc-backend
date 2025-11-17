package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ficcc-backend/db"
	"ficcc-backend/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := routes.SetupRoutes()

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel to receive server errors
	errChan := make(chan error, 1)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for either interrupt signal or server error
	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
		log.Println("Shutting down server...")
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down server...")
	}

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Database will be closed by defer at line 27
	log.Println("Server stopped")
}
