package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnalyticsResponse(t *testing.T) {
	// Arrange
	mealTypeDistribution := map[string]int{
		"breakfast": 7,
		"lunch":     7,
		"dinner":    7,
		"snack":     3,
	}
	
	cuisineDistribution := map[string]int{
		"North Indian": 10,
		"South Indian": 8,
		"Chinese":      4,
		"Italian":      2,
	}

	topDishes := []DishPopularity{
		{
			DishID:   "1",
			DishName: "Chicken Biryani",
			Count:    5,
			Calories: 450,
			Cuisine:  "North Indian",
		},
		{
			DishID:   "2",
			DishName: "Masala Dosa",
			Count:    4,
			Calories: 300,
			Cuisine:  "South Indian",
		},
	}

	nutritionSummary := NutritionAnalytics{
		TotalCalories: 8400,
		AvgCalories:   350.0,
		TotalProtein:  600,
		TotalCarbs:    1200,
		TotalFat:      280,
		TotalFiber:    180,
		AvgProtein:    25.0,
		AvgCarbs:      50.0,
		AvgFat:        11.7,
	}

	weeklyTrend := []DailyMealCount{
		{
			Date:      time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			MealCount: 3,
			Calories:  1200,
		},
		{
			Date:      time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			MealCount: 4,
			Calories:  1400,
		},
	}

	analytics := AnalyticsResponse{
		TotalMeals:           24,
		AvgCaloriesPerDay:    1200.0,
		MealTypeDistribution: mealTypeDistribution,
		CuisineDistribution:  cuisineDistribution,
		TopDishes:            topDishes,
		NutritionSummary:     nutritionSummary,
		WeeklyTrend:          weeklyTrend,
		Period:               7,
	}

	// Assert
	assert.Equal(t, 24, analytics.TotalMeals)
	assert.Equal(t, 1200.0, analytics.AvgCaloriesPerDay)
	assert.Equal(t, mealTypeDistribution, analytics.MealTypeDistribution)
	assert.Equal(t, cuisineDistribution, analytics.CuisineDistribution)
	assert.Equal(t, topDishes, analytics.TopDishes)
	assert.Equal(t, nutritionSummary, analytics.NutritionSummary)
	assert.Equal(t, weeklyTrend, analytics.WeeklyTrend)
	assert.Equal(t, 7, analytics.Period)
}

func TestDishPopularity(t *testing.T) {
	// Arrange
	dish := DishPopularity{
		DishID:   "507f1f77bcf86cd799439011",
		DishName: "Butter Chicken",
		Count:    8,
		Calories: 420,
		Cuisine:  "North Indian",
	}

	// Assert
	assert.Equal(t, "507f1f77bcf86cd799439011", dish.DishID)
	assert.Equal(t, "Butter Chicken", dish.DishName)
	assert.Equal(t, 8, dish.Count)
	assert.Equal(t, 420, dish.Calories)
	assert.Equal(t, "North Indian", dish.Cuisine)
}

func TestNutritionAnalytics(t *testing.T) {
	// Arrange
	nutrition := NutritionAnalytics{
		TotalCalories: 10500,
		AvgCalories:   350.0,
		TotalProtein:  750,
		TotalCarbs:    1500,
		TotalFat:      350,
		TotalFiber:    225,
		AvgProtein:    25.0,
		AvgCarbs:      50.0,
		AvgFat:        11.7,
	}

	// Assert
	assert.Equal(t, 10500, nutrition.TotalCalories)
	assert.Equal(t, 350.0, nutrition.AvgCalories)
	assert.Equal(t, 750, nutrition.TotalProtein)
	assert.Equal(t, 1500, nutrition.TotalCarbs)
	assert.Equal(t, 350, nutrition.TotalFat)
	assert.Equal(t, 225, nutrition.TotalFiber)
	assert.Equal(t, 25.0, nutrition.AvgProtein)
	assert.Equal(t, 50.0, nutrition.AvgCarbs)
	assert.Equal(t, 11.7, nutrition.AvgFat)
}

func TestDailyMealCount(t *testing.T) {
	// Arrange
	date := time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)
	dailyCount := DailyMealCount{
		Date:      date,
		MealCount: 3,
		Calories:  1200,
	}

	// Assert
	assert.Equal(t, date, dailyCount.Date)
	assert.Equal(t, 3, dailyCount.MealCount)
	assert.Equal(t, 1200, dailyCount.Calories)
}

func TestShoppingListResponse(t *testing.T) {
	// Arrange
	ingredients := []IngredientItem{
		{
			Name:     "Chicken",
			Quantity: "1 kg",
			Category: "Meat",
			Count:    3,
		},
		{
			Name:     "Rice",
			Quantity: "2 cups",
			Category: "Grains",
			Count:    5,
		},
	}

	shoppingList := ShoppingListResponse{
		Ingredients: ingredients,
		TotalItems:  2,
		DateRange:   "2023-10-15 to 2023-10-21",
	}

	// Assert
	assert.Equal(t, ingredients, shoppingList.Ingredients)
	assert.Equal(t, 2, shoppingList.TotalItems)
	assert.Equal(t, "2023-10-15 to 2023-10-21", shoppingList.DateRange)
}

func TestIngredientItem(t *testing.T) {
	// Arrange
	ingredient := IngredientItem{
		Name:     "Tomatoes",
		Quantity: "500g",
		Category: "Vegetables",
		Count:    4,
	}

	// Assert
	assert.Equal(t, "Tomatoes", ingredient.Name)
	assert.Equal(t, "500g", ingredient.Quantity)
	assert.Equal(t, "Vegetables", ingredient.Category)
	assert.Equal(t, 4, ingredient.Count)
}

