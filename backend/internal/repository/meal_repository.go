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

// MealRepository interface defines meal database operations
type MealRepository interface {
	Create(ctx context.Context, meal *models.Meal) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Meal, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.Meal, int64, error)
	GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*models.Meal, error)
	Update(ctx context.Context, id primitive.ObjectID, meal *models.Meal) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetNutritionByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]NutritionSummary, error)
}

// NutritionSummary represents daily nutrition summary
type NutritionSummary struct {
	Date      time.Time `bson:"_id" json:"date"`
	Calories  int       `bson:"totalCalories" json:"calories"`
	Protein   int       `bson:"totalProtein" json:"protein"`
	Carbs     int       `bson:"totalCarbs" json:"carbs"`
	Fat       int       `bson:"totalFat" json:"fat"`
	Fiber     int       `bson:"totalFiber" json:"fiber"`
	Sodium    int       `bson:"totalSodium" json:"sodium"`
	MealCount int       `bson:"mealCount" json:"mealCount"`
}

// mealRepository implements MealRepository interface
type mealRepository struct {
	collection *mongo.Collection
}

// NewMealRepository creates a new meal repository
func NewMealRepository(db *mongo.Database) MealRepository {
	collection := db.Collection("meals")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Compound index for user queries
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"userId", 1}, {"date", -1}},
	})
	
	// Index for date range queries
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"date", 1}},
	})
	
	return &mealRepository{
		collection: collection,
	}
}

// Create creates a new meal
func (r *mealRepository) Create(ctx context.Context, meal *models.Meal) error {
	meal.CreatedAt = time.Now()
	meal.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, meal)
	if err != nil {
		return err
	}
	
	meal.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a meal by ID
func (r *mealRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Meal, error) {
	var meal models.Meal
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&meal)
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

// GetByUserID retrieves meals for a user with pagination
func (r *mealRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.Meal, int64, error) {
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
		SetSort(bson.D{{"date", -1}})
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var meals []*models.Meal
	if err = cursor.All(ctx, &meals); err != nil {
		return nil, 0, err
	}
	
	return meals, total, nil
}

// GetByUserAndDateRange retrieves meals for a user within a date range
func (r *mealRepository) GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*models.Meal, error) {
	query := bson.M{
		"userId": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	
	opts := options.Find().SetSort(bson.D{{"date", 1}, {"mealType", 1}})
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var meals []*models.Meal
	if err = cursor.All(ctx, &meals); err != nil {
		return nil, err
	}
	
	return meals, nil
}

// Update updates a meal
func (r *mealRepository) Update(ctx context.Context, id primitive.ObjectID, meal *models.Meal) error {
	meal.UpdatedAt = time.Now()
	
	update := bson.M{"$set": meal}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a meal
func (r *mealRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetNutritionByDateRange aggregates nutrition data for a user within a date range
func (r *mealRepository) GetNutritionByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]NutritionSummary, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"userId": userID,
				"date": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "dishes",
				"localField":   "dishId",
				"foreignField": "_id",
				"as":           "dish",
			},
		},
		{
			"$unwind": "$dish",
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$date",
					},
				},
				"totalCalories": bson.M{"$sum": "$dish.calories"},
				"totalProtein":  bson.M{"$sum": "$dish.nutrition.protein"},
				"totalCarbs":    bson.M{"$sum": "$dish.nutrition.carbs"},
				"totalFat":      bson.M{"$sum": "$dish.nutrition.fat"},
				"totalFiber":    bson.M{"$sum": "$dish.nutrition.fiber"},
				"totalSodium":   bson.M{"$sum": "$dish.nutrition.sodium"},
				"mealCount":     bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []NutritionSummary
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}
