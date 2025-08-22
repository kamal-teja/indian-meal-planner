package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nourish-backend/internal/models"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock UserService for User Handler  
type MockUserServiceForUserHandler struct {
	mock.Mock
}

func (m *MockUserServiceForUserHandler) GetByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForUserHandler) UpdateProfile(ctx context.Context, userID primitive.ObjectID, profile models.UserProfile) error {
	args := m.Called(ctx, userID, profile)
	return args.Error(0)
}

func (m *MockUserServiceForUserHandler) UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, req models.ProfileUpdateRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserServiceForUserHandler) AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserServiceForUserHandler) RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserServiceForUserHandler) GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

func (m *MockUserServiceForUserHandler) Delete(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func setupUserHandler() (*UserHandler, *MockUserServiceForUserHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	
	mockService := new(MockUserServiceForUserHandler)
	log := logger.New("info", "json")
	
	handler := NewUserHandler(mockService, log)
	router := gin.New()
	
	// Add middleware to set user and userID for tests that need auth
	router.Use(func(c *gin.Context) {
		userID := primitive.NewObjectID()
		user := &models.User{
			ID:   userID,
			Name: "Test User",
			Email: "test@example.com",
			Profile: models.UserProfile{
				DietaryPreferences: []string{"vegetarian"},
				SpiceLevel:         "medium",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		c.Set("userID", userID)
		c.Set("user", user)
		c.Next()
	})
	
	return handler, mockService, router
}

func TestUserHandler_GetProfile_Success(t *testing.T) {
	// Arrange
	handler, _, router := setupUserHandler()
	router.GET("/profile", handler.GetProfile)

	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
}

func TestUserHandler_GetProfile_Unauthorized(t *testing.T) {
	// Arrange
	handler, _, _ := setupUserHandler()
	
	// Create router without auth middleware
	router := gin.New()
	router.GET("/profile", handler.GetProfile)

	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Authentication required", response.Error)
}

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupUserHandler()
	router.PUT("/profile", handler.UpdateProfile)

	req := models.ProfileUpdateRequest{
		Name: "Updated Name",
		Profile: models.UserProfile{
			DietaryPreferences: []string{"vegan"},
			SpiceLevel:         "hot",
		},
	}

	mockService.On("UpdateUserProfile", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), req).Return(nil)

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Profile updated successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile_InvalidJSON(t *testing.T) {
	// Arrange
	handler, _, router := setupUserHandler()
	router.PUT("/profile", handler.UpdateProfile)

	request := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBufferString("invalid json"))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Invalid request format")
}

func TestUserHandler_UpdateProfile_ServiceError(t *testing.T) {
	// Arrange
	handler, mockService, router := setupUserHandler()
	router.PUT("/profile", handler.UpdateProfile)

	req := models.ProfileUpdateRequest{
		Name: "Updated Name",
		Profile: models.UserProfile{
			DietaryPreferences: []string{"vegan"},
			SpiceLevel:         "hot",
		},
	}

	mockService.On("UpdateUserProfile", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), req).Return(errors.New("service error"))

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile_Unauthorized(t *testing.T) {
	// Arrange
	handler, _, _ := setupUserHandler()
	
	// Create router without auth middleware
	router := gin.New()
	router.PUT("/profile", handler.UpdateProfile)

	req := models.ProfileUpdateRequest{
		Name: "Updated Name",
	}

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Authentication required", response.Error)
}

func TestUserHandler_DeleteAccount_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupUserHandler()
	router.DELETE("/account", handler.DeleteAccount)

	mockService.On("Delete", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	request := httptest.NewRequest(http.MethodDelete, "/account", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Account deleted successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteAccount_ServiceError(t *testing.T) {
	// Arrange
	handler, mockService, router := setupUserHandler()
	router.DELETE("/account", handler.DeleteAccount)

	mockService.On("Delete", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(errors.New("service error"))

	request := httptest.NewRequest(http.MethodDelete, "/account", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteAccount_Unauthorized(t *testing.T) {
	// Arrange
	handler, _, _ := setupUserHandler()
	
	// Create router without auth middleware
	router := gin.New()
	router.DELETE("/account", handler.DeleteAccount)

	request := httptest.NewRequest(http.MethodDelete, "/account", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Authentication required", response.Error)
}