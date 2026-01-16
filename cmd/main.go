package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"server/http/helper"
	"server/http/router"
	db "server/init"

	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func CreateHttpServer(port string, router http.Handler) *http.Server {
	server := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}

func start() int {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR loading env :- ", zap.Error(err))
	}
	logEnv := helper.GetEnv("LOG_ENV", "development")
	log, err := createLogger(logEnv)
	if err != nil {
		log.Info("Error setting up the logger :- ", zap.Error(err))
		return 1
	}
	defer func() {
		// If we cannot sync, there's probably something wrong with outputting logs,
		// so we probably cannot write using fmt.Println either. So just ignore the error.
		_ = log.Sync()
	}()

	err = db.ConnectDB()
	if err != nil {
		log.Sugar().Panicf("Failed to connect to database: %v", err)
	}
	defer db.DisconnectDB()

	// extract PORT
	port := helper.GetEnv("PORT", "8080")

	// Create Http.Server
	srv := CreateHttpServer(port, router.InitRouter(log))

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run Server in an Go Routine
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Info("Error starting server", zap.Error(err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go func() {
		helper.GracefulShutdown(srv, done, db.DB)
	}()

	// Wait for the graceful shutdown to complete
	<-done
	log.Info("Graceful shutdown complete.")
	return 0
}

func main() {
	os.Exit(start())
}
