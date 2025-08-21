import React, { useState, useEffect } from 'react';
import { Lightbulb, Clock, Users, Flame, ChefHat, Star, Plus } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import CustomDropdown from './ui/CustomDropdown';

const MealRecommendations = ({ onAddMeal, selectedDate, mealType = 'lunch' }) => {
  const [recommendations, setRecommendations] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedMealType, setSelectedMealType] = useState(mealType);
  const { user } = useAuth();

  useEffect(() => {
    if (user) {
      fetchRecommendations();
    }
  }, [selectedMealType, selectedDate, user]);

  const fetchRecommendations = async () => {
    try {
      setLoading(true);
      const response = await mealPlannerAPI.getRecommendations(selectedMealType, selectedDate);
      console.log('Recommendations response:', response.data);
      // Backend returns { success: true, data: { recommendations: [...], reason: "..." } }
      if (response.data && response.data.success && response.data.data) {
        setRecommendations(response.data.data.recommendations || []);
      } else {
        console.error('Invalid recommendations response format:', response.data);
        setRecommendations([]);
      }
    } catch (error) {
      console.error('Error fetching recommendations:', error);
      setRecommendations([]);
    } finally {
      setLoading(false);
    }
  };

  const handleAddMeal = async (dish) => {
    if (onAddMeal) {
      await onAddMeal({
        date: selectedDate,
        mealType: selectedMealType,
        dishId: dish.dishId || dish.id  // Handle both frontend and backend property names
      });
      // Refresh recommendations after adding a meal
      fetchRecommendations();
    }
  };

  const getDifficultyColor = (difficulty) => {
    switch (difficulty) {
      case 'easy': return 'text-sage-600 bg-sage-100';
      case 'medium': return 'text-warm-600 bg-warm-100';
      case 'hard': return 'text-secondary-600 bg-secondary-100';
      default: return 'text-neutral-600 bg-neutral-100';
    }
  };

  const getSpiceLevelEmoji = (level) => {
    switch (level) {
      case 'mild': return '🌶️';
      case 'medium': return '🌶️🌶️';
      case 'hot': return '🌶️🌶️🌶️';
      case 'extra-hot': return '🌶️🌶️🌶️🌶️';
      default: return '🌶️🌶️';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-semibold flex items-center">
          <Lightbulb className="h-6 w-6 text-secondary-600 mr-2" />
          Meal Recommendations
        </h3>
        <CustomDropdown
          value={selectedMealType}
          onChange={setSelectedMealType}
          options={[
            { value: 'breakfast', label: 'Breakfast' },
            { value: 'lunch', label: 'Lunch' },
            { value: 'dinner', label: 'Dinner' },
            { value: 'snack', label: 'Snack' }
          ]}
          placeholder="Select Meal Type"
        />
      </div>

      {loading ? (
        <div className="space-y-4">
          {[1, 2, 3].map(i => (
            <div key={i} className="animate-pulse">
              <div className="bg-gray-200 rounded-lg h-32"></div>
            </div>
          ))}
        </div>
      ) : recommendations.length === 0 ? (
        <div className="text-center py-8">
          <ChefHat className="h-12 w-12 text-neutral-400 mx-auto mb-4" />
          <p className="text-neutral-600">No recommendations available.</p>
          <p className="text-sm text-neutral-500 mt-2">
            Update your preferences in your profile to get personalized suggestions!
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {recommendations.map((dish, index) => (
            <div key={dish.dishId || dish.id || index} className="border border-accent-200 rounded-lg p-4 hover:shadow-md transition-shadow bg-white">
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1">
                  <h4 className="font-medium text-neutral-900 mb-1">{dish.dishName || dish.name}</h4>
                  <p className="text-sm text-neutral-600 mb-2">{dish.cuisine} • {dish.type || 'Main Course'}</p>
                  
                  {/* Recommendation Score */}
                  <div className="flex items-center space-x-2 mb-2">
                    <Star className="h-4 w-4 text-warm-500" />
                    <span className="text-sm text-neutral-600">
                      Score: {Math.round((dish.score || dish.recommendationScore || 0) * 10) / 10}/10
                    </span>
                    {index < 3 && (
                      <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs bg-accent-100 text-accent-800">
                        Top Pick
                      </span>
                    )}
                  </div>
                </div>
                
                {dish.image && (
                  <img 
                    src={dish.image} 
                    alt={dish.dishName || dish.name}
                    className="w-16 h-16 rounded-lg object-cover ml-4"
                  />
                )}
              </div>

              {/* Nutrition & Details */}
              <div className="grid grid-cols-2 gap-4 text-sm text-neutral-600 mb-3">
                <div className="flex items-center">
                  <Flame className="h-4 w-4 text-warm-500 mr-1" />
                  {dish.calories || 'N/A'} cal
                </div>
                <div className="flex items-center">
                  <Clock className="h-4 w-4 text-sage-500 mr-1" />
                  {dish.prepTime || 30} min
                </div>
                <div className="flex items-center">
                  <Users className="h-4 w-4 text-secondary-500 mr-1" />
                  {dish.servings || 2} servings
                </div>
                <div className="flex items-center">
                  <ChefHat className={`h-4 w-4 mr-1 ${getDifficultyColor(dish.difficulty)}`} />
                  <span className={`capitalize text-xs px-2 py-0.5 rounded-full ${getDifficultyColor(dish.difficulty)}`}>
                    {dish.difficulty || 'medium'}
                  </span>
                </div>
              </div>

              {/* Dietary Tags */}
              {dish.dietaryTags && dish.dietaryTags.length > 0 && (
                <div className="mb-3">
                  <div className="flex flex-wrap gap-1">
                    {dish.dietaryTags.slice(0, 3).map(tag => (
                      <span 
                        key={tag}
                        className="inline-flex items-center px-2 py-0.5 rounded-md text-xs bg-lavender-100 text-lavender-800"
                      >
                        {tag}
                      </span>
                    ))}
                    {dish.dietaryTags.length > 3 && (
                      <span className="text-xs text-neutral-500">
                        +{dish.dietaryTags.length - 3} more
                      </span>
                    )}
                  </div>
                </div>
              )}

              {/* Difficulty */}
              {dish.difficulty && (
                <div className="mb-3">
                  <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getDifficultyColor(dish.difficulty)}`}>
                    {dish.difficulty} difficulty
                  </span>
                </div>
              )}

              {/* Nutrition Info */}
              {dish.nutrition && (
                <div className="text-xs text-gray-500 mb-3">
                  Protein: {dish.nutrition.protein || 0}g • 
                  Carbs: {dish.nutrition.carbs || 0}g • 
                  Fat: {dish.nutrition.fat || 0}g
                </div>
              )}

              {/* Add Button */}
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  {dish.isFavorite && (
                    <span className="text-red-500 text-sm">❤️ Favorite</span>
                  )}
                </div>
                <button
                  onClick={() => handleAddMeal(dish)}
                  className="flex items-center px-3 py-1.5 bg-accent-600 text-white rounded-md hover:bg-accent-700 transition-colors text-sm shadow-sm"
                >
                  <Plus className="h-4 w-4 mr-1" />
                  Add to {selectedMealType}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Why These Recommendations? */}
      {recommendations.length > 0 && user && (
        <div className="mt-6 p-4 bg-accent-50 rounded-lg border border-accent-200">
          <h4 className="font-medium text-accent-900 mb-2">Why these recommendations?</h4>
          <div className="text-sm text-accent-700 space-y-1">
            {user.profile.dietaryPreferences.length > 0 && (
              <p>• Matches your dietary preferences: {user.profile.dietaryPreferences.join(', ')}</p>
            )}
            {user.profile.favoriteRegions.length > 0 && (
              <p>• Includes your favorite cuisines: {user.profile.favoriteRegions.join(', ')}</p>
            )}
            <p>• Avoids dishes you've had recently</p>
            <p>• Appropriate calories for {selectedMealType}</p>
            <p>• Considers your spice tolerance</p>
          </div>
        </div>
      )}
    </div>
  );
};

export default MealRecommendations;
