package helper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	log.Println("Server Stopped 🔴")

	// Close DB Connection
	err := db.Close(ctx)
	if err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}
	log.Println("Database Closed 🔴")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func FloatToNumeric(f float64, precision int) (pgtype.Numeric, error) {
	var numericValue pgtype.Numeric

	// Format the float to a string with the desired precision
	// 'f' format specifier, precision specifies the number of digits after the decimal point
	str := strconv.FormatFloat(f, 'f', precision, 64)

	// Use Scan to parse the string into the pgtype.Numeric struct
	if err := numericValue.Scan(str); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("failed to scan string to pgtype.Numeric: %w", err)
	}

	return numericValue, nil
}

// CalculateNetSalary computes the take-home pay after all deductions.
// Formula: Net Salary = Gross Salary - Total Deductions
func CalculateNetSalary(grossSalary, totalDeductions float64) float64 {
	// Ensure net salary isn't negative in edge cases
	if totalDeductions > grossSalary {
		return 0
	}
	return grossSalary - totalDeductions
}

// CalculatePercentage computes what 'percent' of 'total' is.
func CalculatePercentage(percent, total float64) float64 {
	return (total * percent) / 100
}

func GetTaxRatePerCountry(country string) float64 {
	// extract country tax from ".env"
	countryCleaned := strings.ToLower(country)
	taxRate := GetEnv(countryCleaned, "")

	taxRateFloat, err := strconv.ParseFloat(taxRate, 64)
	if err != nil {
		return 0.0
	}

	return taxRateFloat
}

// getEnv retrieves environment variables or returns a default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves integer environment variables or returns a default value
func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
