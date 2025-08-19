import React from 'react';
import { Clock, Users, Flame } from 'lucide-react';

const MealCard = ({ meal }) => {
  const { dish } = meal;

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

  const getTypeColor = (type) => {
    return type === 'Veg' 
      ? 'bg-green-100 text-green-800 border border-green-200' 
      : 'bg-red-100 text-red-800 border border-red-200';
  };

  return (
    <div className="meal-card group">
      <div className="flex items-start space-x-4">
        {/* Dish Image */}
        <div className="flex-shrink-0">
          <div className="relative">
            <img
              src={dish.image}
              alt={dish.name}
              className="w-20 h-20 object-cover rounded-xl shadow-md group-hover:shadow-lg transition-shadow"
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
        </div>

        {/* Dish Details */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between">
            <div>
              <h4 className="text-lg font-semibold text-gray-900 group-hover:text-primary-600 transition-colors">
                {dish.name}
              </h4>
              <span className={`inline-block px-2 py-1 text-xs font-medium rounded-full mt-1 ${getCuisineColor(dish.cuisine)}`}>
                {dish.cuisine}
              </span>
            </div>
          </div>

          {/* Calories and Info */}
          <div className="flex items-center space-x-4 mt-3 text-sm text-gray-600">
            <div className="flex items-center space-x-1">
              <Flame className="h-4 w-4 text-orange-500" />
              <span className="font-medium">{dish.calories} cal</span>
            </div>
            <div className="flex items-center space-x-1">
              <Clock className="h-4 w-4 text-blue-500" />
              <span>30 min</span>
            </div>
            <div className="flex items-center space-x-1">
              <Users className="h-4 w-4 text-green-500" />
              <span>2-3 servings</span>
            </div>
          </div>

          {/* Ingredients Preview */}
          {dish.ingredients && dish.ingredients.length > 0 && (
            <div className="mt-3">
              <p className="text-xs text-gray-500 mb-1">Main ingredients:</p>
              <div className="flex flex-wrap gap-1">
                {dish.ingredients.slice(0, 3).map((ingredient, index) => (
                  <span
                    key={index}
                    className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded-md"
                  >
                    {ingredient}
                  </span>
                ))}
                {dish.ingredients.length > 3 && (
                  <span className="px-2 py-1 bg-gray-100 text-gray-500 text-xs rounded-md">
                    +{dish.ingredients.length - 3} more
                  </span>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MealCard;
