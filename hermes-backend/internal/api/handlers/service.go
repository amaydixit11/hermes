// internal/api/handlers/service.go
package handlers

import (
	"net/http"

	"github.com/amaydixit11/hermes/hermes-backend/internal/domain/models"
	"github.com/amaydixit11/hermes/hermes-backend/internal/service"
	"github.com/amaydixit11/hermes/hermes-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

// ServiceHandler handles HTTP requests for services
type ServiceHandler struct {
	service *service.ServiceService
}

// NewServiceHandler creates a new ServiceHandler
func NewServiceHandler(service *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{
		service: service,
	}
}

// RegisterService handles service registration requests
func (h *ServiceHandler) RegisterService(c *gin.Context) {
	var registration models.ServiceRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Add validation for registration fields

	service, err := h.service.RegisterService(c.Request.Context(), registration)
	if err != nil {
		if errors.Is(err, errors.New("service with this name already exists")) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) BulkRegisterService(c *gin.Context) {
	var bulkRegistration models.BulkServiceRegistration
	if err := c.ShouldBindJSON(&bulkRegistration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bulk registration data"})
		return
	}

	type SuccessResult struct {
		ID     string               `json:"id"`
		Name   string               `json:"name"`
		Status models.ServiceStatus `json:"status"`
	}

	type FailureResult struct {
		Name  string `json:"name"`
		Error string `json:"error"`
	}

	successful := []SuccessResult{}
	failed := []FailureResult{}

	for _, registration := range bulkRegistration.Services {
		// Merge common metadata
		if registration.Metadata == nil {
			registration.Metadata = map[string]string{}
		}
		for k, v := range bulkRegistration.CommonMetadata {
			if _, exists := registration.Metadata[k]; !exists {
				registration.Metadata[k] = v
			}
		}

		// Merge common tags
		tagMap := map[string]bool{}
		for _, tag := range registration.Tags {
			tagMap[tag] = true
		}
		for _, tag := range bulkRegistration.CommonTags {
			tagMap[tag] = true
		}
		registration.Tags = make([]string, 0, len(tagMap))
		for tag := range tagMap {
			registration.Tags = append(registration.Tags, tag)
		}

		// Register service
		service, err := h.service.RegisterService(c.Request.Context(), registration)
		if err != nil {
			failed = append(failed, FailureResult{
				Name:  registration.Name,
				Error: err.Error(),
			})
			continue
		}

		successful = append(successful, SuccessResult{
			ID:     service.ID,
			Name:   service.Name,
			Status: service.Status,
		})
	}

	c.JSON(http.StatusMultiStatus, gin.H{
		"successful":      successful,
		"failed":          failed,
		"total_submitted": len(bulkRegistration.Services),
		"total_succeeded": len(successful),
		"total_failed":    len(failed),
	})

}

// GetServiceByID handles requests to get a service by ID
func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	id := c.Param("id")
	service, err := h.service.GetServiceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service"})
		return
	}
	if service == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

// GetServiceByName handles requests to get a service by its name
func (h *ServiceHandler) GetServiceByName(c *gin.Context) {
	name := c.Param("name")
	service, err := h.service.GetServiceByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service"})
		return
	}
	if service == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

// ListServices handles requests to list services with filters and pagination
func (h *ServiceHandler) ListServices(c *gin.Context) {
	var params models.ServiceQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services, total, err := h.service.ListServices(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    total,
	})
}

// UpdateService handles requests to update an existing service
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	var updateRequest models.ServiceUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	service, err := h.service.UpdateService(c.Request.Context(), id, updateRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}
	if service == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

// DeleteService handles requests to delete a service by ID
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteService(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// UpdateServiceStatus handles requests to update the status of a service
func (h *ServiceHandler) UpdateServiceStatus(c *gin.Context) {
	id := c.Param("id")
	var statusRequest models.ServiceHealthUpdate
	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateServiceStatus(c.Request.Context(), id, statusRequest.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service status updated successfully"})
}

// UpdateServiceLastSeen handles requests to update the LastSeen time of a service
func (h *ServiceHandler) UpdateServiceLastSeen(c *gin.Context) {
	id := c.Param("id")

	err := h.service.UpdateServiceLastSeen(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service last seen"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service last seen updated successfully"})
}
