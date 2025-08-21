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

// MealPlanRepository interface defines meal plan database operations
type MealPlanRepository interface {
	Create(ctx context.Context, mealPlan *models.MealPlan) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.MealPlan, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.MealPlan, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, mealPlan *models.MealPlan) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetActivePlans(ctx context.Context, userID primitive.ObjectID) ([]*models.MealPlan, error)
}

// mealPlanRepository implements MealPlanRepository interface
type mealPlanRepository struct {
	collection *mongo.Collection
}

// NewMealPlanRepository creates a new meal plan repository
func NewMealPlanRepository(db *mongo.Database) MealPlanRepository {
	collection := db.Collection("mealplans")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Compound index for user queries
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"userId", 1}, {"startDate", -1}},
	})

	// Index for date range queries
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"startDate", 1}, {"endDate", 1}},
	})

	return &mealPlanRepository{
		collection: collection,
	}
}

// Create creates a new meal plan
func (r *mealPlanRepository) Create(ctx context.Context, mealPlan *models.MealPlan) error {
	mealPlan.CreatedAt = time.Now()
	mealPlan.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, mealPlan)
	if err != nil {
		return err
	}

	mealPlan.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a meal plan by ID
func (r *mealPlanRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.MealPlan, error) {
	var mealPlan models.MealPlan
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&mealPlan)
	if err != nil {
		return nil, err
	}
	return &mealPlan, nil
}

// GetByUserID retrieves meal plans for a user with pagination
func (r *mealPlanRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.MealPlan, int64, error) {
	query := bson.M{"userId": userID}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Calculate pagination
	skip := (page - 1) * limit

	// Find options
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"startDate", -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var mealPlans []*models.MealPlan
	if err = cursor.All(ctx, &mealPlans); err != nil {
		return nil, 0, err
	}

	return mealPlans, total, nil
}

// GetActivePlans retrieves active meal plans for a user
func (r *mealPlanRepository) GetActivePlans(ctx context.Context, userID primitive.ObjectID) ([]*models.MealPlan, error) {
	now := time.Now()
	query := bson.M{
		"userId":    userID,
		"startDate": bson.M{"$lte": now},
		"endDate":   bson.M{"$gte": now},
	}

	opts := options.Find().SetSort(bson.D{{"startDate", 1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mealPlans []*models.MealPlan
	if err = cursor.All(ctx, &mealPlans); err != nil {
		return nil, err
	}

	return mealPlans, nil
}

// Update updates a meal plan
func (r *mealPlanRepository) Update(ctx context.Context, id primitive.ObjectID, mealPlan *models.MealPlan) error {
	mealPlan.UpdatedAt = time.Now()

	update := bson.M{"$set": mealPlan}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a meal plan
func (r *mealPlanRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
