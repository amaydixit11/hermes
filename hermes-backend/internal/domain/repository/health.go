// repository/health_repository.go
package repository

import (
	"context"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
)

type HealthRepository interface {
	// Health checks
	CreateHealthCheck(ctx context.Context, check *models.HealthCheck) error
	GetHealthChecks(ctx context.Context, serviceID string) ([]*models.HealthCheck, error)
	GetHealthCheck(ctx context.Context, id uint) (*models.HealthCheck, error)
	UpdateHealthCheck(ctx context.Context, check *models.HealthCheck) error
	DeleteHealthCheck(ctx context.Context, id uint) error
	GetAllActiveHealthChecks(ctx context.Context) ([]*models.HealthCheck, error)

	// Health history
	RecordHealthHistory(ctx context.Context, history *models.HealthHistory) error
	GetHealthHistory(ctx context.Context, serviceID string, params models.HealthHistoryQueryParams) ([]*models.HealthHistory, int64, error)

	// Custom metrics
	UpdateCustomMetric(ctx context.Context, metric *models.CustomHealthMetric) error
	GetCustomMetrics(ctx context.Context, serviceID string) ([]*models.CustomHealthMetric, error)

	// Health thresholds
	CreateHealthThreshold(ctx context.Context, threshold *models.HealthThreshold) error
	GetHealthThresholds(ctx context.Context, serviceID string) ([]*models.HealthThreshold, error)
	UpdateHealthThreshold(ctx context.Context, threshold *models.HealthThreshold) error
	DeleteHealthThreshold(ctx context.Context, id uint) error

	// Health status management
	UpdateHealthStatus(ctx context.Context, serviceID string, status models.ServiceStatus, message string) error
	IncrementHealthCheckFailures(ctx context.Context, checkID uint) error
	ResetHealthCheckFailures(ctx context.Context, checkID uint) error
}
