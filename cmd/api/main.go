package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"point-system-api/config"
	"point-system-api/internal/database"
	"point-system-api/internal/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the database
	db := database.New()
	defer db.Close()

	// Run database migrations
	if err := database.MigrateDB(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if err := database.InitializeViewDB(); err != nil {
		log.Fatalf("Failed to creating view in database: %v", err)
	}

	// Create a new server instance
	server := server.NewServer()

	// Create a channel to listen for interrupt signals
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server is running on port %d", cfg.ServerPort)

	// Wait for an interrupt signal
	<-quit
	log.Println("Server is shutting down...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Notify the main goroutine that shutdown is complete
	close(done)
}
