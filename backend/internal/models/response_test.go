package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPaginationResponse(t *testing.T) {
	// Arrange
	pagination := PaginationResponse{
		Page:       2,
		Limit:      10,
		Total:      25,
		TotalPages: 3,
		HasNext:    true,
		HasPrev:    true,
	}

	// Assert
	assert.Equal(t, 2, pagination.Page)
	assert.Equal(t, 10, pagination.Limit)
	assert.Equal(t, 25, pagination.Total)
	assert.Equal(t, 3, pagination.TotalPages)
	assert.True(t, pagination.HasNext)
	assert.True(t, pagination.HasPrev)
}

func TestPaginationResponse_FirstPage(t *testing.T) {
	// Arrange
	pagination := PaginationResponse{
		Page:       1,
		Limit:      10,
		Total:      25,
		TotalPages: 3,
		HasNext:    true,
		HasPrev:    false,
	}

	// Assert
	assert.Equal(t, 1, pagination.Page)
	assert.True(t, pagination.HasNext)
	assert.False(t, pagination.HasPrev)
}

func TestPaginationResponse_LastPage(t *testing.T) {
	// Arrange
	pagination := PaginationResponse{
		Page:       3,
		Limit:      10,
		Total:      25,
		TotalPages: 3,
		HasNext:    false,
		HasPrev:    true,
	}

	// Assert
	assert.Equal(t, 3, pagination.Page)
	assert.False(t, pagination.HasNext)
	assert.True(t, pagination.HasPrev)
}

func TestSuccessResponse(t *testing.T) {
	// Test with data
	t.Run("with data", func(t *testing.T) {
		data := map[string]interface{}{
			"id":   1,
			"name": "test",
		}
		response := SuccessResponse{
			Success: true,
			Message: "Operation successful",
			Data:    data,
		}

		assert.True(t, response.Success)
		assert.Equal(t, "Operation successful", response.Message)
		assert.Equal(t, data, response.Data)
	})

	// Test without data
	t.Run("without data", func(t *testing.T) {
		response := SuccessResponse{
			Success: true,
			Message: "Operation successful",
		}

		assert.True(t, response.Success)
		assert.Equal(t, "Operation successful", response.Message)
		assert.Nil(t, response.Data)
	})

	// Test minimal response
	t.Run("minimal", func(t *testing.T) {
		response := SuccessResponse{
			Success: true,
		}

		assert.True(t, response.Success)
		assert.Empty(t, response.Message)
		assert.Nil(t, response.Data)
	})
}

func TestErrorResponse(t *testing.T) {
	// Test with details and code
	t.Run("with details and code", func(t *testing.T) {
		details := map[string]interface{}{
			"field": "email",
			"issue": "invalid format",
		}
		response := ErrorResponse{
			Success: false,
			Error:   "Validation failed",
			Details: details,
			Code:    "VALIDATION_ERROR",
		}

		assert.False(t, response.Success)
		assert.Equal(t, "Validation failed", response.Error)
		assert.Equal(t, details, response.Details)
		assert.Equal(t, "VALIDATION_ERROR", response.Code)
	})

	// Test minimal error
	t.Run("minimal", func(t *testing.T) {
		response := ErrorResponse{
			Success: false,
			Error:   "Something went wrong",
		}

		assert.False(t, response.Success)
		assert.Equal(t, "Something went wrong", response.Error)
		assert.Nil(t, response.Details)
		assert.Empty(t, response.Code)
	})
}

func TestAuthResponse(t *testing.T) {
	// Arrange
	userResponse := UserResponse{
		ID:    primitive.NewObjectID().Hex(),
		Name:  "John Doe",
		Email: "john@example.com",
		Profile: UserProfile{
			DietaryPreferences: []string{"vegetarian"},
			SpiceLevel:         "medium",
		},
	}

	authResponse := AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   "jwt-token-123",
		User:    userResponse,
	}

	// Assert
	assert.True(t, authResponse.Success)
	assert.Equal(t, "Login successful", authResponse.Message)
	assert.Equal(t, "jwt-token-123", authResponse.Token)
	assert.Equal(t, userResponse, authResponse.User)
	assert.Equal(t, "John Doe", authResponse.User.Name)
	assert.Equal(t, "john@example.com", authResponse.User.Email)
}

func TestHealthResponse(t *testing.T) {
	// Arrange
	healthResponse := HealthResponse{
		Status:    "healthy",
		Timestamp: "2023-10-15T14:30:00Z",
		Version:   "1.0.0",
		Database:  "connected",
	}

	// Assert
	assert.Equal(t, "healthy", healthResponse.Status)
	assert.Equal(t, "2023-10-15T14:30:00Z", healthResponse.Timestamp)
	assert.Equal(t, "1.0.0", healthResponse.Version)
	assert.Equal(t, "connected", healthResponse.Database)
}

func TestHealthResponse_Unhealthy(t *testing.T) {
	// Arrange
	healthResponse := HealthResponse{
		Status:    "unhealthy",
		Timestamp: "2023-10-15T14:30:00Z",
		Version:   "1.0.0",
		Database:  "disconnected",
	}

	// Assert
	assert.Equal(t, "unhealthy", healthResponse.Status)
	assert.Equal(t, "disconnected", healthResponse.Database)
}

func TestResponseInterfaces(t *testing.T) {
	// Test that response structures can hold different data types
	t.Run("success response with different data types", func(t *testing.T) {
		// String data
		stringResponse := SuccessResponse{
			Success: true,
			Data:    "test string",
		}
		assert.Equal(t, "test string", stringResponse.Data)

		// Slice data
		sliceResponse := SuccessResponse{
			Success: true,
			Data:    []string{"item1", "item2"},
		}
		assert.Equal(t, []string{"item1", "item2"}, sliceResponse.Data)

		// Map data
		mapResponse := SuccessResponse{
			Success: true,
			Data:    map[string]int{"count": 5},
		}
		assert.Equal(t, map[string]int{"count": 5}, mapResponse.Data)

		// Struct data
		type CustomData struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		structResponse := SuccessResponse{
			Success: true,
			Data:    CustomData{ID: 1, Name: "test"},
		}
		customData, ok := structResponse.Data.(CustomData)
		assert.True(t, ok)
		assert.Equal(t, 1, customData.ID)
		assert.Equal(t, "test", customData.Name)
	})

	t.Run("error response with different detail types", func(t *testing.T) {
		// String details
		stringErrResponse := ErrorResponse{
			Success: false,
			Error:   "error",
			Details: "additional info",
		}
		assert.Equal(t, "additional info", stringErrResponse.Details)

		// Map details
		mapErrResponse := ErrorResponse{
			Success: false,
			Error:   "validation error",
			Details: map[string]string{"field": "email", "message": "invalid"},
		}
		details, ok := mapErrResponse.Details.(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "email", details["field"])
		assert.Equal(t, "invalid", details["message"])
	})
}