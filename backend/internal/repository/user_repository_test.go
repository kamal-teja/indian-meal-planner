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

func TestUserRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		user := &models.User{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "hashedpassword",
			Profile: models.UserProfile{
				DietaryPreferences: []string{"vegetarian"},
				SpiceLevel:         "medium",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Create(testContext(), user)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, user.ID)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{"_id", userID},
			{"name", "John Doe"},
			{"email", "john@example.com"},
			{"profile", bson.D{
				{"dietaryPreferences", bson.A{"vegetarian"}},
				{"spiceLevel", "medium"},
			}},
		}))

		// Act
		user, err := repo.GetByID(testContext(), userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "medium", user.Profile.SpiceLevel)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))

		// Act
		user, err := repo.GetByID(testContext(), userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()
		email := "john@example.com"

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{"_id", userID},
			{"name", "John Doe"},
			{"email", email},
			{"password", "hashedpassword"},
		}))

		// Act
		user, err := repo.GetByEmail(testContext(), email)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, "John Doe", user.Name)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		email := "nonexistent@example.com"

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))

		// Act
		user, err := repo.GetByEmail(testContext(), email)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()
		user := &models.User{
			Name: "John Updated",
			Profile: models.UserProfile{
				DietaryPreferences: []string{"vegan"},
				SpiceLevel:         "hot",
			},
			UpdatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.Update(testContext(), userID, user)

		// Assert
		assert.NoError(t, err)
	})
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.UpdateLastLogin(testContext(), userID)

		// Assert
		assert.NoError(t, err)
	})
}

func TestUserRepository_AddToFavorites(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()
		dishID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.AddToFavorites(testContext(), userID, dishID)

		// Assert
		assert.NoError(t, err)
	})
}

func TestUserRepository_RemoveFromFavorites(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()
		dishID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := repo.RemoveFromFavorites(testContext(), userID, dishID)

		// Assert
		assert.NoError(t, err)
	})
}

func TestUserRepository_GetFavorites(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()
		dishID1 := primitive.NewObjectID()
		dishID2 := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{"_id", userID},
			{"favorites", bson.A{dishID1, dishID2}},
		}))

		// Act
		favorites, err := repo.GetFavorites(testContext(), userID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, favorites, 2)
		assert.Contains(t, favorites, dishID1)
		assert.Contains(t, favorites, dishID2)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))

		// Act
		favorites, err := repo.GetFavorites(testContext(), userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, favorites)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 1}})

		// Act
		err := repo.Delete(testContext(), userID)

		// Assert
		assert.NoError(t, err)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		// Arrange
		repo := NewUserRepository(mt.DB)
		userID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 0}})

		// Act
		err := repo.Delete(testContext(), userID)

		// Assert
		assert.Error(t, err)
	})
}

func TestNewRepositories(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("creates all repositories", func(mt *mtest.T) {
		// Act
		repos := NewRepositories(mt.DB)

		// Assert
		assert.NotNil(t, repos)
		assert.NotNil(t, repos.User)
		assert.NotNil(t, repos.Dish)
		assert.NotNil(t, repos.Meal)
		assert.NotNil(t, repos.MealPlan)
	})
}