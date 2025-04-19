// cmd/api/main.go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/api"
	"github.com/amaydixit11/hermes/hermes-backend/internal/config"
	"github.com/amaydixit11/hermes/hermes-backend/internal/monitoring"
	repoPostgres "github.com/amaydixit11/hermes/hermes-backend/internal/repository/postgres"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Consistent logging instead of panic
		log := logger.New("error")
		log.Fatal("Failed to load configuration", "error", err)
		return
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	log.Info("Configuration loaded successfully", "env", cfg.Environment)

	// Connect to database
	db, err := gorm.Open(gormpostgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	// Get raw SQL DB to close on shutdown
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to retrieve SQL DB from GORM", "error", err)
	}
	defer sqlDB.Close()

	// Initialize repositories
	serviceRepo := repoPostgres.NewServiceRepository(db)

	// Initialize services
	serviceService := service.NewServiceService(serviceRepo, log)

	// Initialize health checker (interval assumed in config)
	healthCheckInterval := cfg.HealthCheck.Interval
	if healthCheckInterval == 0 {
		healthCheckInterval = int(30 * time.Second / time.Second)
	}
	healthChecker := monitoring.NewHealthChecker(serviceRepo, log, time.Duration(healthCheckInterval)*time.Second)

	// Start health checker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	healthChecker.Start(ctx)

	// Set up HTTP router
	router := api.SetupRouter(cfg, log, serviceService)

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

	// Stop health checker
	healthChecker.Stop()

	log.Info("Server exited gracefully")
	defer log.Sync()

}
