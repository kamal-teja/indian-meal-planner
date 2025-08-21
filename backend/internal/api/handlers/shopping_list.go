package handlers

import (
	"net/http"
	"time"

	"nourish-backend/internal/api/middleware"
	"nourish-backend/internal/models"
	"nourish-backend/internal/service"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ShoppingListHandler handles shopping list requests
type ShoppingListHandler struct {
	mealService service.MealService
	logger      *logger.Logger
}

// NewShoppingListHandler creates a new shopping list handler
func NewShoppingListHandler(mealService service.MealService, log *logger.Logger) *ShoppingListHandler {
	return &ShoppingListHandler{
		mealService: mealService,
		logger:      log,
	}
}

// GetShoppingList handles GET /api/shopping-list
func (h *ShoppingListHandler) GetShoppingList(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Parse date parameters
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "startDate and endDate parameters are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid startDate format. Use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid endDate format. Use YYYY-MM-DD",
		})
		return
	}

	shoppingList, err := h.mealService.GetShoppingList(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    shoppingList,
	})
}
