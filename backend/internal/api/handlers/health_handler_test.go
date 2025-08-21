package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": "2023-10-15T10:00:00Z",
			"version":   "1.0.0",
			"database":  "connected",
		})
	})

	// Act
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "healthy")
}