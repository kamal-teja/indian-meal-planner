package service

import (
	"context"
	"testing"

	"nourish-backend/internal/models"
	"nourish-backend/internal/repository"
	"nourish-backend/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock DishRepository
type MockDishRepository struct {
	mock.Mock
}

func (m *MockDishRepository) Create(ctx context.Context, dish *models.Dish) error {
	args := m.Called(ctx, dish)
	return args.Error(0)
}

func (m *MockDishRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Dish, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dish), args.Error(1)
}

func (m *MockDishRepository) GetAll(ctx context.Context, filter repository.DishFilter, page, limit int) ([]*models.Dish, int64, error) {
	args := m.Called(ctx, filter, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Dish), args.Get(1).(int64), args.Error(2)
}

func (m *MockDishRepository) Search(ctx context.Context, query string, filter repository.DishFilter, page, limit int) ([]*models.Dish, int64, error) {
	args := m.Called(ctx, query, filter, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Dish), args.Get(1).(int64), args.Error(2)
}

func (m *MockDishRepository) Update(ctx context.Context, id primitive.ObjectID, dish *models.Dish) error {
	args := m.Called(ctx, id, dish)
	return args.Error(0)
}

func (m *MockDishRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDishRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDishRepository) GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Dish, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Dish), args.Error(1)
}

func TestDishService_GetByID_Success(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	dish := &models.Dish{
		ID:       dishID,
		Name:     "Test Dish",
		Type:     "Veg",
		Cuisine:  "Indian",
		Calories: 300,
	}

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(dish, nil)
	mockUserRepo.On("GetFavorites", mock.Anything, userID).Return([]primitive.ObjectID{dishID}, nil)

	// Act
	result, err := service.GetByID(context.Background(), dishID, &userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dish.ID.Hex(), result.ID)
	assert.Equal(t, dish.Name, result.Name)
	assert.True(t, result.IsFavorite)
	mockDishRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestDishService_GetByID_DishNotFound(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishID := primitive.NewObjectID()

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(nil, mongo.ErrNoDocuments)

	// Act
	result, err := service.GetByID(context.Background(), dishID, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "dish not found")
	mockDishRepo.AssertExpectations(t)
}

func TestDishService_GetByID_WithoutUser(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishID := primitive.NewObjectID()
	dish := &models.Dish{
		ID:       dishID,
		Name:     "Test Dish",
		Type:     "Veg",
		Cuisine:  "Indian",
		Calories: 300,
	}

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(dish, nil)

	// Act
	result, err := service.GetByID(context.Background(), dishID, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dish.ID.Hex(), result.ID)
	assert.Equal(t, dish.Name, result.Name)
	assert.False(t, result.IsFavorite)
	mockDishRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "GetFavorites")
}

