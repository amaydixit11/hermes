// internal/service/service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/repository"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/errors"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
	"github.com/google/uuid"
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

	var registeredBy string
	if reg.RegisteredBy == "" {
		registeredBy = "self"
	} else {
		registeredBy = reg.RegisteredBy
	}

	// Create new service
	service := &models.Service{
		ID:           "svc-" + uuid.New().String()[:8],
		Name:         reg.Name,
		Description:  reg.Description,
		Status:       models.ServiceStatusUnknown,
		Type:         reg.Type,
		Endpoint:     reg.Endpoint,
		Metadata:     reg.Metadata,
		Tags:         reg.Tags,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastSeen:     time.Now(),
		RegisteredBy: registeredBy,
	}

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, errors.Wrap(err, "failed to create service")
	}

	// TODO: Trigger initial health check (async)
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

// AdvancedDiscovery provides advanced service discovery capabilities
func (s *ServiceService) AdvancedDiscovery(ctx context.Context, params models.AdvancedDiscoveryParams) ([]*models.Service, int64, error) {
	services, total, err := s.repo.AdvancedSearch(ctx, params)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to perform advanced service discovery")
	}

	// If dependency details are requested, we can enhance the response
	if params.DependencyID != "" || params.DependencyOf != "" {
		for _, service := range services {
			// Fetch and enrich with dependency information if needed
			dependencies, err := s.repo.GetDependencies(ctx, service.ID)
			if err != nil {
				s.log.Error("Failed to fetch dependencies", "error", err, "serviceID", service.ID)
				continue
			}

			// Add dependency info to metadata
			if service.Metadata == nil {
				service.Metadata = make(map[string]string)
			}

			for i, dep := range dependencies {
				service.Metadata[fmt.Sprintf("dependency_%d", i)] = dep.DependencyID
				service.Metadata[fmt.Sprintf("dependency_%d_type", i)] = dep.DependencyType
			}
		}
	}

	return services, total, nil
}

// AddServiceVersion adds a new version to a service
func (s *ServiceService) AddServiceVersion(ctx context.Context, serviceID string, versionReq models.ServiceVersionRequest) (*models.ServiceVersion, error) {
	// Check if service exists
	_, err := s.repo.GetByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, errors.Wrap(err, "failed to retrieve service")
	}

	// Check if version already exists
	existingVersion, err := s.repo.GetVersion(ctx, serviceID, versionReq.Version)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "failed to check existing version")
	}

	if existingVersion != nil {
		return nil, errors.New("version already exists for this service")
	}

	// Create new version
	version := &models.ServiceVersion{
		ServiceID:   serviceID,
		Version:     versionReq.Version,
		IsActive:    versionReq.IsActive,
		Endpoint:    versionReq.Endpoint,
		Description: versionReq.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// If this is the first version or marked as active, ensure it's the only active one
	if versionReq.IsActive {
		// Deactivate all other versions first
		versions, err := s.repo.GetVersions(ctx, serviceID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get existing versions")
		}

		if len(versions) == 0 {
			// This is the first version, make it active by default
			version.IsActive = true
		} else if versionReq.IsActive {
			// We need to deactivate other versions
			if err := s.repo.ActivateVersion(ctx, serviceID, versionReq.Version); err != nil {
				return nil, errors.Wrap(err, "failed to activate version")
			}
		}
	}

	if err := s.repo.CreateVersion(ctx, version); err != nil {
		return nil, errors.Wrap(err, "failed to create service version")
	}

	s.log.Info("Service version added", "serviceID", serviceID, "version", version.Version)
	return version, nil
}

// GetServiceVersions retrieves all versions of a service
func (s *ServiceService) GetServiceVersions(ctx context.Context, serviceID string) ([]*models.ServiceVersion, error) {
	// Check if service exists
	service, err := s.repo.GetByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, errors.Wrap(err, "failed to retrieve service")
	}

	versions, err := s.repo.GetVersions(ctx, service.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve service versions")
	}

	return versions, nil
}

