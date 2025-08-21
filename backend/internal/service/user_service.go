package service

import (
	"context"
	"errors"

	"meal-planner-backend/internal/models"
	"meal-planner-backend/internal/repository"
	"meal-planner-backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService interface defines user operations
type UserService interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	UpdateProfile(ctx context.Context, userID primitive.ObjectID, profile models.UserProfile) error
	UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, req models.ProfileUpdateRequest) error
	AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error
	RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error
	GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error)
	Delete(ctx context.Context, userID primitive.ObjectID) error
}

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, log *logger.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   log,
	}
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		s.logger.Error("Failed to get user by ID", "error", err, "userID", id.Hex())
		return nil, errors.New("internal server error")
	}

	return user, nil
}

// UpdateProfile updates user profile
func (s *userService) UpdateProfile(ctx context.Context, userID primitive.ObjectID, profile models.UserProfile) error {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		s.logger.Error("Failed to get user for profile update", "error", err, "userID", userID.Hex())
		return errors.New("internal server error")
	}

	// Update profile
	user.Profile = profile

	if err := s.userRepo.Update(ctx, userID, user); err != nil {
		s.logger.Error("Failed to update user profile", "error", err, "userID", userID.Hex())
		return errors.New("failed to update profile")
	}

	return nil
}

// UpdateUserProfile updates user profile and optionally the name
func (s *userService) UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, req models.ProfileUpdateRequest) error {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		s.logger.Error("Failed to get user for profile update", "error", err, "userID", userID.Hex())
		return errors.New("internal server error")
	}

	// Update profile
	user.Profile = req.Profile

	// Update name if provided
	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.userRepo.Update(ctx, userID, user); err != nil {
		s.logger.Error("Failed to update user profile", "error", err, "userID", userID.Hex())
		return errors.New("failed to update profile")
	}

	return nil
}

// AddToFavorites adds a dish to user's favorites
func (s *userService) AddToFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	if err := s.userRepo.AddToFavorites(ctx, userID, dishID); err != nil {
		s.logger.Error("Failed to add dish to favorites", "error", err, "userID", userID.Hex(), "dishID", dishID.Hex())
		return errors.New("failed to add to favorites")
	}

	return nil
}

// RemoveFromFavorites removes a dish from user's favorites
func (s *userService) RemoveFromFavorites(ctx context.Context, userID, dishID primitive.ObjectID) error {
	if err := s.userRepo.RemoveFromFavorites(ctx, userID, dishID); err != nil {
		s.logger.Error("Failed to remove dish from favorites", "error", err, "userID", userID.Hex(), "dishID", dishID.Hex())
		return errors.New("failed to remove from favorites")
	}

	return nil
}

// GetFavorites gets user's favorite dish IDs
func (s *userService) GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	favorites, err := s.userRepo.GetFavorites(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user favorites", "error", err, "userID", userID.Hex())
		return nil, errors.New("failed to get favorites")
	}

	return favorites, nil
}

// Delete deletes a user account
func (s *userService) Delete(ctx context.Context, userID primitive.ObjectID) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Error("Failed to delete user", "error", err, "userID", userID.Hex())
		return errors.New("failed to delete user")
	}

	return nil
}
