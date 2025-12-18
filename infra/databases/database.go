package databases

import (
	"context"
	"fmt"
	"log/slog"

	// "medisuite/pkg/logs"
	"os"

	"github.com/jackc/pgx/v5"
)

// InitDB initializes a new database connection
func InitDB() (*pgx.Conn, error) {
	// Get database configuration from environment variables with defaults
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// Log database configuration (without password for security)
	logMsg := fmt.Sprintf("Connecting to database: %s@%s:%s/%s", user, host, port, dbname)
	slog.Info(logMsg)

	// Create connection string in DSN format
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=Asia/Jakarta",
		user, password, host, port, dbname,
	)

	// Create connection
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("Successfully connected to database")
	return conn, nil
}

// CloseDB closes the database connection
func CloseDB(conn *pgx.Conn) {
	if conn != nil {
		conn.Close(context.Background())
	}
}
