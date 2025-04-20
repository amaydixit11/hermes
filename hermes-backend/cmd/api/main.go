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
	repoPostgres "github.com/amaydixit11/hermes/hermes-backend/internal/repository/postgres"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/amaydixit11/hermes/hermes-backend/internal/worker"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log := logger.New("error")
		log.Fatal("Failed to load configuration", "error", err)
		return
	}

	log := logger.New(cfg.LogLevel)
	log.Info("Configuration loaded successfully", "env", cfg.Environment)

	// Connect to database
	db, err := gorm.Open(gormpostgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to retrieve SQL DB from GORM", "error", err)
	}
	defer sqlDB.Close()

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
