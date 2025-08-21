package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"nourish-backend/internal/api/middleware"
	"nourish-backend/internal/models"
	"nourish-backend/internal/service"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DishHandler handles dish-related requests
type DishHandler struct {
	dishService service.DishService
	userService service.UserService
	validator   *validator.Validate
	logger      *logger.Logger
}

// NewDishHandler creates a new dish handler
func NewDishHandler(dishService service.DishService, userService service.UserService, log *logger.Logger) *DishHandler {
	return &DishHandler{
		dishService: dishService,
		userService: userService,
		validator:   validator.New(),
		logger:      log,
	}
}

// GetDishes handles GET /api/dishes
func (h *DishHandler) GetDishes(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse search parameter (support both 'search' and 'q' parameters)
	search := strings.TrimSpace(c.Query("search"))
	if search == "" {
		search = strings.TrimSpace(c.Query("q"))
	}

	// Parse filter parameters
	filter := service.DishFilter{
		Type:       c.Query("type"),
		Cuisine:    c.Query("cuisine"),
		SpiceLevel: c.Query("spiceLevel"),
	}

	// Parse dietary tags
	if tags := c.Query("dietaryTags"); tags != "" {
		filter.DietaryTags = strings.Split(tags, ",")
	}

	// Parse calorie range
	if maxCal := c.Query("maxCalories"); maxCal != "" {
		if val, err := strconv.Atoi(maxCal); err == nil {
			filter.MaxCalories = val
		}
	}
	if minCal := c.Query("minCalories"); minCal != "" {
		if val, err := strconv.Atoi(minCal); err == nil {
			filter.MinCalories = val
		}
	}

	// Parse ingredients
	if ingredients := c.Query("ingredients"); ingredients != "" {
		filter.Ingredients = strings.Split(ingredients, ",")
	}

	// Get user ID from context (optional)
	var userID *primitive.ObjectID
	if id, exists := middleware.GetUserIDFromContext(c); exists {
		userID = &id
	}

	var dishes []*models.DishResponse
	var pagination *models.PaginationResponse
	var err error

	// Search or get all dishes
	if search != "" {
		dishes, pagination, err = h.dishService.Search(c.Request.Context(), search, filter, page, limit, userID)
	} else {
		dishes, pagination, err = h.dishService.GetAll(c.Request.Context(), filter, page, limit, userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response := gin.H{
		"success":    true,
		"dishes":     dishes,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
}

// GetDish handles GET /api/dishes/:id
func (h *DishHandler) GetDish(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid dish ID",
		})
		return
	}

	// Get user ID from context (optional)
	var userID *primitive.ObjectID
	if uid, exists := middleware.GetUserIDFromContext(c); exists {
		userID = &uid
	}

	dish, err := h.dishService.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "dish not found" {
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
		Data:    dish,
	})
}

// GetFavorites handles GET /api/dishes/favorites
func (h *DishHandler) GetFavorites(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	dishes, pagination, err := h.dishService.GetFavorites(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response := gin.H{
		"success":    true,
		"dishes":     dishes,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
}

// AddToFavorites handles POST /api/dishes/:id/favorite
func (h *DishHandler) AddToFavorites(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	dishID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid dish ID",
		})
		return
	}

	// Check if dish exists
	_, err = h.dishService.GetByID(c.Request.Context(), dishID, nil)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "dish not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Add to favorites
	if err := h.userService.AddToFavorites(c.Request.Context(), userID, dishID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Dish added to favorites",
	})
}

// RemoveFromFavorites handles DELETE /api/dishes/:id/favorite
func (h *DishHandler) RemoveFromFavorites(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	dishID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid dish ID",
		})
		return
	}

	// Remove from favorites
	if err := h.userService.RemoveFromFavorites(c.Request.Context(), userID, dishID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Dish removed from favorites",
	})
}

// CreateDish handles POST /api/dishes
func (h *DishHandler) CreateDish(c *gin.Context) {
	var req models.DishCreateRequest
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

	// Convert request to dish model
	dish := &models.Dish{
		Name:        req.Name,
		Type:        req.Type,
		Cuisine:     req.Cuisine,
		Image:       req.Image,
		Ingredients: req.Ingredients,
		Calories:    req.Calories,
		Nutrition:   req.Nutrition,
		DietaryTags: req.DietaryTags,
		SpiceLevel:  req.SpiceLevel,
		PrepTime:    req.PrepTime,
		CookTime:    req.CookTime,
		Servings:    req.Servings,
		Difficulty:  req.Difficulty,
		Description: req.Description,
	}

	// Set defaults if not provided
	if dish.SpiceLevel == "" {
		dish.SpiceLevel = "medium"
	}
	if dish.Difficulty == "" {
		dish.Difficulty = "medium"
	}
	if dish.Servings == 0 {
		dish.Servings = 2
	}

	// Create dish
	if err := h.dishService.Create(c.Request.Context(), dish); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "Dish created successfully",
		Data:    dish.ToResponse(),
	})
}
