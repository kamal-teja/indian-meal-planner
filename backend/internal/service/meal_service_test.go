package service

import (
	"context"
	"testing"
	"time"

	"nourish-backend/internal/models"
	"nourish-backend/internal/repository"
	"nourish-backend/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock MealRepository
type MockMealRepository struct {
	mock.Mock
}

func (m *MockMealRepository) Create(ctx context.Context, meal *models.Meal) error {
	args := m.Called(ctx, meal)
	return args.Error(0)
}

func (m *MockMealRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Meal, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Meal), args.Error(1)
}

func (m *MockMealRepository) GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*models.Meal, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Meal), args.Error(1)
}

func (m *MockMealRepository) GetNutritionByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]repository.NutritionSummary, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.NutritionSummary), args.Error(1)
}

func (m *MockMealRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.Meal, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Meal), args.Get(1).(int64), args.Error(2)
}

func (m *MockMealRepository) GetByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*models.Meal, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Meal), args.Error(1)
}

func (m *MockMealRepository) Update(ctx context.Context, id primitive.ObjectID, meal *models.Meal) error {
	args := m.Called(ctx, id, meal)
	return args.Error(0)
}

func (m *MockMealRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMealRepository) Count(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMealRepository) GetNutritionSummary(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]repository.NutritionSummary, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.NutritionSummary), args.Error(1)
}

func TestMealService_Create_Success(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	userID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()
	
	dish := &models.Dish{
		ID:       dishID,
		Name:     "Test Dish",
		Type:     "Veg",
		Calories: 300,
	}

	req := models.MealRequest{
		DishID:   dishID.Hex(),
		MealType: "breakfast",
		Date:     models.FlexibleDate(time.Now().Format("2006-01-02")),
		Rating:   5,
		Notes:    "Delicious",
	}

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(dish, nil)
	mockMealRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Meal")).Return(nil)

	// Act
	result, err := service.Create(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID.Hex(), result.User)
	assert.Equal(t, dishID.Hex(), result.Dish.ID)
	assert.Equal(t, req.MealType, result.MealType)
	assert.Equal(t, dish.Name, result.Dish.Name)
	mockDishRepo.AssertExpectations(t)
	mockMealRepo.AssertExpectations(t)
}

func TestMealService_Create_DishNotFound(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	userID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()

	req := models.MealRequest{
		DishID:   dishID.Hex(),
		MealType: "breakfast",
		Date:     models.FlexibleDate(time.Now().Format("2006-01-02")),
	}

	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(nil, mongo.ErrNoDocuments)

	// Act
	result, err := service.Create(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "dish not found")
	mockDishRepo.AssertExpectations(t)
	mockMealRepo.AssertNotCalled(t, "Create")
}

