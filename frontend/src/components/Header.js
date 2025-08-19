import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Calendar, CalendarDays, Utensils } from 'lucide-react';

const Header = ({ currentView, onViewChange }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const handleViewChange = (view) => {
    onViewChange(view);
    navigate(`/${view}`);
  };

  const isActive = (path) => location.pathname === `/${path}`;

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

          {/* Navigation Tabs */}
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
        </div>
      </div>
    </header>
  );
};

export default Header;
