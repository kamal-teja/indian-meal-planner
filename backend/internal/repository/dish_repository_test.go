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

func TestDishRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dish := &models.Dish{
			Name:        "Chicken Biryani",
			Type:        "Non-Veg",
			Cuisine:     "North Indian",
			Ingredients: []string{"chicken", "rice", "spices"},
			Calories:    450,
			Nutrition: models.Nutrition{
				Protein: 25,
				Carbs:   50,
				Fat:     15,
			},
			SpiceLevel: "medium",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Create(testContext(), dish)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, dish.ID)
	})
}

func TestDishRepository_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dishID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.FirstBatch, bson.D{
			{"_id", dishID},
			{"name", "Chicken Biryani"},
			{"type", "Non-Veg"},
			{"cuisine", "North Indian"},
			{"ingredients", bson.A{"chicken", "rice", "spices"}},
			{"calories", 450},
			{"nutrition", bson.D{
				{"protein", 25},
				{"carbs", 50},
				{"fat", 15},
			}},
			{"spiceLevel", "medium"},
		}))

		// Act
		dish, err := repo.GetByID(testContext(), dishID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, dish)
		assert.Equal(t, dishID, dish.ID)
		assert.Equal(t, "Chicken Biryani", dish.Name)
		assert.Equal(t, "Non-Veg", dish.Type)
		assert.Equal(t, "North Indian", dish.Cuisine)
		assert.Equal(t, 450, dish.Calories)
		assert.Equal(t, "medium", dish.SpiceLevel)
	})

	mt.Run("dish not found", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dishID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.dishes", mtest.FirstBatch))

		// Act
		dish, err := repo.GetByID(testContext(), dishID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, dish)
	})
}

func TestDishRepository_GetAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success with filter", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		filter := DishFilter{
			Type:    "Veg",
			Cuisine: "North Indian",
		}

		dishID1 := primitive.NewObjectID()
		dishID2 := primitive.NewObjectID()

		// Mock count response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.FirstBatch,
			bson.D{{"n", 2}}))

		// Mock find response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.FirstBatch,
			bson.D{
				{"_id", dishID1},
				{"name", "Paneer Makhani"},
				{"type", "Veg"},
				{"cuisine", "North Indian"},
				{"calories", 350},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.NextBatch,
			bson.D{
				{"_id", dishID2},
				{"name", "Dal Makhani"},
				{"type", "Veg"},
				{"cuisine", "North Indian"},
				{"calories", 250},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.dishes", mtest.NextBatch))

		// Act
		dishes, total, err := repo.GetAll(testContext(), filter, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, dishes, 2)
		if len(dishes) >= 2 {
			assert.Equal(t, "Paneer Makhani", dishes[0].Name)
			assert.Equal(t, "Dal Makhani", dishes[1].Name)
		}
	})

	mt.Run("success empty filter", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		filter := DishFilter{}

		// Mock count response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.FirstBatch,
			bson.D{{"n", 0}}))

		// Mock find response
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.dishes", mtest.FirstBatch))

		// Act
		dishes, total, err := repo.GetAll(testContext(), filter, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, dishes, 0)
	})
}

func TestDishRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dishID := primitive.NewObjectID()
		dish := &models.Dish{
			Name:      "Updated Dish",
			Calories:  400,
			UpdatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Update(testContext(), dishID, dish)

		// Assert
		assert.NoError(t, err)
	})
}

func TestDishRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dishID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 1}})

		// Act
		err := repo.Delete(testContext(), dishID)

		// Assert
		assert.NoError(t, err)
	})

	mt.Run("dish not found", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dishID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 0}})

		// Act
		err := repo.Delete(testContext(), dishID)

		// Assert
		assert.Error(t, err)
	})
}

func TestDishRepository_GetByIDs(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dishID1 := primitive.NewObjectID()
		dishID2 := primitive.NewObjectID()
		dishIDs := []primitive.ObjectID{dishID1, dishID2}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.FirstBatch,
			bson.D{
				{"_id", dishID1},
				{"name", "Dish 1"},
				{"type", "Veg"},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.NextBatch,
			bson.D{
				{"_id", dishID2},
				{"name", "Dish 2"},
				{"type", "Non-Veg"},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.dishes", mtest.NextBatch))

		// Act
		dishes, err := repo.GetByIDs(testContext(), dishIDs)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, dishes, 2)
		assert.Equal(t, dishID1, dishes[0].ID)
		assert.Equal(t, dishID2, dishes[1].ID)
	})

	mt.Run("empty ids", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		dishIDs := []primitive.ObjectID{}

		// Act
		dishes, err := repo.GetByIDs(testContext(), dishIDs)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, dishes, 0)
	})
}

func TestDishRepository_Search(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		query := "chicken"
		filter := DishFilter{Type: "Non-Veg"}

		dishID := primitive.NewObjectID()

		// Mock aggregate response for count
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.FirstBatch,
			bson.D{{"count", 1}}))

		// Mock aggregate response for search results
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.dishes", mtest.FirstBatch,
			bson.D{
				{"_id", dishID},
				{"name", "Chicken Curry"},
				{"type", "Non-Veg"},
				{"cuisine", "Indian"},
				{"score", 1.5},
			}))
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.dishes", mtest.NextBatch))

		// Act
		dishes, total, err := repo.Search(testContext(), query, filter, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, dishes, 1)
		assert.Equal(t, "Chicken Curry", dishes[0].Name)
		assert.Equal(t, "Non-Veg", dishes[0].Type)
	})

	mt.Run("empty query", func(mt *mtest.T) {
		// Arrange
		repo := NewDishRepository(mt.DB)
		query := ""
		filter := DishFilter{}

		// Act
		dishes, total, err := repo.Search(testContext(), query, filter, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, dishes, 0)
	})
}

func TestDishFilter(t *testing.T) {
	// Test that DishFilter struct works correctly
	filter := DishFilter{
		Type:        "Veg",
		Cuisine:     "North Indian",
		DietaryTags: []string{"vegetarian", "high-protein"},
		SpiceLevel:  "medium",
		MaxCalories: 500,
		MinCalories: 100,
		Ingredients: []string{"paneer", "tomatoes"},
	}

	assert.Equal(t, "Veg", filter.Type)
	assert.Equal(t, "North Indian", filter.Cuisine)
	assert.Contains(t, filter.DietaryTags, "vegetarian")
	assert.Contains(t, filter.DietaryTags, "high-protein")
	assert.Equal(t, "medium", filter.SpiceLevel)
	assert.Equal(t, 500, filter.MaxCalories)
	assert.Equal(t, 100, filter.MinCalories)
	assert.Contains(t, filter.Ingredients, "paneer")
	assert.Contains(t, filter.Ingredients, "tomatoes")
}