package database

import (
	"fmt"
	"os"
)

// firstEnv returns the first non-empty environment variable value from the provided keys
func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// Exported DB env variables. Call LoadDBEnv() after .env is loaded to populate.
var (
	DBUser string
	DBPass string
	DBHost string
	DBName string
	DBPort string
)

// LoadDBEnv populates the exported DB* variables from environment with fallbacks
func LoadDBEnv() {
	DBUser = firstEnv("POSTGRES_USER", "DB_USER")
	DBPass = firstEnv("POSTGRES_PASSWORD", "DB_PASS")
	DBHost = firstEnv("POSTGRES_HOST", "DB_HOST")
	DBName = firstEnv("POSTGRES_DB", "DB_NAME")
	DBPort = firstEnv("POSTGRES_PORT", "DB_PORT")
	if DBPort == "" {
		DBPort = "5432"
	}
}

// PostgresURL builds the DSN for postgres using envs and sane defaults
func PostgresURL() string {
	// assumes LoadDBEnv() has been called
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DBUser, DBPass, DBHost, DBPort, DBName,
	)
}
