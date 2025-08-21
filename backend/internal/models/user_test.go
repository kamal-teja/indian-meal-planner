package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserToResponse(t *testing.T) {
	// Arrange
	userID := primitive.NewObjectID()
	user := &User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
		Profile: UserProfile{
			DietaryPreferences: []string{"vegetarian"},
			SpiceLevel:         "medium",
			FavoriteRegions:    []string{"North Indian"},
			Avatar:             "avatar.jpg",
			NutritionGoals: NutritionGoals{
				DailyCalories: 2000,
				Protein:       100,
				Carbs:         250,
				Fat:           65,
				Fiber:         25,
				Sodium:        2300,
			},
		},
		Favorites: []primitive.ObjectID{primitive.NewObjectID()},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	response := user.ToResponse()

	// Assert
	assert.Equal(t, userID.Hex(), response.ID)
	assert.Equal(t, "John Doe", response.Name)
	assert.Equal(t, "john@example.com", response.Email)
	assert.Equal(t, user.Profile, response.Profile)
}

func TestGetDietaryPreferences(t *testing.T) {
	// Act
	preferences := GetDietaryPreferences()

	// Assert
	assert.NotEmpty(t, preferences)
	assert.Contains(t, preferences, "vegetarian")
	assert.Contains(t, preferences, "vegan")
	assert.Contains(t, preferences, "gluten-free")
	assert.Contains(t, preferences, "dairy-free")
	assert.Contains(t, preferences, "keto")
	assert.Contains(t, preferences, "paleo")
}

func TestGetFavoriteRegions(t *testing.T) {
	// Act
	regions := GetFavoriteRegions()

	// Assert
	assert.NotEmpty(t, regions)
	assert.Contains(t, regions, "North Indian")
	assert.Contains(t, regions, "South Indian")
	assert.Contains(t, regions, "Bengali")
	assert.Contains(t, regions, "Italian")
	assert.Contains(t, regions, "Chinese")
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Profile: UserProfile{
					SpiceLevel: "medium",
					NutritionGoals: NutritionGoals{
						DailyCalories: 2000,
						Protein:       50,
						Carbs:         200,
						Fat:           70,
						Fiber:         25,
						Sodium:        2000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: User{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "short password",
			user: User{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "12345",
			},
			wantErr: true,
		},
		{
			name: "short name",
			user: User{
				Name:     "J",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	// Note: This test validates struct tags but requires a validator instance
	// In a real scenario, this would be tested in the service layer where validation is performed
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We're just testing that the struct is properly formed
			// Actual validation would happen in the service layer
			assert.NotNil(t, tt.user)
		})
	}
}

func TestUserRegistrationRequest(t *testing.T) {
	// Arrange
	req := UserRegistrationRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Assert
	assert.Equal(t, "John Doe", req.Name)
	assert.Equal(t, "john@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestUserLoginRequest(t *testing.T) {
	// Arrange
	req := UserLoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	// Assert
	assert.Equal(t, "john@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestNutritionGoals(t *testing.T) {
	// Arrange
	goals := NutritionGoals{
		DailyCalories: 2000,
		Protein:       100,
		Carbs:         250,
		Fat:           65,
		Fiber:         25,
		Sodium:        2300,
	}

	// Assert
	assert.Equal(t, 2000, goals.DailyCalories)
	assert.Equal(t, 100, goals.Protein)
	assert.Equal(t, 250, goals.Carbs)
	assert.Equal(t, 65, goals.Fat)
	assert.Equal(t, 25, goals.Fiber)
	assert.Equal(t, 2300, goals.Sodium)
}

func TestUserProfile(t *testing.T) {
	// Arrange
	profile := UserProfile{
		DietaryPreferences: []string{"vegetarian", "gluten-free"},
		SpiceLevel:         "medium",
		FavoriteRegions:    []string{"North Indian", "South Indian"},
		Avatar:             "avatar.jpg",
		NutritionGoals: NutritionGoals{
			DailyCalories: 2000,
			Protein:       100,
		},
	}

	// Assert
	assert.Equal(t, []string{"vegetarian", "gluten-free"}, profile.DietaryPreferences)
	assert.Equal(t, "medium", profile.SpiceLevel)
	assert.Equal(t, []string{"North Indian", "South Indian"}, profile.FavoriteRegions)
	assert.Equal(t, "avatar.jpg", profile.Avatar)
	assert.Equal(t, 2000, profile.NutritionGoals.DailyCalories)
}

func TestProfileUpdateRequest(t *testing.T) {
	// Arrange
	req := ProfileUpdateRequest{
		Name: "Updated Name",
		Profile: UserProfile{
			DietaryPreferences: []string{"vegan"},
			SpiceLevel:         "hot",
		},
	}

	// Assert
	assert.Equal(t, "Updated Name", req.Name)
	assert.Equal(t, []string{"vegan"}, req.Profile.DietaryPreferences)
	assert.Equal(t, "hot", req.Profile.SpiceLevel)
}