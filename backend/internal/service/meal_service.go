package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"meal-planner-backend/internal/models"
	"meal-planner-backend/internal/repository"
	"meal-planner-backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MealService interface defines meal operations
type MealService interface {
	Create(ctx context.Context, userID primitive.ObjectID, req models.MealRequest) (*models.MealWithDish, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.MealWithDish, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.MealWithDish, *models.PaginationResponse, error)
	GetByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*models.MealWithDish, error)
	Update(ctx context.Context, id primitive.ObjectID, req models.MealRequest) (*models.MealWithDish, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetNutritionSummary(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]repository.NutritionSummary, error)
	GetAnalytics(ctx context.Context, userID primitive.ObjectID, period int) (*models.AnalyticsResponse, error)
	GetShoppingList(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*models.ShoppingListResponse, error)
	GetRecommendations(ctx context.Context, userID primitive.ObjectID, mealType string, date time.Time) (*models.RecommendationsResponse, error)
	GetNutritionProgress(ctx context.Context, userID primitive.ObjectID, period int) (*models.NutritionProgressResponse, error)
	GetNutritionGoals(ctx context.Context, userID primitive.ObjectID) (*models.NutritionGoals, error)
	UpdateNutritionGoals(ctx context.Context, userID primitive.ObjectID, req models.NutritionGoalsRequest) (*models.NutritionGoals, error)
}

// mealService implements MealService interface
type mealService struct {
	mealRepo repository.MealRepository
	dishRepo repository.DishRepository
	logger   *logger.Logger
}

// NewMealService creates a new meal service
func NewMealService(mealRepo repository.MealRepository, dishRepo repository.DishRepository, log *logger.Logger) MealService {
	return &mealService{
		mealRepo: mealRepo,
		dishRepo: dishRepo,
		logger:   log,
	}
}

// Create creates a new meal
func (s *mealService) Create(ctx context.Context, userID primitive.ObjectID, req models.MealRequest) (*models.MealWithDish, error) {
	// Validate dish ID
	dishID, err := primitive.ObjectIDFromHex(req.DishID)
	if err != nil {
		return nil, errors.New("invalid dish ID")
	}

	// Check if dish exists
	dish, err := s.dishRepo.GetByID(ctx, dishID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("dish not found")
		}
		s.logger.Error("Failed to get dish", "error", err, "dishID", req.DishID)
		return nil, errors.New("internal server error")
	}

	// Create meal
	meal := &models.Meal{
		Date:     req.Date.Time,
		MealType: req.MealType,
		DishID:   dishID,
		UserID:   userID,
		Notes:    req.Notes,
		Rating:   req.Rating,
	}

	if err := s.mealRepo.Create(ctx, meal); err != nil {
		s.logger.Error("Failed to create meal", "error", err)
		return nil, errors.New("failed to create meal")
	}

	// Return meal with dish info
	return &models.MealWithDish{
		ID:        meal.ID.Hex(),
		Date:      meal.Date,
		MealType:  meal.MealType,
		Dish:      dish.ToResponse(),
		User:      userID.Hex(),
		Notes:     meal.Notes,
		Rating:    meal.Rating,
		CreatedAt: meal.CreatedAt,
	}, nil
}

// GetByID retrieves a meal by ID with dish information
func (s *mealService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.MealWithDish, error) {
	meal, err := s.mealRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("meal not found")
		}
		s.logger.Error("Failed to get meal by ID", "error", err, "mealID", id.Hex())
		return nil, errors.New("internal server error")
	}

	// Get dish information
	dish, err := s.dishRepo.GetByID(ctx, meal.DishID)
	if err != nil {
		s.logger.Error("Failed to get dish for meal", "error", err, "dishID", meal.DishID.Hex())
		return nil, errors.New("internal server error")
	}

	return &models.MealWithDish{
		ID:        meal.ID.Hex(),
		Date:      meal.Date,
		MealType:  meal.MealType,
		Dish:      dish.ToResponse(),
		User:      meal.UserID.Hex(),
		Notes:     meal.Notes,
		Rating:    meal.Rating,
		CreatedAt: meal.CreatedAt,
	}, nil
}

