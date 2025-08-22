package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"nourish-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock AuthService for middleware tests
type MockAuthServiceForMiddleware struct {
	mock.Mock
}

func (m *MockAuthServiceForMiddleware) Register(ctx context.Context, req models.UserRegistrationRequest) (*models.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthServiceForMiddleware) Login(ctx context.Context, req models.UserLoginRequest) (*models.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthServiceForMiddleware) GenerateToken(userID primitive.ObjectID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthServiceForMiddleware) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockAuthServiceForMiddleware) GetUserFromToken(tokenString string) (*models.User, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func setupAuthMiddleware() (*MockAuthServiceForMiddleware, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	
	mockAuthService := new(MockAuthServiceForMiddleware)
	router := gin.New()
	
	return mockAuthService, router
}

func TestAuthMiddleware_Success(t *testing.T) {
	// Arrange
	mockAuthService, router := setupAuthMiddleware()
	
	user := &models.User{
		ID:    primitive.NewObjectID(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	
	mockAuthService.On("GetUserFromToken", "valid-token").Return(user, nil)
	
	router.Use(AuthMiddleware(mockAuthService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()
	
	// Act
	router.ServeHTTP(recorder, request)
	
	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_NoAuthorizationHeader(t *testing.T) {
	// Arrange
	mockAuthService, router := setupAuthMiddleware()
	
	router.Use(AuthMiddleware(mockAuthService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	recorder := httptest.NewRecorder()
	
	// Act
	router.ServeHTTP(recorder, request)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Access denied. No token provided.", response.Error)
}

func TestAuthMiddleware_InvalidAuthorizationFormat(t *testing.T) {
	// Arrange
	mockAuthService, router := setupAuthMiddleware()
	
	router.Use(AuthMiddleware(mockAuthService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "InvalidFormat")
	recorder := httptest.NewRecorder()
	
	// Act
	router.ServeHTTP(recorder, request)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid authorization header format.", response.Error)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Arrange
	mockAuthService, router := setupAuthMiddleware()
	
	mockAuthService.On("GetUserFromToken", "invalid-token").Return(nil, errors.New("invalid token"))
	
	router.Use(AuthMiddleware(mockAuthService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")
	recorder := httptest.NewRecorder()
	
	// Act
	router.ServeHTTP(recorder, request)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid token.", response.Error)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_WrongBearerPrefix(t *testing.T) {
	// Arrange
	mockAuthService, router := setupAuthMiddleware()
	
	router.Use(AuthMiddleware(mockAuthService))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Basic some-token")
	recorder := httptest.NewRecorder()
	
	// Act
	router.ServeHTTP(recorder, request)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid authorization header format.", response.Error)
}

func TestAuthMiddleware_ContextSetCorrectly(t *testing.T) {
	// Arrange
	mockAuthService, router := setupAuthMiddleware()
	
	userID := primitive.NewObjectID()
	user := &models.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}
	
	mockAuthService.On("GetUserFromToken", "valid-token").Return(user, nil)
	
	router.Use(AuthMiddleware(mockAuthService))
	router.GET("/protected", func(c *gin.Context) {
		// Check if user and userID are set in context
		contextUser, userExists := GetUserFromContext(c)
		contextUserID, idExists := GetUserIDFromContext(c)
		
		assert.True(t, userExists)
		assert.True(t, idExists)
		assert.Equal(t, user.ID, contextUser.ID)
		assert.Equal(t, userID, contextUserID)
		
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()
	
	// Act
	router.ServeHTTP(recorder, request)
	
	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockAuthService.AssertExpectations(t)
}

func TestGetUserFromContext_UserExists(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	user := &models.User{
		ID:    primitive.NewObjectID(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	c.Set("user", user)
	
	// Act
	result, exists := GetUserFromContext(c)
	
	// Assert
	assert.True(t, exists)
	assert.Equal(t, user, result)
}

func TestGetUserFromContext_UserDoesNotExist(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	// Act
	result, exists := GetUserFromContext(c)
	
	// Assert
	assert.False(t, exists)
	assert.Nil(t, result)
}

func TestGetUserIDFromContext_UserIDExists(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	userID := primitive.NewObjectID()
	c.Set("userID", userID)
	
	// Act
	result, exists := GetUserIDFromContext(c)
	
	// Assert
	assert.True(t, exists)
	assert.Equal(t, userID, result)
}

func TestGetUserIDFromContext_UserIDDoesNotExist(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	
	// Act
	result, exists := GetUserIDFromContext(c)
	
	// Assert
	assert.False(t, exists)
	assert.Equal(t, primitive.NilObjectID, result)
}