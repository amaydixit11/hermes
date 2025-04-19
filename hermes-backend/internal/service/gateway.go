package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/repository"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/logger"
)

// GatewayService handles business logic for API gateway management
type GatewayService interface {
	// CreateRoute creates a new route
	CreateRoute(ctx context.Context, req *models.RouteCreationRequest) (*models.Route, error)

	// GetRoute retrieves a route by ID
	GetRoute(ctx context.Context, id string) (*models.Route, error)

	// ListRoutes lists all routes with optional filtering
	ListRoutes(ctx context.Context, filters map[string]interface{}) ([]*models.Route, error)

	// UpdateRoute updates an existing route
	UpdateRoute(ctx context.Context, id string, req *models.RouteUpdateRequest) (*models.Route, error)

	// DeleteRoute removes a route
	DeleteRoute(ctx context.Context, id string) error

	// CreateLoadBalancer creates a new load balancer
	CreateLoadBalancer(ctx context.Context, req *models.LoadBalancerCreationRequest) (*models.LoadBalancer, error)

	// GetLoadBalancer retrieves a load balancer by ID
	GetLoadBalancer(ctx context.Context, id string) (*models.LoadBalancer, error)

	// ListLoadBalancers lists all load balancers with optional filtering
	ListLoadBalancers(ctx context.Context, filters map[string]interface{}) ([]*models.LoadBalancer, error)
}

type gatewayService struct {
	routeRepo        repository.RouteRepository
	loadBalancerRepo repository.LoadBalancerRepository
	serviceRepo      repository.ServiceRepository
	log              logger.Logger
}

// NewGatewayService creates a new gateway service
func NewGatewayService(
	routeRepo repository.RouteRepository,
	loadBalancerRepo repository.LoadBalancerRepository,
	serviceRepo repository.ServiceRepository,
	log logger.Logger,
) GatewayService {
	return &gatewayService{
		routeRepo:        routeRepo,
		loadBalancerRepo: loadBalancerRepo,
		serviceRepo:      serviceRepo,
		log:              log,
	}
}

// CreateRoute creates a new route
func (s *gatewayService) CreateRoute(ctx context.Context, req *models.RouteCreationRequest) (*models.Route, error) {
	// Validate service exists
	service, err := s.serviceRepo.Get(ctx, req.ServiceID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, fmt.Errorf("service with ID %s not found", req.ServiceID)
		}
		return nil, fmt.Errorf("failed to check service: %w", err)
	}

	// Validate load balancer if specified
	if req.LoadBalancerID != "" {
		lb, err := s.loadBalancerRepo.Get(ctx, req.LoadBalancerID)
		if err != nil {
			if err == repository.ErrNotFound {
				return nil, fmt.Errorf("load balancer with ID %s not found", req.LoadBalancerID)
			}
			return nil, fmt.Errorf("failed to check load balancer: %w", err)
		}

		// Additional validation for load balancer if needed
		_ = lb // Using the lb variable to avoid unused variable warning
	}

	// Create route
	route := &models.Route{
		ID:             uuid.New().String(),
		Path:           req.Path,
		Description:    req.Description,
		ServiceID:      req.ServiceID,
		LoadBalancerID: req.LoadBalancerID,
		Targets:        req.Targets,
		Active:         true,
		Headers:        req.Headers,
		RateLimit:      req.RateLimit,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// If no targets specified, use service endpoint
	if len(route.Targets) == 0 && service != nil {
		route.Targets = []string{service.Endpoint}
	}

	if err := s.routeRepo.Create(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	return route, nil
}

// GetRoute retrieves a route by ID
func (s *gatewayService) GetRoute(ctx context.Context, id string) (*models.Route, error) {
	route, err := s.routeRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get route: %w", err)
	}

	return route, nil
}

// ListRoutes lists all routes with optional filtering
func (s *gatewayService) ListRoutes(ctx context.Context, filters map[string]interface{}) ([]*models.Route, error) {
	routes, err := s.routeRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}

	return routes, nil
}

// UpdateRoute updates an existing route
func (s *gatewayService) UpdateRoute(ctx context.Context, id string, req *models.RouteUpdateRequest) (*models.Route, error) {
	// Check if route exists
	_, err := s.routeRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get route: %w", err)
	}

	// Validate service if specified
	if req.ServiceID != nil {
		_, err := s.serviceRepo.Get(ctx, *req.ServiceID)
		if err != nil {
			if err == repository.ErrNotFound {
				return nil, fmt.Errorf("service with ID %s not found", *req.ServiceID)
			}
			return nil, fmt.Errorf("failed to check service: %w", err)
		}
	}

	// Validate load balancer if specified
	if req.LoadBalancerID != nil {
		_, err := s.loadBalancerRepo.Get(ctx, *req.LoadBalancerID)
		if err != nil {
			if err == repository.ErrNotFound {
				return nil, fmt.Errorf("load balancer with ID %s not found", *req.LoadBalancerID)
			}
			return nil, fmt.Errorf("failed to check load balancer: %w", err)
		}
	}

	// Update route
	if err := s.routeRepo.Update(ctx, id, req); err != nil {
		return nil, fmt.Errorf("failed to update route: %w", err)
	}

	// Get updated route
	route, err := s.routeRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated route: %w", err)
	}

	return route, nil
}

// DeleteRoute removes a route
func (s *gatewayService) DeleteRoute(ctx context.Context, id string) error {
	if err := s.routeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}

	return nil
}

// Additional methods would be implemented here...
