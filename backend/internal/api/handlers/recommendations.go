package handlers

import (
	"net/http"
	"time"

	"meal-planner-backend/internal/api/middleware"
	"meal-planner-backend/internal/models"
	"meal-planner-backend/internal/service"
	"meal-planner-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RecommendationsHandler handles recommendation requests
type RecommendationsHandler struct {
	dishService service.DishService
	mealService service.MealService
	userService service.UserService
	logger      *logger.Logger
}

// NewRecommendationsHandler creates a new recommendations handler
func NewRecommendationsHandler(dishService service.DishService, mealService service.MealService, userService service.UserService, log *logger.Logger) *RecommendationsHandler {
	return &RecommendationsHandler{
		dishService: dishService,
		mealService: mealService,
		userService: userService,
		logger:      log,
	}
}

// GetRecommendations handles GET /api/recommendations
func (h *RecommendationsHandler) GetRecommendations(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	mealType := c.Query("mealType")
	dateStr := c.Query("date")

	if mealType == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "mealType and date parameters are required",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	recommendations, err := h.mealService.GetRecommendations(c.Request.Context(), userID, mealType, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    recommendations,
	})
}
