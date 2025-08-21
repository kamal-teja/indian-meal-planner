package service

import (
	"meal-planner-backend/internal/config"
	"meal-planner-backend/internal/repository"
	"meal-planner-backend/pkg/logger"
)

// Services holds all service instances
type Services struct {
	Auth     AuthService
	User     UserService
	Dish     DishService
	Meal     MealService
	MealPlan MealPlanService
}

// NewServices creates and returns all service instances
func NewServices(repos *repository.Repositories, cfg *config.Config, log *logger.Logger) *Services {
	return &Services{
		Auth:     NewAuthService(repos.User, cfg, log),
		User:     NewUserService(repos.User, log),
		Dish:     NewDishService(repos.Dish, repos.User, log),
		Meal:     NewMealService(repos.Meal, repos.Dish, log),
		MealPlan: NewMealPlanService(repos.MealPlan, repos.Dish, log),
	}
}
