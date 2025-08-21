package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"nourish-backend/internal/config"
	"nourish-backend/internal/models"
	"nourish-backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id primitive.ObjectID, user *models.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, skip int) ([]*models.User, error) {
	args := m.Called(ctx, limit, skip)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

func (m *MockUserRepository) AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) AddFavorite(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveFavorite(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func TestAuthService_Register_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	req := models.UserRegistrationRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Mock: User doesn't exist
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)
	// Mock: Create user succeeds
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	// Act
	result, err := service.Register(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.User.Name)
	assert.Equal(t, req.Email, result.User.Email)
	assert.NotEmpty(t, result.Token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_UserExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	req := models.UserRegistrationRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	existingUser := &models.User{
		ID:    primitive.NewObjectID(),
		Email: req.Email,
	}

	// Mock: User already exists
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	// Act
	result, err := service.Register(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user already exists")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		ID:        primitive.NewObjectID(),
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := models.UserLoginRequest{
		Email:    user.Email,
		Password: password,
	}

	// Mock: User exists
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(user, nil)

	// Act
	result, err := service.Login(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Name, result.User.Name)
	assert.Equal(t, user.Email, result.User.Email)
	assert.NotEmpty(t, result.Token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	req := models.UserLoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Mock: User doesn't exist
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)

	// Act
	result, err := service.Login(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	user := &models.User{
		ID:       primitive.NewObjectID(),
		Email:    "john@example.com",
		Password: string(hashedPassword),
	}

	req := models.UserLoginRequest{
		Email:    user.Email,
		Password: "wrongpassword",
	}

	// Mock: User exists
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(user, nil)

	// Act
	result, err := service.Login(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GenerateToken_Success(t *testing.T) {
	// Arrange
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(nil, cfg, log)

	userID := primitive.NewObjectID()

	// Act
	token, err := service.GenerateToken(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token can be parsed
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID.Hex(), claims["user_id"])
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Arrange
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(nil, cfg, log)

	userID := primitive.NewObjectID()
	token, _ := service.GenerateToken(userID)

	// Act
	parsedToken, err := service.ValidateToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)
	assert.True(t, parsedToken.Valid)
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(nil, cfg, log)

	invalidToken := "invalid.token.here"

	// Act
	parsedToken, err := service.ValidateToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, parsedToken)
}

func TestAuthService_GetUserFromToken_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	userID := primitive.NewObjectID()
	user := &models.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	token, _ := service.GenerateToken(userID)

	// Mock: User exists
	mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Act
	result, err := service.GetUserFromToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserFromToken_InvalidToken(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	invalidToken := "invalid.token.here"

	// Act
	result, err := service.GetUserFromToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestAuthService_GetUserFromToken_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}
	log := logger.New("info", "json")
	service := NewAuthService(mockRepo, cfg, log)

	userID := primitive.NewObjectID()
	token, _ := service.GenerateToken(userID)

	// Mock: User doesn't exist
	mockRepo.On("GetByID", mock.Anything, userID).Return(nil, mongo.ErrNoDocuments)

	// Act
	result, err := service.GetUserFromToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}