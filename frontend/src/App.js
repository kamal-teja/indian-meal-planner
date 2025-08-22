import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { NotificationProvider } from './contexts/NotificationContext';
import NotificationSnackbar from './components/NotificationSnackbar';
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
  const [currentView, setCurrentView] = useState('day');

  const loadDishes = async (params = {}) => {
    try {
      const response = await mealPlannerAPI.getDishes(params);
      return response.data;
    } catch (error) {
      console.error('Error loading dishes:', error);
      throw error;
    }
  };

  const addDish = async (dishData) => {
    try {
      const response = await mealPlannerAPI.addDish(dishData);
      const newDish = response.data;
      return newDish;
    } catch (error) {
      console.error('Error adding dish:', error);
      throw error;
    }
  };

  return (
    <Router>
      <AuthProvider>
        <NotificationProvider>
          <div className="min-h-screen bg-neutral-50">
          <Header currentView={currentView} onViewChange={setCurrentView} />
          
          <main className="container mx-auto px-6 py-8">
            <Routes>
              {/* Authentication Routes */}
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              
              {/* Protected Routes */}
              <Route 
                path="/day" 
                element={
                  <ProtectedRoute>
                    <DayView loadDishes={loadDishes} onAddDish={addDish} />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/month" 
                element={
                  <ProtectedRoute>
                    <MonthView loadDishes={loadDishes} onAddDish={addDish} />
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
          <NotificationSnackbar />
        </NotificationProvider>
      </AuthProvider>
    </Router>
  );
}

export default App;
