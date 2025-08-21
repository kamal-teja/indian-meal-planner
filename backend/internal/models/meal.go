package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Meal represents a meal entry in the system
type Meal struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date     time.Time          `bson:"date" json:"date" validate:"required"`
	MealType string             `bson:"mealType" json:"mealType" validate:"required,oneof=breakfast lunch dinner snack"`
	DishID   primitive.ObjectID `bson:"dishId" json:"dishId" validate:"required"`
	UserID   primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`

	// Optional fields
	Notes  string `bson:"notes" json:"notes"`
	Rating int    `bson:"rating" json:"rating" validate:"min=0,max=5"` // 0 means no rating

	// Metadata
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// MealWithDish represents a meal with populated dish information
type MealWithDish struct {
	ID        string       `json:"id"`
	Date      time.Time    `json:"date"`
	MealType  string       `json:"mealType"`
	Dish      DishResponse `json:"dish"`
	User      string       `json:"user"`
	Notes     string       `json:"notes"`
	Rating    int          `json:"rating"`
	CreatedAt time.Time    `json:"createdAt"`
}

// MealRequest represents the request for creating/updating a meal
type MealRequest struct {
	Date     FlexibleDate `json:"date" validate:"required"`
	MealType string       `json:"mealType" validate:"required,oneof=breakfast lunch dinner snack"`
	DishID   string       `json:"dishId" validate:"required"`
	Notes    string       `json:"notes"`
	Rating   int          `json:"rating" validate:"min=0,max=5"` // 0 means no rating
}

// MealPlan represents a meal plan for a user
type MealPlan struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`
	Name   string             `bson:"name" json:"name" validate:"required"`

	// Plan details
	StartDate   time.Time `bson:"startDate" json:"startDate" validate:"required"`
	EndDate     time.Time `bson:"endDate" json:"endDate" validate:"required"`
	Description string    `bson:"description" json:"description"`

	// Meals in the plan
	Meals []MealPlanMeal `bson:"meals" json:"meals"`

	// Metadata
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// MealPlanMeal represents a meal within a meal plan
type MealPlanMeal struct {
	Date     time.Time          `bson:"date" json:"date"`
	MealType string             `bson:"mealType" json:"mealType"`
	DishID   primitive.ObjectID `bson:"dishId" json:"dishId"`
	Notes    string             `bson:"notes" json:"notes"`
}

// MealPlanRequest represents the request for creating/updating a meal plan
type MealPlanRequest struct {
	Name        string    `json:"name" validate:"required"`
	StartDate   time.Time `json:"startDate" validate:"required"`
	EndDate     time.Time `json:"endDate" validate:"required"`
	Description string    `json:"description"`
}

// GetValidMealTypes returns the list of valid meal types
func GetValidMealTypes() []string {
	return []string{"breakfast", "lunch", "dinner", "snack"}
}
