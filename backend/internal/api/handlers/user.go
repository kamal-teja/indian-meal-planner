package handlers

import (
	"net/http"

	"nourish-backend/internal/api/middleware"
	"nourish-backend/internal/models"
	"nourish-backend/internal/service"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService service.UserService
	validator   *validator.Validate
	logger      *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
		logger:      log,
	}
}

// GetProfile handles GET /api/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    user.ToResponse(),
	})
}

// UpdateProfile handles PUT /api/user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	var req models.ProfileUpdateRequest
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

	// Update profile using the new service method
	if err := h.userService.UpdateUserProfile(c.Request.Context(), userID, req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Get updated user to return in response
	updatedUser, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to retrieve updated user",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    updatedUser.ToResponse(),
	})
}

// DeleteAccount handles DELETE /api/user/account
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	if err := h.userService.Delete(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Account deleted successfully",
	})
}
