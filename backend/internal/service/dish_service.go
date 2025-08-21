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

// DishService interface defines dish operations
type DishService interface {
	GetByID(ctx context.Context, id primitive.ObjectID, userID *primitive.ObjectID) (*models.DishResponse, error)
	GetAll(ctx context.Context, filter DishFilter, page, limit int, userID *primitive.ObjectID) ([]*models.DishResponse, *models.PaginationResponse, error)
	Search(ctx context.Context, query string, filter DishFilter, page, limit int, userID *primitive.ObjectID) ([]*models.DishResponse, *models.PaginationResponse, error)
	GetFavorites(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.DishResponse, *models.PaginationResponse, error)
	Create(ctx context.Context, dish *models.Dish) error
	Update(ctx context.Context, id primitive.ObjectID, dish *models.Dish) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// DishFilter represents filters for dish queries
type DishFilter struct {
	Type         string
	Cuisine      string
	DietaryTags  []string
	SpiceLevel   string
	MaxCalories  int
	MinCalories  int
	Ingredients  []string
}

// dishService implements DishService interface
type dishService struct {
	dishRepo repository.DishRepository
	userRepo repository.UserRepository
	logger   *logger.Logger
}

// NewDishService creates a new dish service
func NewDishService(dishRepo repository.DishRepository, userRepo repository.UserRepository, log *logger.Logger) DishService {
	return &dishService{
		dishRepo: dishRepo,
		userRepo: userRepo,
		logger:   log,
	}
}

// GetByID retrieves a dish by ID with favorite status
func (s *dishService) GetByID(ctx context.Context, id primitive.ObjectID, userID *primitive.ObjectID) (*models.DishResponse, error) {
	dish, err := s.dishRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("dish not found")
		}
		s.logger.Error("Failed to get dish by ID", "error", err, "dishID", id.Hex())
		return nil, errors.New("internal server error")
	}

	dishResponse := dish.ToResponse()

	// Set favorite status if user is provided
	if userID != nil {
		favorites, err := s.userRepo.GetFavorites(ctx, *userID)
		if err == nil {
			dishResponse.IsFavorite = s.isDishInFavorites(id, favorites)
		}
	}

	return &dishResponse, nil
}

