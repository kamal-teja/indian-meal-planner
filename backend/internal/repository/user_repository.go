package repository

import (
	"context"
	"time"

	"nourish-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository interface defines user database operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id primitive.ObjectID, user *models.User) error
	UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error
	AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error
	RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error
	GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) UserRepository {
	collection := db.Collection("users")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Email index (unique)
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	})
	
	return &userRepository{
		collection: collection,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, id primitive.ObjectID, user *models.User) error {
	user.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"email":     user.Email,
			"profile":   user.Profile,
			"updatedAt": user.UpdatedAt,
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateLastLogin updates the last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// AddToFavorites adds a dish to user's favorites
func (r *userRepository) AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	update := bson.M{
		"$addToSet": bson.M{"favorites": dishID},
		"$set":      bson.M{"updatedAt": time.Now()},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	return err
}

// RemoveFromFavorites removes a dish from user's favorites
func (r *userRepository) RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	update := bson.M{
		"$pull": bson.M{"favorites": dishID},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	return err
}

// GetFavorites gets user's favorite dish IDs
func (r *userRepository) GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	var user struct {
		Favorites []primitive.ObjectID `bson:"favorites"`
	}
	
	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": userID},
		options.FindOne().SetProjection(bson.M{"favorites": 1}),
	).Decode(&user)
	
	if err != nil {
		return nil, err
	}
	
	return user.Favorites, nil
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