func TestRecommendationsResponse(t *testing.T) {
	// Arrange
	recommendations := []RecommendedDish{
		{
			DishID:     "1",
			DishName:   "Paneer Makhani",
			Cuisine:    "North Indian",
			Calories:   380,
			Score:      0.85,
			Reason:     "Matches your vegetarian preference",
			Image:      "paneer.jpg",
			PrepTime:   25,
			Difficulty: "medium",
		},
	}

	response := RecommendationsResponse{
		Recommendations: recommendations,
		Reason:          "Based on your dietary preferences and past meals",
	}

	// Assert
	assert.Equal(t, recommendations, response.Recommendations)
	assert.Equal(t, "Based on your dietary preferences and past meals", response.Reason)
	assert.Len(t, response.Recommendations, 1)
	assert.Equal(t, "Paneer Makhani", response.Recommendations[0].DishName)
}

func TestRecommendedDish(t *testing.T) {
	// Arrange
	dish := RecommendedDish{
		DishID:     "507f1f77bcf86cd799439011",
		DishName:   "Dal Makhani",
		Cuisine:    "North Indian",
		Calories:   250,
		Score:      0.92,
		Reason:     "High protein content",
		Image:      "dal.jpg",
		PrepTime:   60,
		Difficulty: "medium",
	}

	// Assert
	assert.Equal(t, "507f1f77bcf86cd799439011", dish.DishID)
	assert.Equal(t, "Dal Makhani", dish.DishName)
	assert.Equal(t, "North Indian", dish.Cuisine)
	assert.Equal(t, 250, dish.Calories)
	assert.Equal(t, 0.92, dish.Score)
	assert.Equal(t, "High protein content", dish.Reason)
	assert.Equal(t, "dal.jpg", dish.Image)
	assert.Equal(t, 60, dish.PrepTime)
	assert.Equal(t, "medium", dish.Difficulty)
}

func TestNutritionProgressResponse(t *testing.T) {
	// Arrange
	progress := []DailyNutrition{
		{
			Date:      time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			Calories:  1200,
			Protein:   60,
			Carbs:     150,
			Fat:       40,
			Fiber:     25,
			Sodium:    1800,
			MealCount: 3,
		},
	}

	goals := NutritionGoals{
		DailyCalories: 2000,
		Protein:       100,
		Carbs:         250,
		Fat:           67,
		Fiber:         25,
		Sodium:        2300,
	}

	summary := NutritionSummary{
		AvgCalories:    1150.0,
		AvgProtein:     58.0,
		AvgCarbs:       145.0,
		AvgFat:         38.0,
		AvgFiber:       23.0,
		TotalDays:      7,
		CalorieGoalMet: 4,
		ProteinGoalMet: 3,
		GoalPercentage: 57.1,
	}

	response := NutritionProgressResponse{
		Period:   7,
		Progress: progress,
		Goals:    goals,
		Summary:  summary,
	}

	// Assert
	assert.Equal(t, 7, response.Period)
	assert.Equal(t, progress, response.Progress)
	assert.Equal(t, goals, response.Goals)
	assert.Equal(t, summary, response.Summary)
}

func TestDailyNutrition(t *testing.T) {
	// Arrange
	date := time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)
	nutrition := DailyNutrition{
		Date:      date,
		Calories:  1350,
		Protein:   70,
		Carbs:     180,
		Fat:       45,
		Fiber:     28,
		Sodium:    2100,
		MealCount: 4,
	}

	// Assert
	assert.Equal(t, date, nutrition.Date)
	assert.Equal(t, 1350, nutrition.Calories)
	assert.Equal(t, 70, nutrition.Protein)
	assert.Equal(t, 180, nutrition.Carbs)
	assert.Equal(t, 45, nutrition.Fat)
	assert.Equal(t, 28, nutrition.Fiber)
	assert.Equal(t, 2100, nutrition.Sodium)
	assert.Equal(t, 4, nutrition.MealCount)
}

func TestNutritionSummary(t *testing.T) {
	// Arrange
	summary := NutritionSummary{
		AvgCalories:    1250.5,
		AvgProtein:     65.2,
		AvgCarbs:       165.8,
		AvgFat:         42.3,
		AvgFiber:       26.1,
		TotalDays:      14,
		CalorieGoalMet: 9,
		ProteinGoalMet: 11,
		GoalPercentage: 64.3,
	}

	// Assert
	assert.Equal(t, 1250.5, summary.AvgCalories)
	assert.Equal(t, 65.2, summary.AvgProtein)
	assert.Equal(t, 165.8, summary.AvgCarbs)
	assert.Equal(t, 42.3, summary.AvgFat)
	assert.Equal(t, 26.1, summary.AvgFiber)
	assert.Equal(t, 14, summary.TotalDays)
	assert.Equal(t, 9, summary.CalorieGoalMet)
	assert.Equal(t, 11, summary.ProteinGoalMet)
	assert.Equal(t, 64.3, summary.GoalPercentage)
}

func TestNutritionGoalsRequest(t *testing.T) {
	// Arrange
	request := NutritionGoalsRequest{
		DailyCalories: 2200,
		Protein:       110,
		Carbs:         275,
		Fat:           73,
		Fiber:         30,
		Sodium:        2000,
	}

	// Assert
	assert.Equal(t, 2200, request.DailyCalories)
	assert.Equal(t, 110, request.Protein)
	assert.Equal(t, 275, request.Carbs)
	assert.Equal(t, 73, request.Fat)
	assert.Equal(t, 30, request.Fiber)
	assert.Equal(t, 2000, request.Sodium)
}