// GetByUserID retrieves meals for a user with pagination
func (s *mealService) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.MealWithDish, *models.PaginationResponse, error) {
	meals, total, err := s.mealRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		s.logger.Error("Failed to get meals by user ID", "error", err, "userID", userID.Hex())
		return nil, nil, errors.New("failed to get meals")
	}

	// Get dish information for each meal
	mealsWithDish := make([]*models.MealWithDish, len(meals))
	for i, meal := range meals {
		dish, err := s.dishRepo.GetByID(ctx, meal.DishID)
		if err != nil {
			s.logger.Error("Failed to get dish for meal", "error", err, "dishID", meal.DishID.Hex())
			continue
		}

		mealsWithDish[i] = &models.MealWithDish{
			ID:        meal.ID.Hex(),
			Date:      meal.Date,
			MealType:  meal.MealType,
			Dish:      dish.ToResponse(),
			User:      meal.UserID.Hex(),
			Notes:     meal.Notes,
			Rating:    meal.Rating,
			CreatedAt: meal.CreatedAt,
		}
	}

	// Create pagination response
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	pagination := &models.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return mealsWithDish, pagination, nil
}

// GetByDateRange retrieves meals for a user within a date range
func (s *mealService) GetByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*models.MealWithDish, error) {
	meals, err := s.mealRepo.GetByUserAndDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get meals by date range", "error", err, "userID", userID.Hex())
		return nil, errors.New("failed to get meals")
	}

	// Get dish information for each meal
	mealsWithDish := make([]*models.MealWithDish, len(meals))
	for i, meal := range meals {
		dish, err := s.dishRepo.GetByID(ctx, meal.DishID)
		if err != nil {
			s.logger.Error("Failed to get dish for meal", "error", err, "dishID", meal.DishID.Hex())
			continue
		}

		mealsWithDish[i] = &models.MealWithDish{
			ID:        meal.ID.Hex(),
			Date:      meal.Date,
			MealType:  meal.MealType,
			Dish:      dish.ToResponse(),
			User:      meal.UserID.Hex(),
			Notes:     meal.Notes,
			Rating:    meal.Rating,
			CreatedAt: meal.CreatedAt,
		}
	}

	return mealsWithDish, nil
}

// Update updates a meal
func (s *mealService) Update(ctx context.Context, id primitive.ObjectID, req models.MealRequest) (*models.MealWithDish, error) {
	// Get existing meal
	existingMeal, err := s.mealRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("meal not found")
		}
		s.logger.Error("Failed to get meal for update", "error", err, "mealID", id.Hex())
		return nil, errors.New("internal server error")
	}

	// Validate dish ID
	dishID, err := primitive.ObjectIDFromHex(req.DishID)
	if err != nil {
		return nil, errors.New("invalid dish ID")
	}

	// Check if dish exists
	dish, err := s.dishRepo.GetByID(ctx, dishID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("dish not found")
		}
		s.logger.Error("Failed to get dish", "error", err, "dishID", req.DishID)
		return nil, errors.New("internal server error")
	}

	// Update meal
	existingMeal.Date = req.Date.Time
	existingMeal.MealType = req.MealType
	existingMeal.DishID = dishID
	existingMeal.Notes = req.Notes
	existingMeal.Rating = req.Rating

	if err := s.mealRepo.Update(ctx, id, existingMeal); err != nil {
		s.logger.Error("Failed to update meal", "error", err, "mealID", id.Hex())
		return nil, errors.New("failed to update meal")
	}

	// Return updated meal with dish info
	return &models.MealWithDish{
		ID:        existingMeal.ID.Hex(),
		Date:      existingMeal.Date,
		MealType:  existingMeal.MealType,
		Dish:      dish.ToResponse(),
		User:      existingMeal.UserID.Hex(),
		Notes:     existingMeal.Notes,
		Rating:    existingMeal.Rating,
		CreatedAt: existingMeal.CreatedAt,
	}, nil
}

