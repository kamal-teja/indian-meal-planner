import React, { useState, useEffect } from 'react';
import { BarChart3, TrendingUp, Calendar, Utensils, Star, Clock } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';

const Analytics = () => {
  const [analytics, setAnalytics] = useState(null);
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState(30);

  useEffect(() => {
    loadAnalytics();
  }, [period]);

  const loadAnalytics = async () => {
    try {
      setLoading(true);
      const response = await mealPlannerAPI.getMealAnalytics(period);
      // Handle the response format from backend: { success: true, data: analytics }
      setAnalytics(response.data.data);
    } catch (error) {
      console.error('Error loading analytics:', error);
      setAnalytics(null);
    } finally {
      setLoading(false);
    }
  };

  const StatCard = ({ icon: Icon, title, value, subtitle, color = 'primary' }) => (
    <div className="bg-white rounded-xl shadow-md p-6 hover:shadow-lg transition-shadow duration-300">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-600 mb-1">{title}</p>
          <p className="text-3xl font-bold text-gray-900">{value}</p>
          {subtitle && <p className="text-sm text-gray-500 mt-1">{subtitle}</p>}
        </div>
        <div className={`p-3 bg-${color}-100 rounded-xl`}>
          <Icon className={`h-6 w-6 text-${color}-600`} />
        </div>
      </div>
    </div>
  );

  const DistributionChart = ({ data, title }) => {
    const total = Object.values(data).reduce((sum, value) => sum + value, 0);
    const colors = {
      breakfast: 'bg-yellow-400',
      lunch: 'bg-green-400',
      dinner: 'bg-blue-400',
      snack: 'bg-purple-400'
    };

    return (
      <div className="bg-white rounded-xl shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">{title}</h3>
        <div className="space-y-3">
          {Object.entries(data).map(([key, value]) => {
            const percentage = total > 0 ? (value / total) * 100 : 0;
            return (
              <div key={key} className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className={`w-4 h-4 rounded-full ${colors[key] || 'bg-gray-400'}`}></div>
                  <span className="text-sm font-medium text-gray-700 capitalize">
                    {key}
                  </span>
                </div>
                <div className="flex items-center space-x-2">
                  <div className="w-24 bg-gray-200 rounded-full h-2">
                    <div 
                      className={`h-2 rounded-full ${colors[key] || 'bg-gray-400'}`}
                      style={{ width: `${percentage}%` }}
                    ></div>
                  </div>
                  <span className="text-sm text-gray-500 w-8 text-right">
                    {value}
                  </span>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  const TopCuisinesChart = ({ cuisines }) => {
    const sortedCuisines = Object.entries(cuisines)
      .sort(([,a], [,b]) => b - a)
      .slice(0, 5);
    const maxValue = Math.max(...sortedCuisines.map(([, value]) => value));

    const colors = [
      'bg-red-400',
      'bg-orange-400', 
      'bg-amber-400',
      'bg-emerald-400',
      'bg-blue-400'
    ];

    return (
      <div className="bg-white rounded-xl shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Cuisines</h3>
        <div className="space-y-3">
          {sortedCuisines.map(([cuisine, count], index) => {
            const percentage = maxValue > 0 ? (count / maxValue) * 100 : 0;
            return (
              <div key={cuisine} className="flex items-center justify-between">
                <div className="flex items-center space-x-3 flex-1">
                  <div className={`w-4 h-4 rounded-full ${colors[index]}`}></div>
                  <span className="text-sm font-medium text-gray-700 truncate">
                    {cuisine}
                  </span>
                </div>
                <div className="flex items-center space-x-3 ml-4">
                  <div className="w-20 bg-gray-200 rounded-full h-2">
                    <div 
                      className={`h-2 rounded-full ${colors[index]}`}
                      style={{ width: `${percentage}%` }}
                    ></div>
                  </div>
                  <span className="text-sm text-gray-500 w-6 text-right">
                    {count}
                  </span>
                </div>
              </div>
            );
          })}
        </div>
        {sortedCuisines.length === 0 && (
          <p className="text-gray-500 text-center py-4">No cuisine data available</p>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-4 border-primary-200 border-t-primary-600 mx-auto mb-4"></div>
          <h2 className="text-2xl font-display font-semibold gradient-text">
            Loading Analytics...
          </h2>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl shadow-lg">
              <BarChart3 className="h-8 w-8 text-white" />
            </div>
            <div>
              <h1 className="text-3xl font-display font-bold gradient-text">
                Meal Analytics
              </h1>
              <p className="text-gray-600">
                Insights into your eating patterns
              </p>
            </div>
          </div>

          {/* Period Selector */}
          <div className="flex items-center space-x-2">
            <label className="text-sm font-medium text-gray-700">Period:</label>
            <select
              value={period}
              onChange={(e) => setPeriod(Number(e.target.value))}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            >
              <option value={7}>Last 7 days</option>
              <option value={30}>Last 30 days</option>
              <option value={90}>Last 90 days</option>
            </select>
          </div>
        </div>
      </div>

      {analytics ? (
        <div className="space-y-6">
          {/* Stats Overview */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatCard
              icon={Utensils}
              title="Total Meals"
              value={analytics.totalMeals}
              subtitle={`in last ${period} days`}
              color="primary"
            />
            <StatCard
              icon={TrendingUp}
              title="Avg Calories/Day"
              value={analytics.averageCaloriesPerDay}
              subtitle="calories"
              color="green"
            />
            <StatCard
              icon={Star}
              title="Average Rating"
              value={analytics.averageRating || 'N/A'}
              subtitle={analytics.averageRating ? 'out of 5 stars' : 'No ratings yet'}
              color="yellow"
            />
            <StatCard
              icon={Calendar}
              title="Total Calories"
              value={analytics.totalCalories.toLocaleString()}
              subtitle="calories consumed"
              color="red"
            />
          </div>

          {/* Charts */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <DistributionChart
              data={analytics.mealTypeDistribution}
              title="Meal Distribution"
            />
            <TopCuisinesChart cuisines={analytics.topCuisines} />
          </div>

          {/* Health Insights */}
          <div className="bg-white rounded-xl shadow-md p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center space-x-2">
              <TrendingUp className="h-5 w-5" />
              <span>Health Insights</span>
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div className="flex items-center justify-between p-4 bg-green-50 rounded-lg">
                  <div>
                    <p className="font-medium text-green-800">Daily Calorie Average</p>
                    <p className="text-sm text-green-600">
                      {analytics.averageCaloriesPerDay < 2000 ? 'Below recommended' : 
                       analytics.averageCaloriesPerDay > 2500 ? 'Above recommended' : 'Within range'}
                    </p>
                  </div>
                  <div className="text-2xl font-bold text-green-700">
                    {analytics.averageCaloriesPerDay}
                  </div>
                </div>
                
                <div className="flex items-center justify-between p-4 bg-blue-50 rounded-lg">
                  <div>
                    <p className="font-medium text-blue-800">Meal Frequency</p>
                    <p className="text-sm text-blue-600">
                      {analytics.totalMeals / period < 2 ? 'Consider more meals' : 
                       analytics.totalMeals / period > 4 ? 'Great meal frequency' : 'Good frequency'}
                    </p>
                  </div>
                  <div className="text-2xl font-bold text-blue-700">
                    {(analytics.totalMeals / period).toFixed(1)}/day
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <div className="p-4 bg-purple-50 rounded-lg">
                  <h4 className="font-medium text-purple-800 mb-2">Cuisine Diversity</h4>
                  <p className="text-sm text-purple-600 mb-2">
                    You've tried {Object.keys(analytics.topCuisines).length} different cuisines
                  </p>
                  <div className="text-lg font-semibold text-purple-700">
                    {Object.keys(analytics.topCuisines).length < 3 ? 'Try more variety!' :
                     Object.keys(analytics.topCuisines).length > 5 ? 'Excellent variety!' : 'Good diversity!'}
                  </div>
                </div>

                {analytics.averageRating > 0 && (
                  <div className="p-4 bg-yellow-50 rounded-lg">
                    <h4 className="font-medium text-yellow-800 mb-2">Meal Satisfaction</h4>
                    <p className="text-sm text-yellow-600 mb-2">
                      Average rating of your meals
                    </p>
                    <div className="flex items-center space-x-1">
                      {[1, 2, 3, 4, 5].map((star) => (
                        <Star
                          key={star}
                          className={`h-5 w-5 ${
                            star <= Math.round(analytics.averageRating)
                              ? 'text-yellow-400 fill-current'
                              : 'text-gray-300'
                          }`}
                        />
                      ))}
                      <span className="ml-2 text-lg font-semibold text-yellow-700">
                        {analytics.averageRating}
                      </span>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div className="text-center py-16">
          <BarChart3 className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-2xl font-display font-semibold text-gray-700 mb-2">
            No Data Available
          </h3>
          <p className="text-gray-500">
            Start planning and eating meals to see your analytics here.
          </p>
        </div>
      )}
    </div>
  );
};

export default Analytics;
