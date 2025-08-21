package service

import (
	"context"
	"testing"

	"nourish-backend/internal/models"
	"nourish-backend/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock UserRepository for UserService tests
type MockUserRepositoryForUserService struct {
	mock.Mock
}

func (m *MockUserRepositoryForUserService) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) Update(ctx context.Context, id primitive.ObjectID, user *models.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) List(ctx context.Context, limit, skip int) ([]*models.User, error) {
	args := m.Called(ctx, limit, skip)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepositoryForUserService) GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

func (m *MockUserRepositoryForUserService) AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) AddFavorite(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func (m *MockUserRepositoryForUserService) RemoveFavorite(ctx context.Context, userID, dishID primitive.ObjectID) error {
	args := m.Called(ctx, userID, dishID)
	return args.Error(0)
}

func TestUserService_GetProfile_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()
	user := &models.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
		Profile: models.UserProfile{
			DietaryPreferences: []string{"vegetarian"},
			SpiceLevel:         "medium",
		},
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Act
	result, err := service.GetProfile(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID.Hex(), result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetProfile_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()

	mockRepo.On("GetByID", mock.Anything, userID).Return(nil, mongo.ErrNoDocuments)

	// Act
	result, err := service.GetProfile(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()
	existingUser := &models.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	req := models.UserProfileUpdateRequest{
		Name: "John Smith",
		Profile: models.UserProfile{
			DietaryPreferences: []string{"vegan"},
			SpiceLevel:         "high",
		},
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, userID, mock.AnythingOfType("*models.User")).Return(nil)

	// Act
	result, err := service.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()
	req := models.UserProfileUpdateRequest{
		Name: "John Smith",
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(nil, mongo.ErrNoDocuments)

	// Act
	result, err := service.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetFavorites_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()
	dishID1 := primitive.NewObjectID()
	dishID2 := primitive.NewObjectID()
	favorites := []primitive.ObjectID{dishID1, dishID2}

	mockRepo.On("GetFavorites", mock.Anything, userID).Return(favorites, nil)

	// Act
	result, err := service.GetFavorites(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, dishID1.Hex(), result[0])
	assert.Equal(t, dishID2.Hex(), result[1])
	mockRepo.AssertExpectations(t)
}

func TestUserService_AddFavorite_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()

	mockRepo.On("AddFavorite", mock.Anything, userID, dishID).Return(nil)

	// Act
	err := service.AddFavorite(context.Background(), userID, dishID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_RemoveFavorite_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()

	mockRepo.On("RemoveFavorite", mock.Anything, userID, dishID).Return(nil)

	// Act
	err := service.RemoveFavorite(context.Background(), userID, dishID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteAccount_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	userID := primitive.NewObjectID()

	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	// Act
	err := service.DeleteAccount(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewUserService(mockRepo, log)

	users := []*models.User{
		{
			ID:    primitive.NewObjectID(),
			Name:  "User 1",
			Email: "user1@example.com",
		},
		{
			ID:    primitive.NewObjectID(),
			Name:  "User 2",
			Email: "user2@example.com",
		},
	}

	limit := 10
	skip := 0

	mockRepo.On("List", mock.Anything, limit, skip).Return(users, nil)

	// Act
	result, err := service.ListUsers(context.Background(), limit, skip)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].Name, result[0].Name)
	assert.Equal(t, users[1].Name, result[1].Name)
	mockRepo.AssertExpectations(t)
}