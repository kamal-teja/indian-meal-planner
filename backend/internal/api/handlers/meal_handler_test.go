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
	"nourish-backend/internal/repository"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock MealService
type MockMealService struct {
	mock.Mock
}

func (m *MockMealService) Create(ctx context.Context, userID primitive.ObjectID, req models.MealRequest) (*models.MealWithDish, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MealWithDish), args.Error(1)
}

func (m *MockMealService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.MealWithDish, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MealWithDish), args.Error(1)
}

func (m *MockMealService) GetByDate(ctx context.Context, userID primitive.ObjectID, date string) ([]*models.MealWithDish, error) {
	args := m.Called(ctx, userID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MealWithDish), args.Error(1)
}

func (m *MockMealService) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.MealWithDish, *models.PaginationResponse, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.MealWithDish), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockMealService) GetByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*models.MealWithDish, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MealWithDish), args.Error(1)
}

func (m *MockMealService) Update(ctx context.Context, id primitive.ObjectID, req models.MealRequest) (*models.MealWithDish, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MealWithDish), args.Error(1)
}

func (m *MockMealService) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMealService) GetAnalytics(ctx context.Context, userID primitive.ObjectID, period int) (*models.AnalyticsResponse, error) {
	args := m.Called(ctx, userID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AnalyticsResponse), args.Error(1)
}

func (m *MockMealService) GetShoppingList(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*models.ShoppingListResponse, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShoppingListResponse), args.Error(1)
}

func (m *MockMealService) GetRecommendations(ctx context.Context, userID primitive.ObjectID, mealType string, date time.Time) (*models.RecommendationsResponse, error) {
	args := m.Called(ctx, userID, mealType, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RecommendationsResponse), args.Error(1)
}

func (m *MockMealService) GetNutritionProgress(ctx context.Context, userID primitive.ObjectID, period int) (*models.NutritionProgressResponse, error) {
	args := m.Called(ctx, userID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NutritionProgressResponse), args.Error(1)
}

func (m *MockMealService) GetNutritionGoals(ctx context.Context, userID primitive.ObjectID) (*models.NutritionGoals, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NutritionGoals), args.Error(1)
}

func (m *MockMealService) UpdateNutritionGoals(ctx context.Context, userID primitive.ObjectID, req models.NutritionGoalsRequest) (*models.NutritionGoals, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NutritionGoals), args.Error(1)
}

func (m *MockMealService) GetNutritionSummary(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]repository.NutritionSummary, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.NutritionSummary), args.Error(1)
}

func setupMealHandler() (*MealHandler, *MockMealService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	
	mockService := new(MockMealService)
	log := logger.New("info", "json")
	
	handler := NewMealHandler(mockService, log)
	router := gin.New()
	
	// Add middleware to set user ID for tests that need auth
	router.Use(func(c *gin.Context) {
		userID := primitive.NewObjectID()
		c.Set("userID", userID)
		c.Next()
	})
	
	return handler, mockService, router
}