// GetAll retrieves dishes with pagination and filtering
func (s *dishService) GetAll(ctx context.Context, filter DishFilter, page, limit int, userID *primitive.ObjectID) ([]*models.DishResponse, *models.PaginationResponse, error) {
	// Convert service filter to repository filter
	repoFilter := repository.DishFilter{
		Type:        filter.Type,
		Cuisine:     filter.Cuisine,
		DietaryTags: filter.DietaryTags,
		SpiceLevel:  filter.SpiceLevel,
		MaxCalories: filter.MaxCalories,
		MinCalories: filter.MinCalories,
		Ingredients: filter.Ingredients,
	}

	dishes, total, err := s.dishRepo.GetAll(ctx, repoFilter, page, limit)
	if err != nil {
		s.logger.Error("Failed to get dishes", "error", err)
		return nil, nil, errors.New("failed to get dishes")
	}

	// Get user favorites if user is provided
	var favorites []primitive.ObjectID
	if userID != nil {
		favorites, _ = s.userRepo.GetFavorites(ctx, *userID)
	}

	// Convert to response format
	dishResponses := make([]*models.DishResponse, len(dishes))
	for i, dish := range dishes {
		dishResponse := dish.ToResponse()
		if userID != nil {
			dishResponse.IsFavorite = s.isDishInFavorites(dish.ID, favorites)
		}
		dishResponses[i] = &dishResponse
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

	return dishResponses, pagination, nil
}

// Search searches dishes with text search and filtering
func (s *dishService) Search(ctx context.Context, query string, filter DishFilter, page, limit int, userID *primitive.ObjectID) ([]*models.DishResponse, *models.PaginationResponse, error) {
	// Convert service filter to repository filter
	repoFilter := repository.DishFilter{
		Type:        filter.Type,
		Cuisine:     filter.Cuisine,
		DietaryTags: filter.DietaryTags,
		SpiceLevel:  filter.SpiceLevel,
		MaxCalories: filter.MaxCalories,
		MinCalories: filter.MinCalories,
		Ingredients: filter.Ingredients,
	}

	dishes, total, err := s.dishRepo.Search(ctx, query, repoFilter, page, limit)
	if err != nil {
		s.logger.Error("Failed to search dishes", "error", err, "query", query)
		return nil, nil, errors.New("failed to search dishes")
	}

	// Get user favorites if user is provided
	var favorites []primitive.ObjectID
	if userID != nil {
		favorites, _ = s.userRepo.GetFavorites(ctx, *userID)
	}

	// Convert to response format
	dishResponses := make([]*models.DishResponse, len(dishes))
	for i, dish := range dishes {
		dishResponse := dish.ToResponse()
		if userID != nil {
			dishResponse.IsFavorite = s.isDishInFavorites(dish.ID, favorites)
		}
		dishResponses[i] = &dishResponse
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

	return dishResponses, pagination, nil
}

// GetFavorites retrieves user's favorite dishes
func (s *dishService) GetFavorites(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.DishResponse, *models.PaginationResponse, error) {
	// Get user's favorite dish IDs
	favoriteIDs, err := s.userRepo.GetFavorites(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user favorites", "error", err, "userID", userID.Hex())
		return nil, nil, errors.New("failed to get favorites")
	}

	if len(favoriteIDs) == 0 {
		return []*models.DishResponse{}, &models.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      0,
			TotalPages: 0,
			HasNext:    false,
			HasPrev:    false,
		}, nil
	}

	// Calculate pagination for favorites
	total := len(favoriteIDs)
	start := (page - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}
	if start >= total {
		start = total
	}

	paginatedIDs := favoriteIDs[start:end]

	// Get dishes by IDs
	dishes, err := s.dishRepo.GetByIDs(ctx, paginatedIDs)
	if err != nil {
		s.logger.Error("Failed to get favorite dishes", "error", err, "userID", userID.Hex())
		return nil, nil, errors.New("failed to get favorite dishes")
	}

	// Convert to response format
	dishResponses := make([]*models.DishResponse, len(dishes))
	for i, dish := range dishes {
		dishResponse := dish.ToResponse()
		dishResponse.IsFavorite = true // All dishes here are favorites
		dishResponses[i] = &dishResponse
	}

	// Create pagination response
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	pagination := &models.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return dishResponses, pagination, nil
}

// Create creates a new dish
func (s *dishService) Create(ctx context.Context, dish *models.Dish) error {
	if err := s.dishRepo.Create(ctx, dish); err != nil {
		s.logger.Error("Failed to create dish", "error", err)
		return errors.New("failed to create dish")
	}

	return nil
}

// Update updates a dish
func (s *dishService) Update(ctx context.Context, id primitive.ObjectID, dish *models.Dish) error {
	// Check if dish exists
	_, err := s.dishRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("dish not found")
		}
		s.logger.Error("Failed to get dish for update", "error", err, "dishID", id.Hex())
		return errors.New("internal server error")
	}

	if err := s.dishRepo.Update(ctx, id, dish); err != nil {
		s.logger.Error("Failed to update dish", "error", err, "dishID", id.Hex())
		return errors.New("failed to update dish")
	}

	return nil
}

// Delete deletes a dish
func (s *dishService) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Check if dish exists
	_, err := s.dishRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("dish not found")
		}
		s.logger.Error("Failed to get dish for deletion", "error", err, "dishID", id.Hex())
		return errors.New("internal server error")
	}

	if err := s.dishRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete dish", "error", err, "dishID", id.Hex())
		return errors.New("failed to delete dish")
	}

	return nil
}

// isDishInFavorites checks if a dish ID is in the favorites list
func (s *dishService) isDishInFavorites(dishID primitive.ObjectID, favorites []primitive.ObjectID) bool {
	for _, fav := range favorites {
		if fav == dishID {
			return true
		}
	}
	return false
}