// Delete deletes a meal
func (s *mealService) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Check if meal exists
	_, err := s.mealRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("meal not found")
		}
		s.logger.Error("Failed to get meal for deletion", "error", err, "mealID", id.Hex())
		return errors.New("internal server error")
	}

	if err := s.mealRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete meal", "error", err, "mealID", id.Hex())
		return errors.New("failed to delete meal")
	}

	return nil
}

// GetNutritionSummary gets nutrition summary for a user within a date range
func (s *mealService) GetNutritionSummary(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]repository.NutritionSummary, error) {
	summary, err := s.mealRepo.GetNutritionByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get nutrition summary", "error", err, "userID", userID.Hex())
		return nil, errors.New("failed to get nutrition summary")
	}

	return summary, nil
}

// GetAnalytics gets meal analytics for a user for the specified period
func (s *mealService) GetAnalytics(ctx context.Context, userID primitive.ObjectID, period int) (*models.AnalyticsResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -period)

	// Get meals for the period
	meals, err := s.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get meals for analytics", "error", err, "userID", userID.Hex())
		return nil, errors.New("failed to get analytics data")
	}

	analytics := &models.AnalyticsResponse{
		TotalMeals:           len(meals),
		MealTypeDistribution: make(map[string]int),
		CuisineDistribution:  make(map[string]int),
		TopDishes:            []models.DishPopularity{},
		Period:               period,
	}

	// Calculate distributions and analytics
	dishCount := make(map[string]*models.DishPopularity)
	totalCalories := 0
	totalProtein := 0
	totalCarbs := 0
	totalFat := 0
	totalFiber := 0
	dailyCounts := make(map[string]*models.DailyMealCount)

	for _, meal := range meals {
		// Meal type distribution
		analytics.MealTypeDistribution[meal.MealType]++

		// Cuisine distribution
		analytics.CuisineDistribution[meal.Dish.Cuisine]++

		// Dish popularity
		dishKey := meal.Dish.ID
		if dish, exists := dishCount[dishKey]; exists {
			dish.Count++
		} else {
			dishCount[dishKey] = &models.DishPopularity{
				DishID:   meal.Dish.ID,
				DishName: meal.Dish.Name,
				Count:    1,
				Calories: meal.Dish.Calories,
				Cuisine:  meal.Dish.Cuisine,
			}
		}

		// Nutrition totals
		totalCalories += meal.Dish.Calories
		totalProtein += meal.Dish.Nutrition.Protein
		totalCarbs += meal.Dish.Nutrition.Carbs
		totalFat += meal.Dish.Nutrition.Fat
		totalFiber += meal.Dish.Nutrition.Fiber

		// Daily counts
		dateKey := meal.Date.Format("2006-01-02")
		if daily, exists := dailyCounts[dateKey]; exists {
			daily.MealCount++
			daily.Calories += meal.Dish.Calories
		} else {
			dailyCounts[dateKey] = &models.DailyMealCount{
				Date:      meal.Date,
				MealCount: 1,
				Calories:  meal.Dish.Calories,
			}
		}
	}

	// Convert dish counts to sorted slice (top 10)
	for _, dish := range dishCount {
		analytics.TopDishes = append(analytics.TopDishes, *dish)
	}

	// Sort by count (simple bubble sort for small datasets)
	for i := 0; i < len(analytics.TopDishes)-1; i++ {
		for j := i + 1; j < len(analytics.TopDishes); j++ {
			if analytics.TopDishes[i].Count < analytics.TopDishes[j].Count {
				analytics.TopDishes[i], analytics.TopDishes[j] = analytics.TopDishes[j], analytics.TopDishes[i]
			}
		}
	}

	// Limit to top 10
	if len(analytics.TopDishes) > 10 {
		analytics.TopDishes = analytics.TopDishes[:10]
	}

	// Calculate averages
	if period > 0 {
		analytics.AvgCaloriesPerDay = float64(totalCalories) / float64(period)
	}

	analytics.NutritionSummary = models.NutritionAnalytics{
		TotalCalories: totalCalories,
		TotalProtein:  totalProtein,
		TotalCarbs:    totalCarbs,
		TotalFat:      totalFat,
		TotalFiber:    totalFiber,
		AvgCalories:   float64(totalCalories) / float64(len(meals)),
		AvgProtein:    float64(totalProtein) / float64(len(meals)),
		AvgCarbs:      float64(totalCarbs) / float64(len(meals)),
		AvgFat:        float64(totalFat) / float64(len(meals)),
	}

	// Convert daily counts to slice
	for _, daily := range dailyCounts {
		analytics.WeeklyTrend = append(analytics.WeeklyTrend, *daily)
	}

	return analytics, nil
}

