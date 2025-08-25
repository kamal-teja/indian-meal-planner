import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Calendar, CalendarDays, Utensils, User, LogOut, ChevronDown, Settings, Heart, BarChart3, ShoppingCart, Activity } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

const Header = ({ currentView, onViewChange }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout, isAuthenticated } = useAuth();

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
    <header className="glass-header border-b border-neutral-200 sticky top-0 z-50">
      <div className="container mx-auto px-6">
        <div className="flex items-center justify-between h-16">
          {/* Logo and Title */}
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-secondary-600 rounded-lg shadow-sm">
              <Utensils className="h-8 w-8 text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-display font-bold gradient-text">
                Nourish
              </h1>
              <p className="text-sm text-accent-600 hidden sm:block">
                Plan your healthy, delicious meals
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {/* Navigation Tabs (only if authenticated) */}
            {isAuthenticated && (
              <div className="flex items-center space-x-2 bg-neutral-100 rounded-lg p-1">
                <button
                  onClick={() => handleViewChange('day')}
                  className={`nav-tab ${
                    isActive('day') ? 'nav-tab-active' : 'nav-tab-inactive'
                  }`}
                >
                  <Calendar className="h-4 w-4" />
                  <span className="hidden sm:inline">Day View</span>
                </button>
                
                <button
                  onClick={() => handleViewChange('month')}
                  className={`nav-tab ${
                    isActive('month') ? 'nav-tab-active' : 'nav-tab-inactive'
                  }`}
                >
                  <CalendarDays className="h-4 w-4" />
                  <span className="hidden sm:inline">Month View</span>
                </button>
              </div>
            )}

            {/* User Menu */}
            {isAuthenticated ? (
              <div className="flex items-center space-x-3">
                <div className="flex items-center space-x-2 bg-neutral-100 rounded-lg px-3 py-2 text-accent-700">
                  <div className="w-8 h-8 bg-sage-500 rounded-lg flex items-center justify-center">
                    <User className="h-4 w-4 text-white" />
                  </div>
                  <span className="hidden sm:inline font-medium">{user?.name}</span>
                </div>
                <button
                  onClick={handleLogout}
                  className="flex items-center space-x-2 px-3 py-2 text-sm text-warm-600 hover:bg-warm-50 rounded-lg transition-colors bg-neutral-100"
                >
                  <LogOut className="h-4 w-4" />
                  <span className="hidden sm:inline">Sign Out</span>
                </button>
              </div>
            ) : (
              <div className="flex items-center space-x-3">
                <button
                  onClick={() => navigate('/login')}
                  className="px-4 py-2 text-sm font-medium text-accent-700 hover:text-secondary-600 transition-colors bg-neutral-100 rounded-lg hover:bg-neutral-200"
                >
                  Sign In
                </button>
                <button
                  onClick={() => navigate('/register')}
                  className="btn-primary text-sm"
                >
                  Sign Up
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

  {/* No profile dropdown — simplified header */}
    </header>
  );
};

export default Header;
