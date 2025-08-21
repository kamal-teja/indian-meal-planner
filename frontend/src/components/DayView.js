import React, { useState, useEffect } from 'react';
import { format, addDays, subDays } from 'date-fns';
import { ChevronLeft, ChevronRight, Plus, X, Edit3, Trash2, ShoppingCart } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import MealCard from './MealCard';
import DishSelector from './DishSelector';
import IngredientsList from './IngredientsList';
import MealRecommendations from './MealRecommendations';

const MEAL_TYPES = [
  { id: 'breakfast', name: 'Breakfast', icon: '🌅', color: 'bg-warm-100 border-warm-200', bgColor: 'bg-warm-50', textColor: 'text-warm-700' },
  { id: 'lunch', name: 'Lunch', icon: '☀️', color: 'bg-sage-100 border-sage-200', bgColor: 'bg-sage-50', textColor: 'text-sage-700' },
  { id: 'dinner', name: 'Dinner', icon: '🌙', color: 'bg-lavender-100 border-lavender-200', bgColor: 'bg-lavender-50', textColor: 'text-lavender-700' },
  { id: 'snack', name: 'Snacks', icon: '🍿', color: 'bg-secondary-100 border-secondary-200', bgColor: 'bg-secondary-50', textColor: 'text-secondary-700' }
];

const DayView = ({ loadDishes, onAddDish }) => {
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
      // Backend returns { success: true, data: [meals] }, so we need response.data.data
      setMeals(response.data.data || []);
    } catch (error) {
      console.error('Error loading meals:', error);
      setMeals([]); // Set empty array on error to prevent reduce error
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

  const addMeal = async (mealData) => {
    try {
      await mealPlannerAPI.addMeal(mealData);
      loadMealsForDate(); // Refresh meals after adding
    } catch (error) {
      console.error('Error adding meal:', error);
    }
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
    // Safety check: ensure meals is an array before calling filter
    if (!Array.isArray(meals)) {
      console.warn('meals is not an array:', meals);
      return [];
    }
    return meals.filter(meal => meal.mealType === mealType);
  };

  const getTotalCalories = () => {
    // Safety check: ensure meals is an array before calling reduce
    if (!Array.isArray(meals)) {
      console.warn('meals is not an array:', meals);
      return 0;
    }
    return meals.reduce((total, meal) => total + (meal.dish?.calories || 0), 0);
  };

  return (
    <div className="max-w-7xl mx-auto p-6">
      <div className={`grid transition-all duration-300 gap-6 ${showIngredientsPanel ? 'grid-cols-1 lg:grid-cols-3' : 'grid-cols-1'}`}>
        {/* Main Content */}
        <div className={`space-y-6 ${showIngredientsPanel ? 'lg:col-span-2' : ''}`}>
          {/* Date Navigation */}
          <div className="card-elevated p-6">
            <div className="flex items-center justify-between">
              <button
                onClick={() => navigateDate('prev')}
                className="p-3 rounded-lg hover:bg-accent-50 transition-all duration-200 hover:shadow-md"
              >
                <ChevronLeft className="h-6 w-6 text-accent-600" />
              </button>
              
              <div className="text-center">
                <h2 className="text-3xl font-display font-bold text-neutral-800">
                  {format(selectedDate, 'EEEE')}
                </h2>
                <p className="text-lg text-accent-600">
                  {format(selectedDate, 'MMMM d, yyyy')}
                </p>
                <div className="mt-3 px-6 py-2 bg-accent-100 rounded-full inline-block border border-accent-200">
                  <span className="text-sm font-medium text-accent-700">
                    {getTotalCalories()} calories today
                  </span>
                </div>
              </div>
              
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setShowRecommendations(!showRecommendations)}
                  className={`p-3 rounded-lg transition-all duration-200 flex items-center space-x-2 ${
                    showRecommendations 
                      ? 'bg-lavender-100 text-lavender-700 hover:bg-lavender-200 shadow-md' 
                      : 'hover:bg-accent-50 text-accent-600 hover:shadow-md'
                  }`}
                  title="Toggle Meal Recommendations"
                >
                  ✨
                  <span className="hidden sm:inline text-sm font-medium">Recommendations</span>
                </button>
                <button
                  onClick={() => setShowIngredientsPanel(!showIngredientsPanel)}
                  className={`p-3 rounded-lg transition-all duration-200 flex items-center space-x-2 ${
                    showIngredientsPanel 
                      ? 'bg-sage-100 text-sage-700 hover:bg-sage-200 shadow-md' 
                      : 'hover:bg-accent-50 text-accent-600 hover:shadow-md'
                  }`}
                  title="Toggle Shopping List"
                >
                  <ShoppingCart className="h-5 w-5" />
                  {meals.length > 0 && (
                    <span className="bg-accent-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center shadow-sm">
                      {meals.reduce((count, meal) => count + (meal.dish?.ingredients?.length || 0), 0)}
                    </span>
                  )}
                </button>
                <button
                  onClick={() => navigateDate('next')}
                  className="p-3 rounded-lg hover:bg-accent-50 transition-all duration-200 hover:shadow-md"
                >
                  <ChevronRight className="h-6 w-6 text-accent-600" />
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
                  <div key={mealType.id} className="card-elevated p-6">
                    <div className="flex items-center justify-between mb-6">
                      <div className="flex items-center space-x-3">
                        <div className={`p-3 ${mealType.color} rounded-lg border`}>
                          <span className="text-2xl">{mealType.icon}</span>
                        </div>
                        <div>
                          <h3 className={`text-xl font-display font-semibold ${mealType.textColor}`}>
                            {mealType.name}
                          </h3>
                          <p className="text-sm text-accent-600">
                            {mealTypemeals.length} meal{mealTypemeals.length !== 1 ? 's' : ''}
                          </p>
                        </div>
                      </div>
                      
                      <button
                        onClick={() => handleAddMeal(mealType.id)}
                        className={`p-3 ${mealType.bgColor} hover:shadow-md text-accent-600 rounded-lg transition-all duration-200 hover:scale-105`}
                      >
                        <Plus className="h-5 w-5" />
                      </button>
                    </div>

                    <div className="space-y-4">
                      {mealTypemeals.length === 0 ? (
                        <div className="text-center py-8 text-accent-500">
                          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-accent-100 flex items-center justify-center">
                            <span className="text-2xl">{mealType.icon}</span>
                          </div>
                          <p className="text-sm mb-3">No meals planned</p>
                          <button
                            onClick={() => handleAddMeal(mealType.id)}
                            className="btn-secondary text-sm"
                          >
                            Add your first meal
                          </button>
                        </div>
                      ) : (
                        mealTypemeals.map((meal) => (
                          <div key={meal.id} className="relative group">
                            <MealCard meal={meal} />
                            <div className="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity">
                              <div className="flex space-x-2">
                                <button
                                  onClick={() => handleEditMeal(meal)}
                                  className="p-2 bg-secondary-500 hover:bg-secondary-600 text-white rounded-lg shadow-md transition-colors"
                                >
                                  <Edit3 className="h-4 w-4" />
                                </button>
                                <button
                                  onClick={() => handleDeleteMeal(meal.id)}
                                  className="p-2 bg-rose-500 hover:bg-rose-600 text-white rounded-lg shadow-md transition-colors"
                                >
                                  <Trash2 className="h-4 w-4" />
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
          loadDishes={loadDishes}
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
