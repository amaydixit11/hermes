// internal/repository/postgres/service.go
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceRepository implements the repository.ServiceRepository interface
type ServiceRepository struct {
	db *gorm.DB
}

// NewServiceRepository creates a new ServiceRepository
func NewServiceRepository(db *gorm.DB) repository.ServiceRepository {
	return &ServiceRepository{db: db}
}

// Create adds a new service to the database
func (r *ServiceRepository) Create(ctx context.Context, service *models.Service) error {
	if service.ID == "" {
		service.ID = uuid.New().String()
	}

	return r.db.WithContext(ctx).Create(service).Error
}

// GetByID retrieves a service by its ID
func (r *ServiceRepository) GetByID(ctx context.Context, id string) (*models.Service, error) {
	var service models.Service
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if service is not found
		}
		return nil, err // Return error if something went wrong
	}
	return &service, nil
}

// GetByName retrieves a service by its name
func (r *ServiceRepository) GetByName(ctx context.Context, name string) (*models.Service, error) {
	var service models.Service
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if service is not found
		}
		return nil, err // Return error if something went wrong
	}
	return &service, nil
}

// List retrieves services with filtering and pagination
func (r *ServiceRepository) List(ctx context.Context, params models.ServiceQueryParams) ([]*models.Service, int64, error) {
	var services []*models.Service
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Service{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if len(params.Tags) > 0 {
		query = query.Where("tags @> ?", params.Tags) // Array contains condition (if tags are stored as an array in DB)
	}
	if params.Search != "" {
		query = query.Where("name LIKE ?", "%"+params.Search+"%")
	}

	// Count total number of records (for pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if err := query.Offset(params.Offset).Limit(params.Limit).Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

// Update modifies an existing service's details
func (r *ServiceRepository) Update(ctx context.Context, service *models.Service) error {
	return r.db.WithContext(ctx).Save(service).Error
}

// Delete removes a service by its ID
func (r *ServiceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Service{}).Error
}

// UpdateStatus updates the status of a service
func (r *ServiceRepository) UpdateStatus(ctx context.Context, id string, status models.ServiceStatus) error {
	return r.db.WithContext(ctx).Model(&models.Service{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateLastSeen updates the LastSeen timestamp of a service
func (r *ServiceRepository) UpdateLastSeen(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Service{}).Where("id = ?", id).Update("last_seen", time.Now()).Error
}
