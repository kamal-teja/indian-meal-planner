import React, { useState, useEffect } from 'react';
import { format, startOfMonth, endOfMonth, eachDayOfInterval, isSameMonth, isToday, addMonths, subMonths } from 'date-fns';
import { ChevronLeft, ChevronRight, Plus, Calendar } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import DishSelector from './DishSelector';

const WEEKDAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

const MonthView = ({ loadDishes, onAddDish }) => {
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [monthMeals, setMonthMeals] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedDay, setSelectedDay] = useState(null);
  const [showDishSelector, setShowDishSelector] = useState(false);
  const [selectedMealType, setSelectedMealType] = useState('lunch');

  useEffect(() => {
    loadMonthMeals();
  }, [selectedDate]);

  const loadMonthMeals = async () => {
    try {
      setLoading(true);
      const year = selectedDate.getFullYear();
      const month = selectedDate.getMonth() + 1;
      const response = await mealPlannerAPI.getMealsByMonth(year, month);
      // Backend returns { success: true, data: [meals] }, so we need response.data.data
      const meals = response.data.data || [];
      console.log(`Loaded ${meals.length} meals for ${year}-${month}:`, meals);
      setMonthMeals(meals);
    } catch (error) {
      console.error('Error loading month meals:', error);
      setMonthMeals([]); // Set empty array on error
    } finally {
      setLoading(false);
    }
  };

  const navigateMonth = (direction) => {
    if (direction === 'prev') {
      setSelectedDate(subMonths(selectedDate, 1));
    } else {
      setSelectedDate(addMonths(selectedDate, 1));
    }
  };

  const getMealsForDate = (date) => {
    const dateStr = format(date, 'yyyy-MM-dd');
    const filteredMeals = monthMeals.filter(meal => {
      // Handle both date string formats: "2025-08-21" and "2025-08-21T10:30:00Z"
      const mealDate = meal.date instanceof Date ? meal.date : new Date(meal.date);
      const mealDateStr = format(mealDate, 'yyyy-MM-dd');
      return mealDateStr === dateStr;
    });
    
    // Debug logging (can be removed later)
    if (filteredMeals.length > 0) {
      console.log(`Found ${filteredMeals.length} meals for ${dateStr}:`, filteredMeals);
    }
    
    return filteredMeals;
  };

  const getCalendarDays = () => {
    const start = startOfMonth(selectedDate);
    const end = endOfMonth(selectedDate);
    
    // Get the start of the week for the first day of the month
    const calendarStart = new Date(start);
    calendarStart.setDate(calendarStart.getDate() - calendarStart.getDay());
    
    // Get the end of the week for the last day of the month
    const calendarEnd = new Date(end);
    calendarEnd.setDate(calendarEnd.getDate() + (6 - calendarEnd.getDay()));
    
    return eachDayOfInterval({ start: calendarStart, end: calendarEnd });
  };

  const handleDayClick = (date) => {
    if (isSameMonth(date, selectedDate)) {
      setSelectedDay(date);
      setShowDishSelector(true);
    }
  };

  const handleDishSelect = async (dish) => {
    try {
      const dateStr = format(selectedDay, 'yyyy-MM-dd');
      await mealPlannerAPI.addMeal({
        date: dateStr,
        mealType: selectedMealType,
        dishId: dish.id
      });
      
      await loadMonthMeals();
      setShowDishSelector(false);
      setSelectedDay(null);
    } catch (error) {
      console.error('Error adding meal:', error);
    }
  };

  const getTotalCaloriesForDate = (date) => {
    const meals = getMealsForDate(date);
    return meals.reduce((total, meal) => total + (meal.dish.calories || 0), 0);
  };

  const getMealTypeColors = () => ({
    breakfast: 'bg-orange-100 border-orange-300',
    lunch: 'bg-green-100 border-green-300',
    dinner: 'bg-purple-100 border-purple-300',
    snack: 'bg-pink-100 border-pink-300'
  });

  const calendarDays = getCalendarDays();

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Month Navigation */}
      <div className="card p-6">
        <div className="flex items-center justify-between">
          <button
            onClick={() => navigateMonth('prev')}
            className="p-3 rounded-xl hover:bg-gray-100 transition-colors"
          >
            <ChevronLeft className="h-6 w-6 text-gray-600" />
          </button>
          
          <div className="text-center">
            <h2 className="text-4xl font-display font-bold gradient-text">
              {format(selectedDate, 'MMMM yyyy')}
            </h2>
            <p className="text-gray-600 mt-1">
              Plan your monthly meals
            </p>
          </div>
          
          <button
            onClick={() => navigateMonth('next')}
            className="p-3 rounded-xl hover:bg-gray-100 transition-colors"
          >
            <ChevronRight className="h-6 w-6 text-gray-600" />
          </button>
        </div>
      </div>

      {/* Calendar Grid */}
      <div className="card p-6">
        {loading ? (
          <div className="flex justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-4 border-primary-200 border-t-primary-600"></div>
          </div>
        ) : (
          <>
            {/* Weekday Headers */}
            <div className="grid grid-cols-7 gap-2 mb-4">
              {WEEKDAYS.map(day => (
                <div key={day} className="text-center py-3">
                  <span className="text-sm font-semibold text-gray-600">{day}</span>
                </div>
              ))}
            </div>

            {/* Calendar Days */}
            <div className="grid grid-cols-7 gap-2">
              {calendarDays.map((date, index) => {
                const dayMeals = getMealsForDate(date);
                const isCurrentMonth = isSameMonth(date, selectedDate);
                const isCurrentDay = isToday(date);
                const totalCalories = getTotalCaloriesForDate(date);
                const mealColors = getMealTypeColors();

                return (
                  <div
                    key={index}
                    onClick={() => handleDayClick(date)}
                    className={`
                      min-h-[120px] p-2 rounded-xl border-2 transition-all duration-200 cursor-pointer
                      ${isCurrentMonth 
                        ? 'bg-white hover:bg-gray-50 border-gray-200 hover:border-primary-300 hover:shadow-md' 
                        : 'bg-gray-50 border-gray-100 text-gray-400'
                      }
                      ${isCurrentDay 
                        ? 'ring-2 ring-primary-500 border-primary-500 bg-primary-50' 
                        : ''
                      }
                    `}
                  >
                    {/* Date Number */}
                    <div className="flex items-center justify-between mb-2">
                      <span className={`
                        text-sm font-medium
                        ${isCurrentDay 
                          ? 'text-primary-700 font-bold' 
                          : isCurrentMonth 
                            ? 'text-gray-900' 
                            : 'text-gray-400'
                        }
                      `}>
                        {format(date, 'd')}
                      </span>
                      
                      {isCurrentMonth && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDayClick(date);
                          }}
                          className="p-1 rounded-md hover:bg-primary-100 text-primary-600 opacity-0 group-hover:opacity-100 transition-opacity"
                        >
                          <Plus className="h-3 w-3" />
                        </button>
                      )}
                    </div>

                    {/* Meals for the day */}
                    {isCurrentMonth && (
                      <div className="space-y-1">
                        {dayMeals.slice(0, 3).map((meal, mealIndex) => (
                          <div
                            key={mealIndex}
                            className={`
                              px-2 py-1 rounded-md text-xs font-medium border
                              ${mealColors[meal.mealType] || 'bg-gray-100 border-gray-300'}
                            `}
                          >
                            <div className="truncate">
                              {meal.dish.name}
                            </div>
                          </div>
                        ))}
                        
                        {dayMeals.length > 3 && (
                          <div className="text-xs text-gray-500 px-2">
                            +{dayMeals.length - 3} more
                          </div>
                        )}
                        
                        {/* Total calories */}
                        {totalCalories > 0 && (
                          <div className="text-xs text-gray-600 px-2 mt-2 border-t border-gray-200 pt-1">
                            {totalCalories} cal
                          </div>
                        )}
                        
                        {/* Empty state */}
                        {dayMeals.length === 0 && (
                          <div className="text-center py-4">
                            <Calendar className="h-6 w-6 text-gray-300 mx-auto mb-1" />
                            <p className="text-xs text-gray-500">No meals</p>
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </>
        )}
      </div>

      {/* Month Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6 text-center">
          <div className="text-3xl font-bold text-primary-600 mb-2">
            {monthMeals.length}
          </div>
          <p className="text-gray-600">Total Planned Meals</p>
        </div>
        
        <div className="card p-6 text-center">
          <div className="text-3xl font-bold text-secondary-600 mb-2">
            {monthMeals.reduce((total, meal) => total + (meal.dish.calories || 0), 0)}
          </div>
          <p className="text-gray-600">Total Calories</p>
        </div>
        
        <div className="card p-6 text-center">
          <div className="text-3xl font-bold text-purple-600 mb-2">
            {new Set(monthMeals.map(meal => meal.date)).size}
          </div>
          <p className="text-gray-600">Days with Meals</p>
        </div>
      </div>

      {/* Quick Add Section */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Add Meal</h3>
        <div className="flex flex-wrap gap-3">
          <select
            value={selectedMealType}
            onChange={(e) => setSelectedMealType(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          >
            <option value="breakfast">Breakfast</option>
            <option value="lunch">Lunch</option>
            <option value="dinner">Dinner</option>
            <option value="snack">Snack</option>
          </select>
          <p className="text-gray-600 py-2">
            Click on any day in the calendar to add a {selectedMealType} meal
          </p>
        </div>
      </div>

      {/* Dish Selector Modal */}
      {showDishSelector && selectedDay && (
        <DishSelector
          loadDishes={loadDishes}
          onSelect={handleDishSelect}
          onClose={() => {
            setShowDishSelector(false);
            setSelectedDay(null);
          }}
          mealType={selectedMealType}
          isEditing={false}
          onAddDish={onAddDish}
        />
      )}
    </div>
  );
};

export default MonthView;
