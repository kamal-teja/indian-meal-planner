package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name" validate:"required,min=2,max=50"`
	Email    string             `bson:"email" json:"email" validate:"required,email"`
	Password string             `bson:"password" json:"-" validate:"required,min=6"`
	Profile  UserProfile        `bson:"profile" json:"profile"`

	// User metadata
	Favorites   []primitive.ObjectID `bson:"favorites" json:"favorites"`
	LastLoginAt *time.Time           `bson:"lastLoginAt" json:"lastLoginAt"`
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// UserProfile contains user's dietary preferences and nutrition goals
type UserProfile struct {
	DietaryPreferences []string       `bson:"dietaryPreferences" json:"dietaryPreferences"`
	SpiceLevel         string         `bson:"spiceLevel" json:"spiceLevel" validate:"omitempty,oneof=mild medium hot extra-hot"`
	FavoriteRegions    []string       `bson:"favoriteRegions" json:"favoriteRegions"`
	Avatar             string         `bson:"avatar" json:"avatar"`
	NutritionGoals     NutritionGoals `bson:"nutritionGoals" json:"nutritionGoals"`
}

// NutritionGoals represents daily nutrition targets
type NutritionGoals struct {
	DailyCalories int `bson:"dailyCalories" json:"dailyCalories" validate:"omitempty,min=0,max=10000"`
	Protein       int `bson:"protein" json:"protein" validate:"min=0"` // grams
	Carbs         int `bson:"carbs" json:"carbs" validate:"min=0"`     // grams
	Fat           int `bson:"fat" json:"fat" validate:"min=0"`         // grams
	Fiber         int `bson:"fiber" json:"fiber" validate:"min=0"`     // grams
	Sodium        int `bson:"sodium" json:"sodium" validate:"min=0"`   // milligrams
}

// UserRegistrationRequest represents the request for user registration
type UserRegistrationRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLoginRequest represents the request for user login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Email   string      `json:"email"`
	Profile UserProfile `json:"profile"`
}

// ProfileUpdateRequest represents the request for updating user profile
type ProfileUpdateRequest struct {
	Name    string      `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Profile UserProfile `json:"profile"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:      u.ID.Hex(),
		Name:    u.Name,
		Email:   u.Email,
		Profile: u.Profile,
	}
}

// GetDietaryPreferences returns the list of valid dietary preferences
func GetDietaryPreferences() []string {
	return []string{
		"vegetarian", "vegan", "gluten-free", "dairy-free", "nut-free",
		"keto", "paleo", "low-carb", "high-protein", "low-sodium", "sugar-free",
	}
}

// GetFavoriteRegions returns the list of valid favorite regions
func GetFavoriteRegions() []string {
	return []string{
		"North Indian", "South Indian", "Bengali", "Gujarati", "Punjabi",
		"Rajasthani", "Maharashtrian", "Italian", "Chinese", "Thai",
		"French", "Korean", "Continental",
	}
}