// GetShoppingList generates a shopping list for the specified date range
func (s *mealService) GetShoppingList(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*models.ShoppingListResponse, error) {
	// Get meals for the date range
	meals, err := s.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get meals for shopping list", "error", err, "userID", userID.Hex())
		return nil, errors.New("failed to get meals for shopping list")
	}

	// Aggregate ingredients
	ingredientMap := make(map[string]*models.IngredientItem)

	for _, meal := range meals {
		for _, ingredient := range meal.Dish.Ingredients {
			if item, exists := ingredientMap[ingredient]; exists {
				item.Count++
			} else {
				ingredientMap[ingredient] = &models.IngredientItem{
					Name:     ingredient,
					Quantity: "1 unit", // Default quantity
					Category: categorizeIngredient(ingredient),
					Count:    1,
				}
			}
		}
	}

	// Convert map to slice
	ingredients := make([]models.IngredientItem, 0, len(ingredientMap))
	for _, item := range ingredientMap {
		// Adjust quantity based on count
		if item.Count > 1 {
			item.Quantity = string(rune(item.Count+'0')) + " units"
		}
		ingredients = append(ingredients, *item)
	}

	dateRange := startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02")

	response := &models.ShoppingListResponse{
		Ingredients: ingredients,
		TotalItems:  len(ingredients),
		DateRange:   dateRange,
	}

	return response, nil
}

// categorizeIngredient provides basic categorization for ingredients
func categorizeIngredient(ingredient string) string {
	// Basic categorization logic - can be enhanced
	ingredient = ingredient
	switch {
	case contains(ingredient, "rice", "wheat", "flour", "bread"):
		return "Grains"
	case contains(ingredient, "onion", "tomato", "potato", "carrot", "peas", "beans"):
		return "Vegetables"
	case contains(ingredient, "chicken", "fish", "meat", "egg"):
		return "Protein"
	case contains(ingredient, "milk", "cheese", "yogurt", "butter"):
		return "Dairy"
	case contains(ingredient, "oil", "ghee", "salt", "sugar", "spice"):
		return "Pantry"
	default:
		return "Others"
	}
}

