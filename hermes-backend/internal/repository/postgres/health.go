// repository/health_repository_gorm.go
package postgres

import (
	"context"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"gorm.io/gorm"
)

type HealthRepositoryGorm struct {
	db *gorm.DB
}

func NewHealthRepositoryGorm(db *gorm.DB) *HealthRepositoryGorm {
	return &HealthRepositoryGorm{db: db}
}

// Health checks
func (r *HealthRepositoryGorm) CreateHealthCheck(ctx context.Context, check *models.HealthCheck) error {
	return r.db.WithContext(ctx).Create(check).Error
}

func (r *HealthRepositoryGorm) GetHealthChecks(ctx context.Context, serviceID string) ([]*models.HealthCheck, error) {
	var checks []*models.HealthCheck
	err := r.db.WithContext(ctx).Where("service_id = ?", serviceID).Find(&checks).Error
	return checks, err
}

func (r *HealthRepositoryGorm) GetHealthCheck(ctx context.Context, id uint) (*models.HealthCheck, error) {
	var check models.HealthCheck
	err := r.db.WithContext(ctx).First(&check, id).Error
	if err != nil {
		return nil, err
	}
	return &check, nil
}

func (r *HealthRepositoryGorm) UpdateHealthCheck(ctx context.Context, check *models.HealthCheck) error {
	return r.db.WithContext(ctx).Save(check).Error
}

func (r *HealthRepositoryGorm) DeleteHealthCheck(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.HealthCheck{}, id).Error
}

// Health history
func (r *HealthRepositoryGorm) RecordHealthHistory(ctx context.Context, history *models.HealthHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *HealthRepositoryGorm) GetHealthHistory(ctx context.Context, serviceID string, params models.HealthHistoryQueryParams) ([]*models.HealthHistory, int64, error) {
	var histories []*models.HealthHistory
	var count int64

	query := r.db.WithContext(ctx).Model(&models.HealthHistory{}).Where("service_id = ?", serviceID)

	if !params.StartTime.IsZero() {
		query = query.Where("timestamp >= ?", params.StartTime)
	}

	if !params.EndTime.IsZero() {
		query = query.Where("timestamp <= ?", params.EndTime)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Get total count
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset(params.Offset).Limit(params.Limit).Order("timestamp desc").Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, count, nil
}

// Custom metrics
func (r *HealthRepositoryGorm) UpdateCustomMetric(ctx context.Context, metric *models.CustomHealthMetric) error {
	// Try to update existing metric, create if not exists
	result := r.db.WithContext(ctx).Where("service_id = ? AND name = ?", metric.ServiceID, metric.Name).
		Updates(map[string]interface{}{
			"value":      metric.Value,
			"unit":       metric.Unit,
			"updated_at": time.Now(),
		})

	if result.RowsAffected == 0 {
		// Create new metric
		return r.db.WithContext(ctx).Create(metric).Error
	}

	return result.Error
}

func (r *HealthRepositoryGorm) GetCustomMetrics(ctx context.Context, serviceID string) ([]*models.CustomHealthMetric, error) {
	var metrics []*models.CustomHealthMetric
	err := r.db.WithContext(ctx).Where("service_id = ?", serviceID).Find(&metrics).Error
	return metrics, err
}

// Health thresholds
func (r *HealthRepositoryGorm) CreateHealthThreshold(ctx context.Context, threshold *models.HealthThreshold) error {
	return r.db.WithContext(ctx).Create(threshold).Error
}

func (r *HealthRepositoryGorm) GetHealthThresholds(ctx context.Context, serviceID string) ([]*models.HealthThreshold, error) {
	var thresholds []*models.HealthThreshold
	err := r.db.WithContext(ctx).Where("service_id = ?", serviceID).Find(&thresholds).Error
	return thresholds, err
}

func (r *HealthRepositoryGorm) UpdateHealthThreshold(ctx context.Context, threshold *models.HealthThreshold) error {
	return r.db.WithContext(ctx).Save(threshold).Error
}

func (r *HealthRepositoryGorm) DeleteHealthThreshold(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.HealthThreshold{}, id).Error
}

// Health status management
func (r *HealthRepositoryGorm) UpdateHealthStatus(ctx context.Context, serviceID string, status models.ServiceStatus, message string) error {
	return r.db.WithContext(ctx).Model(&models.Service{}).
		Where("id = ?", serviceID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *HealthRepositoryGorm) IncrementHealthCheckFailures(ctx context.Context, checkID uint) error {
	return r.db.WithContext(ctx).Model(&models.HealthCheck{}).
		Where("id = ?", checkID).
		UpdateColumn("timeout_count", gorm.Expr("timeout_count + 1")).
		Error
}

func (r *HealthRepositoryGorm) ResetHealthCheckFailures(ctx context.Context, checkID uint) error {
	return r.db.WithContext(ctx).Model(&models.HealthCheck{}).
		Where("id = ?", checkID).
		UpdateColumn("timeout_count", 0).
		Error
}

func (r *HealthRepositoryGorm) GetAllActiveHealthChecks(ctx context.Context) ([]*models.HealthCheck, error) {
	var checks []*models.HealthCheck
	err := r.db.WithContext(ctx).Where("enabled = ?", true).Find(&checks).Error
	return checks, err
}
