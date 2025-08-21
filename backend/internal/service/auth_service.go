package service

import (
	"context"
	"errors"
	"time"

	"meal-planner-backend/internal/config"
	"meal-planner-backend/internal/models"
	"meal-planner-backend/internal/repository"
	"meal-planner-backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// AuthService interface defines authentication operations
type AuthService interface {
	Register(ctx context.Context, req models.UserRegistrationRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.UserLoginRequest) (*models.AuthResponse, error)
	GenerateToken(userID primitive.ObjectID) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserFromToken(tokenString string) (*models.User, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
	logger   *logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config, log *logger.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   cfg,
		logger:   log,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req models.UserRegistrationRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		s.logger.Error("Failed to check existing user", "error", err)
		return nil, errors.New("internal server error")
	}

	if existingUser != nil {
		return nil, errors.New("an account with this email address already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, errors.New("internal server error")
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Profile: models.UserProfile{
			DietaryPreferences: []string{},
			SpiceLevel:         "medium",
			FavoriteRegions:    []string{},
			NutritionGoals: models.NutritionGoals{
				DailyCalories: 2000,
				Protein:       150,
				Carbs:         250,
				Fat:           65,
				Fiber:         25,
				Sodium:        2300,
			},
		},
		Favorites: []primitive.ObjectID{},
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, errors.New("failed to create user account")
	}

	// Generate token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token", "error", err)
		return nil, errors.New("internal server error")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Warn("Failed to update last login", "error", err)
	}

	return &models.AuthResponse{
		Success: true,
		Message: "User registered successfully",
		Token:   token,
		User:    user.ToResponse(),
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, req models.UserLoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("invalid email or password")
		}
		s.logger.Error("Failed to get user by email", "error", err)
		return nil, errors.New("internal server error")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token", "error", err)
		return nil, errors.New("internal server error")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Warn("Failed to update last login", "error", err)
	}

	return &models.AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    user.ToResponse(),
	}, nil
}

// GenerateToken generates a JWT token for a user
func (s *authService) GenerateToken(userID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID.Hex(),
		"exp":    time.Now().Add(s.config.JWTExpiresIn).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// ValidateToken validates a JWT token
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// GetUserFromToken extracts user from JWT token
func (s *authService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["userId"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	return user, nil
}
