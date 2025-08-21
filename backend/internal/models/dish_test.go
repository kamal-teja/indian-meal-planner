package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDishToResponse(t *testing.T) {
	// Arrange
	dishID := primitive.NewObjectID()
	dish := &Dish{
		ID:          dishID,
		Name:        "Chicken Biryani",
		Type:        "Non-Veg",
		Cuisine:     "North Indian",
		Image:       "biryani.jpg",
		Ingredients: []string{"chicken", "rice", "spices"},
		Calories:    450,
		Nutrition: Nutrition{
			Protein: 25,
			Carbs:   50,
			Fat:     15,
			Fiber:   3,
			Sugar:   5,
			Sodium:  800,
		},
		DietaryTags: []string{"high-protein"},
		SpiceLevel:  "medium",
		PrepTime:    30,
		CookTime:    45,
		Servings:    4,
		Difficulty:  "medium",
		Description: "Delicious chicken biryani",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	response := dish.ToResponse()

	// Assert
	assert.Equal(t, dishID.Hex(), response.ID)
	assert.Equal(t, "Chicken Biryani", response.Name)
	assert.Equal(t, "Non-Veg", response.Type)
	assert.Equal(t, "North Indian", response.Cuisine)
	assert.Equal(t, "biryani.jpg", response.Image)
	assert.Equal(t, []string{"chicken", "rice", "spices"}, response.Ingredients)
	assert.Equal(t, 450, response.Calories)
	assert.Equal(t, dish.Nutrition, response.Nutrition)
	assert.Equal(t, []string{"high-protein"}, response.DietaryTags)
	assert.Equal(t, "medium", response.SpiceLevel)
	assert.Equal(t, 30, response.PrepTime)
	assert.Equal(t, 45, response.CookTime)
	assert.Equal(t, 4, response.Servings)
	assert.Equal(t, "medium", response.Difficulty)
	assert.Equal(t, "Delicious chicken biryani", response.Description)
	assert.False(t, response.IsFavorite) // Default value
}

func TestGetValidDietaryTags(t *testing.T) {
	// Act
	tags := GetValidDietaryTags()

	// Assert
	assert.NotEmpty(t, tags)
	assert.Contains(t, tags, "vegetarian")
	assert.Contains(t, tags, "vegan")
	assert.Contains(t, tags, "gluten-free")
	assert.Contains(t, tags, "dairy-free")
	assert.Contains(t, tags, "nut-free")
	assert.Contains(t, tags, "keto")
	assert.Contains(t, tags, "paleo")
	assert.Contains(t, tags, "low-carb")
	assert.Contains(t, tags, "high-protein")
	assert.Contains(t, tags, "low-sodium")
	assert.Contains(t, tags, "sugar-free")
}

func TestGetValidCuisines(t *testing.T) {
	// Act
	cuisines := GetValidCuisines()

	// Assert
	assert.NotEmpty(t, cuisines)
	assert.Contains(t, cuisines, "North Indian")
	assert.Contains(t, cuisines, "South Indian")
	assert.Contains(t, cuisines, "Bengali")
	assert.Contains(t, cuisines, "Gujarati")
	assert.Contains(t, cuisines, "Punjabi")
	assert.Contains(t, cuisines, "Italian")
	assert.Contains(t, cuisines, "Chinese")
	assert.Contains(t, cuisines, "Thai")
}

func TestNutrition(t *testing.T) {
	// Arrange
	nutrition := Nutrition{
		Protein: 25,
		Carbs:   50,
		Fat:     15,
		Fiber:   3,
		Sugar:   5,
		Sodium:  800,
	}

	// Assert
	assert.Equal(t, 25, nutrition.Protein)
	assert.Equal(t, 50, nutrition.Carbs)
	assert.Equal(t, 15, nutrition.Fat)
	assert.Equal(t, 3, nutrition.Fiber)
	assert.Equal(t, 5, nutrition.Sugar)
	assert.Equal(t, 800, nutrition.Sodium)
}

func TestDishCreateRequest(t *testing.T) {
	// Arrange
	req := DishCreateRequest{
		Name:        "Paneer Butter Masala",
		Type:        "Veg",
		Cuisine:     "North Indian",
		Image:       "paneer.jpg",
		Ingredients: []string{"paneer", "tomatoes", "cream", "spices"},
		Calories:    350,
		Nutrition: Nutrition{
			Protein: 20,
			Carbs:   25,
			Fat:     20,
			Fiber:   4,
			Sugar:   8,
			Sodium:  600,
		},
		DietaryTags: []string{"vegetarian", "high-protein"},
		SpiceLevel:  "medium",
		PrepTime:    15,
		CookTime:    30,
		Servings:    4,
		Difficulty:  "easy",
		Description: "Creamy paneer curry",
	}

	// Assert
	assert.Equal(t, "Paneer Butter Masala", req.Name)
	assert.Equal(t, "Veg", req.Type)
	assert.Equal(t, "North Indian", req.Cuisine)
	assert.Equal(t, "paneer.jpg", req.Image)
	assert.Equal(t, []string{"paneer", "tomatoes", "cream", "spices"}, req.Ingredients)
	assert.Equal(t, 350, req.Calories)
	assert.Equal(t, []string{"vegetarian", "high-protein"}, req.DietaryTags)
	assert.Equal(t, "medium", req.SpiceLevel)
	assert.Equal(t, 15, req.PrepTime)
	assert.Equal(t, 30, req.CookTime)
	assert.Equal(t, 4, req.Servings)
	assert.Equal(t, "easy", req.Difficulty)
	assert.Equal(t, "Creamy paneer curry", req.Description)
}

func TestDishValidationFields(t *testing.T) {
	tests := []struct {
		name string
		dish Dish
		desc string
	}{
		{
			name: "valid vegetarian dish",
			dish: Dish{
				Name:        "Dal Tadka",
				Type:        "Veg",
				Cuisine:     "North Indian",
				Ingredients: []string{"lentils", "spices"},
				Calories:    200,
				SpiceLevel:  "mild",
				PrepTime:    10,
				CookTime:    20,
				Servings:    2,
				Difficulty:  "easy",
			},
			desc: "Should be valid vegetarian dish",
		},
		{
			name: "valid non-vegetarian dish",
			dish: Dish{
				Name:        "Chicken Curry",
				Type:        "Non-Veg",
				Cuisine:     "South Indian",
				Ingredients: []string{"chicken", "coconut", "spices"},
				Calories:    400,
				SpiceLevel:  "hot",
				PrepTime:    20,
				CookTime:    40,
				Servings:    4,
				Difficulty:  "medium",
			},
			desc: "Should be valid non-vegetarian dish",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that dish struct is properly initialized
			assert.NotEmpty(t, tt.dish.Name)
			assert.Contains(t, []string{"Veg", "Non-Veg"}, tt.dish.Type)
			assert.NotEmpty(t, tt.dish.Cuisine)
			assert.NotEmpty(t, tt.dish.Ingredients)
			assert.GreaterOrEqual(t, tt.dish.Calories, 0)
			assert.GreaterOrEqual(t, tt.dish.PrepTime, 0)
			assert.GreaterOrEqual(t, tt.dish.CookTime, 0)
			assert.GreaterOrEqual(t, tt.dish.Servings, 1)
		})
	}
}

