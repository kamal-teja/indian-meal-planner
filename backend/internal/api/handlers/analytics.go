package handlers

import (
	"net/http"
	"strconv"

	"meal-planner-backend/internal/api/middleware"
	"meal-planner-backend/internal/models"
	"meal-planner-backend/internal/service"
	"meal-planner-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles analytics-related requests
type AnalyticsHandler struct {
	mealService service.MealService
	logger      *logger.Logger
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(mealService service.MealService, log *logger.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		mealService: mealService,
		logger:      log,
	}
}

// GetMealAnalytics handles GET /api/analytics/meals
func (h *AnalyticsHandler) GetMealAnalytics(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Parse period parameter (default 30 days)
	period, err := strconv.Atoi(c.DefaultQuery("period", "30"))
	if err != nil || period < 1 || period > 365 {
		period = 30
	}

	analytics, err := h.mealService.GetAnalytics(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    analytics,
	})
}
