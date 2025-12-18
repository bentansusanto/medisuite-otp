package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"medisuite-api/api/handler"
	"medisuite-api/api/routes"
	"medisuite-api/app/repo"
	"medisuite-api/app/services"
	"medisuite-api/common/middlewares"
	"medisuite-api/infra/databases"

	dbconsts "medisuite-api/constants/databases"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	goose "github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Medisuite API server",
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Load environment variables
	err := godotenv.Load(".env." + env)
	if err != nil {
		slog.Error("Warning: Error loading .env." + env + " file")
		// Try to load .env as fallback
		err = godotenv.Load(".env")
		if err != nil {
			slog.Error("Error: No .env file found")
		}
	}

	// Initialize database
	db, err := databases.InitDB()
	if err != nil {
		slog.Error("Failed to initialize database", "err", err)
		return
	}
	defer databases.CloseDB(db)

	store := repo.NewStore(db)
	repo := repo.NewRepo(store)
	service := services.NewService(repo)
	handler := handler.NewHandler(service)

	// Run migrations
	if err := runMigrations(); err != nil {
		slog.Error("Migration failed", "err", err)
		return
	}

	// Initialize Gin
	r := gin.Default()
	r.Use(middlewares.HandlePanic())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Test route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Medisuite API is running - Welcome to the backend! Environment: "+env)
	})

	// Add your routes here
	group := r.Group("/api/v1")
	route := routes.NewRoutes(handler, group, repo)
	route.Serve()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Debug("Server running on port " + port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Failed to start server", "err", err)
	}
}

func runMigrations() error {
	// Allow disabling migrations via env
	if os.Getenv("MIGRATE_ENABLED") == "false" {
		slog.Info("Migrations are disabled via MIGRATE_ENABLED=false")
		return nil
	}

	// Build Postgres DSN from central DB constants
	dbconsts.LoadDBEnv()
	dsn := dbconsts.PostgresURL()

	// Open *sql.DB using pgx stdlib driver (no pgx Pool)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB for migrations: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB for migrations: %w", err)
	}

	// Configure Goose
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	migrationDir := os.Getenv("MIGRATION_PATH")
	if migrationDir == "" {
		migrationDir = "infra/databases/migrations"
	}

	// Apply all pending migrations
	if err := goose.Up(db, migrationDir); err != nil {
		if errors.Is(err, goose.ErrNoMigrations) {
			slog.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("Goose migrations applied successfully")
	return nil
}
