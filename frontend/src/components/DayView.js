import React, { useState, useEffect } from 'react';
import { format, addDays, subDays } from 'date-fns';
import { ChevronLeft, ChevronRight, Plus, X, Edit3, Trash2, ShoppingCart } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import MealCard from './MealCard';
import DishSelector from './DishSelector';
import IngredientsList from './IngredientsList';
import MealRecommendations from './MealRecommendations';

const MEAL_TYPES = [
  { id: 'breakfast', name: 'Breakfast', icon: '🌅', color: 'from-orange-400 to-yellow-400' },
  { id: 'lunch', name: 'Lunch', icon: '☀️', color: 'from-green-400 to-blue-400' },
  { id: 'dinner', name: 'Dinner', icon: '🌙', color: 'from-purple-400 to-pink-400' },
  { id: 'snack', name: 'Snacks', icon: '🍿', color: 'from-pink-400 to-red-400' }
];

const DayView = ({ dishes, onAddDish }) => {
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [meals, setMeals] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showDishSelector, setShowDishSelector] = useState(false);
  const [selectedMealType, setSelectedMealType] = useState(null);
  const [editingMeal, setEditingMeal] = useState(null);
  const [showIngredientsPanel, setShowIngredientsPanel] = useState(false);
  const [showRecommendations, setShowRecommendations] = useState(false);

  useEffect(() => {
    loadMealsForDate();
  }, [selectedDate]);

  const loadMealsForDate = async () => {
    try {
      setLoading(true);
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      const response = await mealPlannerAPI.getMealsByDate(dateStr);
      setMeals(response.data);
    } catch (error) {
      console.error('Error loading meals:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddMeal = (mealType) => {
    setSelectedMealType(mealType);
    setEditingMeal(null);
    setShowDishSelector(true);
  };

  const handleEditMeal = (meal) => {
    setSelectedMealType(meal.mealType);
    setEditingMeal(meal);
    setShowDishSelector(true);
  };

  const handleDishSelect = async (dish) => {
    try {
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      
      if (editingMeal) {
        // Update existing meal
        await mealPlannerAPI.updateMeal(editingMeal.id, dish.id);
      } else {
        // Add new meal
        await mealPlannerAPI.addMeal({
          date: dateStr,
          mealType: selectedMealType,
          dishId: dish.id
        });
      }
      
      await loadMealsForDate();
      setShowDishSelector(false);
      setEditingMeal(null);
    } catch (error) {
      console.error('Error saving meal:', error);
    }
  };

  const handleDeleteMeal = async (mealId) => {
    if (window.confirm('Are you sure you want to delete this meal?')) {
      try {
        await mealPlannerAPI.deleteMeal(mealId);
        await loadMealsForDate();
      } catch (error) {
        console.error('Error deleting meal:', error);
      }
    }
  };

  const navigateDate = (direction) => {
    if (direction === 'prev') {
      setSelectedDate(subDays(selectedDate, 1));
    } else {
      setSelectedDate(addDays(selectedDate, 1));
    }
  };

  const getMealsForType = (mealType) => {
    return meals.filter(meal => meal.mealType === mealType);
  };

  const getTotalCalories = () => {
    return meals.reduce((total, meal) => total + (meal.dish.calories || 0), 0);
  };

  return (
    <div className="max-w-7xl mx-auto">
      <div className={`grid transition-all duration-300 gap-6 ${showIngredientsPanel ? 'grid-cols-1 lg:grid-cols-3' : 'grid-cols-1'}`}>
        {/* Main Content */}
        <div className={`space-y-6 ${showIngredientsPanel ? 'lg:col-span-2' : ''}`}>
          {/* Date Navigation */}
          <div className="card p-6">
            <div className="flex items-center justify-between">
              <button
                onClick={() => navigateDate('prev')}
                className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
              >
                <ChevronLeft className="h-6 w-6 text-gray-600" />
              </button>
              
              <div className="text-center">
                <h2 className="text-3xl font-display font-bold gradient-text">
                  {format(selectedDate, 'EEEE')}
                </h2>
                <p className="text-lg text-gray-600">
                  {format(selectedDate, 'MMMM d, yyyy')}
                </p>
                <div className="mt-2 px-4 py-2 bg-gradient-to-r from-primary-100 to-secondary-100 rounded-full inline-block">
                  <span className="text-sm font-medium text-gray-700">
                    Total Calories: <span className="font-bold text-primary-600">{getTotalCalories()}</span>
                  </span>
                </div>
              </div>
              
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setShowRecommendations(!showRecommendations)}
                  className={`p-2 rounded-lg transition-colors flex items-center space-x-2 ${
                    showRecommendations 
                      ? 'bg-purple-100 text-purple-600 hover:bg-purple-200' 
                      : 'hover:bg-gray-100 text-gray-600'
                  }`}
                  title="Toggle Meal Recommendations"
                >
                  ✨
                  <span className="hidden sm:inline text-sm">Recommendations</span>
                </button>
                <button
                  onClick={() => setShowIngredientsPanel(!showIngredientsPanel)}
                  className={`p-2 rounded-lg transition-colors flex items-center space-x-2 ${
                    showIngredientsPanel 
                      ? 'bg-green-100 text-green-600 hover:bg-green-200' 
                      : 'hover:bg-gray-100 text-gray-600'
                  }`}
                  title="Toggle Shopping List"
                >
                  <ShoppingCart className="h-5 w-5" />
                  {meals.length > 0 && (
                    <span className="bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                      {meals.reduce((count, meal) => count + (meal.dish?.ingredients?.length || 0), 0)}
                    </span>
                  )}
                </button>
                <button
                  onClick={() => navigateDate('next')}
                  className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <ChevronRight className="h-6 w-6 text-gray-600" />
                </button>
              </div>
            </div>
          </div>

          {/* Meal Sections */}
          {loading ? (
            <div className="flex justify-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-4 border-primary-200 border-t-primary-600"></div>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {MEAL_TYPES.map((mealType) => {
                const mealTypemeals = getMealsForType(mealType.id);
                
                return (
                  <div key={mealType.id} className="card p-6">
                    <div className="flex items-center justify-between mb-4">
                      <div className="flex items-center space-x-3">
                        <div className={`p-3 bg-gradient-to-r ${mealType.color} rounded-xl shadow-lg`}>
                          <span className="text-2xl">{mealType.icon}</span>
                        </div>
                        <div>
                          <h3 className="text-xl font-display font-semibold text-gray-800">
                            {mealType.name}
                          </h3>
                          <p className="text-sm text-gray-600">
                            {mealTypemeals.length} meal{mealTypemeals.length !== 1 ? 's' : ''}
                          </p>
                        </div>
                      </div>
                      
                      <button
                        onClick={() => handleAddMeal(mealType.id)}
                        className="p-2 bg-primary-100 hover:bg-primary-200 text-primary-600 rounded-lg transition-colors"
                      >
                        <Plus className="h-5 w-5" />
                      </button>
                    </div>

                    <div className="space-y-3">
                      {mealTypemeals.length === 0 ? (
                        <div className="text-center py-8 text-gray-500">
                          <p>No meals planned</p>
                          <button
                            onClick={() => handleAddMeal(mealType.id)}
                            className="mt-2 text-primary-600 hover:text-primary-700 font-medium"
                          >
                            Add your first meal
                          </button>
                        </div>
                      ) : (
                        mealTypemeals.map((meal) => (
                          <div key={meal.id} className="relative group">
                            <MealCard meal={meal} />
                            <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                              <div className="flex space-x-1">
                                <button
                                  onClick={() => handleEditMeal(meal)}
                                  className="p-1.5 bg-blue-500 hover:bg-blue-600 text-white rounded-lg shadow-md transition-colors"
                                >
                                  <Edit3 className="h-3 w-3" />
                                </button>
                                <button
                                  onClick={() => handleDeleteMeal(meal.id)}
                                  className="p-1.5 bg-red-500 hover:bg-red-600 text-white rounded-lg shadow-md transition-colors"
                                >
                                  <Trash2 className="h-3 w-3" />
                                </button>
                              </div>
                            </div>
                          </div>
                        ))
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {/* Recommendations Panel */}
        {showRecommendations && (
          <div className="lg:col-span-1">
            <div className="sticky top-6">
              <MealRecommendations
                onAddMeal={addMeal}
                selectedDate={format(selectedDate, 'yyyy-MM-dd')}
                mealType="lunch"
              />
            </div>
          </div>
        )}

        {/* Ingredients Panel */}
        {showIngredientsPanel && (
          <div className="lg:col-span-1">
            <div className="sticky top-6">
              <IngredientsList
                meals={meals}
                selectedDate={selectedDate}
                onClose={() => setShowIngredientsPanel(false)}
              />
            </div>
          </div>
        )}
      </div>

      {/* Dish Selector Modal */}
      {showDishSelector && (
        <DishSelector
          dishes={dishes}
          onSelect={handleDishSelect}
          onClose={() => setShowDishSelector(false)}
          mealType={selectedMealType}
          isEditing={!!editingMeal}
          onAddDish={onAddDish}
        />
      )}
    </div>
  );
};

export default DayView;
