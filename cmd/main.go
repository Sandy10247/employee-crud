package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"server/http/helper"
	router "server/http/router"
	db "server/init"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR loading env :", err)
	}

	err = db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.DisconnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	srv := serve(port)
	println("Server running on port", port)

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Printf("ERROR starting server : %v\n", err)
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go func() {
		helper.GracefulShutdown(srv, done, db.DB)
	}()

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func serve(port string) *http.Server {
	router := router.InitRouter()

	server := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