// ActivateServiceVersion sets a specific version as the active one
func (s *ServiceService) ActivateServiceVersion(ctx context.Context, serviceID string, version string) error {
	// Check if service exists
	_, service_err := s.repo.GetByID(ctx, serviceID)
	if service_err != nil {
		if errors.Is(service_err, gorm.ErrRecordNotFound) {
			return errors.New("service not found")
		}
		return errors.Wrap(service_err, "failed to retrieve service")
	}

	// Check if version exists
	_, ver_err := s.repo.GetVersion(ctx, serviceID, version)
	if ver_err != nil {
		if errors.Is(ver_err, gorm.ErrRecordNotFound) {
			return errors.New("version not found")
		}
		return errors.Wrap(ver_err, "failed to retrieve version")
	}

	// Activate the version (repository handles deactivating others)
	if err := s.repo.ActivateVersion(ctx, serviceID, version); err != nil {
		return errors.Wrap(err, "failed to activate version")
	}

	s.log.Info("Service version activated", "serviceID", serviceID, "version", version)
	return nil
}

// AddServiceDependency creates a dependency relationship between services
func (s *ServiceService) AddServiceDependency(ctx context.Context, serviceID string, depReq models.ServiceDependencyRequest) (*models.ServiceDependency, error) {
	// Check if source service exists
	_, source_err := s.repo.GetByID(ctx, serviceID)
	if source_err != nil {
		if errors.Is(source_err, gorm.ErrRecordNotFound) {
			return nil, errors.New("source service not found")
		}
		return nil, errors.Wrap(source_err, "failed to retrieve source service")
	}

	// Check if dependency service exists
	_, dep_err := s.repo.GetByID(ctx, depReq.DependencyID)
	if dep_err != nil {
		if errors.Is(dep_err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dependency service not found")
		}
		return nil, errors.Wrap(dep_err, "failed to retrieve dependency service")
	}

	// Validate dependency type
	validTypes := []string{"REQUIRED", "OPTIONAL"}
	validType := false
	for _, t := range validTypes {
		if depReq.DependencyType == t {
			validType = true
			break
		}
	}

	if !validType {
		return nil, errors.New("invalid dependency type, must be one of: REQUIRED, OPTIONAL")
	}

	// Check for circular dependency
	// Simple case: A depends on B, B depends on A
	dependencies, err := s.repo.GetDependencies(ctx, depReq.DependencyID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check for circular dependencies")
	}

	for _, dep := range dependencies {
		if dep.DependencyID == serviceID {
			return nil, errors.New("circular dependency detected")
		}
	}

	// Create the dependency
	dependency := &models.ServiceDependency{
		ServiceID:      serviceID,
		DependencyID:   depReq.DependencyID,
		DependencyType: depReq.DependencyType,
		Description:    depReq.Description,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.AddDependency(ctx, dependency); err != nil {
		return nil, errors.Wrap(err, "failed to add service dependency")
	}

	s.log.Info("Service dependency added",
		"serviceID", serviceID,
		"dependencyID", depReq.DependencyID,
		"type", depReq.DependencyType)

	return dependency, nil
}

// GetServiceDependencies retrieves all dependencies of a service
func (s *ServiceService) GetServiceDependencies(ctx context.Context, serviceID string) ([]*models.ServiceDependency, error) {
	// Check if service exists
	service, err := s.repo.GetByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, errors.Wrap(err, "failed to retrieve service")
	}

	dependencies, err := s.repo.GetDependencies(ctx, service.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve service dependencies")
	}

	return dependencies, nil
}

// GetServiceDependents retrieves all services that depend on this service
func (s *ServiceService) GetServiceDependents(ctx context.Context, serviceID string) ([]*models.ServiceDependency, error) {
	// Check if service exists
	service, err := s.repo.GetByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, errors.Wrap(err, "failed to retrieve service")
	}

	dependents, err := s.repo.GetDependencyOf(ctx, service.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve service dependents")
	}

	return dependents, nil
}

// RemoveServiceDependency removes a dependency relationship
func (s *ServiceService) RemoveServiceDependency(ctx context.Context, serviceID string, dependencyID string) error {
	// Find the dependency record
	dependencies, err := s.repo.GetDependencies(ctx, serviceID)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve dependencies")
	}

	var depID uint
	found := false

	for _, dep := range dependencies {
		if dep.DependencyID == dependencyID {
			depID = dep.ID
			found = true
			break
		}
	}

	if !found {
		return errors.New("dependency relationship not found")
	}

	if err := s.repo.RemoveDependency(ctx, depID); err != nil {
		return errors.Wrap(err, "failed to remove dependency")
	}

	s.log.Info("Service dependency removed", "serviceID", serviceID, "dependencyID", dependencyID)
	return nil
}
