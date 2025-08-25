import React from 'react';
import { Clock, Users, Flame } from 'lucide-react';

const MealCard = ({ meal, count = 1, animateDecrease = false }) => {
  // Handle null or undefined meal
  if (!meal) {
    return null;
  }

  const { dish } = meal;

  const getCuisineColor = (cuisine) => {
    const colors = {
      'North Indian': 'bg-warm-100 text-warm-800 border border-warm-200',
      'South Indian': 'bg-sage-100 text-sage-800 border border-sage-200',
      'Punjabi': 'bg-secondary-100 text-secondary-800 border border-secondary-200',
      'Hyderabadi': 'bg-lavender-100 text-lavender-800 border border-lavender-200',
      'Bengali': 'bg-accent-100 text-accent-800 border border-accent-200',
      'Gujarati': 'bg-neutral-100 text-neutral-800 border border-neutral-200',
    };
    return colors[cuisine] || 'bg-neutral-100 text-neutral-800 border border-neutral-200';
  };

  const getTypeColor = (type) => {
    return type === 'Veg' 
      ? 'badge-veg' 
      : 'badge-non-veg';
  };

  return (
    <div className={`meal-card group relative transition-all ${animateDecrease ? 'animate-count-decrease' : ''} overflow-hidden`}> 
      {count > 1 && (
        <div className="absolute -top-2 -right-2 bg-warm-600 text-white text-xs rounded-full h-7 w-7 flex items-center justify-center shadow-md font-semibold transform translate-z-0">
          <span className="text-xs">×{count}</span>
        </div>
      )}
      <div className="flex items-start space-x-4">
        {/* Dish Image */}
        <div className="flex-shrink-0">
          <div className="relative">
            <img
              src={dish.image}
              alt={dish.name}
              className="w-20 h-20 object-cover rounded-lg shadow-sm group-hover:shadow-md transition-shadow border border-neutral-200"
              onError={(e) => {
                e.target.src = `https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=300&h=200&fit=crop`;
              }}
            />
            <div className="absolute -top-2 -right-2">
              <span className={`px-2 py-1 text-xs font-medium rounded-md ${getTypeColor(dish.type)} shadow-sm`}>
                {dish.type}
              </span>
            </div>
          </div>
        </div>

        {/* Dish Details */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between min-w-0">
            <div className="min-w-0">
              <h4 className="text-lg font-semibold text-accent-900 group-hover:text-secondary-600 transition-colors">
                {dish.name}
              </h4>
              <span className={`inline-block px-3 py-1 text-xs font-medium rounded-md mt-2 ${getCuisineColor(dish.cuisine)} shadow-sm`}>
                {dish.cuisine}
              </span>
            </div>
          </div>

          {/* Calories and Info */}
          <div className="flex items-center space-x-3 mt-4 text-sm text-accent-600 flex-wrap">
            <div className="flex items-center space-x-1 surface-warm px-2 py-1 rounded-md min-w-0">
              <Flame className="h-4 w-4 text-warm-600 flex-shrink-0" />
              <span className="font-medium truncate">{dish.calories} cal</span>
            </div>
            <div className="flex items-center space-x-1 surface-secondary px-2 py-1 rounded-md min-w-0">
              <Clock className="h-4 w-4 text-secondary-600 flex-shrink-0" />
              <span className="truncate">30 min</span>
            </div>
            <div className="flex items-center space-x-1 surface-primary px-2 py-1 rounded-md min-w-0">
              <Users className="h-4 w-4 text-secondary-600 flex-shrink-0" />
              <span className="truncate">2-3 servings</span>
            </div>
          </div>

          {/* Ingredients Preview */}
          {dish.ingredients && dish.ingredients.length > 0 && (
            <div className="mt-4">
              <p className="text-xs text-accent-500 mb-2 font-medium">Main ingredients:</p>
              <div className="flex flex-wrap gap-2 content-start items-start">
                {dish.ingredients.slice(0, 3).map((ingredient, index) => (
                  <span
                    key={index}
                    className="inline-flex items-center px-3 py-1 bg-neutral-100 text-neutral-700 text-xs rounded-md border border-neutral-200 max-w-[220px] truncate"
                    title={ingredient}
                  >
                    <span className="truncate">{ingredient}</span>
                  </span>
                ))}
                {dish.ingredients.length > 3 && (
                  <span className="inline-flex items-center px-3 py-1 bg-lavender-100 text-lavender-700 text-xs rounded-md border border-lavender-200 max-w-[220px] truncate">
                    <span className="truncate">+{dish.ingredients.length - 3} more</span>
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
