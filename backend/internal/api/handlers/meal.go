package handlers

import (
	"net/http"
	"strconv"
	"time"

	"nourish-backend/internal/api/middleware"
	"nourish-backend/internal/models"
	"nourish-backend/internal/service"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MealHandler handles meal-related requests
type MealHandler struct {
	mealService service.MealService
	validator   *validator.Validate
	logger      *logger.Logger
}

// NewMealHandler creates a new meal handler
func NewMealHandler(mealService service.MealService, log *logger.Logger) *MealHandler {
	return &MealHandler{
		mealService: mealService,
		validator:   validator.New(),
		logger:      log,
	}
}

// CreateMeal handles POST /api/meals
func (h *MealHandler) CreateMeal(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	var req models.MealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to parse meal request",
			"error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid request format - check required fields and date format",
			Details: err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Meal request validation failed",
			"error", err.Error(),
			"request", req)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Validation failed - check required fields and valid values",
			Details: err.Error(),
		})
		return
	}

	meal, err := h.mealService.Create(c.Request.Context(), userID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "dish not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid dish ID" {
			status = http.StatusBadRequest
		}

		c.JSON(status, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "Meal created successfully",
		Data:    meal,
	})
}

// GetMeals handles GET /api/meals
func (h *MealHandler) GetMeals(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Check if date range is provided
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr != "" && endDateStr != "" {
		// Handle date range query
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Invalid start date format. Use YYYY-MM-DD",
			})
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Invalid end date format. Use YYYY-MM-DD",
			})
			return
		}

		meals, err := h.mealService.GetByDateRange(c.Request.Context(), userID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
			Data:    meals,
		})
		return
	}

	// Handle pagination query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	meals, pagination, err := h.mealService.GetByUserID(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response := gin.H{
		"success":    true,
		"meals":      meals,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
}

// GetMeal handles GET /api/meals/:id (can be ObjectID or date in YYYY-MM-DD format)
func (h *MealHandler) GetMeal(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	idParam := c.Param("id")

	// Try to parse as date first (YYYY-MM-DD format)
	if date, err := time.Parse("2006-01-02", idParam); err == nil {
		// Handle as date - get meals for the specific date
		startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

		meals, err := h.mealService.GetByDateRange(c.Request.Context(), userID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
			Data:    meals,
		})
		return
	}

	// Try to parse as ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid meal ID or date format. Use ObjectID or YYYY-MM-DD",
		})
		return
	}

	meal, err := h.mealService.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "meal not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    meal,
	})
}

// UpdateMeal handles PUT /api/meals/:id
func (h *MealHandler) UpdateMeal(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid meal ID",
		})
		return
	}

	var req models.MealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	meal, err := h.mealService.Update(c.Request.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "meal not found" || err.Error() == "dish not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid dish ID" {
			status = http.StatusBadRequest
		}

		c.JSON(status, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Meal updated successfully",
		Data:    meal,
	})
}

// DeleteMeal handles DELETE /api/meals/:id
func (h *MealHandler) DeleteMeal(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid meal ID",
		})
		return
	}

	// Soft-delete the meal by setting DeletedAt
	if err := h.mealService.SoftDeleteByIDs(c.Request.Context(), []primitive.ObjectID{id}); err != nil {
		h.logger.Error("Failed to soft-delete meal", "error", err, "mealID", id.Hex())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "failed to delete meal",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Meal deleted successfully",
	})
}

// DeleteMealsBulk handles DELETE /api/meals with body { ids: ["id1","id2"] }
func (h *MealHandler) DeleteMealsBulk(c *gin.Context) {
	// Auth is enforced by middleware; proceed to parse IDs from body

	var body struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Error("Failed to parse bulk delete request", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid request format"})
		return
	}

	var objIDs []primitive.ObjectID
	for _, idStr := range body.IDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			h.logger.Error("Invalid meal ID in bulk delete", "id", idStr)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid meal ID in request"})
			return
		}
		objIDs = append(objIDs, id)
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Success: false, Error: "Authentication required"})
		return
	}

	// Soft-delete the provided IDs and create undo token
	token, err := h.mealService.CreateUndoableSoftDelete(c.Request.Context(), userID, objIDs, 5*time.Minute)
	if err != nil {
		h.logger.Error("Failed to soft-delete meals bulk", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to delete meals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Meals deleted successfully",
		"undoToken": token,
	})
}

// DeleteByDateAndDish handles DELETE /api/meals/by-date-dish with body { date: "YYYY-MM-DD", dishId: "..." }
func (h *MealHandler) DeleteByDateAndDish(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Success: false, Error: "Authentication required"})
		return
	}

	var body struct {
		Date   string `json:"date"`
		DishID string `json:"dishId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Error("Failed to parse delete-by-date-dish request", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid request format"})
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid date format"})
		return
	}

	dishObjID, err := primitive.ObjectIDFromHex(body.DishID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid dish ID"})
		return
	}

	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	if err := h.mealService.DeleteByUserDateAndDish(c.Request.Context(), userID, startDate, endDate, dishObjID); err != nil {
		h.logger.Error("Failed to delete meals by date and dish", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to delete meals"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true, Message: "Meals deleted successfully"})
}

// UndoDeleteByDateAndDish handles POST /api/meals/undo-by-date-dish with body { date: "YYYY-MM-DD", dishId: "..." }
func (h *MealHandler) UndoDeleteByDateAndDish(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Success: false, Error: "Authentication required"})
		return
	}

	var body struct {
		Date   string `json:"date"`
		DishID string `json:"dishId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Error("Failed to parse undo-delete request", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid request format"})
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid date format"})
		return
	}

	dishObjID, err := primitive.ObjectIDFromHex(body.DishID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid dish ID"})
		return
	}

	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	if err := h.mealService.UndoDeleteByUserDateAndDish(c.Request.Context(), userID, startDate, endDate, dishObjID); err != nil {
		h.logger.Error("Failed to undo delete by date and dish", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to undo delete"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true, Message: "Undo successful"})
}

// UndoByToken handles POST /api/meals/undo with body { token: "..." }
func (h *MealHandler) UndoByToken(c *gin.Context) {
	_, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Success: false, Error: "Authentication required"})
		return
	}

	var body struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Error("Failed to parse undo-by-token request", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid request format"})
		return
	}

	if body.Token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "token is required"})
		return
	}

	if err := h.mealService.UndoByToken(c.Request.Context(), body.Token); err != nil {
		h.logger.Error("Failed to undo by token", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to undo"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true, Message: "Undo successful"})
}

// GetNutritionSummary handles GET /api/meals/nutrition-summary
func (h *MealHandler) GetNutritionSummary(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "startDate and endDate query parameters are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid start date format. Use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid end date format. Use YYYY-MM-DD",
		})
		return
	}

	summary, err := h.mealService.GetNutritionSummary(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    summary,
	})
}

// GetMealsByMonth handles GET /api/meals/month/:year/:month
func (h *MealHandler) GetMealsByMonth(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Parse year and month parameters
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid year parameter",
		})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid month parameter (1-12)",
		})
		return
	}

	// Calculate start and end dates for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond) // Last nanosecond of the month

	// Get meals for the date range
	meals, err := h.mealService.GetByDateRange(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get meals for month",
			"error", err.Error(),
			"year", year,
			"month", month,
			"userID", userID.Hex())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to retrieve meals for the month",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    meals,
	})
}
