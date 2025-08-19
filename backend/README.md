# 🍛 Indian Meal Planner App

A modern, beautiful React-based meal planning application featuring authentic Indian cuisine with both day and month views. Built with React, Node.js, Express, and styled with Tailwind CSS.

## ✨ Features

- **📅 Day View**: Plan meals for specific days with detailed meal cards
- **🗓️ Month View**: Calendar-based monthly meal planning
- **🍽️ Indian Cuisine**: Curated collection of authentic Indian dishes
- **🎨 Modern UI**: Beautiful, responsive design with Tailwind CSS
- **📱 Mobile Responsive**: Works seamlessly on all devices
- **🔍 Smart Search**: Filter dishes by cuisine, type, and ingredients
- **📊 Nutrition Tracking**: Track calories and meal information
- **⚡ Real-time Updates**: Add, edit, and delete meals instantly

## 🛠️ Tech Stack

### Frontend
- **React 18** - Modern React with hooks
- **Tailwind CSS** - Utility-first CSS framework
- **React Router** - Client-side routing
- **Lucide React** - Beautiful icons
- **Axios** - HTTP client
- **date-fns** - Date manipulation library

### Backend
- **Node.js** - JavaScript runtime
- **Express.js** - Web framework
- **CORS** - Cross-origin resource sharing
- **UUID** - Unique identifier generation

## 🚀 Getting Started

### Prerequisites
- Node.js 16+ and npm/yarn
- Git

### Installation

1. **Install Backend Dependencies**
   ```bash
   cd backend
   npm install
   ```

2. **Install Frontend Dependencies**
   ```bash
   cd frontend
   npm install
   ```

3. **Start the Backend Server**
   ```bash
   cd backend
   npm run dev
   ```
   Server will start on http://localhost:5000

4. **Start the Frontend Development Server**
   ```bash
   cd frontend
   npm start
   ```
   App will open on http://localhost:3000

## 📁 Project Structure

```
indian-meal-planner/
├── backend/
│   ├── package.json
│   └── server.js
├── frontend/
│   ├── public/
│   │   └── index.html
│   ├── src/
│   │   ├── components/
│   │   │   ├── Header.js
│   │   │   ├── DayView.js
│   │   │   ├── MonthView.js
│   │   │   ├── MealCard.js
│   │   │   └── DishSelector.js
│   │   ├── services/
│   │   │   └── api.js
│   │   ├── App.js
│   │   ├── index.js
│   │   └── index.css
│   ├── package.json
│   ├── tailwind.config.js
│   └── postcss.config.js
└── README.md
```

## 🍛 Available Dishes

The app includes a variety of authentic Indian dishes:

- **North Indian**: Butter Chicken, Dal Tadka, Palak Paneer, Samosa
- **South Indian**: Masala Dosa, Fish Curry
- **Punjabi**: Tandoori Chicken
- **Hyderabadi**: Biryani

Each dish includes:
- High-quality food images from Unsplash
- Nutritional information (calories)
- Cuisine type and vegetarian/non-vegetarian classification
- Main ingredients list

## 🎨 Design Features

- **Glass Morphism**: Modern glassmorphism effects throughout the UI
- **Gradient Backgrounds**: Beautiful gradient backgrounds and button effects
- **Responsive Design**: Mobile-first design that works on all screen sizes
- **Color-coded Meals**: Different colors for breakfast, lunch, dinner, and snacks
- **Smooth Animations**: Hover effects and transitions for better UX
- **Indian-inspired Theme**: Warm colors inspired by Indian spices and culture

---

Made with ❤️ and 🌶️ for Indian food lovers!
