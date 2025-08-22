package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"nourish-backend/internal/models"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req models.UserRegistrationRequest) (*models.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req models.UserLoginRequest) (*models.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthService) GenerateToken(userID primitive.ObjectID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Mock UserService for AuthHandler
type MockUserServiceForAuth struct {
	mock.Mock
}

func (m *MockUserServiceForAuth) GetByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForAuth) UpdateProfile(ctx context.Context, userID primitive.ObjectID, profile models.UserProfile) error {
	args := m.Called(ctx, userID, profile)
	return args.Error(0)
}

func (m *MockUserServiceForAuth) UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, req models.ProfileUpdateRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserServiceForAuth) GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

func (m *MockUserServiceForAuth) AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserServiceForAuth) RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserServiceForAuth) Delete(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func setupAuthHandler() (*AuthHandler, *MockAuthService, *MockUserServiceForAuth, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserServiceForAuth)
	log := logger.New("info", "json")
	
	handler := NewAuthHandler(mockAuthService, mockUserService, log)
	router := gin.New()
	
	return handler, mockAuthService, mockUserService, router
}

func TestAuthHandler_Register_Success(t *testing.T) {
	// Arrange
	handler, mockAuthService, _, router := setupAuthHandler()
	router.POST("/register", handler.Register)

	req := models.UserRegistrationRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	authResponse := &models.AuthResponse{
		User: models.UserResponse{
			ID:    primitive.NewObjectID().Hex(),
			Name:  req.Name,
			Email: req.Email,
		},
		Token: "jwt-token-here",
	}

	mockAuthService.On("Register", mock.Anything, req).Return(authResponse, nil)

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusCreated, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "User registered successfully", response.Message)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	// Arrange
	handler, _, _, router := setupAuthHandler()
	router.POST("/register", handler.Register)

	// Invalid request - missing required fields
	req := models.UserRegistrationRequest{
		Email: "invalid-email", // Invalid email format
	}

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
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
	assert.Contains(t, response.Error, "Validation failed")
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	// Arrange
	handler, mockAuthService, _, router := setupAuthHandler()
	router.POST("/register", handler.Register)

	req := models.UserRegistrationRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockAuthService.On("Register", mock.Anything, req).Return(nil, errors.New("user already exists"))

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusConflict, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "user already exists")
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	handler, mockAuthService, _, router := setupAuthHandler()
	router.POST("/login", handler.Login)

	req := models.UserLoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	authResponse := &models.AuthResponse{
		User: models.UserResponse{
			ID:    primitive.NewObjectID().Hex(),
			Email: req.Email,
		},
		Token: "jwt-token-here",
	}

	mockAuthService.On("Login", mock.Anything, req).Return(authResponse, nil)

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
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
	assert.Equal(t, "Login successful", response.Message)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	handler, mockAuthService, _, router := setupAuthHandler()
	router.POST("/login", handler.Login)

	req := models.UserLoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	mockAuthService.On("Login", mock.Anything, req).Return(nil, errors.New("invalid credentials"))

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
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
	assert.Contains(t, response.Error, "invalid credentials")
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// Arrange
	handler, _, _, router := setupAuthHandler()
	router.POST("/login", handler.Login)

	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
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
	assert.Equal(t, "Invalid request format", response.Error)
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	// Arrange
	handler, _, _, router := setupAuthHandler()
	router.POST("/login", handler.Login)

	// Invalid request - missing required fields
	req := models.UserLoginRequest{
		Email: "invalid-email-format",
		// Missing password
	}

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
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
	assert.Contains(t, response.Error, "Validation failed")
}