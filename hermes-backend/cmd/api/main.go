// cmd/api/main.go
package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/api"
	"github.com/amaydixit11/hermes/hermes-backend/internal/config"
	"github.com/amaydixit11/hermes/hermes-backend/internal/database"
	repoPostgres "github.com/amaydixit11/hermes/hermes-backend/internal/repository/postgres"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/amaydixit11/hermes/hermes-backend/internal/worker"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Define command line flags
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	migrationsPath := flag.String("migrations", "./migrations", "Path to migrations directory")
	rollbackFlag := flag.Bool("rollback", false, "Rollback the last migration")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log := logger.New("error")
		log.Fatal("Failed to load configuration", "error", err)
		return
	}
	log := logger.New(cfg.LogLevel)
	log.Info("Configuration loaded successfully", "env", cfg.Environment)

	// Handle migration commands if requested
	if *migrateFlag || *rollbackFlag {
		absPath, err := filepath.Abs(*migrationsPath)
		if err != nil {
			log.Fatal("Failed to get absolute path for migrations", "error", err)
		}

		dsn := cfg.DSN()
		log.Info("Database DSN", "dsn", dsn)

		if *migrateFlag {
			log.Info("Running migrations", "path", absPath)
			if err := database.RunMigrations(dsn, absPath); err != nil {
				log.Fatal("Migration failed", "error", err)
			}
			log.Info("Migrations completed successfully")
		} else if *rollbackFlag {
			log.Info("Rolling back last migration", "path", absPath)
			if err := database.RollbackLastMigration(dsn, absPath); err != nil {
				log.Fatal("Rollback failed", "error", err)
			}
			log.Info("Rollback completed successfully")
		}

		return
	}

	// Connect to database for normal operation
	db, err := gorm.Open(gormpostgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to retrieve SQL DB from GORM", "error", err)
	}
	defer sqlDB.Close()

	// Initialize repositories and services
	serviceRepo := repoPostgres.NewServiceRepository(db)
	healthRepo := repoPostgres.NewHealthRepositoryGorm(db)
	serviceService := service.NewServiceService(serviceRepo, log)
	healthService := service.NewHealthService(healthRepo, serviceRepo, log)
	healthCheckManager := worker.NewHealthCheckManager(healthRepo, healthService, log)
	go healthCheckManager.Start()

	// Set up HTTP router
	router := api.SetupRouter(cfg, log, serviceService, healthService)

	// Start HTTP server with proper timeouts
	srv := &http.Server{
		Addr:         "localhost:" + strconv.Itoa(cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.TimeoutRead) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.TimeoutWrite) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.TimeoutIdle) * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info("Starting server", "address", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Shutdown server with timeout
	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.TimeoutShutdown)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}
	healthCheckManager.Stop()
	log.Info("Server exited gracefully")
	defer log.Sync()
}
