// internal/domain/repository/service.go
package repository

import (
	"context"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
	GetByID(ctx context.Context, id string) (*models.Service, error)
	GetByName(ctx context.Context, name string) (*models.Service, error)
	List(ctx context.Context, params models.ServiceQueryParams) ([]*models.Service, int64, error)
	Update(ctx context.Context, service *models.Service) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status models.ServiceStatus) error
	UpdateLastSeen(ctx context.Context, id string) error

	AdvancedSearch(ctx context.Context, params models.AdvancedDiscoveryParams) ([]*models.Service, int64, error)

	// Version management
	CreateVersion(ctx context.Context, version *models.ServiceVersion) error
	GetVersions(ctx context.Context, serviceID string) ([]*models.ServiceVersion, error)
	GetVersion(ctx context.Context, serviceID string, version string) (*models.ServiceVersion, error)
	UpdateVersion(ctx context.Context, version *models.ServiceVersion) error
	DeleteVersion(ctx context.Context, id uint) error
	ActivateVersion(ctx context.Context, serviceID string, version string) error

	// Dependency management
	AddDependency(ctx context.Context, dependency *models.ServiceDependency) error
	GetDependencies(ctx context.Context, serviceID string) ([]*models.ServiceDependency, error)
	GetDependencyOf(ctx context.Context, dependencyID string) ([]*models.ServiceDependency, error)
	RemoveDependency(ctx context.Context, id uint) error
}
