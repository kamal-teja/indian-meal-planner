package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Repositories holds all repository instances
type Repositories struct {
	User     UserRepository
	Dish     DishRepository
	Meal     MealRepository
	MealPlan MealPlanRepository
	Undo     UndoRepository
}

// NewRepositories creates and returns all repository instances
func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		User:     NewUserRepository(db),
		Dish:     NewDishRepository(db),
		Meal:     NewMealRepository(db),
		MealPlan: NewMealPlanRepository(db),
		Undo:     NewUndoRepository(db),
	}
}
