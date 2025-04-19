// internal/domain/repository/service.go
package repository

import (
	"context"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
)

// ServiceRepository defines the interface for service data access
type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
	GetByID(ctx context.Context, id string) (*models.Service, error)
	GetByName(ctx context.Context, name string) (*models.Service, error)
	List(ctx context.Context, params models.ServiceQueryParams) ([]*models.Service, int64, error)
	Update(ctx context.Context, service *models.Service) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status models.ServiceStatus) error
	UpdateLastSeen(ctx context.Context, id string) error
}
