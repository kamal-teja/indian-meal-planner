package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetValidMealTypes(t *testing.T) {
	// Act
	mealTypes := GetValidMealTypes()

	// Assert
	assert.NotEmpty(t, mealTypes)
	assert.Contains(t, mealTypes, "breakfast")
	assert.Contains(t, mealTypes, "lunch")
	assert.Contains(t, mealTypes, "dinner")
	assert.Contains(t, mealTypes, "snack")
	assert.Equal(t, 4, len(mealTypes))
}

func TestMeal(t *testing.T) {
	// Arrange
	mealID := primitive.NewObjectID()
	dishID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	now := time.Now()

	meal := Meal{
		ID:        mealID,
		Date:      now,
		MealType:  "breakfast",
		DishID:    dishID,
		UserID:    userID,
		Notes:     "Delicious breakfast",
		Rating:    5,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Assert
	assert.Equal(t, mealID, meal.ID)
	assert.Equal(t, now, meal.Date)
	assert.Equal(t, "breakfast", meal.MealType)
	assert.Equal(t, dishID, meal.DishID)
	assert.Equal(t, userID, meal.UserID)
	assert.Equal(t, "Delicious breakfast", meal.Notes)
	assert.Equal(t, 5, meal.Rating)
	assert.Equal(t, now, meal.CreatedAt)
	assert.Equal(t, now, meal.UpdatedAt)
}

func TestMealWithDish(t *testing.T) {
	// Arrange
	dishResponse := DishResponse{
		ID:          "507f1f77bcf86cd799439011",
		Name:        "Oatmeal",
		Type:        "Veg",
		Cuisine:     "Continental",
		Calories:    150,
		IsFavorite:  false,
	}

	now := time.Now()
	mealWithDish := MealWithDish{
		ID:        "507f1f77bcf86cd799439012",
		Date:      now,
		MealType:  "breakfast",
		Dish:      dishResponse,
		User:      "507f1f77bcf86cd799439013",
		Notes:     "Healthy start to the day",
		Rating:    4,
		CreatedAt: now,
	}

	// Assert
	assert.Equal(t, "507f1f77bcf86cd799439012", mealWithDish.ID)
	assert.Equal(t, now, mealWithDish.Date)
	assert.Equal(t, "breakfast", mealWithDish.MealType)
	assert.Equal(t, dishResponse, mealWithDish.Dish)
	assert.Equal(t, "507f1f77bcf86cd799439013", mealWithDish.User)
	assert.Equal(t, "Healthy start to the day", mealWithDish.Notes)
	assert.Equal(t, 4, mealWithDish.Rating)
	assert.Equal(t, now, mealWithDish.CreatedAt)
}

func TestMealRequest(t *testing.T) {
	// Arrange
	flexDate := FlexibleDate{Time: time.Now()}
	req := MealRequest{
		Date:     flexDate,
		MealType: "lunch",
		DishID:   "507f1f77bcf86cd799439011",
		Notes:    "Lunch with colleagues",
		Rating:   4,
	}

	// Assert
	assert.Equal(t, flexDate, req.Date)
	assert.Equal(t, "lunch", req.MealType)
	assert.Equal(t, "507f1f77bcf86cd799439011", req.DishID)
	assert.Equal(t, "Lunch with colleagues", req.Notes)
	assert.Equal(t, 4, req.Rating)
}

func TestMealPlan(t *testing.T) {
	// Arrange
	planID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 7) // 7 days later
	now := time.Now()

	mealPlanMeal := MealPlanMeal{
		Date:     startDate,
		MealType: "breakfast",
		DishID:   primitive.NewObjectID(),
		Notes:    "Start the week right",
	}

	mealPlan := MealPlan{
		ID:          planID,
		UserID:      userID,
		Name:        "Weekly Meal Plan",
		StartDate:   startDate,
		EndDate:     endDate,
		Description: "Healthy week ahead",
		Meals:       []MealPlanMeal{mealPlanMeal},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Assert
	assert.Equal(t, planID, mealPlan.ID)
	assert.Equal(t, userID, mealPlan.UserID)
	assert.Equal(t, "Weekly Meal Plan", mealPlan.Name)
	assert.Equal(t, startDate, mealPlan.StartDate)
	assert.Equal(t, endDate, mealPlan.EndDate)
	assert.Equal(t, "Healthy week ahead", mealPlan.Description)
	assert.Len(t, mealPlan.Meals, 1)
	assert.Equal(t, mealPlanMeal, mealPlan.Meals[0])
	assert.Equal(t, now, mealPlan.CreatedAt)
	assert.Equal(t, now, mealPlan.UpdatedAt)
}

func TestMealPlanMeal(t *testing.T) {
	// Arrange
	date := time.Now()
	dishID := primitive.NewObjectID()

	mealPlanMeal := MealPlanMeal{
		Date:     date,
		MealType: "dinner",
		DishID:   dishID,
		Notes:    "Special dinner",
	}

	// Assert
	assert.Equal(t, date, mealPlanMeal.Date)
	assert.Equal(t, "dinner", mealPlanMeal.MealType)
	assert.Equal(t, dishID, mealPlanMeal.DishID)
	assert.Equal(t, "Special dinner", mealPlanMeal.Notes)
}

func TestMealPlanRequest(t *testing.T) {
	// Arrange
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 14) // 14 days later

	req := MealPlanRequest{
		Name:        "Two Week Plan",
		StartDate:   startDate,
		EndDate:     endDate,
		Description: "Balanced nutrition for two weeks",
	}

	// Assert
	assert.Equal(t, "Two Week Plan", req.Name)
	assert.Equal(t, startDate, req.StartDate)
	assert.Equal(t, endDate, req.EndDate)
	assert.Equal(t, "Balanced nutrition for two weeks", req.Description)
}

func TestMealValidation(t *testing.T) {
	tests := []struct {
		name     string
		meal     Meal
		desc     string
	}{
		{
			name: "valid breakfast meal",
			meal: Meal{
				Date:     time.Now(),
				MealType: "breakfast",
				DishID:   primitive.NewObjectID(),
				UserID:   primitive.NewObjectID(),
				Rating:   5,
			},
			desc: "Should be valid breakfast meal",
		},
		{
			name: "valid lunch meal",
			meal: Meal{
				Date:     time.Now(),
				MealType: "lunch",
				DishID:   primitive.NewObjectID(),
				UserID:   primitive.NewObjectID(),
				Rating:   3,
			},
			desc: "Should be valid lunch meal",
		},
		{
			name: "valid dinner meal",
			meal: Meal{
				Date:     time.Now(),
				MealType: "dinner",
				DishID:   primitive.NewObjectID(),
				UserID:   primitive.NewObjectID(),
				Rating:   0, // No rating
			},
			desc: "Should be valid dinner meal without rating",
		},
		{
			name: "valid snack meal",
			meal: Meal{
				Date:     time.Now(),
				MealType: "snack",
				DishID:   primitive.NewObjectID(),
				UserID:   primitive.NewObjectID(),
				Notes:    "Healthy snack",
				Rating:   4,
			},
			desc: "Should be valid snack meal with notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that meal struct is properly initialized
			assert.NotZero(t, tt.meal.Date)
			assert.Contains(t, GetValidMealTypes(), tt.meal.MealType)
			assert.NotEqual(t, primitive.NilObjectID, tt.meal.DishID)
			assert.NotEqual(t, primitive.NilObjectID, tt.meal.UserID)
			assert.GreaterOrEqual(t, tt.meal.Rating, 0)
			assert.LessOrEqual(t, tt.meal.Rating, 5)
		})
	}
}