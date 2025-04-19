// internal/service/service.go
package service

import (
	"context"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/repository"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/errors"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
	"gorm.io/gorm"
)

// ServiceService handles business logic for services
type ServiceService struct {
	repo repository.ServiceRepository
	log  *logger.Logger
}

// NewServiceService creates a new ServiceService
func NewServiceService(repo repository.ServiceRepository, log *logger.Logger) *ServiceService {
	return &ServiceService{
		repo: repo,
		log:  log,
	}
}

// RegisterService handles the registration of a new service
func (s *ServiceService) RegisterService(ctx context.Context, reg models.ServiceRegistration) (*models.Service, error) {
	// Check if service with the same name already exists
	existing, err := s.repo.GetByName(ctx, reg.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "failed to check for existing service")
	}

	if existing != nil {
		return nil, errors.New("service with this name already exists")
	}
	// tags := pq.Array(reg.Tags)

	// Create new service
	service := &models.Service{
		Name:        reg.Name,
		Description: reg.Description,
		Status:      models.ServiceStatusUnknown,
		Type:        reg.Type,
		Endpoint:    reg.Endpoint,
		Metadata:    reg.Metadata,
		Tags:        reg.Tags,
		LastSeen:    time.Now(),
	}

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, errors.Wrap(err, "failed to create service")
	}

	s.log.Info("Service registered", "id", service.ID, "name", service.Name)
	return service, nil
}

// GetServiceByID retrieves a service by its ID
func (s *ServiceService) GetServiceByID(ctx context.Context, id string) (*models.Service, error) {
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Service not found
		}
		return nil, errors.Wrap(err, "failed to retrieve service")
	}
	return service, nil
}

// GetServiceByName retrieves a service by its name
func (s *ServiceService) GetServiceByName(ctx context.Context, name string) (*models.Service, error) {
	service, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Service not found
		}
		return nil, errors.Wrap(err, "failed to retrieve service")
	}
	return service, nil
}

// ListServices lists all services with optional filters and pagination
func (s *ServiceService) ListServices(ctx context.Context, params models.ServiceQueryParams) ([]*models.Service, int64, error) {
	services, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list services")
	}
	return services, total, nil
}

// UpdateService updates an existing service
func (s *ServiceService) UpdateService(ctx context.Context, id string, update models.ServiceUpdateRequest) (*models.Service, error) {
	// Check if service exists
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, errors.Wrap(err, "failed to retrieve service")
	}

	// Update service fields
	if update.Description != nil {
		service.Description = *update.Description
	}
	if update.Endpoint != nil {
		service.Endpoint = *update.Endpoint
	}
	if update.Metadata != nil {
		service.Metadata = update.Metadata
	}
	if update.Tags != nil {
		service.Tags = update.Tags
	}

	if err := s.repo.Update(ctx, service); err != nil {
		return nil, errors.Wrap(err, "failed to update service")
	}

	s.log.Info("Service updated", "id", service.ID, "name", service.Name)
	return service, nil
}

// DeleteService deletes a service by its ID
func (s *ServiceService) DeleteService(ctx context.Context, id string) error {
	// Check if service exists
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("service not found")
		}
		return errors.Wrap(err, "failed to retrieve service")
	}

	if err := s.repo.Delete(ctx, service.ID); err != nil {
		return errors.Wrap(err, "failed to delete service")
	}

	s.log.Info("Service deleted", "id", service.ID, "name", service.Name)
	return nil
}

// UpdateServiceStatus updates the health status of a service
func (s *ServiceService) UpdateServiceStatus(ctx context.Context, id string, status models.ServiceStatus) error {
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("service not found")
		}
		return errors.Wrap(err, "failed to retrieve service")
	}

	service.Status = status
	if err := s.repo.Update(ctx, service); err != nil {
		return errors.Wrap(err, "failed to update service status")
	}

	s.log.Info("Service status updated", "id", service.ID, "name", service.Name, "status", status)
	return nil
}

// UpdateServiceLastSeen updates the LastSeen timestamp of a service
func (s *ServiceService) UpdateServiceLastSeen(ctx context.Context, id string) error {
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("service not found")
		}
		return errors.Wrap(err, "failed to retrieve service")
	}

	service.LastSeen = time.Now()
	if err := s.repo.Update(ctx, service); err != nil {
		return errors.Wrap(err, "failed to update service last seen")
	}

	s.log.Info("Service LastSeen updated", "id", service.ID, "name", service.Name)
	return nil
}