// GetRecommendations gets meal recommendations based on user preferences and history
func (s *mealService) GetRecommendations(ctx context.Context, userID primitive.ObjectID, mealType string, date time.Time) (*models.RecommendationsResponse, error) {
	// Get user's recent meals to understand preferences
	endDate := date.AddDate(0, 0, 1)
	startDate := date.AddDate(0, 0, -30) // Last 30 days

	recentMeals, err := s.mealRepo.GetByUserAndDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get all available dishes with empty filter
	filter := repository.DishFilter{}
	dishes, _, err := s.dishRepo.GetAll(ctx, filter, 1, 20) // Get top 20 dishes
	if err != nil {
		return nil, err
	}

	// Analyze user preferences
	preferredCuisines := make(map[string]int)
	preferredDishIDs := make(map[primitive.ObjectID]int)

	for _, meal := range recentMeals {
		if meal.MealType == mealType {
			preferredDishIDs[meal.DishID]++

			// Get dish to analyze cuisine preferences
			dish, err := s.dishRepo.GetByID(ctx, meal.DishID)
			if err == nil {
				preferredCuisines[dish.Cuisine]++
			}
		}
	}

	// Generate recommendations
	var recommendations []models.RecommendedDish
	maxRecommendations := 5

	// First, add dishes from preferred cuisines that user hasn't had recently
	for _, dish := range dishes {
		if len(recommendations) >= maxRecommendations {
			break
		}

		// Skip if user had this dish recently
		if preferredDishIDs[dish.ID] > 0 {
			continue
		}

		// All dishes can be recommended for any meal type
		// Enhanced recommendation logic based on user preferences
		score := 0.6 // Base score

		// Increase score if cuisine is popular with user
		if preferredCuisines[dish.Cuisine] > 0 {
			score += 0.3
		}

		// Boost score for highly rated dishes (if available)
		if dish.Calories > 0 && dish.Calories < 600 { // Light meals
			score += 0.1
		}

		// Ensure score doesn't exceed 1.0
		if score > 1.0 {
			score = 1.0
		}

		reason := fmt.Sprintf("Based on your %s preferences", dish.Cuisine)
		if preferredCuisines[dish.Cuisine] > 0 {
			reason = fmt.Sprintf("You enjoyed %s cuisine recently", dish.Cuisine)
		} else {
			reason = "Trying something new based on your dietary patterns"
		}

		recommendation := models.RecommendedDish{
			DishID:     dish.ID.Hex(),
			DishName:   dish.Name,
			Cuisine:    dish.Cuisine,
			Calories:   dish.Calories, // Use the main calories field from dish
			Score:      score,
			Reason:     reason,
			Image:      dish.Image,
			PrepTime:   dish.PrepTime,
			Difficulty: dish.Difficulty,
		}

		recommendations = append(recommendations, recommendation)
	}

	// If no recommendations found based on preferences, add some popular dishes
	if len(recommendations) == 0 {
		for i, dish := range dishes {
			if i >= maxRecommendations {
				break
			}

			recommendation := models.RecommendedDish{
				DishID:     dish.ID.Hex(),
				DishName:   dish.Name,
				Cuisine:    dish.Cuisine,
				Calories:   dish.Calories,
				Score:      0.5, // Lower score for fallback recommendations
				Reason:     "Popular dish to try",
				Image:      dish.Image,
				PrepTime:   dish.PrepTime,
				Difficulty: dish.Difficulty,
			}
			recommendations = append(recommendations, recommendation)
		}
	}

	reason := fmt.Sprintf("Recommendations for %s", mealType)
	if len(recentMeals) > 0 {
		reason = fmt.Sprintf("Recommendations for %s based on your recent meal history", mealType)
	} else {
		reason = fmt.Sprintf("Popular %s recommendations to get you started", mealType)
	}

	return &models.RecommendationsResponse{
		Recommendations: recommendations,
		Reason:          reason,
	}, nil
}

