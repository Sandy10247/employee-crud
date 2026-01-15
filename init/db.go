package init

import (
	// "context"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	sqlc "server/sql/database" // Adjust this import path as per your project structure

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// Global variables to hold the database and queries
var (
	DB      *pgx.Conn
	Queries *sqlc.Queries
)

// Config holds the configuration for database connection
type Config struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig loads database configuration from environment variables or another source
func LoadConfig() Config {
	return Config{
		Driver:   getEnv("DB_DRIVER", "postgres"),
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "root"),
		DBName:   getEnv("DB_NAME", "attempt2"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv retrieves environment variables or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves integer environment variables or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func DSN(c *Config) string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

func Connect(c *Config) *pgx.Conn {
	// Use conf.Database to constrct the connection string and connect to the database
	connectionDsn := DSN(c)

	// Example using pgx to connect to PostgreSQL
	conn, err := pgx.Connect(context.Background(), connectionDsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Test the connection
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}

	return conn
}

// ConnectDB initializes the database connection and sqlc Queries
func ConnectDB() error {
	config := LoadConfig()

	conn := Connect(&config)

	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFn()

	if err := conn.Ping(ctx); err != nil {
		return fmt.Errorf("could not connect to db: %w", err)
	}

	DB = conn
	Queries = sqlc.New(conn) // Initialize sqlc Queries with the database connection

	log.Println("Successfully connected to the database")
	return nil
}

// DisconnectDB closes the database connection
func DisconnectDB() {
	if DB != nil {
		ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancelFn()
		DB.Close(ctx)
	}
}
