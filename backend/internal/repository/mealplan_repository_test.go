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

func TestMealPlanRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlan := &models.MealPlan{
			UserID:      primitive.NewObjectID(),
			Name:        "Test Plan",
			Description: "Test meal plan",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 7),
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Create(testContext(), mealPlan)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, mealPlan.ID)
	})

	mt.Run("database error", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlan := &models.MealPlan{
			UserID:    primitive.NewObjectID(),
			Name:      "Test Plan",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		// Act
		err := repo.Create(testContext(), mealPlan)

		// Assert
		assert.Error(t, err)
	})
}

func TestMealPlanRepository_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlanID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.mealplans", mtest.FirstBatch,
			bson.D{
				{"_id", mealPlanID},
				{"name", "Test Plan"},
				{"description", "Test meal plan"},
				{"isActive", true},
			}))

		// Act
		mealPlan, err := repo.GetByID(testContext(), mealPlanID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, mealPlan)
		assert.Equal(t, mealPlanID, mealPlan.ID)
		assert.Equal(t, "Test Plan", mealPlan.Name)
		assert.True(t, mealPlan.IsActive)
	})

	mt.Run("not found", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlanID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.mealplans", mtest.FirstBatch))

		// Act
		mealPlan, err := repo.GetByID(testContext(), mealPlanID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mealPlan)
	})
}

func TestMealPlanRepository_GetByUserID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		userID := primitive.NewObjectID()
		mealPlanID1 := primitive.NewObjectID()
		mealPlanID2 := primitive.NewObjectID()

		// Mock count response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.mealplans", mtest.FirstBatch,
			bson.D{{"n", 2}}))

		// Mock find response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.mealplans", mtest.FirstBatch,
			bson.D{
				{"_id", mealPlanID1},
				{"userId", userID},
				{"name", "Plan 1"},
				{"isActive", true},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.mealplans", mtest.NextBatch,
			bson.D{
				{"_id", mealPlanID2},
				{"userId", userID},
				{"name", "Plan 2"},
				{"isActive", false},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.mealplans", mtest.NextBatch))

		// Act
		mealPlans, total, err := repo.GetByUserID(testContext(), userID, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, mealPlans, 2)
		if len(mealPlans) >= 2 {
			assert.Equal(t, "Plan 1", mealPlans[0].Name)
			assert.Equal(t, "Plan 2", mealPlans[1].Name)
		}
	})

	mt.Run("empty result", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		userID := primitive.NewObjectID()

		// Mock count response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.mealplans", mtest.FirstBatch,
			bson.D{{"n", 0}}))

		// Mock find response
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.mealplans", mtest.FirstBatch))

		// Act
		mealPlans, total, err := repo.GetByUserID(testContext(), userID, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, mealPlans, 0)
	})
}

func TestMealPlanRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlanID := primitive.NewObjectID()
		mealPlan := &models.MealPlan{
			ID:          mealPlanID,
			Name:        "Updated Plan",
			Description: "Updated description",
			IsActive:    false,
			UpdatedAt:   time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Update(testContext(), mealPlanID, mealPlan)

		// Assert
		assert.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlanID := primitive.NewObjectID()
		mealPlan := &models.MealPlan{
			ID:   mealPlanID,
			Name: "Updated Plan",
		}

		// Mock no documents matched
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 0},
			{"nModified", 0},
		})

		// Act
		err := repo.Update(testContext(), mealPlanID, mealPlan)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no documents updated")
	})
}

func TestMealPlanRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlanID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 1},
		})

		// Act
		err := repo.Delete(testContext(), mealPlanID)

		// Assert
		assert.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		mealPlanID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 0},
		})

		// Act
		err := repo.Delete(testContext(), mealPlanID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no documents deleted")
	})
}

func TestMealPlanRepository_GetActivePlans(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		userID := primitive.NewObjectID()
		mealPlanID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.mealplans", mtest.FirstBatch,
			bson.D{
				{"_id", mealPlanID},
				{"userId", userID},
				{"name", "Active Plan"},
				{"isActive", true},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.mealplans", mtest.NextBatch))

		// Act
		mealPlans, err := repo.GetActivePlans(testContext(), userID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, mealPlans, 1)
		if len(mealPlans) > 0 {
			assert.Equal(t, "Active Plan", mealPlans[0].Name)
			assert.True(t, mealPlans[0].IsActive)
		}
	})

	mt.Run("no active plans", func(mt *mtest.T) {
		// Arrange
		repo := NewMealPlanRepository(mt.DB)
		userID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.mealplans", mtest.FirstBatch))

		// Act
		mealPlans, err := repo.GetActivePlans(testContext(), userID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, mealPlans, 0)
	})
}