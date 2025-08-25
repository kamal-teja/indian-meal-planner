import React from 'react';
import { NavLink } from 'react-router-dom';
import { Home, Calendar, CalendarDays, BarChart3, Activity, ShoppingCart, Heart, User } from 'lucide-react';

const LeftNav = () => {
  const NavItem = ({ to, icon: Icon, label }) => (
    <NavLink to={to} end={to !== '/home'}>
      {({ isActive }) => (
        <div className={`flex items-center space-x-3 px-4 py-2 rounded-md transition-colors ${isActive ? 'bg-accent-600 text-white' : 'text-accent-800 hover:bg-accent-50'}`}>
          <Icon className={`h-4 w-4 ${isActive ? 'text-white' : 'text-accent-700'}`} />
          <span className="hidden sm:inline truncate">{label}</span>
        </div>
      )}
    </NavLink>
  );

  return (
    <aside className="w-56 hidden md:block pr-6">
      <nav className="sticky top-20 bg-white rounded-lg p-4 shadow-lg text-accent-900 ring-1 ring-accent-50">
        <div className="space-y-1">
          <NavItem to="/home" icon={Home} label="Home" />
          <NavItem to="/analytics" icon={BarChart3} label="Analytics" />
          <NavItem to="/nutrition" icon={Activity} label="Nutrition" />
          <NavItem to="/shopping-list" icon={ShoppingCart} label="Planner" />
          <NavItem to="/favorites" icon={Heart} label="Favorites" />
          <NavItem to="/profile" icon={User} label="Profile" />
        </div>
      </nav>
    </aside>
  );
};

export default LeftNav;
