import React, { useState, useEffect } from 'react';
import { Target, TrendingUp, Award, AlertCircle, Calendar, Activity, BarChart3, Settings, Star } from 'lucide-react';
import mealPlannerAPI from '../services/api';

const NutritionDashboard = () => {
  const [nutritionData, setNutritionData] = useState(null);
  const [goals, setGoals] = useState({
    dailyCalories: 2000,
    protein: 150,
    carbs: 250,
    fat: 65,
    fiber: 25,
    sodium: 2300
  });
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState(7); // days

  useEffect(() => {
    fetchNutritionData();
  }, [period]);

  const fetchNutritionData = async () => {
    try {
      setLoading(true);
      
      // Get nutrition progress and goals
      const [progressResponse, goalsResponse] = await Promise.all([
        mealPlannerAPI.getNutritionProgress(period),
        mealPlannerAPI.getNutritionGoals()
      ]);
      
      // Debug: Log the actual response data
      console.log('Nutrition Progress Response:', progressResponse.data);
      console.log('Nutrition Goals Response:', goalsResponse.data);
      
      // Handle the response format from backend: { success: true, data: ... }
      const progressData = progressResponse.data.data;
      const goalsData = goalsResponse.data.data;
      
      console.log('Processed Progress Data:', progressData);
      console.log('Processed Goals Data:', goalsData);
      
      setNutritionData(progressData);
      
      // Ensure goalsData has all required fields with defaults
      const completeGoals = {
        dailyCalories: goalsData?.dailyCalories || 2000,
        protein: goalsData?.protein || 150,
        carbs: goalsData?.carbs || 250,
        fat: goalsData?.fat || 65,
        fiber: goalsData?.fiber || 25,
        sodium: goalsData?.sodium || 2300
      };
      
      console.log('Setting goals to:', completeGoals);
      setGoals(completeGoals);
    } catch (error) {
      console.error('Error fetching nutrition data:', error);
      // Log more details about the error
      if (error.response) {
        console.error('Error Response Status:', error.response.status);
        console.error('Error Response Data:', error.response.data);
      }
      setNutritionData(null);
    } finally {
      setLoading(false);
    }
  };

  const updateGoals = async () => {
    try {
      console.log('Updating goals with data:', goals);
      const response = await mealPlannerAPI.updateNutritionGoals(goals);
      console.log('Update response:', response.data);
      
      // Handle the response format from backend: { success: true, data: goals }
      const updatedGoals = response.data.data;
      console.log('Updated goals received:', updatedGoals);
      
      setGoals(updatedGoals);
      setIsEditing(false);
      
      // Reload data to see the changes
      fetchNutritionData();
    } catch (error) {
      console.error('Error updating goals:', error);
      if (error.response) {
        console.error('Error response:', error.response.data);
        console.error('Error status:', error.response.status);
      }
    }
  };

  const getProgressColor = (percentage) => {
    if (percentage >= 90 && percentage <= 110) return 'text-green-600 bg-green-100';
    if (percentage >= 70 && percentage <= 130) return 'text-yellow-600 bg-yellow-100';
    return 'text-red-600 bg-red-100';
  };

  const getProgressWidth = (percentage) => {
    return Math.min(percentage, 100);
  };

  if (loading) {
    return (
      <div className="p-6 bg-white rounded-lg shadow-md">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded mb-4"></div>
          <div className="space-y-4">
            {[1, 2, 3, 4].map(i => (
              <div key={i} className="h-16 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!nutritionData) {
    return (
      <div className="p-6 bg-white rounded-lg shadow-md text-center">
        <AlertCircle className="h-12 w-12 text-gray-400 mx-auto mb-4" />
        <h3 className="text-lg font-semibold text-gray-900 mb-2">No Nutrition Data Available</h3>
        <p className="text-gray-600 mb-4">
          Start logging meals to see your nutrition progress and track your daily goals!
        </p>
        <div className="text-sm text-gray-500">
          <p>• Add meals from the Day View or Month View</p>
          <p>• Track calories, protein, carbs, and more</p>
          <p>• Set personalized nutrition goals</p>
        </div>
      </div>
    );
  }

  // Handle the backend response format: { period, progress, goals, summary }
  const dailyData = nutritionData.progress || [];
  const todayData = dailyData.length > 0 ? dailyData[dailyData.length - 1] : null;
  const avgData = nutritionData.summary;

  // If there's no progress data, show a helpful message
  if (dailyData.length === 0) {
    return (
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-900 flex items-center">
            <Activity className="h-8 w-8 text-blue-600 mr-3" />
            Nutrition Dashboard
          </h2>
        </div>
        
        <div className="p-8 bg-white rounded-lg shadow-md text-center">
          <Activity className="h-16 w-16 text-blue-400 mx-auto mb-4" />
          <h3 className="text-xl font-semibold text-gray-900 mb-3">Start Your Nutrition Journey!</h3>
          <p className="text-gray-600 mb-6 max-w-md mx-auto">
            You haven't logged any meals yet. Add some meals to your meal plan to start tracking your nutrition progress.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-2xl mx-auto text-sm text-gray-600">
            <div className="p-4 bg-blue-50 rounded-lg">
              <h4 className="font-medium text-blue-800 mb-2">📅 Plan Meals</h4>
              <p>Use the Day View to add meals to your daily plan</p>
            </div>
            <div className="p-4 bg-green-50 rounded-lg">
              <div className="flex items-center space-x-2 mb-2">
                <BarChart3 className="h-5 w-5 text-secondary-600" />
                <h4 className="font-medium text-secondary-800">Track Progress</h4>
              </div>
              <p>Monitor calories, protein, carbs, and other nutrients</p>
            </div>
            <div className="p-4 bg-purple-50 rounded-lg">
              <div className="flex items-center space-x-2 mb-2">
                <Target className="h-5 w-5 text-sage-600" />
                <h4 className="font-medium text-sage-800">Set Goals</h4>
              </div>
              <p>Customize your daily nutrition targets</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900 flex items-center">
          <Activity className="h-8 w-8 text-blue-600 mr-3" />
          Nutrition Dashboard
        </h2>
        <div className="flex items-center space-x-4">
          <select 
            value={period} 
            onChange={(e) => setPeriod(parseInt(e.target.value))}
            className="dropdown-elegant"
          >
            <option value={7}>Last 7 days</option>
            <option value={14}>Last 14 days</option>
            <option value={30}>Last 30 days</option>
          </select>
          <button
            onClick={() => setIsEditing(!isEditing)}
            className="btn-secondary flex items-center space-x-2"
          >
            <Settings className="h-4 w-4" />
            <span>{isEditing ? 'Cancel' : 'Edit Goals'}</span>
          </button>
        </div>
      </div>

      {/* Today's Progress */}
      {todayData ? (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold mb-4 flex items-center">
            <Calendar className="h-5 w-5 text-green-600 mr-2" />
            Today's Progress ({new Date(todayData.date).toLocaleDateString()})
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
            {[
              { name: 'calories', value: todayData.calories, goal: goals.dailyCalories, unit: '' },
              { name: 'protein', value: todayData.protein, goal: goals.protein, unit: 'g' },
              { name: 'carbs', value: todayData.carbs, goal: goals.carbs, unit: 'g' },
              { name: 'fat', value: todayData.fat, goal: goals.fat, unit: 'g' },
              { name: 'fiber', value: todayData.fiber, goal: goals.fiber, unit: 'g' },
              { name: 'sodium', value: todayData.sodium, goal: goals.sodium, unit: 'mg' }
            ].map(({ name, value, goal, unit }) => {
              const percentage = goal > 0 ? (value / goal) * 100 : 0;
              return (
                <div key={name} className="text-center">
                  <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getProgressColor(percentage)}`}>
                    {Math.round(percentage)}%
                  </div>
                  <p className="text-sm text-gray-600 mt-1 capitalize">{name}</p>
                  <p className="text-xs text-gray-500">{Math.round(value)}{unit} / {goal}{unit}</p>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
                    <div 
                      className={`h-2 rounded-full transition-all ${
                        percentage >= 90 && percentage <= 110 
                          ? 'bg-green-500' 
                          : percentage > 110 
                            ? 'bg-red-500' 
                            : 'bg-yellow-500'
                      }`}
                      style={{ width: `${Math.min(percentage, 100)}%` }}
                    ></div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold mb-4 flex items-center">
            <Calendar className="h-5 w-5 text-orange-600 mr-2" />
            Today's Progress ({new Date().toLocaleDateString()})
          </h3>
          <div className="text-center py-8">
            <Calendar className="h-12 w-12 text-gray-400 mx-auto mb-3" />
            <p className="text-gray-600 mb-2">No meals logged for today yet</p>
            <p className="text-sm text-gray-500">Add some meals to start tracking your daily nutrition!</p>
          </div>
        </div>
      )}

      {/* Goals Setting */}
      {isEditing && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold mb-4 flex items-center">
            <Target className="h-5 w-5 text-blue-600 mr-2" />
            Nutrition Goals
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            {Object.entries(goals).map(([key, value]) => (
              <div key={key}>
                <label className="block text-sm font-medium text-gray-700 mb-1 capitalize">
                  {key.replace(/([A-Z])/g, ' $1').toLowerCase()}
                  {key === 'sodium' ? ' (mg)' : key.includes('Calories') ? '' : ' (g)'}
                </label>
                <input
                  type="number"
                  value={value}
                  onChange={(e) => setGoals({ ...goals, [key]: parseFloat(e.target.value) || 0 })}
                  className="w-full border border-gray-300 rounded-md px-3 py-2"
                  min="0"
                />
              </div>
            ))}
          </div>
          <div className="flex justify-end space-x-3 mt-6">
            <button
              onClick={() => setIsEditing(false)}
              className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={updateGoals}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Save Goals
            </button>
          </div>
        </div>
      )}

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center">
            <TrendingUp className="h-8 w-8 text-green-600" />
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Avg Daily Calories</p>
              <p className="text-2xl font-bold text-gray-900">{Math.round(avgData.avgCalories)}</p>
              <p className="text-sm text-gray-500">Goal: {goals.dailyCalories}</p>
            </div>
          </div>
        </div>
        
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center">
            <Award className="h-8 w-8 text-blue-600" />
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Days On Track</p>
              <p className="text-2xl font-bold text-gray-900">{avgData.daysOnTrack || 0}</p>
              <p className="text-sm text-gray-500">Out of {dailyData.length} days</p>
            </div>
          </div>
        </div>
        
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center">
            <Activity className="h-8 w-8 text-purple-600" />
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Avg Protein</p>
              <p className="text-2xl font-bold text-gray-900">{Math.round(avgData.avgProtein)}g</p>
              <p className="text-sm text-gray-500">Goal: {goals.protein}g</p>
            </div>
          </div>
        </div>
      </div>

      {/* Daily Breakdown */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold mb-4">Daily Breakdown</h3>
        <div className="space-y-4">
          {dailyData.slice().reverse().map((day) => (
            <div key={day.date} className="border border-gray-200 rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <h4 className="font-medium text-gray-900">{new Date(day.date).toLocaleDateString()}</h4>
                <div className="flex space-x-2">
                  {day.calories >= goals.dailyCalories && 
                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-green-100 text-green-800">
                      Calories ✓
                    </span>
                  }
                  {day.protein >= goals.protein && 
                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-blue-100 text-blue-800">
                      Protein ✓
                    </span>
                  }
                  {day.sodium <= goals.sodium && 
                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-purple-100 text-purple-800">
                      Sodium ✓
                    </span>
                  }
                </div>
              </div>
              <div className="grid grid-cols-3 md:grid-cols-6 gap-4 text-sm">
                <div>
                  <span className="text-gray-600">Calories:</span>
                  <p className="font-medium">{Math.round(day.calories || 0)}</p>
                </div>
                <div>
                  <span className="text-gray-600">Protein:</span>
                  <p className="font-medium">{Math.round(day.protein || 0)}g</p>
                </div>
                <div>
                  <span className="text-gray-600">Carbs:</span>
                  <p className="font-medium">{Math.round(day.carbs || 0)}g</p>
                </div>
                <div>
                  <span className="text-gray-600">Fat:</span>
                  <p className="font-medium">{Math.round(day.fat || 0)}g</p>
                </div>
                <div>
                  <span className="text-gray-600">Fiber:</span>
                  <p className="font-medium">{Math.round(day.fiber || 0)}g</p>
                </div>
                <div>
                  <span className="text-gray-600">Sodium:</span>
                  <p className="font-medium">{Math.round(day.sodium || 0)}mg</p>
                </div>
              </div>
              {day.meals && day.meals.length > 0 && (
                <div className="mt-3 pt-3 border-t border-gray-100">
                  <p className="text-sm text-gray-600 mb-2">Meals:</p>
                  <div className="flex flex-wrap gap-2">
                    {day.meals.map((meal, idx) => (
                      <span key={idx} className="inline-flex items-center px-2 py-1 rounded-md text-xs bg-gray-100 text-gray-700">
                        {meal.mealType}: {meal.dish?.name || 'Unknown dish'}
                        {meal.rating && (
                          <span className="ml-1 flex items-center">
                            <Star className="h-3 w-3 text-warm-500 mr-1" />
                            {meal.rating}
                          </span>
                        )}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default NutritionDashboard;
