import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:5000/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token expiration
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const mealPlannerAPI = {
  // Authentication
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData),
  logout: () => api.post('/auth/logout'),
  getCurrentUser: () => api.get('/auth/me'),
  updateProfile: (profileData) => api.put('/user/profile', profileData),
  toggleFavorite: (dishId) => api.post(`/dishes/${dishId}/favorite`),
  
  // Get all available dishes with pagination
  getDishes: (params = {}) => api.get('/dishes', { params }),
  
  // Search dishes
  searchDishes: (params) => api.get('/dishes/search', { params }),
  
  // Get favorite dishes
  getFavoriteDishes: () => api.get('/dishes/favorites'),
  
  // Add a new dish
  addDish: (dishData) => api.post('/dishes', dishData),
  
  // Get meals for a specific date (YYYY-MM-DD format)
  getMealsByDate: (date) => api.get(`/meals/${date}`),
  
  // Get meals for a month
  getMealsByMonth: (year, month) => api.get(`/meals/month/${year}/${month}`),
  
  // Add a new meal
  addMeal: (mealData) => api.post('/meals', mealData),
  
  // Update an existing meal
  updateMeal: (mealId, updateData) => api.put(`/meals/${mealId}`, updateData),
  
  // Delete a meal
  deleteMeal: (mealId) => api.delete(`/meals/${mealId}`),
  
  // Analytics
  getMealAnalytics: (period = 30) => api.get(`/analytics?period=${period}`),
  
  // Shopping list
  getShoppingList: (startDate, endDate) => api.get(`/shopping-list?startDate=${startDate}&endDate=${endDate}`),
  
  // Recommendations
  getRecommendations: (mealType, date) => api.get(`/recommendations?mealType=${mealType}&date=${date}`),
  
  // Nutrition tracking
  getNutritionProgress: (period = 7) => api.get(`/nutrition/progress?period=${period}`),
  getNutritionGoals: () => api.get('/nutrition/goals'),
  updateNutritionGoals: (goals) => api.put('/nutrition/goals', goals),
  
  // Health check
  healthCheck: () => api.get('/health'),
};

export default mealPlannerAPI;
