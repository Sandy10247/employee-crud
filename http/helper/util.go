package helper

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
)

// Gracefull Teriminate Server
func GracefulShutdown(server *http.Server, done chan bool, db *pgx.Conn) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error 'server.Shutdown' error: %v", err)
	}

	// Close DB Connection
	err := db.Close(ctx)
	if err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server Down 🔴")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
