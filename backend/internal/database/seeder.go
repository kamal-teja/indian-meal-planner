package database

import (
	"context"
	"encoding/json"
	"time"

	"meal-planner-backend/internal/models"
	"meal-planner-backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedDefaultDishes seeds the database with default dishes if it's empty
func SeedDefaultDishes(db *mongo.Database, log *logger.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := db.Collection("dishes")

	// Check if dishes already exist
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		log.Info("Database already contains dishes, skipping seeding", "count", count)
		return nil
	}

	// Default dishes data
	defaultDishes := getDefaultDishes()

	// Insert dishes
	var docs []interface{}
	for _, dish := range defaultDishes {
		docs = append(docs, dish)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	log.Info("Successfully seeded database with default dishes", "count", len(result.InsertedIDs))
	return nil
}

// getDefaultDishes returns a slice of default dishes
func getDefaultDishes() []models.Dish {
	dishesJSON := `[
		{
			"name": "Butter Chicken",
			"type": "Non-Veg",
			"cuisine": "North Indian",
			"image": "https://images.unsplash.com/photo-1565557623262-b51c2513a641?w=500&h=300&fit=crop",
			"ingredients": ["chicken", "butter", "tomato", "cream", "garam masala", "ginger", "garlic"],
			"calories": 438,
			"nutrition": {
				"protein": 25,
				"carbs": 8,
				"fat": 35,
				"fiber": 2,
				"sugar": 6,
				"sodium": 875
			},
			"dietaryTags": ["high-protein"],
			"spiceLevel": "medium",
			"prepTime": 45,
			"cookTime": 30,
			"servings": 4,
			"difficulty": "medium",
			"description": "Rich and creamy North Indian curry with tender chicken in a spiced tomato-butter sauce."
		},
		{
			"name": "Dal Tadka",
			"type": "Veg",
			"cuisine": "North Indian",
			"image": "https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=500&h=300&fit=crop",
			"ingredients": ["yellow lentils", "onion", "tomato", "cumin", "turmeric", "ginger", "garlic", "green chilies"],
			"calories": 230,
			"nutrition": {
				"protein": 12,
				"carbs": 35,
				"fat": 6,
				"fiber": 8,
				"sugar": 4,
				"sodium": 450
			},
			"dietaryTags": ["vegetarian", "vegan", "high-protein", "gluten-free"],
			"spiceLevel": "medium",
			"prepTime": 15,
			"cookTime": 25,
			"servings": 4,
			"difficulty": "easy",
			"description": "Comfort food at its best - yellow lentils cooked with aromatic spices and tempered with cumin."
		},
		{
			"name": "Masala Dosa",
			"type": "Veg",
			"cuisine": "South Indian",
			"image": "https://images.unsplash.com/photo-1567188040759-fb8a883dc6d8?w=500&h=300&fit=crop",
			"ingredients": ["rice", "urad dal", "potato", "onion", "mustard seeds", "curry leaves", "turmeric"],
			"calories": 375,
			"nutrition": {
				"protein": 8,
				"carbs": 68,
				"fat": 8,
				"fiber": 4,
				"sugar": 3,
				"sodium": 620
			},
			"dietaryTags": ["vegetarian", "vegan", "gluten-free"],
			"spiceLevel": "mild",
			"prepTime": 480,
			"cookTime": 20,
			"servings": 2,
			"difficulty": "hard",
			"description": "Crispy South Indian crepe made from fermented rice and lentil batter, filled with spiced potatoes."
		},
		{
			"name": "Biryani",
			"type": "Non-Veg",
			"cuisine": "Mughlai",
			"image": "https://images.unsplash.com/photo-1563379091339-03246963d4b5?w=500&h=300&fit=crop",
			"ingredients": ["basmati rice", "chicken", "yogurt", "saffron", "mint", "fried onions", "ghee", "whole spices"],
			"calories": 520,
			"nutrition": {
				"protein": 22,
				"carbs": 65,
				"fat": 18,
				"fiber": 3,
				"sugar": 5,
				"sodium": 780
			},
			"dietaryTags": ["high-protein"],
			"spiceLevel": "medium",
			"prepTime": 60,
			"cookTime": 45,
			"servings": 6,
			"difficulty": "hard",
			"description": "Fragrant and flavorful rice dish layered with marinated meat and aromatic spices."
		},
		{
			"name": "Palak Paneer",
			"type": "Veg",
			"cuisine": "North Indian",
			"image": "https://images.unsplash.com/photo-1601050690597-df0568f70950?w=500&h=300&fit=crop",
			"ingredients": ["spinach", "paneer", "onion", "tomato", "ginger", "garlic", "cream", "garam masala"],
			"calories": 285,
			"nutrition": {
				"protein": 15,
				"carbs": 12,
				"fat": 22,
				"fiber": 4,
				"sugar": 8,
				"sodium": 520
			},
			"dietaryTags": ["vegetarian", "high-protein", "gluten-free"],
			"spiceLevel": "mild",
			"prepTime": 20,
			"cookTime": 25,
			"servings": 4,
			"difficulty": "medium",
			"description": "Creamy spinach curry with chunks of soft paneer cheese in aromatic spices."
		}
	]`

	var dishes []models.Dish
	if err := json.Unmarshal([]byte(dishesJSON), &dishes); err != nil {
		// Return empty slice if JSON parsing fails
		return []models.Dish{}
	}

	return dishes
}
