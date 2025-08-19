import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const mealPlannerAPI = {
  // Get all available dishes
  getDishes: () => api.get('/dishes'),
  
  // Add a new dish
  addDish: (dishData) => api.post('/dishes', dishData),
  
  // Get meals for a specific date (YYYY-MM-DD format)
  getMealsByDate: (date) => api.get(`/meals/${date}`),
  
  // Get meals for a month
  getMealsByMonth: (year, month) => api.get(`/meals/month/${year}/${month}`),
  
  // Add a new meal
  addMeal: (mealData) => api.post('/meals', mealData),
  
  // Update an existing meal
  updateMeal: (mealId, dishId) => api.put(`/meals/${mealId}`, { dishId }),
  
  // Delete a meal
  deleteMeal: (mealId) => api.delete(`/meals/${mealId}`),
  
  // Health check
  healthCheck: () => api.get('/health'),
};

export default mealPlannerAPI;
