package repository

import (
	"testing"
	"time"

	"nourish-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMealRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		meal := &models.Meal{
			Date:      time.Now(),
			MealType:  "breakfast",
			DishID:    primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			Notes:     "Delicious breakfast",
			Rating:    5,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Create(testContext(), meal)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, meal.ID)
	})
}

func TestMealRepository_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		mealID := primitive.NewObjectID()
		dishID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.meals", mtest.FirstBatch, bson.D{
			{"_id", mealID},
			{"date", now},
			{"mealType", "breakfast"},
			{"dishId", dishID},
			{"userId", userID},
			{"notes", "Great breakfast"},
			{"rating", 5},
			{"createdAt", now},
			{"updatedAt", now},
		}))

		// Act
		meal, err := repo.GetByID(testContext(), mealID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, mealID, meal.ID)
		assert.Equal(t, "breakfast", meal.MealType)
		assert.Equal(t, dishID, meal.DishID)
		assert.Equal(t, userID, meal.UserID)
		assert.Equal(t, "Great breakfast", meal.Notes)
		assert.Equal(t, 5, meal.Rating)
	})

	mt.Run("meal not found", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		mealID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.meals", mtest.FirstBatch))

		// Act
		meal, err := repo.GetByID(testContext(), mealID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, meal)
	})
}

func TestMealRepository_GetByUserID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		userID := primitive.NewObjectID()
		mealID1 := primitive.NewObjectID()
		mealID2 := primitive.NewObjectID()

		// Mock count response
		mt.AddMockResponses(
			bson.D{{"ok", 1}, {"n", 2}},
		)

		// Mock find response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.meals", mtest.FirstBatch,
			bson.D{
				{"_id", mealID1},
				{"userId", userID},
				{"mealType", "breakfast"},
				{"date", time.Now()},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.meals", mtest.NextBatch,
			bson.D{
				{"_id", mealID2},
				{"userId", userID},
				{"mealType", "lunch"},
				{"date", time.Now()},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.meals", mtest.NextBatch))

		// Act
		meals, total, err := repo.GetByUserID(testContext(), userID, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, meals, 2)
		assert.Equal(t, "breakfast", meals[0].MealType)
		assert.Equal(t, "lunch", meals[1].MealType)
	})
}

func TestMealRepository_GetByUserAndDateRange(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		userID := primitive.NewObjectID()
		startDate := time.Now().AddDate(0, 0, -7) // 7 days ago
		endDate := time.Now()

		mealID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.meals", mtest.FirstBatch,
			bson.D{
				{"_id", mealID},
				{"userId", userID},
				{"mealType", "dinner"},
				{"date", time.Now().AddDate(0, 0, -3)}, // 3 days ago
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.meals", mtest.NextBatch))

		// Act
		meals, err := repo.GetByUserAndDateRange(testContext(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, meals, 1)
		assert.Equal(t, "dinner", meals[0].MealType)
		assert.Equal(t, userID, meals[0].UserID)
	})

	mt.Run("no meals in range", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		userID := primitive.NewObjectID()
		startDate := time.Now().AddDate(0, 0, -7)
		endDate := time.Now().AddDate(0, 0, -1)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.meals", mtest.FirstBatch))

		// Act
		meals, err := repo.GetByUserAndDateRange(testContext(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, meals, 0)
	})
}

func TestMealRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		mealID := primitive.NewObjectID()
		meal := &models.Meal{
			Notes:     "Updated notes",
			Rating:    4,
			UpdatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Update(testContext(), mealID, meal)

		// Assert
		assert.NoError(t, err)
	})
}

func TestMealRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		mealID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 1}})

		// Act
		err := repo.Delete(testContext(), mealID)

		// Assert
		assert.NoError(t, err)
	})

	mt.Run("meal not found", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		mealID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 0}})

		// Act
		err := repo.Delete(testContext(), mealID)

		// Assert
		assert.Error(t, err)
	})
}

func TestMealRepository_GetNutritionByDateRange(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		userID := primitive.NewObjectID()
		startDate := time.Now().AddDate(0, 0, -7)
		endDate := time.Now()
		
		date1 := time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)
		date2 := time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.meals", mtest.FirstBatch,
			bson.D{
				{"_id", date1},
				{"totalCalories", 1200},
				{"totalProtein", 60},
				{"totalCarbs", 150},
				{"totalFat", 40},
				{"totalFiber", 25},
				{"totalSodium", 1800},
				{"mealCount", 3},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.meals", mtest.NextBatch,
			bson.D{
				{"_id", date2},
				{"totalCalories", 1400},
				{"totalProtein", 70},
				{"totalCarbs", 180},
				{"totalFat", 45},
				{"totalFiber", 30},
				{"totalSodium", 2000},
				{"mealCount", 4},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.meals", mtest.NextBatch))

		// Act
		summaries, err := repo.GetNutritionByDateRange(testContext(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, summaries, 2)
		
		// Verify first summary
		assert.Equal(t, date1, summaries[0].Date)
		assert.Equal(t, 1200, summaries[0].Calories)
		assert.Equal(t, 60, summaries[0].Protein)
		assert.Equal(t, 150, summaries[0].Carbs)
		assert.Equal(t, 40, summaries[0].Fat)
		assert.Equal(t, 25, summaries[0].Fiber)
		assert.Equal(t, 1800, summaries[0].Sodium)
		assert.Equal(t, 3, summaries[0].MealCount)
		
		// Verify second summary
		assert.Equal(t, date2, summaries[1].Date)
		assert.Equal(t, 1400, summaries[1].Calories)
		assert.Equal(t, 70, summaries[1].Protein)
		assert.Equal(t, 4, summaries[1].MealCount)
	})

	mt.Run("no nutrition data", func(mt *mtest.T) {
		// Arrange
		repo := NewMealRepository(mt.DB)
		userID := primitive.NewObjectID()
		startDate := time.Now().AddDate(0, 0, -7)
		endDate := time.Now()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.meals", mtest.FirstBatch))

		// Act
		summaries, err := repo.GetNutritionByDateRange(testContext(), userID, startDate, endDate)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, summaries, 0)
	})
}

func TestNutritionSummary(t *testing.T) {
	// Test NutritionSummary struct
	date := time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)
	summary := NutritionSummary{
		Date:      date,
		Calories:  1500,
		Protein:   75,
		Carbs:     200,
		Fat:       50,
		Fiber:     30,
		Sodium:    2200,
		MealCount: 4,
	}

	assert.Equal(t, date, summary.Date)
	assert.Equal(t, 1500, summary.Calories)
	assert.Equal(t, 75, summary.Protein)
	assert.Equal(t, 200, summary.Carbs)
	assert.Equal(t, 50, summary.Fat)
	assert.Equal(t, 30, summary.Fiber)
	assert.Equal(t, 2200, summary.Sodium)
	assert.Equal(t, 4, summary.MealCount)
}