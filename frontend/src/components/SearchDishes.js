import React, { useState, useEffect, useMemo } from 'react';
import { Search, Filter, X, Heart, Star } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

const SearchDishes = ({ onSelectDish, selectedDish, mealType = 'lunch' }) => {
  const { user, toggleFavorite } = useAuth();
  const [searchTerm, setSearchTerm] = useState('');
  const [filters, setFilters] = useState({
    cuisine: '',
    type: '',
    maxCalories: ''
  });
  const [dishes, setDishes] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showFilters, setShowFilters] = useState(false);

  // Debounced search
  useEffect(() => {
    const delayedSearch = setTimeout(() => {
      searchDishes();
    }, 300);

    return () => clearTimeout(delayedSearch);
  }, [searchTerm, filters]);

  const searchDishes = async () => {
    setLoading(true);
    try {
      const params = {};
      if (searchTerm) params.q = searchTerm;
      if (filters.cuisine) params.cuisine = filters.cuisine;
      if (filters.type) params.type = filters.type;
      if (filters.maxCalories) params.maxCalories = filters.maxCalories;

      const response = Object.keys(params).length > 0
        ? await mealPlannerAPI.searchDishes(params)
        : await mealPlannerAPI.getDishes();
      
      setDishes(response.data);
    } catch (error) {
      console.error('Error searching dishes:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (key, value) => {
    setFilters(prev => ({
      ...prev,
      [key]: value
    }));
  };

  const clearFilters = () => {
    setSearchTerm('');
    setFilters({
      cuisine: '',
      type: '',
      maxCalories: ''
    });
  };

  const handleToggleFavorite = async (e, dishId) => {
    e.stopPropagation();
    const result = await toggleFavorite(dishId);
    if (result.success) {
      setDishes(prev => prev.map(dish => 
        dish.id === dishId 
          ? { ...dish, isFavorite: result.isFavorite }
          : dish
      ));
    }
  };

  const cuisineOptions = [
    'North Indian', 'South Indian', 'Bengali', 'Gujarati', 
    'Punjabi', 'Rajasthani', 'Maharashtrian'
  ];

  const typeOptions = ['breakfast', 'lunch', 'dinner', 'snack'];

  const filteredDishes = useMemo(() => {
    if (!searchTerm && !filters.cuisine && !filters.type && !filters.maxCalories) {
      return dishes;
    }
    return dishes;
  }, [dishes, searchTerm, filters]);

  return (
    <div className="space-y-4">
      {/* Search Bar */}
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          placeholder="Search dishes by name or ingredients..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
        />
        <button
          onClick={() => setShowFilters(!showFilters)}
          className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
        >
          <Filter className="h-5 w-5" />
        </button>
      </div>

      {/* Filters */}
      {showFilters && (
        <div className="bg-gray-50 rounded-lg p-4 space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-medium text-gray-700">Filters</h3>
            <button
              onClick={clearFilters}
              className="text-sm text-primary-600 hover:text-primary-700 flex items-center space-x-1"
            >
              <X className="h-4 w-4" />
              <span>Clear All</span>
            </button>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Cuisine
              </label>
              <select
                value={filters.cuisine}
                onChange={(e) => handleFilterChange('cuisine', e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="">All Cuisines</option>
                {cuisineOptions.map(cuisine => (
                  <option key={cuisine} value={cuisine}>{cuisine}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Meal Type
              </label>
              <select
                value={filters.type}
                onChange={(e) => handleFilterChange('type', e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="">All Types</option>
                {typeOptions.map(type => (
                  <option key={type} value={type}>
                    {type.charAt(0).toUpperCase() + type.slice(1)}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Max Calories
              </label>
              <input
                type="number"
                placeholder="e.g. 500"
                value={filters.maxCalories}
                onChange={(e) => handleFilterChange('maxCalories', e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              />
            </div>
          </div>
        </div>
      )}

      {/* Results */}
      <div className="space-y-2">
        {loading ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-2 border-primary-200 border-t-primary-600 mx-auto"></div>
            <p className="text-gray-500 mt-2">Searching dishes...</p>
          </div>
        ) : filteredDishes.length > 0 ? (
          <div className="max-h-96 overflow-y-auto space-y-2">
            {filteredDishes.map((dish) => (
              <div
                key={dish.id}
                onClick={() => onSelectDish(dish)}
                className={`
                  relative p-4 border rounded-lg cursor-pointer transition-all duration-200 hover:shadow-md
                  ${selectedDish?.id === dish.id 
                    ? 'border-primary-500 bg-primary-50' 
                    : 'border-gray-200 hover:border-gray-300'
                  }
                `}
              >
                <div className="flex items-start space-x-4">
                  <img
                    src={dish.image}
                    alt={dish.name}
                    className="w-16 h-16 object-cover rounded-lg"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h3 className="text-lg font-medium text-gray-900 truncate">
                          {dish.name}
                        </h3>
                        <p className="text-sm text-gray-500">
                          {dish.cuisine} • {dish.type} • {dish.calories} cal
                        </p>
                      </div>
                      {user && (
                        <button
                          onClick={(e) => handleToggleFavorite(e, dish.id)}
                          className={`ml-2 p-1 rounded-full transition-colors ${
                            dish.isFavorite
                              ? 'text-red-500 hover:text-red-600'
                              : 'text-gray-400 hover:text-red-500'
                          }`}
                        >
                          <Heart className={`h-5 w-5 ${dish.isFavorite ? 'fill-current' : ''}`} />
                        </button>
                      )}
                    </div>
                    <div className="mt-2">
                      <p className="text-sm text-gray-600 line-clamp-2">
                        <span className="font-medium">Ingredients:</span>{' '}
                        {dish.ingredients.join(', ')}
                      </p>
                    </div>
                  </div>
                </div>
                
                {selectedDish?.id === dish.id && (
                  <div className="absolute inset-0 border-2 border-primary-500 rounded-lg pointer-events-none"></div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <Search className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">
              {searchTerm || Object.values(filters).some(f => f) 
                ? 'No dishes found matching your criteria'
                : 'Start typing to search for dishes'
              }
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default SearchDishes;