func TestDishResponse(t *testing.T) {
	// Arrange
	response := DishResponse{
		ID:          "507f1f77bcf86cd799439011",
		Name:        "Masala Dosa",
		Type:        "Veg",
		Cuisine:     "South Indian",
		Image:       "dosa.jpg",
		Ingredients: []string{"rice", "lentils", "potatoes"},
		Calories:    300,
		Nutrition: Nutrition{
			Protein: 10,
			Carbs:   55,
			Fat:     8,
			Fiber:   5,
			Sugar:   3,
			Sodium:  400,
		},
		DietaryTags: []string{"vegetarian", "gluten-free"},
		SpiceLevel:  "mild",
		PrepTime:    120, // Including fermentation time
		CookTime:    15,
		Servings:    2,
		Difficulty:  "medium",
		Description: "Crispy fermented crepe with potato filling",
		IsFavorite:  true,
	}

	// Assert
	assert.Equal(t, "507f1f77bcf86cd799439011", response.ID)
	assert.Equal(t, "Masala Dosa", response.Name)
	assert.Equal(t, "Veg", response.Type)
	assert.Equal(t, "South Indian", response.Cuisine)
	assert.True(t, response.IsFavorite)
	assert.Equal(t, 300, response.Calories)
	assert.Equal(t, 120, response.PrepTime)
	assert.Equal(t, 15, response.CookTime)
}