// internal/repository/postgres/service.go
package postgres

import (
	"context"
	"errors"
	"fmt"
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
			return nil, nil
		}
		return nil, err
	}
	return &service, nil
}

// GetByName retrieves a service by its name
func (r *ServiceRepository) GetByName(ctx context.Context, name string) (*models.Service, error) {
	var service models.Service
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
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
		query = query.Where("tags @> ?", params.Tags)
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

// internal/infrastructure/database/service_repository.go
// Implement the new methods in the repository

func (r *ServiceRepository) AdvancedSearch(ctx context.Context, params models.AdvancedDiscoveryParams) ([]*models.Service, int64, error) {
	var services []*models.Service
	var count int64

	query := r.db.WithContext(ctx)

	// Apply basic filters from ServiceQueryParams
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	if len(params.Tags) > 0 {
		for _, tag := range params.Tags {
			query = query.Where("? = ANY(tags)", tag)
		}
	}

	if params.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Apply health-aware filters
	if len(params.HealthStatus) > 0 {
		query = query.Where("status IN ?", params.HealthStatus)
	}

	// Apply metadata filters
	if params.MetadataKey != "" {
		if params.MetadataValue != "" {
			// Filter by both key and value
			query = query.Where("metadata->? = ?", params.MetadataKey, params.MetadataValue)
		} else {
			// Filter by key existence
			query = query.Where("metadata::jsonb ? ?", params.MetadataKey)
		}
	}

	// Apply last seen filter
	if !params.LastSeenSince.IsZero() {
		query = query.Where("last_seen > ?", params.LastSeenSince)
	}

	// Count total matching records
	err := query.Model(&models.Service{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// Execute the query
	err = query.Find(&services).Error
	if err != nil {
		return nil, 0, err
	}

	// Handle dependency filters post-query
	if params.DependencyID != "" || params.DependencyOf != "" {
		// This is more complex, we'll need to filter the results
		filteredServices := []*models.Service{}

		for _, svc := range services {
			if params.DependencyID != "" {
				// Check if this service depends on the specified service
				var depCount int64
				if err := r.db.Model(&models.ServiceDependency{}).
					Where("service_id = ? AND dependency_id = ?", svc.ID, params.DependencyID).
					Count(&depCount).Error; err != nil {
					return nil, 0, err
				}

				if depCount == 0 {
					continue // Skip this service
				}
			}

			if params.DependencyOf != "" {
				// Check if this service is a dependency of the specified service
				var depCount int64
				if err := r.db.Model(&models.ServiceDependency{}).
					Where("service_id = ? AND dependency_id = ?", params.DependencyOf, svc.ID).
					Count(&depCount).Error; err != nil {
					return nil, 0, err
				}

				if depCount == 0 {
					continue // Skip this service
				}
			}

			filteredServices = append(filteredServices, svc)
		}

		services = filteredServices
		count = int64(len(filteredServices))
	}

	// If versions are requested, fetch them (but don't filter by them here)
	if params.IncludeVersions {
		for _, svc := range services {
			var versions []*models.ServiceVersion
			versionQuery := r.db.WithContext(ctx).Where("service_id = ?", svc.ID)

			if params.VersionFilter != "" {
				versionQuery = versionQuery.Where("version = ?", params.VersionFilter)
			}

			if params.ActiveVersionOnly {
				versionQuery = versionQuery.Where("is_active = ?", true)
			}

			if err := versionQuery.Find(&versions).Error; err != nil {
				return nil, 0, err
			}

			// We'll need to add this to the service object
			// For now, we can add it to the metadata (not ideal but works for this example)
			if svc.Metadata == nil {
				svc.Metadata = make(map[string]string)
			}

			// We'd need a proper way to include versions in the response
			// This is just a placeholder approach
			for i, v := range versions {
				svc.Metadata[fmt.Sprintf("version_%d", i)] = v.Version
				svc.Metadata[fmt.Sprintf("version_%d_endpoint", i)] = v.Endpoint
				if v.IsActive {
					svc.Metadata["active_version"] = v.Version
				}
			}
		}
	}

	return services, count, nil
}

func (r *ServiceRepository) CreateVersion(ctx context.Context, version *models.ServiceVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *ServiceRepository) GetVersions(ctx context.Context, serviceID string) ([]*models.ServiceVersion, error) {
	var versions []*models.ServiceVersion
	err := r.db.WithContext(ctx).Where("service_id = ?", serviceID).Find(&versions).Error
	return versions, err
}

func (r *ServiceRepository) GetVersion(ctx context.Context, serviceID string, version string) (*models.ServiceVersion, error) {
	var v models.ServiceVersion
	err := r.db.WithContext(ctx).Where("service_id = ? AND version = ?", serviceID, version).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ServiceRepository) UpdateVersion(ctx context.Context, version *models.ServiceVersion) error {
	return r.db.WithContext(ctx).Save(version).Error
}

func (r *ServiceRepository) DeleteVersion(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ServiceVersion{}, id).Error
}

func (r *ServiceRepository) ActivateVersion(ctx context.Context, serviceID string, version string) error {
	// First, deactivate all versions for this service
	err := r.db.WithContext(ctx).Model(&models.ServiceVersion{}).
		Where("service_id = ?", serviceID).
		Update("is_active", false).Error
	if err != nil {
		return err
	}

	// Then activate the specified version
	return r.db.WithContext(ctx).Model(&models.ServiceVersion{}).
		Where("service_id = ? AND version = ?", serviceID, version).
		Update("is_active", true).Error
}

func (r *ServiceRepository) AddDependency(ctx context.Context, dependency *models.ServiceDependency) error {
	return r.db.WithContext(ctx).Create(dependency).Error
}

func (r *ServiceRepository) GetDependencies(ctx context.Context, serviceID string) ([]*models.ServiceDependency, error) {
	var dependencies []*models.ServiceDependency
	err := r.db.WithContext(ctx).Where("service_id = ?", serviceID).Find(&dependencies).Error
	return dependencies, err
}

func (r *ServiceRepository) GetDependencyOf(ctx context.Context, dependencyID string) ([]*models.ServiceDependency, error) {
	var dependencies []*models.ServiceDependency
	err := r.db.WithContext(ctx).Where("dependency_id = ?", dependencyID).Find(&dependencies).Error
	return dependencies, err
}

func (r *ServiceRepository) RemoveDependency(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ServiceDependency{}, id).Error
}
