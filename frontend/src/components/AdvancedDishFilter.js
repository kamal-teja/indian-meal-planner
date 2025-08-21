import React, { useState, useEffect } from 'react';
import { Filter, X, Search, Clock, Flame, ChefHat, Tag } from 'lucide-react';

const AdvancedDishFilter = ({ onFilterChange, onClose }) => {
  const [filters, setFilters] = useState({
    q: '',
    cuisine: '',
    type: '',
    maxCalories: '',
    dietaryTags: [],
    spiceLevel: '',
    maxPrepTime: '',
    difficulty: ''
  });

  const [isExpanded, setIsExpanded] = useState(false);

  const dietaryOptions = [
    'vegetarian', 'vegan', 'gluten-free', 'dairy-free', 
    'nut-free', 'keto', 'paleo', 'low-carb', 
    'high-protein', 'low-sodium', 'sugar-free'
  ];

  const cuisineOptions = [
    'North Indian', 'South Indian', 'Bengali', 'Gujarati', 
    'Punjabi', 'Rajasthani', 'Maharashtrian', 'Italian', 
    'Chinese', 'Thai', 'French', 'Korean', 'Continental'
  ];

  const spiceLevels = ['mild', 'medium', 'hot', 'extra-hot'];
  const difficulties = ['easy', 'medium', 'hard'];

  useEffect(() => {
    // Debounce the filter changes
    const timeoutId = setTimeout(() => {
      onFilterChange(filters);
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [filters, onFilterChange]);

  const handleInputChange = (key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  };

  const handleDietaryTagToggle = (tag) => {
    setFilters(prev => ({
      ...prev,
      dietaryTags: prev.dietaryTags.includes(tag)
        ? prev.dietaryTags.filter(t => t !== tag)
        : [...prev.dietaryTags, tag]
    }));
  };

  const clearFilters = () => {
    setFilters({
      q: '',
      cuisine: '',
      type: '',
      maxCalories: '',
      dietaryTags: [],
      spiceLevel: '',
      maxPrepTime: '',
      difficulty: ''
    });
  };

  const getActiveFilterCount = () => {
    return Object.entries(filters).reduce((count, [key, value]) => {
      if (key === 'dietaryTags') return count + value.length;
      if (value && value !== '') return count + 1;
      return count;
    }, 0);
  };

  return (
    <div className="bg-white rounded-lg shadow-lg p-6 max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold flex items-center">
          <Filter className="h-6 w-6 text-blue-600 mr-2" />
          Advanced Filters
          {getActiveFilterCount() > 0 && (
            <span className="ml-2 bg-blue-100 text-blue-800 text-sm px-2 py-1 rounded-full">
              {getActiveFilterCount()} active
            </span>
          )}
        </h2>
        <div className="flex items-center space-x-2">
          <button
            onClick={clearFilters}
            className="text-gray-600 hover:text-gray-800 text-sm"
          >
            Clear all
          </button>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
      </div>

      {/* Search */}
      <div className="mb-6">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          <Search className="h-4 w-4 inline mr-1" />
          Search dishes or ingredients
        </label>
        <input
          type="text"
          value={filters.q}
          onChange={(e) => handleInputChange('q', e.target.value)}
          placeholder="e.g., chicken curry, pasta, tomato..."
          className="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      {/* Basic Filters */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Cuisine</label>
          <select
            value={filters.cuisine}
            onChange={(e) => handleInputChange('cuisine', e.target.value)}
            className="dropdown-elegant w-full"
          >
            <option value="">All cuisines</option>
            {cuisineOptions.map(cuisine => (
              <option key={cuisine} value={cuisine}>{cuisine}</option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Dish Type</label>
          <select
            value={filters.type}
            onChange={(e) => handleInputChange('type', e.target.value)}
            className="dropdown-elegant w-full"
          >
            <option value="">All types</option>
            <option value="Veg">Vegetarian</option>
            <option value="Non-Veg">Non-Vegetarian</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            <Flame className="h-4 w-4 inline mr-1" />
            Max Calories
          </label>
          <input
            type="number"
            value={filters.maxCalories}
            onChange={(e) => handleInputChange('maxCalories', e.target.value)}
            placeholder="e.g., 500"
            className="w-full border border-gray-300 rounded-md px-3 py-2"
            min="0"
          />
        </div>
      </div>

      {/* Advanced Filters Toggle */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="flex items-center text-blue-600 hover:text-blue-700 mb-4"
      >
        <ChefHat className="h-4 w-4 mr-1" />
        {isExpanded ? 'Hide' : 'Show'} advanced options
      </button>

      {/* Advanced Filters */}
      {isExpanded && (
        <div className="space-y-6 border-t border-gray-200 pt-6">
          {/* Dietary Tags */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">
              <Tag className="h-4 w-4 inline mr-1" />
              Dietary Preferences
            </label>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
              {dietaryOptions.map(tag => (
                <button
                  key={tag}
                  onClick={() => handleDietaryTagToggle(tag)}
                  className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    filters.dietaryTags.includes(tag)
                      ? 'bg-blue-600 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {tag}
                </button>
              ))}
            </div>
          </div>

          {/* Spice Level & Other Options */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Spice Level</label>
              <select
                value={filters.spiceLevel}
                onChange={(e) => handleInputChange('spiceLevel', e.target.value)}
                className="dropdown-elegant w-full"
              >
                <option value="">Any spice level</option>
                {spiceLevels.map(level => (
                  <option key={level} value={level}>
                    {level.charAt(0).toUpperCase() + level.slice(1)}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <Clock className="h-4 w-4 inline mr-1" />
                Max Prep Time (minutes)
              </label>
              <input
                type="number"
                value={filters.maxPrepTime}
                onChange={(e) => handleInputChange('maxPrepTime', e.target.value)}
                placeholder="e.g., 30"
                className="w-full border border-gray-300 rounded-md px-3 py-2"
                min="0"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Difficulty</label>
              <select
                value={filters.difficulty}
                onChange={(e) => handleInputChange('difficulty', e.target.value)}
                className="dropdown-elegant w-full"
              >
                <option value="">Any difficulty</option>
                {difficulties.map(diff => (
                  <option key={diff} value={diff}>
                    {diff.charAt(0).toUpperCase() + diff.slice(1)}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
      )}

      {/* Active Filters Summary */}
      {getActiveFilterCount() > 0 && (
        <div className="mt-6 p-4 bg-blue-50 rounded-lg">
          <h4 className="font-medium text-blue-900 mb-2">Active Filters:</h4>
          <div className="flex flex-wrap gap-2">
            {filters.q && (
              <span className="inline-flex items-center px-2 py-1 rounded-md text-sm bg-blue-100 text-blue-800">
                Search: "{filters.q}"
                <button 
                  onClick={() => handleInputChange('q', '')}
                  className="ml-1 text-blue-600 hover:text-blue-800"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            )}
            {filters.cuisine && (
              <span className="inline-flex items-center px-2 py-1 rounded-md text-sm bg-blue-100 text-blue-800">
                Cuisine: {filters.cuisine}
                <button 
                  onClick={() => handleInputChange('cuisine', '')}
                  className="ml-1 text-blue-600 hover:text-blue-800"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            )}
            {filters.type && (
              <span className="inline-flex items-center px-2 py-1 rounded-md text-sm bg-blue-100 text-blue-800">
                Type: {filters.type}
                <button 
                  onClick={() => handleInputChange('type', '')}
                  className="ml-1 text-blue-600 hover:text-blue-800"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            )}
            {filters.maxCalories && (
              <span className="inline-flex items-center px-2 py-1 rounded-md text-sm bg-blue-100 text-blue-800">
                Max {filters.maxCalories} cal
                <button 
                  onClick={() => handleInputChange('maxCalories', '')}
                  className="ml-1 text-blue-600 hover:text-blue-800"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            )}
            {filters.dietaryTags.map(tag => (
              <span key={tag} className="inline-flex items-center px-2 py-1 rounded-md text-sm bg-green-100 text-green-800">
                {tag}
                <button 
                  onClick={() => handleDietaryTagToggle(tag)}
                  className="ml-1 text-green-600 hover:text-green-800"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default AdvancedDishFilter;