func TestDishService_GetAll_Success(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishes := []*models.Dish{
		{
			ID:      primitive.NewObjectID(),
			Name:    "Dish 1",
			Type:    "Veg",
			Cuisine: "Indian",
		},
		{
			ID:      primitive.NewObjectID(),
			Name:    "Dish 2",
			Type:    "Non-Veg",
			Cuisine: "Chinese",
		},
	}

	filter := DishFilter{Type: "Veg"}
	page := 1
	limit := 10

	mockDishRepo.On("GetAll", mock.Anything, mock.AnythingOfType("repository.DishFilter"), page, limit).Return(dishes, int64(2), nil)

	// Act
	result, paginationResult, err := service.GetAll(context.Background(), filter, page, limit, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.NotNil(t, paginationResult)
	assert.Equal(t, page, paginationResult.Page)
	assert.Equal(t, 1, paginationResult.TotalPages)
	assert.Equal(t, 2, paginationResult.Total)
	mockDishRepo.AssertExpectations(t)
}

func TestDishService_Search_Success(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishes := []*models.Dish{
		{
			ID:   primitive.NewObjectID(),
			Name: "Paneer Butter Masala",
			Type: "Veg",
		},
	}

	query := "paneer"
	filter := DishFilter{}
	page := 1
	limit := 10

	mockDishRepo.On("Search", mock.Anything, query, page, limit).Return(dishes, int64(1), nil)

	// Act
	result, paginationResult, err := service.Search(context.Background(), query, filter, page, limit, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Paneer Butter Masala", result[0].Name)
	assert.NotNil(t, paginationResult)
	mockDishRepo.AssertExpectations(t)
}

func TestDishService_GetFavorites_Success(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	userID := primitive.NewObjectID()
	dishID1 := primitive.NewObjectID()
	dishID2 := primitive.NewObjectID()

	favoriteIDs := []primitive.ObjectID{dishID1, dishID2}
	dishes := []*models.Dish{
		{
			ID:   dishID1,
			Name: "Favorite Dish 1",
		},
		{
			ID:   dishID2,
			Name: "Favorite Dish 2",
		},
	}

	page := 1
	limit := 10

	mockUserRepo.On("GetFavorites", mock.Anything, userID).Return(favoriteIDs, nil)
	mockDishRepo.On("GetByIDs", mock.Anything, favoriteIDs).Return(dishes, nil)

	// Act
	result, pagination, err := service.GetFavorites(context.Background(), userID, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.True(t, result[0].IsFavorite)
	assert.True(t, result[1].IsFavorite)
	assert.NotNil(t, pagination)
	mockDishRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestDishService_Create_Success(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dish := &models.Dish{
		Name:    "New Dish",
		Type:    "Veg",
		Cuisine: "Indian",
	}

	mockDishRepo.On("Create", mock.Anything, dish).Return(nil)

	// Act
	err := service.Create(context.Background(), dish)

	// Assert
	assert.NoError(t, err)
	mockDishRepo.AssertExpectations(t)
}

func TestDishService_Update_Success(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishID := primitive.NewObjectID()
	existingDish := &models.Dish{
		ID:   dishID,
		Name: "Old Name",
	}
	updatedDish := &models.Dish{
		Name: "New Name",
	}

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(existingDish, nil)
	mockDishRepo.On("Update", mock.Anything, dishID, updatedDish).Return(nil)

	// Act
	err := service.Update(context.Background(), dishID, updatedDish)

	// Assert
	assert.NoError(t, err)
	mockDishRepo.AssertExpectations(t)
}

func TestDishService_Update_DishNotFound(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishID := primitive.NewObjectID()
	updatedDish := &models.Dish{
		Name: "New Name",
	}

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(nil, mongo.ErrNoDocuments)

	// Act
	err := service.Update(context.Background(), dishID, updatedDish)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dish not found")
	mockDishRepo.AssertExpectations(t)
	mockDishRepo.AssertNotCalled(t, "Update")
}

func TestDishService_Delete_Success(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishID := primitive.NewObjectID()
	existingDish := &models.Dish{
		ID:   dishID,
		Name: "Dish to Delete",
	}

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(existingDish, nil)
	mockDishRepo.On("Delete", mock.Anything, dishID).Return(nil)

	// Act
	err := service.Delete(context.Background(), dishID)

	// Assert
	assert.NoError(t, err)
	mockDishRepo.AssertExpectations(t)
}

func TestDishService_Delete_DishNotFound(t *testing.T) {
	// Arrange
	mockDishRepo := new(MockDishRepository)
	mockUserRepo := new(MockUserRepositoryForUserService)
	log := logger.New("info", "json")
	service := NewDishService(mockDishRepo, mockUserRepo, log)

	dishID := primitive.NewObjectID()

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(nil, mongo.ErrNoDocuments)

	// Act
	err := service.Delete(context.Background(), dishID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dish not found")
	mockDishRepo.AssertExpectations(t)
	mockDishRepo.AssertNotCalled(t, "Delete")
}
