import React from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Calendar, CalendarDays } from 'lucide-react';

const HomeLayout = ({ onViewChange }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const handleViewChange = (view) => {
    if (onViewChange) onViewChange(view);
    navigate(`/home/${view}`);
  };

  const isActive = (path) => location.pathname.endsWith(`/${path}`) || location.pathname === `/home/${path}`;

  return (
    <div>
      {/* Inline, right-aligned compact toggle to preserve vertical space and maintain layout */}
      <div className="mb-4 flex justify-end">
        <div className="flex items-center space-x-1 bg-neutral-100 rounded-md p-1 shadow-sm">
          <button
            onClick={() => handleViewChange('day')}
            aria-pressed={isActive('day')}
            className={`flex items-center px-3 py-1 rounded-md text-sm ${isActive('day') ? 'bg-accent-600 text-white' : 'text-accent-800 hover:bg-accent-50'}`}
          >
            <Calendar className="h-4 w-4 mr-2" />
            <span className="hidden sm:inline">Day</span>
          </button>

          <button
            onClick={() => handleViewChange('month')}
            aria-pressed={isActive('month')}
            className={`flex items-center px-3 py-1 rounded-md text-sm ${isActive('month') ? 'bg-accent-600 text-white' : 'text-accent-800 hover:bg-accent-50'}`}
          >
            <CalendarDays className="h-4 w-4 mr-2" />
            <span className="hidden sm:inline">Month</span>
          </button>
        </div>
      </div>

      <Outlet />
    </div>
  );
};

export default HomeLayout;