func TestMealService_Create_InvalidDate(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	userID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()

	req := models.MealRequest{
		DishID:   dishID.Hex(),
		MealType: "breakfast",
		Date:     models.FlexibleDate("invalid-date"),
	}

	// Act
	result, err := service.Create(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid date format")
	mockDishRepo.AssertNotCalled(t, "GetByID")
	mockMealRepo.AssertNotCalled(t, "Create")
}

func TestMealService_GetByID_Success(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	mealID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	meal := &models.Meal{
		ID:       mealID,
		UserID:   userID,
		DishID:   dishID,
		MealType: "breakfast",
		Date:     time.Now(),
		Rating:   5,
	}

	dish := &models.Dish{
		ID:       dishID,
		Name:     "Test Dish",
		Type:     "Veg",
		Calories: 300,
	}

	mockMealRepo.On("GetByID", mock.Anything, mealID).Return(meal, nil)
	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(dish, nil)

	// Act
	result, err := service.GetByID(context.Background(), mealID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, meal.ID, result.ID)
	assert.Equal(t, meal.UserID, result.UserID)
	assert.Equal(t, dish.Name, result.Dish.Name)
	mockMealRepo.AssertExpectations(t)
	mockDishRepo.AssertExpectations(t)
}

func TestMealService_GetByID_MealNotFound(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	mealID := primitive.NewObjectID()

	mockMealRepo.On("GetByID", mock.Anything, mealID).Return(nil, mongo.ErrNoDocuments)

	// Act
	result, err := service.GetByID(context.Background(), mealID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "meal not found")
	mockMealRepo.AssertExpectations(t)
	mockDishRepo.AssertNotCalled(t, "GetByID")
}

func TestMealService_GetByUserID_Success(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	userID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()

	meals := []*models.Meal{
		{
			ID:       primitive.NewObjectID(),
			UserID:   userID,
			DishID:   dishID,
			MealType: "breakfast",
			Date:     time.Now(),
		},
	}

	dish := &models.Dish{
		ID:   dishID,
		Name: "Test Dish",
	}

	page := 1
	limit := 10
	skip := 0

	mockMealRepo.On("GetByUserID", mock.Anything, userID, page, limit).Return(meals, int64(5), nil)
	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(dish, nil)

	// Act
	result, pagination, err := service.GetByUserID(context.Background(), userID, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, meals[0].ID, result[0].ID)
	assert.Equal(t, dish.Name, result[0].Dish.Name)
	assert.NotNil(t, pagination)
	assert.Equal(t, page, pagination.CurrentPage)
	mockMealRepo.AssertExpectations(t)
	mockDishRepo.AssertExpectations(t)
}

func TestMealService_GetByDateRange_Success(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	userID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()
	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()

	meals := []*models.Meal{
		{
			ID:       primitive.NewObjectID(),
			UserID:   userID,
			DishID:   dishID,
			MealType: "breakfast",
			Date:     time.Now(),
		},
	}

	dish := &models.Dish{
		ID:   dishID,
		Name: "Test Dish",
	}

	mockMealRepo.On("GetByUserAndDateRange", mock.Anything, userID, startDate, endDate).Return(meals, nil)
	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(dish, nil)

	// Act
	result, err := service.GetByDateRange(context.Background(), userID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, meals[0].ID, result[0].ID)
	assert.Equal(t, dish.Name, result[0].Dish.Name)
	mockMealRepo.AssertExpectations(t)
	mockDishRepo.AssertExpectations(t)
}

func TestMealService_Update_Success(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	mealID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	existingMeal := &models.Meal{
		ID:       mealID,
		UserID:   userID,
		DishID:   dishID,
		MealType: "breakfast",
		Date:     time.Now(),
	}

	dish := &models.Dish{
		ID:   dishID,
		Name: "Test Dish",
	}

	req := models.MealRequest{
		DishID:   dishID.Hex(),
		MealType: "lunch",
		Date:     models.FlexibleDate(time.Now().Format("2006-01-02")),
		Rating:   4,
		Notes:    "Updated notes",
	}

	mockMealRepo.On("GetByID", mock.Anything, mealID).Return(existingMeal, nil)
	mockDishRepo.On("GetByID", mock.Anything, dishID).Return(dish, nil)
	mockMealRepo.On("Update", mock.Anything, mealID, mock.AnythingOfType("*models.Meal")).Return(nil)

	// Act
	result, err := service.Update(context.Background(), mealID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.MealType, result.MealType)
	assert.Equal(t, req.Rating, result.Rating)
	assert.Equal(t, req.Notes, result.Notes)
	mockMealRepo.AssertExpectations(t)
	mockDishRepo.AssertExpectations(t)
}

func TestMealService_Delete_Success(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	mealID := primitive.NewObjectID()
	existingMeal := &models.Meal{
		ID: mealID,
	}

	mockMealRepo.On("GetByID", mock.Anything, mealID).Return(existingMeal, nil)
	mockMealRepo.On("Delete", mock.Anything, mealID).Return(nil)

	// Act
	err := service.Delete(context.Background(), mealID)

	// Assert
	assert.NoError(t, err)
	mockMealRepo.AssertExpectations(t)
}

func TestMealService_Delete_MealNotFound(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	mealID := primitive.NewObjectID()

	mockMealRepo.On("GetByID", mock.Anything, mealID).Return(nil, mongo.ErrNoDocuments)

	// Act
	err := service.Delete(context.Background(), mealID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "meal not found")
	mockMealRepo.AssertExpectations(t)
	mockMealRepo.AssertNotCalled(t, "Delete")
}

func TestMealService_GetNutritionSummary_Success(t *testing.T) {
	// Arrange
	mockMealRepo := new(MockMealRepository)
	mockDishRepo := new(MockDishRepository)
	log := logger.New("info", "json")
	service := NewMealService(mockMealRepo, mockDishRepo, log)

	userID := primitive.NewObjectID()
	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()

	nutritionSummary := []repository.NutritionSummary{
		{
			Date:     time.Now().Format("2006-01-02"),
			Calories: 1800,
			Protein:  80,
			Carbs:    200,
			Fat:      60,
		},
	}

	mockMealRepo.On("GetNutritionByDateRange", mock.Anything, userID, startDate, endDate).Return(nutritionSummary, nil)

	// Act
	result, err := service.GetNutritionSummary(context.Background(), userID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, nutritionSummary[0].Calories, result[0].Calories)
	mockMealRepo.AssertExpectations(t)
}