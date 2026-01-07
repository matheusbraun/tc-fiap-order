package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server
	ServerPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// External Services
	CustomerServiceURL string
	ProductServiceURL  string

	// HTTP Client
	HTTPClientTimeout      time.Duration
	HTTPClientRetryCount   int
	HTTPClientRetryBackoff time.Duration
}

func Load() (*Config, error) {
	config := &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "order_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// External Services
		CustomerServiceURL: getEnv("CUSTOMER_SERVICE_URL", "http://localhost:8081"),
		ProductServiceURL:  getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082"),

		// HTTP Client
		HTTPClientTimeout:      time.Duration(getEnvAsInt("HTTP_CLIENT_TIMEOUT_SECONDS", 30)) * time.Second,
		HTTPClientRetryCount:   getEnvAsInt("HTTP_CLIENT_RETRY_COUNT", 3),
		HTTPClientRetryBackoff: time.Duration(getEnvAsInt("HTTP_CLIENT_RETRY_BACKOFF_MS", 100)) * time.Millisecond,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: invalid integer for %s, using default %d\n", key, defaultValue)
		return defaultValue
	}
	return value
}
