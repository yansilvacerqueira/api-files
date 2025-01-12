package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig creates a database configuration from environment variables
func NewConfig() (*Config, error) {
	port, err := strconv.Atoi(getEnvOrDefault("PORT_DB", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	return &Config{
		Host:     getEnvOrDefault("HOST_DB", "localhost"),
		Port:     port,
		User:     getEnvOrDefault("USER_DB", "postgres"),
		Password: getEnvOrDefault("PASSWORD_DB", ""),
		DBName:   getEnvOrDefault("NAME_DB", "postgres"),
		SSLMode:  getEnvOrDefault("SSL_MODE", "disable"),
	}, nil
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// NewConnection establishes a new database connection with retry mechanism
func NewConnection(maxRetries int, retryDelay time.Duration) (*sql.DB, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	var db *sql.DB
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			lastErr = fmt.Errorf("failed to open database connection: %w", err)
			log.Printf("Attempt %d: %v", i+1, lastErr)
			time.Sleep(retryDelay)
			continue
		}

		// Configure connection pool
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		// Test connection
		if err = db.Ping(); err != nil {
			lastErr = fmt.Errorf("failed to ping database: %w", err)
			log.Printf("Attempt %d: %v", i+1, lastErr)
			db.Close()
			time.Sleep(retryDelay)
			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, lastErr)
}

// Close safely closes the database connection
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
