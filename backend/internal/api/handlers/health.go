package handlers

import (
	"net/http"
	"time"

	"meal-planner-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// CheckHealth returns the health status of the API
func (h *HealthHandler) CheckHealth(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
		Database:  "connected",
	}

	c.JSON(http.StatusOK, response)
}
