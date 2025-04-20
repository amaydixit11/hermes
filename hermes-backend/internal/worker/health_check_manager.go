// worker/health_check_manager.go
package worker

import (
	"context"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/repository"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
)

type HealthCheckManager struct {
	healthRepo    repository.HealthRepository
	healthService *service.HealthService
	log           *logger.Logger
	stopCh        chan struct{}
}

func NewHealthCheckManager(
	healthRepo repository.HealthRepository,
	healthService *service.HealthService,
	log *logger.Logger,
) *HealthCheckManager {
	return &HealthCheckManager{
		healthRepo:    healthRepo,
		healthService: healthService,
		log:           log,
		stopCh:        make(chan struct{}),
	}
}

// Start begins the health check manager
func (m *HealthCheckManager) Start() {
	m.log.Info("Starting health check manager")

	// Perform initial check to get all health checks
	go m.processAllHealthChecks()

	// Process health checks every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go m.processAllHealthChecks()
		case <-m.stopCh:
			m.log.Info("Stopping health check manager")
			return
		}
	}
}

// Stop gracefully stops the health check manager
func (m *HealthCheckManager) Stop() {
	close(m.stopCh)
}

// processAllHealthChecks fetches all active health checks and processes them
func (m *HealthCheckManager) processAllHealthChecks() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// We'll need to implement this method in the health repository
	checks, err := m.healthRepo.GetAllActiveHealthChecks(ctx)
	if err != nil {
		m.log.Error("Failed to fetch active health checks: %v", err)
		return
	}

	m.log.Info("Processing %d active health checks", len(checks))

	// Create a map of checks by interval
	checksByInterval := make(map[int][]*models.HealthCheck)
	for _, check := range checks {
		if check.Type == models.HealthCheckTypeActive && check.Enabled {
			checksByInterval[check.Interval] = append(checksByInterval[check.Interval], check)
		}
	}

	// Process checks by interval
	for interval, checksForInterval := range checksByInterval {
		m.processChecksWithInterval(checksForInterval, interval)
	}
}

// processChecksWithInterval handles checks with the same interval
func (m *HealthCheckManager) processChecksWithInterval(checks []*models.HealthCheck, interval int) {
	m.log.Debug("Processing %d health checks with interval %d seconds", len(checks), interval)

	// Create a ticker for this interval
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Process all checks immediately
	for _, check := range checks {
		go func(c *models.HealthCheck) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
			defer cancel()

			if err := m.healthService.RunActiveHealthCheck(ctx, c); err != nil {
				m.log.Error("Failed to run active health check for service %s: %v", c.ServiceID, err)
			}
		}(check)
	}

	// Continue processing on the interval
	for {
		select {
		case <-ticker.C:
			for _, check := range checks {
				go func(c *models.HealthCheck) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
					defer cancel()

					if err := m.healthService.RunActiveHealthCheck(ctx, c); err != nil {
						m.log.Error("Failed to run active health check for service %s: %v", c.ServiceID, err)
					}
				}(check)
			}
		case <-m.stopCh:
			return
		}
	}
}
