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
	"nourish-backend/internal/service"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock DishService
type MockDishService struct {
	mock.Mock
}

func (m *MockDishService) GetByID(ctx context.Context, id primitive.ObjectID, userID *primitive.ObjectID) (*models.DishResponse, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DishResponse), args.Error(1)
}

func (m *MockDishService) GetAll(ctx context.Context, filter service.DishFilter, page, limit int, userID *primitive.ObjectID) ([]*models.DishResponse, *models.PaginationResponse, error) {
	args := m.Called(ctx, filter, page, limit, userID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.DishResponse), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockDishService) Search(ctx context.Context, query string, filter service.DishFilter, page, limit int, userID *primitive.ObjectID) ([]*models.DishResponse, *models.PaginationResponse, error) {
	args := m.Called(ctx, query, filter, page, limit, userID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.DishResponse), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockDishService) GetFavorites(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.DishResponse, *models.PaginationResponse, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.DishResponse), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockDishService) Create(ctx context.Context, dish *models.Dish) error {
	args := m.Called(ctx, dish)
	return args.Error(0)
}

func (m *MockDishService) Update(ctx context.Context, id primitive.ObjectID, dish *models.Dish) error {
	args := m.Called(ctx, id, dish)
	return args.Error(0)
}

func (m *MockDishService) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupDishHandler() (*DishHandler, *MockDishService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	
	mockService := new(MockDishService)
	log := logger.New("info", "json")
	
	handler := NewDishHandler(mockService, log)
	router := gin.New()
	
	return handler, mockService, router
}

func TestDishHandler_GetDish_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupDishHandler()
	router.GET("/dishes/:id", handler.GetDish)

	dishID := primitive.NewObjectID()
	dish := &models.DishResponse{
		ID:       dishID.Hex(),
		Name:     "Test Dish",
		Type:     "Veg",
		Cuisine:  "Indian",
		Calories: 300,
	}

	mockService.On("GetByID", mock.Anything, dishID, (*primitive.ObjectID)(nil)).Return(dish, nil)

	request := httptest.NewRequest(http.MethodGet, "/dishes/"+dishID.Hex(), nil)
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

func TestDishHandler_GetDish_InvalidID(t *testing.T) {
	// Arrange
	handler, _, router := setupDishHandler()
	router.GET("/dishes/:id", handler.GetDish)

	request := httptest.NewRequest(http.MethodGet, "/dishes/invalid-id", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Invalid dish ID")
}

func TestDishHandler_GetDish_NotFound(t *testing.T) {
	// Arrange
	handler, mockService, router := setupDishHandler()
	router.GET("/dishes/:id", handler.GetDish)

	dishID := primitive.NewObjectID()

	mockService.On("GetByID", mock.Anything, dishID, (*primitive.ObjectID)(nil)).Return(nil, errors.New("dish not found"))

	request := httptest.NewRequest(http.MethodGet, "/dishes/"+dishID.Hex(), nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestDishHandler_GetDishes_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupDishHandler()
	router.GET("/dishes", handler.GetDishes)

	dishes := []*models.DishResponse{
		{
			ID:      primitive.NewObjectID().Hex(),
			Name:    "Dish 1",
			Type:    "Veg",
			Cuisine: "Indian",
		},
		{
			ID:      primitive.NewObjectID().Hex(),
			Name:    "Dish 2",
			Type:    "Non-Veg",
			Cuisine: "Chinese",
		},
	}

	pagination := &models.PaginationResponse{
		CurrentPage: 1,
		TotalPages:  1,
		TotalItems:  2,
		HasNext:     false,
		HasPrev:     false,
	}

	mockService.On("GetAll", mock.Anything, mock.AnythingOfType("service.DishFilter"), 1, 10, (*primitive.ObjectID)(nil)).Return(dishes, pagination, nil)

	request := httptest.NewRequest(http.MethodGet, "/dishes", nil)
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

func TestDishHandler_SearchDishes_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupDishHandler()
	router.GET("/dishes/search", handler.SearchDishes)

	dishes := []*models.DishResponse{
		{
			ID:   primitive.NewObjectID().Hex(),
			Name: "Paneer Butter Masala",
			Type: "Veg",
		},
	}

	pagination := &models.PaginationResponse{
		CurrentPage: 1,
		TotalPages:  1,
		TotalItems:  1,
		HasNext:     false,
		HasPrev:     false,
	}

	mockService.On("Search", mock.Anything, "paneer", mock.AnythingOfType("service.DishFilter"), 1, 10, (*primitive.ObjectID)(nil)).Return(dishes, pagination, nil)

	request := httptest.NewRequest(http.MethodGet, "/dishes/search?q=paneer", nil)
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

func TestDishHandler_SearchDishes_MissingQuery(t *testing.T) {
	// Arrange
	handler, _, router := setupDishHandler()
	router.GET("/dishes/search", handler.SearchDishes)

	request := httptest.NewRequest(http.MethodGet, "/dishes/search", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Query parameter 'q' is required")
}

func TestDishHandler_CreateDish_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupDishHandler()
	router.POST("/dishes", handler.CreateDish)

	dish := models.Dish{
		Name:     "New Dish",
		Type:     "Veg",
		Cuisine:  "Indian",
		Calories: 300,
	}

	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Dish")).Return(nil)

	requestBody, _ := json.Marshal(dish)
	request := httptest.NewRequest(http.MethodPost, "/dishes", bytes.NewBuffer(requestBody))
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
	assert.Equal(t, "Dish created successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestDishHandler_CreateDish_InvalidJSON(t *testing.T) {
	// Arrange
	handler, _, router := setupDishHandler()
	router.POST("/dishes", handler.CreateDish)

	request := httptest.NewRequest(http.MethodPost, "/dishes", bytes.NewBufferString("invalid json"))
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

func TestDishHandler_UpdateDish_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupDishHandler()
	router.PUT("/dishes/:id", handler.UpdateDish)

	dishID := primitive.NewObjectID()
	dish := models.Dish{
		Name:     "Updated Dish",
		Type:     "Veg",
		Cuisine:  "Indian",
		Calories: 350,
	}

	mockService.On("Update", mock.Anything, dishID, mock.AnythingOfType("*models.Dish")).Return(nil)

	requestBody, _ := json.Marshal(dish)
	request := httptest.NewRequest(http.MethodPut, "/dishes/"+dishID.Hex(), bytes.NewBuffer(requestBody))
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
	assert.Equal(t, "Dish updated successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestDishHandler_DeleteDish_Success(t *testing.T) {
	// Arrange
	handler, mockService, router := setupDishHandler()
	router.DELETE("/dishes/:id", handler.DeleteDish)

	dishID := primitive.NewObjectID()

	mockService.On("Delete", mock.Anything, dishID).Return(nil)

	request := httptest.NewRequest(http.MethodDelete, "/dishes/"+dishID.Hex(), nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Dish deleted successfully", response.Message)
	
	mockService.AssertExpectations(t)
}