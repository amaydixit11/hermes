// internal/api/handlers/service_discovery.go
package handlers

import (
	"net/http"
	"time"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// ServiceDiscoveryHandler handles service discovery HTTP requests
type ServiceDiscoveryHandler struct {
	service *service.ServiceService
}

// NewServiceDiscoveryHandler creates a new ServiceDiscoveryHandler
func NewServiceDiscoveryHandler(service *service.ServiceService) *ServiceDiscoveryHandler {
	return &ServiceDiscoveryHandler{
		service: service,
	}
}

// AdvancedSearch handles advanced service discovery requests
func (h *ServiceDiscoveryHandler) AdvancedSearch(c *gin.Context) {
	var params models.AdvancedDiscoveryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse any time-based parameters
	if lastSeenStr := c.Query("last_seen_since"); lastSeenStr != "" {
		lastSeen, err := time.Parse(time.RFC3339, lastSeenStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid last_seen_since format. Use RFC3339 format."})
			return
		}
		params.LastSeenSince = lastSeen
	}

	services, total, err := h.service.AdvancedDiscovery(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    total,
		"params":   params, // Include the used parameters for debugging/clarity
	})
}

// ServiceVersionHandler handles service versioning HTTP requests
type ServiceVersionHandler struct {
	service *service.ServiceService
}

// NewServiceVersionHandler creates a new ServiceVersionHandler
func NewServiceVersionHandler(service *service.ServiceService) *ServiceVersionHandler {
	return &ServiceVersionHandler{
		service: service,
	}
}

// AddServiceVersion handles requests to add a new service version
func (h *ServiceVersionHandler) AddServiceVersion(c *gin.Context) {
	serviceID := c.Param("id")

	var versionReq models.ServiceVersionRequest
	if err := c.ShouldBindJSON(&versionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	version, err := h.service.AddServiceVersion(c.Request.Context(), serviceID, versionReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, version)
}

// GetServiceVersions handles requests to list all versions of a service
func (h *ServiceVersionHandler) GetServiceVersions(c *gin.Context) {
	serviceID := c.Param("id")

	versions, err := h.service.GetServiceVersions(c.Request.Context(), serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"versions": versions,
		"count":    len(versions),
	})
}

// ActivateServiceVersion handles requests to make a version active
func (h *ServiceVersionHandler) ActivateServiceVersion(c *gin.Context) {
	serviceID := c.Param("id")
	version := c.Param("version")

	err := h.service.ActivateServiceVersion(c.Request.Context(), serviceID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Version activated successfully",
		"service_id": serviceID,
		"version":    version,
	})
}

// ServiceDependencyHandler handles service dependency HTTP requests
type ServiceDependencyHandler struct {
	service *service.ServiceService
}

// NewServiceDependencyHandler creates a new ServiceDependencyHandler
func NewServiceDependencyHandler(service *service.ServiceService) *ServiceDependencyHandler {
	return &ServiceDependencyHandler{
		service: service,
	}
}

// AddServiceDependency handles requests to add a new service dependency
func (h *ServiceDependencyHandler) AddServiceDependency(c *gin.Context) {
	serviceID := c.Param("id")

	var depReq models.ServiceDependencyRequest
	if err := c.ShouldBindJSON(&depReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dependency, err := h.service.AddServiceDependency(c.Request.Context(), serviceID, depReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dependency)
}

// GetServiceDependencies handles requests to list all dependencies of a service
func (h *ServiceDependencyHandler) GetServiceDependencies(c *gin.Context) {
	serviceID := c.Param("id")

	dependencies, err := h.service.GetServiceDependencies(c.Request.Context(), serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enhance the response with dependency service details
	type enhancedDependency struct {
		*models.ServiceDependency
		DependencyName   string               `json:"dependency_name"`
		DependencyStatus models.ServiceStatus `json:"dependency_status"`
	}

	enhancedDeps := make([]enhancedDependency, 0, len(dependencies))

	for _, dep := range dependencies {
		depService, err := h.service.GetServiceByID(c.Request.Context(), dep.DependencyID)
		if err != nil {
			continue // Skip this one if we can't get details
		}

		enhanced := enhancedDependency{
			ServiceDependency: dep,
			DependencyName:    depService.Name,
			DependencyStatus:  depService.Status,
		}

		enhancedDeps = append(enhancedDeps, enhanced)
	}

	c.JSON(http.StatusOK, gin.H{
		"dependencies": enhancedDeps,
		"count":        len(enhancedDeps),
	})
}

// GetServiceDependents handler completion
func (h *ServiceDependencyHandler) GetServiceDependents(c *gin.Context) {
	serviceID := c.Param("id")

	dependents, err := h.service.GetServiceDependents(c.Request.Context(), serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enhance the response with dependent service details
	type enhancedDependent struct {
		*models.ServiceDependency
		ServiceName   string               `json:"service_name"`
		ServiceStatus models.ServiceStatus `json:"service_status"`
	}

	enhancedDeps := make([]enhancedDependent, 0, len(dependents))

	for _, dep := range dependents {
		depService, err := h.service.GetServiceByID(c.Request.Context(), dep.ServiceID)
		if err != nil {
			continue // Skip this one if we can't get details
		}

		enhanced := enhancedDependent{
			ServiceDependency: dep,
			ServiceName:       depService.Name,
			ServiceStatus:     depService.Status,
		}

		enhancedDeps = append(enhancedDeps, enhanced)
	}

	c.JSON(http.StatusOK, gin.H{
		"dependents": enhancedDeps,
		"count":      len(enhancedDeps),
	})
}

// RemoveServiceDependency handles requests to remove a dependency relationship
func (h *ServiceDependencyHandler) RemoveServiceDependency(c *gin.Context) {
	serviceID := c.Param("id")
	dependencyID := c.Param("dependency_id")

	err := h.service.RemoveServiceDependency(c.Request.Context(), serviceID, dependencyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Dependency removed successfully",
		"service_id":    serviceID,
		"dependency_id": dependencyID,
	})
}
