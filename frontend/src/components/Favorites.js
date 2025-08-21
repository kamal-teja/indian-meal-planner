import React, { useState, useEffect } from 'react';
import { Heart, Search, Calendar, Plus } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

const Favorites = () => {
  const { user, toggleFavorite } = useAuth();
  const [favoriteDishes, setFavoriteDishes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    loadFavorites();
  }, []);

  const loadFavorites = async () => {
    try {
      setLoading(true);
      const response = await mealPlannerAPI.getFavoriteDishes();
      // Handle the response format from backend: { success: true, dishes: [...], pagination: {...} }
      setFavoriteDishes(response.data.dishes || []);
    } catch (error) {
      console.error('Error loading favorites:', error);
      setFavoriteDishes([]); // Set empty array on error
    } finally {
      setLoading(false);
    }
  };

  const handleToggleFavorite = async (dishId) => {
    const result = await toggleFavorite(dishId);
    if (result.success) {
      // Remove from favorites since we're toggling it off
      setFavoriteDishes(prev => prev.filter(dish => dish.id !== dishId));
    }
  };

  const filteredDishes = favoriteDishes.filter(dish =>
    dish.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    dish.cuisine.toLowerCase().includes(searchTerm.toLowerCase()) ||
    dish.ingredients.some(ingredient => 
      ingredient.toLowerCase().includes(searchTerm.toLowerCase())
    )
  );

  const addToMealPlan = async (dish, mealType = 'lunch') => {
    try {
      const today = new Date().toISOString().split('T')[0];
      await mealPlannerAPI.addMeal({
        date: today,
        mealType,
        dishId: dish.id
      });
      // Show success notification (you could add a toast notification here)
      alert(`${dish.name} added to today's ${mealType}!`);
    } catch (error) {
      console.error('Error adding to meal plan:', error);
      alert('Failed to add dish to meal plan. Please try again.');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-4 border-primary-200 border-t-primary-600 mx-auto mb-4"></div>
          <h2 className="text-2xl font-display font-semibold gradient-text">
            Loading Favorites...
          </h2>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center space-x-3 mb-4">
          <div className="p-2 bg-gradient-to-r from-red-500 to-pink-600 rounded-xl shadow-lg">
            <Heart className="h-8 w-8 text-white fill-current" />
          </div>
          <div>
            <h1 className="text-3xl font-display font-bold gradient-text">
              Favorite Dishes
            </h1>
            <p className="text-gray-600">
              Your collection of loved dishes
            </p>
          </div>
        </div>

        {/* Search Bar */}
        {favoriteDishes.length > 0 && (
          <div className="relative max-w-md">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              placeholder="Search your favorites..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            />
          </div>
        )}
      </div>

      {/* Content */}
      {favoriteDishes.length === 0 ? (
        <div className="text-center py-16">
          <div className="p-4 bg-gray-100 rounded-full w-24 h-24 mx-auto mb-6 flex items-center justify-center">
            <Heart className="h-12 w-12 text-gray-400" />
          </div>
          <h3 className="text-2xl font-display font-semibold text-gray-700 mb-2">
            No Favorite Dishes Yet
          </h3>
          <p className="text-gray-500 mb-6 max-w-md mx-auto">
            Start exploring dishes and click the heart icon to add them to your favorites.
            You'll see them all here for easy access.
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredDishes.map((dish) => (
            <div key={dish.id} className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
              {/* Image */}
              <div className="relative">
                <img
                  src={dish.image}
                  alt={dish.name}
                  className="w-full h-48 object-cover"
                />
                <button
                  onClick={() => handleToggleFavorite(dish.id)}
                  className="absolute top-3 right-3 p-2 bg-white/90 rounded-full text-red-500 hover:bg-white transition-colors"
                >
                  <Heart className="h-5 w-5 fill-current" />
                </button>
              </div>

              {/* Content */}
              <div className="p-6">
                <div className="mb-3">
                  <h3 className="text-xl font-display font-semibold text-gray-900 mb-1">
                    {dish.name}
                  </h3>
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <span className="bg-primary-100 text-primary-700 px-2 py-1 rounded-full">
                      {dish.cuisine}
                    </span>
                    <span>•</span>
                    <span>{dish.calories} cal</span>
                  </div>
                </div>

                <div className="mb-4">
                  <p className="text-sm text-gray-600 line-clamp-3">
                    <span className="font-medium">Ingredients:</span>{' '}
                    {dish.ingredients.join(', ')}
                  </p>
                </div>

                {/* Quick Add Buttons */}
                <div className="space-y-2">
                  <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                    Quick Add to Meal Plan
                  </p>
                  <div className="grid grid-cols-2 gap-2">
                    <button
                      onClick={() => addToMealPlan(dish, 'lunch')}
                      className="flex items-center justify-center space-x-1 px-3 py-2 bg-primary-100 text-primary-700 rounded-lg hover:bg-primary-200 transition-colors text-sm font-medium"
                    >
                      <Plus className="h-4 w-4" />
                      <span>Lunch</span>
                    </button>
                    <button
                      onClick={() => addToMealPlan(dish, 'dinner')}
                      className="flex items-center justify-center space-x-1 px-3 py-2 bg-secondary-100 text-secondary-700 rounded-lg hover:bg-secondary-200 transition-colors text-sm font-medium"
                    >
                      <Plus className="h-4 w-4" />
                      <span>Dinner</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* No search results */}
      {favoriteDishes.length > 0 && filteredDishes.length === 0 && searchTerm && (
        <div className="text-center py-16">
          <Search className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-xl font-display font-semibold text-gray-700 mb-2">
            No dishes found
          </h3>
          <p className="text-gray-500">
            Try adjusting your search term to find your favorite dishes.
          </p>
        </div>
      )}
    </div>
  );
};

export default Favorites;
