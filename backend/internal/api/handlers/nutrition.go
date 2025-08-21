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

// NutritionHandler handles nutrition-related requests
type NutritionHandler struct {
	mealService service.MealService
	userService service.UserService
	logger      *logger.Logger
}

// NewNutritionHandler creates a new nutrition handler
func NewNutritionHandler(mealService service.MealService, userService service.UserService, log *logger.Logger) *NutritionHandler {
	return &NutritionHandler{
		mealService: mealService,
		userService: userService,
		logger:      log,
	}
}

// GetNutritionProgress handles GET /api/nutrition/progress
func (h *NutritionHandler) GetNutritionProgress(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Parse period parameter (default to 7 days)
	periodStr := c.DefaultQuery("period", "7")
	period, err := strconv.Atoi(periodStr)
	if err != nil || period <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid period parameter",
		})
		return
	}

	progress, err := h.mealService.GetNutritionProgress(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    progress,
	})
}

// GetNutritionGoals handles GET /api/nutrition/goals
func (h *NutritionHandler) GetNutritionGoals(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Get the user to access their nutrition goals from profile
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user for nutrition goals", "error", err, "userID", userID.Hex())
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// If user has no nutrition goals set, return default values
	goals := user.Profile.NutritionGoals
	if goals.DailyCalories == 0 {
		// Set default nutrition goals if none are set
		goals = models.NutritionGoals{
			DailyCalories: 2000,
			Protein:       150,
			Carbs:         250,
			Fat:           65,
			Fiber:         25,
			Sodium:        2300,
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    goals,
	})
}

// UpdateNutritionGoals handles PUT /api/nutrition/goals
func (h *NutritionHandler) UpdateNutritionGoals(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	var req models.NutritionGoalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// Get the current user to update their profile
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user for nutrition goals update", "error", err, "userID", userID.Hex())
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Update the nutrition goals in the user's profile
	goals := models.NutritionGoals{
		DailyCalories: req.DailyCalories,
		Protein:       req.DailyProtein,
		Carbs:         req.DailyCarbs,
		Fat:           req.DailyFat,
		Fiber:         req.DailyFiber,
		Sodium:        req.DailySodium,
	}

	// Update the user's profile with new nutrition goals
	user.Profile.NutritionGoals = goals

	// Save the updated profile using the user service
	err = h.userService.UpdateProfile(c.Request.Context(), userID, user.Profile)
	if err != nil {
		h.logger.Error("Failed to update user profile with nutrition goals", "error", err, "userID", userID.Hex())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to update nutrition goals",
		})
		return
	}

	h.logger.Info("Successfully updated nutrition goals", "userID", userID.Hex())

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    goals,
	})
}
