import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Calendar, CalendarDays, Utensils, User, LogOut, ChevronDown, Settings, Heart, BarChart3, ShoppingCart, Activity } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

const Header = ({ currentView, onViewChange }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout, isAuthenticated } = useAuth();
  const [dropdownOpen, setDropdownOpen] = useState(false);

  const handleViewChange = (view) => {
    onViewChange(view);
    navigate(`/${view}`);
  };

  const isActive = (path) => location.pathname === `/${path}`;

  const handleLogout = async () => {
    await logout();
    navigate('/login');
    setDropdownOpen(false);
  };

  const handleNavigation = (path) => {
    navigate(path);
    setDropdownOpen(false);
  };

  return (
    <header className="glass-effect border-b border-white/30 sticky top-0 z-50">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo and Title */}
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-gradient-to-r from-primary-500 to-primary-600 rounded-xl shadow-lg">
              <Utensils className="h-8 w-8 text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-display font-bold gradient-text">
                Indian Meal Planner
              </h1>
              <p className="text-sm text-gray-600 hidden sm:block">
                Plan your delicious Indian meals
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {/* Navigation Tabs (only if authenticated) */}
            {isAuthenticated && (
              <div className="flex items-center space-x-2 bg-white/50 rounded-xl p-1 backdrop-blur-sm">
                <button
                  onClick={() => handleViewChange('day')}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg font-medium transition-all duration-300 ${
                    isActive('day')
                      ? 'bg-white text-primary-700 shadow-md'
                      : 'text-gray-600 hover:text-primary-600 hover:bg-white/50'
                  }`}
                >
                  <Calendar className="h-4 w-4" />
                  <span className="hidden sm:inline">Day View</span>
                </button>
                
                <button
                  onClick={() => handleViewChange('month')}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg font-medium transition-all duration-300 ${
                    isActive('month')
                      ? 'bg-white text-primary-700 shadow-md'
                      : 'text-gray-600 hover:text-primary-600 hover:bg-white/50'
                  }`}
                >
                  <CalendarDays className="h-4 w-4" />
                  <span className="hidden sm:inline">Month View</span>
                </button>
              </div>
            )}

            {/* User Menu */}
            {isAuthenticated ? (
              <div className="relative">
                <button
                  onClick={() => setDropdownOpen(!dropdownOpen)}
                  className="flex items-center space-x-2 bg-white/50 rounded-xl px-3 py-2 text-gray-700 hover:bg-white/70 transition-all duration-200"
                >
                  <div className="w-8 h-8 bg-gradient-to-r from-primary-500 to-primary-600 rounded-full flex items-center justify-center">
                    <User className="h-4 w-4 text-white" />
                  </div>
                  <span className="hidden sm:inline font-medium">{user?.name}</span>
                  <ChevronDown className={`h-4 w-4 transition-transform ${dropdownOpen ? 'rotate-180' : ''}`} />
                </button>

                {dropdownOpen && (
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg py-1 z-50 border border-gray-200">
                    <button
                      onClick={() => handleNavigation('/profile')}
                      className="flex items-center space-x-2 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                    >
                      <Settings className="h-4 w-4" />
                      <span>Profile Settings</span>
                    </button>
                    <button
                      onClick={() => handleNavigation('/favorites')}
                      className="flex items-center space-x-2 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                    >
                      <Heart className="h-4 w-4" />
                      <span>Favorite Dishes</span>
                    </button>
                    <button
                      onClick={() => handleNavigation('/analytics')}
                      className="flex items-center space-x-2 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                    >
                      <BarChart3 className="h-4 w-4" />
                      <span>Analytics</span>
                    </button>
                    <button
                      onClick={() => handleNavigation('/nutrition')}
                      className="flex items-center space-x-2 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                    >
                      <Activity className="h-4 w-4" />
                      <span>Nutrition</span>
                    </button>
                    <button
                      onClick={() => handleNavigation('/shopping-list')}
                      className="flex items-center space-x-2 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                    >
                      <ShoppingCart className="h-4 w-4" />
                      <span>Shopping List</span>
                    </button>
                    <div className="border-t border-gray-100 my-1"></div>
                    <button
                      onClick={handleLogout}
                      className="flex items-center space-x-2 w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50"
                    >
                      <LogOut className="h-4 w-4" />
                      <span>Sign Out</span>
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => navigate('/login')}
                  className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-primary-600 transition-colors"
                >
                  Sign In
                </button>
                <button
                  onClick={() => navigate('/register')}
                  className="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-primary-500 to-primary-600 rounded-lg hover:from-primary-600 hover:to-primary-700 transition-all duration-200"
                >
                  Sign Up
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Close dropdown when clicking outside */}
      {dropdownOpen && (
        <div
          className="fixed inset-0 z-40"
          onClick={() => setDropdownOpen(false)}
        />
      )}
    </header>
  );
};

export default Header;
