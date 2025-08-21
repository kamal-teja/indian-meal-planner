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

// DishRepository interface defines dish database operations
type DishRepository interface {
	Create(ctx context.Context, dish *models.Dish) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Dish, error)
	GetAll(ctx context.Context, filter DishFilter, page, limit int) ([]*models.Dish, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, dish *models.Dish) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Dish, error)
	Search(ctx context.Context, query string, filter DishFilter, page, limit int) ([]*models.Dish, int64, error)
}

// DishFilter represents filters for dish queries
type DishFilter struct {
	Type         string   // "Veg" or "Non-Veg"
	Cuisine      string   // cuisine type
	DietaryTags  []string // dietary tags
	SpiceLevel   string   // spice level
	MaxCalories  int      // maximum calories
	MinCalories  int      // minimum calories
	Ingredients  []string // must contain these ingredients
}

// dishRepository implements DishRepository interface
type dishRepository struct {
	collection *mongo.Collection
}

// NewDishRepository creates a new dish repository
func NewDishRepository(db *mongo.Database) DishRepository {
	collection := db.Collection("dishes")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Text index for search
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
			{"cuisine", "text"},
			{"ingredients", "text"},
			{"description", "text"},
		},
	})
	
	// Compound indexes for filtering
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"type", 1}, {"cuisine", 1}},
	})
	
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"dietaryTags", 1}},
	})
	
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"calories", 1}},
	})
	
	return &dishRepository{
		collection: collection,
	}
}

// Create creates a new dish
func (r *dishRepository) Create(ctx context.Context, dish *models.Dish) error {
	dish.CreatedAt = time.Now()
	dish.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, dish)
	if err != nil {
		return err
	}
	
	dish.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a dish by ID
func (r *dishRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Dish, error) {
	var dish models.Dish
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&dish)
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

// GetAll retrieves dishes with pagination and filtering
func (r *dishRepository) GetAll(ctx context.Context, filter DishFilter, page, limit int) ([]*models.Dish, int64, error) {
	// Build filter query
	query := r.buildFilterQuery(filter)
	
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
		SetSort(bson.D{{"name", 1}})
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var dishes []*models.Dish
	if err = cursor.All(ctx, &dishes); err != nil {
		return nil, 0, err
	}
	
	return dishes, total, nil
}

// Search searches dishes with text search and filtering
func (r *dishRepository) Search(ctx context.Context, query string, filter DishFilter, page, limit int) ([]*models.Dish, int64, error) {
	// Build search query
	searchQuery := r.buildSearchQuery(query, filter)
	
	// Count total documents
	total, err := r.collection.CountDocuments(ctx, searchQuery)
	if err != nil {
		return nil, 0, err
	}
	
	// Calculate pagination
	skip := (page - 1) * limit
	
	// Find options with text score sorting
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}}).
		SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
	
	cursor, err := r.collection.Find(ctx, searchQuery, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var dishes []*models.Dish
	if err = cursor.All(ctx, &dishes); err != nil {
		return nil, 0, err
	}
	
	return dishes, total, nil
}

// GetByIDs retrieves multiple dishes by their IDs
func (r *dishRepository) GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Dish, error) {
	query := bson.M{"_id": bson.M{"$in": ids}}
	
	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var dishes []*models.Dish
	if err = cursor.All(ctx, &dishes); err != nil {
		return nil, err
	}
	
	return dishes, nil
}

// Update updates a dish
func (r *dishRepository) Update(ctx context.Context, id primitive.ObjectID, dish *models.Dish) error {
	dish.UpdatedAt = time.Now()
	
	update := bson.M{"$set": dish}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a dish
func (r *dishRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// buildFilterQuery builds MongoDB query from DishFilter
func (r *dishRepository) buildFilterQuery(filter DishFilter) bson.M {
	query := bson.M{}
	
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	
	if filter.Cuisine != "" {
		query["cuisine"] = filter.Cuisine
	}
	
	if len(filter.DietaryTags) > 0 {
		query["dietaryTags"] = bson.M{"$in": filter.DietaryTags}
	}
	
	if filter.SpiceLevel != "" {
		query["spiceLevel"] = filter.SpiceLevel
	}
	
	if filter.MaxCalories > 0 || filter.MinCalories > 0 {
		caloriesQuery := bson.M{}
		if filter.MaxCalories > 0 {
			caloriesQuery["$lte"] = filter.MaxCalories
		}
		if filter.MinCalories > 0 {
			caloriesQuery["$gte"] = filter.MinCalories
		}
		query["calories"] = caloriesQuery
	}
	
	if len(filter.Ingredients) > 0 {
		query["ingredients"] = bson.M{"$in": filter.Ingredients}
	}
	
	return query
}

// buildSearchQuery builds MongoDB query for text search with filters
func (r *dishRepository) buildSearchQuery(searchText string, filter DishFilter) bson.M {
	query := bson.M{}
	
	// Add text search
	if searchText != "" {
		query["$text"] = bson.M{"$search": searchText}
	}
	
	// Add filters
	filterQuery := r.buildFilterQuery(filter)
	for k, v := range filterQuery {
		query[k] = v
	}
	
	return query
}
