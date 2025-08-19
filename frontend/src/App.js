import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import Header from './components/Header';
import DayView from './components/DayView';
import MonthView from './components/MonthView';
import Login from './components/auth/Login';
import Register from './components/auth/Register';
import ProtectedRoute from './components/auth/ProtectedRoute';
import Favorites from './components/Favorites';
import Analytics from './components/Analytics';
import ShoppingList from './components/ShoppingList';
import UserProfile from './components/UserProfile';
import NutritionDashboard from './components/NutritionDashboard';
import { mealPlannerAPI } from './services/api';

function App() {
  const [dishes, setDishes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [currentView, setCurrentView] = useState('day');

  useEffect(() => {
    loadDishes();
  }, []);

  const loadDishes = async () => {
    try {
      setLoading(true);
      const response = await mealPlannerAPI.getDishes();
      setDishes(response.data);
    } catch (error) {
      console.error('Error loading dishes:', error);
    } finally {
      setLoading(false);
    }
  };

  const addDish = async (dishData) => {
    try {
      const response = await mealPlannerAPI.addDish(dishData);
      const newDish = response.data;
      setDishes(prev => [...prev, newDish]);
      return newDish;
    } catch (error) {
      console.error('Error adding dish:', error);
      throw error;
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-4 border-primary-200 border-t-primary-600 mx-auto mb-4"></div>
          <h2 className="text-2xl font-display font-semibold gradient-text">
            Loading Delicious Meals...
          </h2>
          <p className="text-gray-600 mt-2">Preparing your Indian cuisine experience</p>
        </div>
      </div>
    );
  }

  return (
    <Router>
      <AuthProvider>
        <div className="min-h-screen bg-gradient-to-br from-primary-50 via-white to-secondary-50 bg-spice-pattern">
          <Header currentView={currentView} onViewChange={setCurrentView} />
          
          <main className="container mx-auto px-4 py-8">
            <Routes>
              {/* Authentication Routes */}
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              
              {/* Protected Routes */}
              <Route 
                path="/day" 
                element={
                  <ProtectedRoute>
                    <DayView dishes={dishes} onAddDish={addDish} />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/month" 
                element={
                  <ProtectedRoute>
                    <MonthView dishes={dishes} onAddDish={addDish} />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/favorites" 
                element={
                  <ProtectedRoute>
                    <Favorites />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/analytics" 
                element={
                  <ProtectedRoute>
                    <Analytics />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/nutrition" 
                element={
                  <ProtectedRoute>
                    <NutritionDashboard />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/shopping-list" 
                element={
                  <ProtectedRoute>
                    <ShoppingList />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/profile" 
                element={
                  <ProtectedRoute>
                    <UserProfile />
                  </ProtectedRoute>
                } 
              />
              
              {/* Default Route */}
              <Route 
                path="/" 
                element={<Navigate to="/day" replace />} 
              />
            </Routes>
          </main>
          
          {/* Decorative background elements */}
          <div className="fixed inset-0 pointer-events-none overflow-hidden">
            <div className="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-primary-200/30 to-secondary-200/30 rounded-full blur-3xl"></div>
            <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-secondary-200/30 to-primary-200/30 rounded-full blur-3xl"></div>
          </div>
        </div>
      </AuthProvider>
    </Router>
  );
}

export default App;
