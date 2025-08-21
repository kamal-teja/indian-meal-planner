package service

import (
	"context"
	"errors"

	"nourish-backend/internal/models"
	"nourish-backend/internal/repository"
	"nourish-backend/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MealPlanService interface defines meal plan operations
type MealPlanService interface {
	Create(ctx context.Context, userID primitive.ObjectID, req models.MealPlanRequest) (*models.MealPlan, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.MealPlan, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.MealPlan, *models.PaginationResponse, error)
	GetActivePlans(ctx context.Context, userID primitive.ObjectID) ([]*models.MealPlan, error)
	Update(ctx context.Context, id primitive.ObjectID, req models.MealPlanRequest) (*models.MealPlan, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// mealPlanService implements MealPlanService interface
type mealPlanService struct {
	mealPlanRepo repository.MealPlanRepository
	dishRepo     repository.DishRepository
	logger       *logger.Logger
}

// NewMealPlanService creates a new meal plan service
func NewMealPlanService(mealPlanRepo repository.MealPlanRepository, dishRepo repository.DishRepository, log *logger.Logger) MealPlanService {
	return &mealPlanService{
		mealPlanRepo: mealPlanRepo,
		dishRepo:     dishRepo,
		logger:       log,
	}
}

// Create creates a new meal plan
func (s *mealPlanService) Create(ctx context.Context, userID primitive.ObjectID, req models.MealPlanRequest) (*models.MealPlan, error) {
	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Create meal plan
	mealPlan := &models.MealPlan{
		UserID:      userID,
		Name:        req.Name,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
		Meals:       []models.MealPlanMeal{},
	}

	if err := s.mealPlanRepo.Create(ctx, mealPlan); err != nil {
		s.logger.Error("Failed to create meal plan", "error", err)
		return nil, errors.New("failed to create meal plan")
	}

	return mealPlan, nil
}

// GetByID retrieves a meal plan by ID
func (s *mealPlanService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.MealPlan, error) {
	mealPlan, err := s.mealPlanRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("meal plan not found")
		}
		s.logger.Error("Failed to get meal plan by ID", "error", err, "mealPlanID", id.Hex())
		return nil, errors.New("internal server error")
	}

	return mealPlan, nil
}

// GetByUserID retrieves meal plans for a user with pagination
func (s *mealPlanService) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.MealPlan, *models.PaginationResponse, error) {
	mealPlans, total, err := s.mealPlanRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		s.logger.Error("Failed to get meal plans by user ID", "error", err, "userID", userID.Hex())
		return nil, nil, errors.New("failed to get meal plans")
	}

	// Create pagination response
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	pagination := &models.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return mealPlans, pagination, nil
}

// GetActivePlans retrieves active meal plans for a user
func (s *mealPlanService) GetActivePlans(ctx context.Context, userID primitive.ObjectID) ([]*models.MealPlan, error) {
	mealPlans, err := s.mealPlanRepo.GetActivePlans(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get active meal plans", "error", err, "userID", userID.Hex())
		return nil, errors.New("failed to get active meal plans")
	}

	return mealPlans, nil
}

// Update updates a meal plan
func (s *mealPlanService) Update(ctx context.Context, id primitive.ObjectID, req models.MealPlanRequest) (*models.MealPlan, error) {
	// Get existing meal plan
	existingPlan, err := s.mealPlanRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("meal plan not found")
		}
		s.logger.Error("Failed to get meal plan for update", "error", err, "mealPlanID", id.Hex())
		return nil, errors.New("internal server error")
	}

	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Update meal plan
	existingPlan.Name = req.Name
	existingPlan.StartDate = req.StartDate
	existingPlan.EndDate = req.EndDate
	existingPlan.Description = req.Description

	if err := s.mealPlanRepo.Update(ctx, id, existingPlan); err != nil {
		s.logger.Error("Failed to update meal plan", "error", err, "mealPlanID", id.Hex())
		return nil, errors.New("failed to update meal plan")
	}

	return existingPlan, nil
}

// Delete deletes a meal plan
func (s *mealPlanService) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Check if meal plan exists
	_, err := s.mealPlanRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("meal plan not found")
		}
		s.logger.Error("Failed to get meal plan for deletion", "error", err, "mealPlanID", id.Hex())
		return errors.New("internal server error")
	}

	if err := s.mealPlanRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete meal plan", "error", err, "mealPlanID", id.Hex())
		return errors.New("failed to delete meal plan")
	}

	return nil
}