// GetNutritionProgress gets nutrition progress for a user over a period
func (s *mealService) GetNutritionProgress(ctx context.Context, userID primitive.ObjectID, period int) (*models.NutritionProgressResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -period)

	meals, err := s.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Group meals by date and calculate daily nutrition
	dailyNutrition := make(map[string]*models.DailyNutrition)
	for _, mealWithDish := range meals {
		dateStr := mealWithDish.Date.Format("2006-01-02")

		if _, exists := dailyNutrition[dateStr]; !exists {
			dailyNutrition[dateStr] = &models.DailyNutrition{
				Date:      mealWithDish.Date,
				Calories:  0,
				Protein:   0,
				Carbs:     0,
				Fat:       0,
				Fiber:     0,
				Sodium:    0,
				MealCount: 0,
			}
		}

		// Add nutrition values
		dailyNutrition[dateStr].Calories += mealWithDish.Dish.Calories
		dailyNutrition[dateStr].Protein += mealWithDish.Dish.Nutrition.Protein
		dailyNutrition[dateStr].Carbs += mealWithDish.Dish.Nutrition.Carbs
		dailyNutrition[dateStr].Fat += mealWithDish.Dish.Nutrition.Fat
		dailyNutrition[dateStr].Fiber += mealWithDish.Dish.Nutrition.Fiber
		dailyNutrition[dateStr].Sodium += mealWithDish.Dish.Nutrition.Sodium
		dailyNutrition[dateStr].MealCount++
	}

	// Convert map to slice
	var progressData []models.DailyNutrition
	for _, daily := range dailyNutrition {
		progressData = append(progressData, *daily)
	}

	// Get user's nutrition goals
	goals, err := s.GetNutritionGoals(ctx, userID)
	if err != nil {
		// Use default goals if user hasn't set any
		goals = &models.NutritionGoals{
			DailyCalories: 2000,
			Protein:       150,
			Carbs:         250,
			Fat:           65,
			Fiber:         25,
			Sodium:        2300,
		}
	}

	// Calculate summary
	var totalCalories, totalProtein, totalCarbs, totalFat, totalFiber, totalSodium int
	for _, daily := range progressData {
		totalCalories += daily.Calories
		totalProtein += daily.Protein
		totalCarbs += daily.Carbs
		totalFat += daily.Fat
		totalFiber += daily.Fiber
		totalSodium += daily.Sodium
	}

	days := len(progressData)
	if days == 0 {
		days = 1 // Avoid division by zero
	}

	summary := models.NutritionSummary{
		AvgCalories:    float64(totalCalories) / float64(days),
		AvgProtein:     float64(totalProtein) / float64(days),
		AvgCarbs:       float64(totalCarbs) / float64(days),
		AvgFat:         float64(totalFat) / float64(days),
		AvgFiber:       float64(totalFiber) / float64(days),
		TotalDays:      days,
		CalorieGoalMet: 0,
		ProteinGoalMet: 0,
		GoalPercentage: 0,
	}

	// Calculate goal achievement
	for _, daily := range progressData {
		if daily.Calories >= goals.DailyCalories {
			summary.CalorieGoalMet++
		}
		if daily.Protein >= goals.Protein {
			summary.ProteinGoalMet++
		}
	}

	if days > 0 {
		summary.GoalPercentage = float64(summary.CalorieGoalMet) / float64(days) * 100
	}

	return &models.NutritionProgressResponse{
		Period:   period,
		Progress: progressData,
		Goals:    *goals,
		Summary:  summary,
	}, nil
}

// GetNutritionGoals gets nutrition goals for a user
func (s *mealService) GetNutritionGoals(ctx context.Context, userID primitive.ObjectID) (*models.NutritionGoals, error) {
	// For now, return default goals. In a real implementation, this would be stored in user profile
	return &models.NutritionGoals{
		DailyCalories: 2000,
		Protein:       150,
		Carbs:         250,
		Fat:           65,
		Fiber:         25,
		Sodium:        2300,
	}, nil
}

// UpdateNutritionGoals updates nutrition goals for a user
func (s *mealService) UpdateNutritionGoals(ctx context.Context, userID primitive.ObjectID, req models.NutritionGoalsRequest) (*models.NutritionGoals, error) {
	// Just return the goals object. The actual persistence will be handled by the nutrition handler
	// using the user service directly since meal service shouldn't depend on user service
	goals := &models.NutritionGoals{
		DailyCalories: req.DailyCalories,
		Protein:       req.DailyProtein,
		Carbs:         req.DailyCarbs,
		Fat:           req.DailyFat,
		Fiber:         req.DailyFiber,
		Sodium:        req.DailySodium,
	}
	return goals, nil
}

// contains checks if any of the keywords are present in the text (case-insensitive)
func contains(text string, keywords ...string) bool {
	text = strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}
