package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Dish represents a dish/recipe in the system
type Dish struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Type        string             `bson:"type" json:"type" validate:"required,oneof=Veg Non-Veg"`
	Cuisine     string             `bson:"cuisine" json:"cuisine" validate:"required"`
	Image       string             `bson:"image" json:"image"`
	Ingredients []string           `bson:"ingredients" json:"ingredients" validate:"required,min=1"`
	Calories    int                `bson:"calories" json:"calories" validate:"min=0"`

	// Nutritional information
	Nutrition Nutrition `bson:"nutrition" json:"nutrition"`

	// Tags and metadata
	DietaryTags []string `bson:"dietaryTags" json:"dietaryTags"`
	SpiceLevel  string   `bson:"spiceLevel" json:"spiceLevel" validate:"oneof=mild medium hot extra-hot"`

	// Cooking information
	PrepTime   int    `bson:"prepTime" json:"prepTime" validate:"min=0"` // minutes
	CookTime   int    `bson:"cookTime" json:"cookTime" validate:"min=0"` // minutes
	Servings   int    `bson:"servings" json:"servings" validate:"min=1"`
	Difficulty string `bson:"difficulty" json:"difficulty" validate:"oneof=easy medium hard"`

	// Additional information
	Description string `bson:"description" json:"description"`

	// Metadata
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// Nutrition represents nutritional information for a dish
type Nutrition struct {
	Protein int `bson:"protein" json:"protein" validate:"min=0"` // grams
	Carbs   int `bson:"carbs" json:"carbs" validate:"min=0"`     // grams
	Fat     int `bson:"fat" json:"fat" validate:"min=0"`         // grams
	Fiber   int `bson:"fiber" json:"fiber" validate:"min=0"`     // grams
	Sugar   int `bson:"sugar" json:"sugar" validate:"min=0"`     // grams
	Sodium  int `bson:"sodium" json:"sodium" validate:"min=0"`   // milligrams
}

// DishResponse represents the dish data returned in API responses with favorites info
type DishResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Cuisine     string    `json:"cuisine"`
	Image       string    `json:"image"`
	Ingredients []string  `json:"ingredients"`
	Calories    int       `json:"calories"`
	Nutrition   Nutrition `json:"nutrition"`
	DietaryTags []string  `json:"dietaryTags"`
	SpiceLevel  string    `json:"spiceLevel"`
	PrepTime    int       `json:"prepTime"`
	CookTime    int       `json:"cookTime"`
	Servings    int       `json:"servings"`
	Difficulty  string    `json:"difficulty"`
	Description string    `json:"description"`
	IsFavorite  bool      `json:"isFavorite,omitempty"`
}

// DishCreateRequest represents the request for creating a dish
type DishCreateRequest struct {
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Type        string    `json:"type" validate:"required,oneof=Veg Non-Veg"`
	Cuisine     string    `json:"cuisine" validate:"required"`
	Image       string    `json:"image"`
	Ingredients []string  `json:"ingredients" validate:"required,min=1"`
	Calories    int       `json:"calories" validate:"omitempty,min=0"`
	Nutrition   Nutrition `json:"nutrition"`
	DietaryTags []string  `json:"dietaryTags"`
	SpiceLevel  string    `json:"spiceLevel" validate:"omitempty,oneof=mild medium hot extra-hot"`
	PrepTime    int       `json:"prepTime" validate:"omitempty,min=0"`
	CookTime    int       `json:"cookTime" validate:"omitempty,min=0"`
	Servings    int       `json:"servings" validate:"omitempty,min=1"`
	Difficulty  string    `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Description string    `json:"description"`
}

// ToResponse converts Dish model to DishResponse
func (d *Dish) ToResponse() DishResponse {
	return DishResponse{
		ID:          d.ID.Hex(),
		Name:        d.Name,
		Type:        d.Type,
		Cuisine:     d.Cuisine,
		Image:       d.Image,
		Ingredients: d.Ingredients,
		Calories:    d.Calories,
		Nutrition:   d.Nutrition,
		DietaryTags: d.DietaryTags,
		SpiceLevel:  d.SpiceLevel,
		PrepTime:    d.PrepTime,
		CookTime:    d.CookTime,
		Servings:    d.Servings,
		Difficulty:  d.Difficulty,
		Description: d.Description,
		IsFavorite:  false, // Will be set by service layer
	}
}

// GetValidDietaryTags returns the list of valid dietary tags
func GetValidDietaryTags() []string {
	return []string{
		"vegetarian", "vegan", "gluten-free", "dairy-free", "nut-free",
		"keto", "paleo", "low-carb", "high-protein", "low-sodium", "sugar-free",
	}
}

// GetValidCuisines returns the list of valid cuisines
func GetValidCuisines() []string {
	return []string{
		"North Indian", "South Indian", "Bengali", "Gujarati", "Punjabi",
		"Rajasthani", "Maharashtrian", "Italian", "Chinese", "Thai",
		"French", "Korean", "Continental", "Mughlai", "Kashmiri",
	}
}
