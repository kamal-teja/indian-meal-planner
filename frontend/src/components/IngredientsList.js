import React, { useState, useMemo } from 'react';
import { ShoppingCart, Check, X, ChefHat, Utensils, Sun, Sunset, Moon, Coffee } from 'lucide-react';

const IngredientsList = ({ meals, selectedDate, onClose }) => {
  const [checkedIngredients, setCheckedIngredients] = useState(new Set());

  // Extract and consolidate ingredients from all meals
  const consolidatedIngredients = useMemo(() => {
    const ingredientMap = new Map();
    
    meals.forEach(meal => {
      if (meal.dish && meal.dish.ingredients) {
        meal.dish.ingredients.forEach(ingredient => {
          const normalizedIngredient = ingredient.toLowerCase().trim();
          if (ingredientMap.has(normalizedIngredient)) {
            const existing = ingredientMap.get(normalizedIngredient);
            existing.meals.push({
              mealType: meal.mealType,
              dishName: meal.dish.name
            });
          } else {
            ingredientMap.set(normalizedIngredient, {
              name: ingredient,
              meals: [{
                mealType: meal.mealType,
                dishName: meal.dish.name
              }]
            });
          }
        });
      }
    });

    return Array.from(ingredientMap.values()).sort((a, b) => 
      a.name.localeCompare(b.name)
    );
  }, [meals]);

  const toggleIngredient = (ingredient) => {
    const newChecked = new Set(checkedIngredients);
    if (newChecked.has(ingredient)) {
      newChecked.delete(ingredient);
    } else {
      newChecked.add(ingredient);
    }
    setCheckedIngredients(newChecked);
  };

  const getMealTypeIcon = (mealType) => {
    const iconMap = {
      breakfast: <Coffee className="h-4 w-4 text-accent-600" />,
      lunch: <Sun className="h-4 w-4 text-accent-600" />,
      dinner: <Moon className="h-4 w-4 text-accent-600" />,
      snack: <Sunset className="h-4 w-4 text-accent-600" />
    };
    return iconMap[mealType] || <Coffee className="h-4 w-4 text-accent-600" />;
  };

  const getMealTypeColor = (mealType) => {
    const colors = {
      breakfast: 'bg-warm-100 text-warm-700 border border-warm-200',
      lunch: 'bg-secondary-100 text-secondary-700 border border-secondary-200',
      dinner: 'bg-lavender-100 text-lavender-700 border border-lavender-200',
      snack: 'bg-sage-100 text-sage-700 border border-sage-200'
    };
    return colors[mealType] || 'bg-neutral-100 text-accent-700 border border-neutral-200';
  };

  const formatDate = (date) => {
    return new Intl.DateTimeFormat('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    }).format(date);
  };

  const checkedCount = checkedIngredients.size;
  const totalCount = consolidatedIngredients.length;

  return (
  <div className="bg-white rounded-2xl shadow-sm border border-gray-200 h-full flex flex-col">
    {/* Header (softer theme) */}
  <div className="bg-primary-50 p-6 rounded-t-2xl text-accent-900">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
              <ShoppingCart className="w-6 h-6 text-white" />
            <div>
                <h2 className="text-xl font-bold">Groceries</h2>
                <p className="text-white/90 text-sm">
                  {formatDate(selectedDate)}
                </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-white/20 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>
        
        {/* Progress */}
        {totalCount > 0 && (
          <div className="mt-4">
            <div className="flex items-center justify-between text-sm mb-2">
              <span>Progress</span>
              <span>{checkedCount}/{totalCount} items</span>
            </div>
            <div className="w-full bg-white/20 rounded-full h-2">
              <div 
                className="bg-white rounded-full h-2 transition-all duration-300"
                style={{ width: `${(checkedCount / totalCount) * 100}%` }}
              />
            </div>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        {consolidatedIngredients.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full p-8 text-gray-500">
            <ChefHat className="w-16 h-16 mb-4 text-gray-300" />
            <h3 className="text-lg font-medium mb-2">No meals planned</h3>
            <p className="text-center text-sm">
              Add some meals to see your shopping list for this day
            </p>
          </div>
        ) : (
          <div className="p-6">
            <div className="mb-4">
              <h3 className="text-lg font-semibold text-gray-800 mb-2">
                Ingredients Needed
              </h3>
              <p className="text-sm text-gray-600">
                Tap items to mark them as purchased
              </p>
            </div>

            <div className="space-y-3">
              {consolidatedIngredients.map((item, index) => (
                <div
                  key={index}
                  className={`p-4 rounded-xl border-2 transition-all duration-200 cursor-pointer ${
                    checkedIngredients.has(item.name)
                      ? 'bg-green-50 border-green-200 opacity-60'
                      : 'bg-white border-gray-200 hover:border-green-300 hover:shadow-sm'
                  }`}
                  onClick={() => toggleIngredient(item.name)}
                >
                  <div className="flex items-start space-x-3">
                    <div className={`flex-shrink-0 w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all ${
                      checkedIngredients.has(item.name)
                        ? 'bg-green-500 border-green-500'
                        : 'border-gray-300 hover:border-green-400'
                    }`}>
                      {checkedIngredients.has(item.name) && (
                        <Check className="w-4 h-4 text-white" />
                      )}
                    </div>
                    
                    <div className="flex-1 min-w-0 overflow-hidden">
                      <h4 className={`font-medium truncate ${
                        checkedIngredients.has(item.name)
                          ? 'line-through text-gray-500'
                          : 'text-gray-800'
                      }`} title={item.name}>
                        {item.name}
                      </h4>

                      <div className="mt-2 flex flex-wrap gap-1">
                        {item.meals.map((meal, mealIndex) => (
                          <span
                            key={mealIndex}
                            className={`inline-flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium max-w-[160px] truncate overflow-hidden ${getMealTypeColor(meal.mealType)}`}
                          >
                            <span className="flex-shrink-0">{getMealTypeIcon(meal.mealType)}</span>
                            <span className="truncate">{meal.dishName}</span>
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Footer */}
      {totalCount > 0 && (
        <div className="p-4 border-t border-gray-200 bg-gray-50 rounded-b-2xl">
          <div className="flex items-center justify-between text-sm">
            <div className="flex items-center space-x-2 text-gray-600">
              <Utensils className="w-4 h-4" />
              <span>{meals.length} meal{meals.length !== 1 ? 's' : ''} planned</span>
            </div>
            <div className="text-gray-600">
              {checkedCount > 0 && (
                <span className="text-green-600 font-medium">
                  {checkedCount} purchased
                </span>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default IngredientsList;
