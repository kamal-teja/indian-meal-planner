package models

import "time"

// AnalyticsResponse represents analytics data for meals
type AnalyticsResponse struct {
	TotalMeals           int                `json:"totalMeals"`
	AvgCaloriesPerDay    float64            `json:"avgCaloriesPerDay"`
	MealTypeDistribution map[string]int     `json:"mealTypeDistribution"`
	CuisineDistribution  map[string]int     `json:"cuisineDistribution"`
	TopDishes            []DishPopularity   `json:"topDishes"`
	NutritionSummary     NutritionAnalytics `json:"nutritionSummary"`
	WeeklyTrend          []DailyMealCount   `json:"weeklyTrend"`
	Period               int                `json:"period"`
}

// DishPopularity represents a dish with its popularity count
type DishPopularity struct {
	DishID   string `json:"dishId"`
	DishName string `json:"dishName"`
	Count    int    `json:"count"`
	Calories int    `json:"calories"`
	Cuisine  string `json:"cuisine"`
}

// NutritionAnalytics represents nutrition analytics
type NutritionAnalytics struct {
	TotalCalories int     `json:"totalCalories"`
	AvgCalories   float64 `json:"avgCalories"`
	TotalProtein  int     `json:"totalProtein"`
	TotalCarbs    int     `json:"totalCarbs"`
	TotalFat      int     `json:"totalFat"`
	TotalFiber    int     `json:"totalFiber"`
	AvgProtein    float64 `json:"avgProtein"`
	AvgCarbs      float64 `json:"avgCarbs"`
	AvgFat        float64 `json:"avgFat"`
}

// DailyMealCount represents meal count for a specific day
type DailyMealCount struct {
	Date      time.Time `json:"date"`
	MealCount int       `json:"mealCount"`
	Calories  int       `json:"calories"`
}

// ShoppingListResponse represents a shopping list
type ShoppingListResponse struct {
	Ingredients []IngredientItem `json:"ingredients"`
	TotalItems  int              `json:"totalItems"`
	DateRange   string           `json:"dateRange"`
}

// IngredientItem represents an ingredient in shopping list
type IngredientItem struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Category string `json:"category"`
	Count    int    `json:"count"` // How many dishes use this ingredient
}

// RecommendationsResponse represents meal recommendations
type RecommendationsResponse struct {
	Recommendations []RecommendedDish `json:"recommendations"`
	Reason          string            `json:"reason"`
}

// RecommendedDish represents a recommended dish
type RecommendedDish struct {
	DishID     string  `json:"dishId"`
	DishName   string  `json:"dishName"`
	Cuisine    string  `json:"cuisine"`
	Calories   int     `json:"calories"`
	Score      float64 `json:"score"`
	Reason     string  `json:"reason"`
	Image      string  `json:"image"`
	PrepTime   int     `json:"prepTime"`
	Difficulty string  `json:"difficulty"`
}

// NutritionProgressResponse represents nutrition progress
type NutritionProgressResponse struct {
	Period   int              `json:"period"`
	Progress []DailyNutrition `json:"progress"`
	Goals    NutritionGoals   `json:"goals"`
	Summary  NutritionSummary `json:"summary"`
}

// DailyNutrition represents nutrition for a specific day
type DailyNutrition struct {
	Date      time.Time `json:"date"`
	Calories  int       `json:"calories"`
	Protein   int       `json:"protein"`
	Carbs     int       `json:"carbs"`
	Fat       int       `json:"fat"`
	Fiber     int       `json:"fiber"`
	Sodium    int       `json:"sodium"`
	MealCount int       `json:"mealCount"`
}

// NutritionSummary represents aggregated nutrition data
type NutritionSummary struct {
	AvgCalories    float64 `json:"avgCalories"`
	AvgProtein     float64 `json:"avgProtein"`
	AvgCarbs       float64 `json:"avgCarbs"`
	AvgFat         float64 `json:"avgFat"`
	AvgFiber       float64 `json:"avgFiber"`
	TotalDays      int     `json:"totalDays"`
	CalorieGoalMet int     `json:"calorieGoalMet"`
	ProteinGoalMet int     `json:"proteinGoalMet"`
	GoalPercentage float64 `json:"goalPercentage"`
}

// NutritionGoalsRequest represents the request to update nutrition goals
type NutritionGoalsRequest struct {
	DailyCalories int `json:"dailyCalories" validate:"min=0"`
	DailyProtein  int `json:"dailyProtein" validate:"min=0"`
	DailyCarbs    int `json:"dailyCarbs" validate:"min=0"`
	DailyFat      int `json:"dailyFat" validate:"min=0"`
	DailyFiber    int `json:"dailyFiber" validate:"min=0"`
	DailySodium   int `json:"dailySodium" validate:"min=0"`
}
