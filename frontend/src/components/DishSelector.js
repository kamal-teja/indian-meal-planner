import React, { useState } from 'react';
import { X, Search, Filter, Flame, Clock, Users, Plus } from 'lucide-react';
import AddDishForm from './AddDishForm';

const DishSelector = ({ dishes, onSelect, onClose, mealType, isEditing, onAddDish }) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedType, setSelectedType] = useState('all');
  const [selectedCuisine, setSelectedCuisine] = useState('all');
  const [showAddDishForm, setShowAddDishForm] = useState(false);

  const getMealTypeIcon = (mealType) => {
    const icons = {
      breakfast: '🌅',
      lunch: '☀️',
      dinner: '🌙',
      snack: '🍿'
    };
    return icons[mealType] || '🍽️';
  };

  const getCuisines = () => {
    const cuisines = [...new Set(dishes.map(dish => dish.cuisine))];
    return cuisines.sort();
  };

  const getTypeColor = (type) => {
    return type === 'Veg' 
      ? 'bg-green-100 text-green-800 border border-green-200' 
      : 'bg-red-100 text-red-800 border border-red-200';
  };

  const getCuisineColor = (cuisine) => {
    const colors = {
      'North Indian': 'bg-orange-100 text-orange-800',
      'South Indian': 'bg-green-100 text-green-800',
      'Punjabi': 'bg-yellow-100 text-yellow-800',
      'Hyderabadi': 'bg-purple-100 text-purple-800',
      'Bengali': 'bg-blue-100 text-blue-800',
      'Gujarati': 'bg-pink-100 text-pink-800',
    };
    return colors[cuisine] || 'bg-gray-100 text-gray-800';
  };

  const handleAddDish = async (dishData) => {
    try {
      const newDish = await onAddDish(dishData);
      setShowAddDishForm(false);
      // Optionally select the newly added dish
      if (newDish) {
        onSelect(newDish);
      }
    } catch (error) {
      console.error('Error adding dish:', error);
      throw error; // Re-throw to let AddDishForm handle the error display
    }
  };

  const filteredDishes = dishes.filter(dish => {
    const matchesSearch = dish.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         dish.cuisine.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         dish.ingredients.some(ing => ing.toLowerCase().includes(searchTerm.toLowerCase()));
    
    const matchesType = selectedType === 'all' || dish.type === selectedType;
    const matchesCuisine = selectedCuisine === 'all' || dish.cuisine === selectedCuisine;
    
    return matchesSearch && matchesType && matchesCuisine;
  });

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="bg-gradient-to-r from-primary-500 to-primary-600 p-6 text-white">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <span className="text-3xl">{getMealTypeIcon(mealType)}</span>
              <div>
                <h2 className="text-2xl font-display font-bold">
                  {isEditing ? 'Change' : 'Select'} {mealType.charAt(0).toUpperCase() + mealType.slice(1)}
                </h2>
                <p className="text-primary-100">Choose from our delicious Indian dishes</p>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <button
                onClick={() => setShowAddDishForm(true)}
                className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30 text-white rounded-xl transition-colors"
              >
                <Plus className="w-4 h-4" />
                <span>Add New Dish</span>
              </button>
              <button
                onClick={onClose}
                className="p-2 hover:bg-white/20 rounded-lg transition-colors"
              >
                <X className="h-6 w-6" />
              </button>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className="p-6 border-b border-gray-200 bg-gray-50">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              <input
                type="text"
                placeholder="Search dishes..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
            </div>

            {/* Type Filter */}
            <select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            >
              <option value="all">All Types</option>
              <option value="Veg">Vegetarian</option>
              <option value="Non-Veg">Non-Vegetarian</option>
            </select>

            {/* Cuisine Filter */}
            <select
              value={selectedCuisine}
              onChange={(e) => setSelectedCuisine(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            >
              <option value="all">All Cuisines</option>
              {getCuisines().map(cuisine => (
                <option key={cuisine} value={cuisine}>{cuisine}</option>
              ))}
            </select>
          </div>

          <div className="flex items-center justify-between mt-4">
            <span className="text-sm text-gray-600">
              {filteredDishes.length} dish{filteredDishes.length !== 1 ? 'es' : ''} found
            </span>
            <button
              onClick={() => {
                setSearchTerm('');
                setSelectedType('all');
                setSelectedCuisine('all');
              }}
              className="text-sm text-primary-600 hover:text-primary-700 font-medium"
            >
              Clear filters
            </button>
          </div>
        </div>

        {/* Dishes Grid */}
        <div className="p-6 overflow-y-auto max-h-96">
          {filteredDishes.length === 0 ? (
            <div className="text-center py-12">
              <Filter className="h-16 w-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No dishes found</h3>
              <p className="text-gray-600">Try adjusting your search or filters</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {filteredDishes.map(dish => (
                <div
                  key={dish.id}
                  onClick={() => onSelect(dish)}
                  className="p-4 border border-gray-200 rounded-xl hover:border-primary-300 hover:shadow-md cursor-pointer transition-all duration-200 group"
                >
                  <div className="flex items-start space-x-4">
                    <div className="relative flex-shrink-0">
                      <img
                        src={dish.image}
                        alt={dish.name}
                        className="w-20 h-20 object-cover rounded-lg group-hover:scale-105 transition-transform"
                        onError={(e) => {
                          e.target.src = `https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=300&h=200&fit=crop`;
                        }}
                      />
                      <div className="absolute -top-2 -right-2">
                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${getTypeColor(dish.type)}`}>
                          {dish.type}
                        </span>
                      </div>
                    </div>

                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-gray-900 group-hover:text-primary-600 transition-colors">
                        {dish.name}
                      </h3>
                      <span className={`inline-block px-2 py-1 text-xs font-medium rounded-full mt-1 ${getCuisineColor(dish.cuisine)}`}>
                        {dish.cuisine}
                      </span>

                      <div className="flex items-center space-x-3 mt-2 text-sm text-gray-600">
                        <div className="flex items-center space-x-1">
                          <Flame className="h-4 w-4 text-orange-500" />
                          <span className="font-medium">{dish.calories} cal</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <Clock className="h-4 w-4 text-blue-500" />
                          <span>30 min</span>
                        </div>
                      </div>

                      {dish.ingredients && dish.ingredients.length > 0 && (
                        <div className="mt-2">
                          <div className="flex flex-wrap gap-1">
                            {dish.ingredients.slice(0, 2).map((ingredient, index) => (
                              <span
                                key={index}
                                className="px-1.5 py-0.5 bg-gray-100 text-gray-600 text-xs rounded"
                              >
                                {ingredient}
                              </span>
                            ))}
                            {dish.ingredients.length > 2 && (
                              <span className="px-1.5 py-0.5 bg-gray-100 text-gray-500 text-xs rounded">
                                +{dish.ingredients.length - 2}
                              </span>
                            )}
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
      
      {/* Add Dish Form Modal */}
      {showAddDishForm && (
        <AddDishForm
          onSubmit={handleAddDish}
          onClose={() => setShowAddDishForm(false)}
        />
      )}
    </div>
  );
};

export default DishSelector;