func TestMealHandler_CreateMeal_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.POST("/meals", handler.CreateMeal)

	req := models.MealRequest{
		DishID:   primitive.NewObjectID().Hex(),
		MealType: "breakfast",
		Date:     models.FlexibleDate{},
		Rating:   5,
		Notes:    "Delicious",
	}

	expectedMeal := &models.MealWithDish{
		ID:       primitive.NewObjectID().Hex(),
		MealType: "breakfast",
		Rating:   5,
		Notes:    "Delicious",
	}

	mockService.On("Create", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), req).Return(expectedMeal, nil)

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/meals", bytes.NewBuffer(requestBody))
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
	assert.Equal(t, "Meal created successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestMealHandler_CreateMeal_InvalidJSON(t *testing.T) {
	// Arrange
	handler, _, router := setupMealHandler()
	router.POST("/meals", handler.CreateMeal)

	request := httptest.NewRequest(http.MethodPost, "/meals", bytes.NewBufferString("invalid json"))
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

func TestMealHandler_CreateMeal_ServiceError(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.POST("/meals", handler.CreateMeal)

	req := models.MealRequest{
		DishID:   primitive.NewObjectID().Hex(),
		MealType: "breakfast",
		Date:     models.FlexibleDate{},
		Rating:   5,
	}

	mockService.On("Create", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), req).Return(nil, errors.New("service error"))

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/meals", bytes.NewBuffer(requestBody))
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

func TestMealHandler_GetMeal_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.GET("/meals/:id", handler.GetMeal)

	mealID := primitive.NewObjectID()
	expectedMeal := &models.MealWithDish{
		ID:       mealID.Hex(),
		MealType: "breakfast",
		Rating:   5,
		Notes:    "Delicious",
	}

	mockService.On("GetByID", mock.Anything, mealID).Return(expectedMeal, nil)

	request := httptest.NewRequest(http.MethodGet, "/meals/"+mealID.Hex(), nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestMealHandler_GetMeal_InvalidID(t *testing.T) {
	// Arrange
	handler, _, router := setupMealHandler()
	router.GET("/meals/:id", handler.GetMeal)

	request := httptest.NewRequest(http.MethodGet, "/meals/invalid-id", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Invalid meal ID")
}

func TestMealHandler_GetMeal_NotFound(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.GET("/meals/:id", handler.GetMeal)

	mealID := primitive.NewObjectID()
	mockService.On("GetByID", mock.Anything, mealID).Return(nil, errors.New("meal not found"))

	request := httptest.NewRequest(http.MethodGet, "/meals/"+mealID.Hex(), nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	
	mockService.AssertExpectations(t)
}

func TestMealHandler_GetMeals_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.GET("/meals", handler.GetMeals)

	date := "2023-10-15"
	meals := []*models.MealWithDish{
		{
			ID:       primitive.NewObjectID().Hex(),
			MealType: "breakfast",
			Rating:   5,
		},
		{
			ID:       primitive.NewObjectID().Hex(),
			MealType: "lunch",
			Rating:   4,
		},
	}

	// Mock GetByDateRange which is what the handler actually calls for date queries
	mockService.On("GetByDateRange", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(meals, nil)

	request := httptest.NewRequest(http.MethodGet, "/meals?date="+date, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestMealHandler_UpdateMeal_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.PUT("/meals/:id", handler.UpdateMeal)

	mealID := primitive.NewObjectID()
	req := models.MealRequest{
		DishID:   primitive.NewObjectID().Hex(),
		MealType: "lunch",
		Rating:   4,
		Notes:    "Updated notes",
	}

	expectedMeal := &models.MealWithDish{
		ID:       mealID.Hex(),
		MealType: "lunch",
		Rating:   4,
		Notes:    "Updated notes",
	}

	mockService.On("Update", mock.Anything, mealID, req).Return(expectedMeal, nil)

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPut, "/meals/"+mealID.Hex(), bytes.NewBuffer(requestBody))
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
	assert.Equal(t, "Meal updated successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestMealHandler_DeleteMeal_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.DELETE("/meals/:id", handler.DeleteMeal)

	mealID := primitive.NewObjectID()
	mockService.On("Delete", mock.Anything, mealID).Return(nil)

	request := httptest.NewRequest(http.MethodDelete, "/meals/"+mealID.Hex(), nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Meal deleted successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestMealHandler_GetMealsByMonth_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupMealHandler()
	router.GET("/meals/month/:year/:month", handler.GetMealsByMonth)

	meals := []*models.MealWithDish{
		{
			ID:       primitive.NewObjectID().Hex(),
			MealType: "breakfast",
			Rating:   5,
		},
	}

	// GetMealsByMonth uses GetByDateRange internally
	mockService.On("GetByDateRange", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(meals, nil)

	request := httptest.NewRequest(http.MethodGet, "/meals/month/2023/10", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	mockService.AssertExpectations(t)
}