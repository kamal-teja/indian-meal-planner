import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Header from './components/Header';
import DayView from './components/DayView';
import MonthView from './components/MonthView';
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
      <div className="min-h-screen bg-gradient-to-br from-primary-50 via-white to-secondary-50 bg-spice-pattern">
        <Header currentView={currentView} onViewChange={setCurrentView} />
        
        <main className="container mx-auto px-4 py-8">
          <Routes>
            <Route 
              path="/day" 
              element={<DayView dishes={dishes} onAddDish={addDish} />} 
            />
            <Route 
              path="/month" 
              element={<MonthView dishes={dishes} onAddDish={addDish} />} 
            />
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
    </Router>
  );
}

export default App;
