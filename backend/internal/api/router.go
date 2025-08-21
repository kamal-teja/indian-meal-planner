package api

import (
	"nourish-backend/internal/api/handlers"
	"nourish-backend/internal/api/middleware"
	"nourish-backend/internal/config"
	"nourish-backend/internal/service"
	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the API router
func NewRouter(services *service.Services, cfg *config.Config, log *logger.Logger) *gin.Engine {
	// Set gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))
	router.Use(middleware.LoggingMiddleware(log))

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(services.Auth, services.User, log)
	userHandler := handlers.NewUserHandler(services.User, log)
	dishHandler := handlers.NewDishHandler(services.Dish, services.User, log)
	mealHandler := handlers.NewMealHandler(services.Meal, log)
	analyticsHandler := handlers.NewAnalyticsHandler(services.Meal, log)
	shoppingListHandler := handlers.NewShoppingListHandler(services.Meal, log)
	recommendationsHandler := handlers.NewRecommendationsHandler(services.Dish, services.Meal, services.User, log)
	nutritionHandler := handlers.NewNutritionHandler(services.Meal, services.User, log)

	// Public routes
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", healthHandler.CheckHealth)

		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)

			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(services.Auth))
			{
				authProtected.GET("/me", authHandler.GetCurrentUser)
			}
		}

		// Dishes routes (with optional auth for favorites)
		dishes := api.Group("/dishes")
		dishes.Use(middleware.OptionalAuthMiddleware(services.Auth))
		{
			dishes.GET("", dishHandler.GetDishes)
			dishes.POST("", dishHandler.CreateDish)      // Add new dish
			dishes.GET("/search", dishHandler.GetDishes) // Alias for search functionality
			dishes.GET("/:id", dishHandler.GetDish)

			// Protected dish routes
			protected := dishes.Group("")
			protected.Use(middleware.AuthMiddleware(services.Auth))
			{
				protected.GET("/favorites", dishHandler.GetFavorites)
				protected.POST("/:id/favorite", dishHandler.AddToFavorites)
				protected.DELETE("/:id/favorite", dishHandler.RemoveFromFavorites)
			}
		}
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(services.Auth))
	{
		// User routes
		user := protected.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.DELETE("/account", userHandler.DeleteAccount)
		}

		// Meals routes
		meals := protected.Group("/meals")
		{
			meals.POST("", mealHandler.CreateMeal)
			meals.GET("", mealHandler.GetMeals)
			meals.GET("/nutrition-summary", mealHandler.GetNutritionSummary)
			meals.GET("/month/:year/:month", mealHandler.GetMealsByMonth) // Add month endpoint
			meals.GET("/:id", mealHandler.GetMeal)                        // This will handle both ObjectIDs and dates
			meals.PUT("/:id", mealHandler.UpdateMeal)
			meals.DELETE("/:id", mealHandler.DeleteMeal)
		}

		// Analytics routes
		analytics := protected.Group("/analytics")
		{
			analytics.GET("", analyticsHandler.GetMealAnalytics)
		}

		// Shopping List routes
		shoppingList := protected.Group("/shopping-list")
		{
			shoppingList.GET("", shoppingListHandler.GetShoppingList)
		}

		// Recommendations routes
		recommendations := protected.Group("/recommendations")
		{
			recommendations.GET("", recommendationsHandler.GetRecommendations)
		}

		// Nutrition routes
		nutrition := protected.Group("/nutrition")
		{
			nutrition.GET("/progress", nutritionHandler.GetNutritionProgress)
			nutrition.GET("/goals", nutritionHandler.GetNutritionGoals)
			nutrition.PUT("/goals", nutritionHandler.UpdateNutritionGoals)
		}
	}

	return router
}
