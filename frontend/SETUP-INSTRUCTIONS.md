# 🚀 Setup Instructions for Indian Meal Planner App

## Quick Start Guide

Follow these steps to get your Indian Meal Planner app up and running:

### Step 1: Install Dependencies

**Backend Setup:**
```bash
cd backend
npm install
```

**Frontend Setup:**
```bash
cd frontend
npm install
```

### Step 2: Start the Applications

**Start Backend Server (Terminal 1):**
```bash
cd backend
npm run dev
```
✅ Backend will run on: http://localhost:5000

**Start Frontend App (Terminal 2):**
```bash
cd frontend
npm start
```
✅ Frontend will open in browser at: http://localhost:3000

### Step 3: Enjoy Your Meal Planner!

🎉 Your app is now ready! You can:
- Plan daily meals in Day View
- View monthly meal calendar in Month View
- Browse and select from 8+ Indian dishes
- Track calories and nutrition
- Add, edit, and delete meals

## 📦 What's Included

### Backend Features:
- REST API with Express.js
- In-memory meal storage
- Pre-loaded Indian dishes database
- CORS enabled for frontend communication

### Frontend Features:
- React 18 with modern hooks
- Tailwind CSS for beautiful styling
- Day and Month views
- Responsive design
- Indian food images from Unsplash

## 🛠️ Development Scripts

### Backend Commands:
```bash
npm start          # Start production server
npm run dev        # Start development server with nodemon
```

### Frontend Commands:
```bash
npm start          # Start development server
npm run build      # Build for production
npm test           # Run tests
```

## 🌟 Features to Explore

1. **Day View**: Click on the calendar icon to switch to day view
2. **Add Meals**: Click the "+" button on any meal type
3. **Browse Dishes**: Search and filter by cuisine type and vegetarian/non-veg
4. **Month View**: See your entire month's meal plan at a glance
5. **Edit/Delete**: Hover over meals to edit or delete them

## 📱 Responsive Design

The app works beautifully on:
- 📱 Mobile phones
- 📱 Tablets
- 💻 Desktop computers

## 🍛 Sample Dishes Available

- Butter Chicken (Non-Veg, North Indian)
- Dal Tadka (Veg, North Indian)
- Biryani (Non-Veg, Hyderabadi)
- Masala Dosa (Veg, South Indian)
- Palak Paneer (Veg, North Indian)
- Tandoori Chicken (Non-Veg, Punjabi)
- Samosa (Veg, North Indian)
- Fish Curry (Non-Veg, South Indian)

## ❓ Troubleshooting

**Port already in use?**
- Backend: Change PORT in backend/server.js
- Frontend: It will prompt to use a different port

**Images not loading?**
- Images are loaded from Unsplash CDN
- Check your internet connection

**API not connecting?**
- Make sure backend is running on port 5000
- Check frontend/package.json proxy setting

## 🎨 Customization

Want to add your own dishes? Edit the `indianDishes` array in `backend/server.js`:

```javascript
{
  id: 9,
  name: "Your Dish Name",
  type: "Veg", // or "Non-Veg"
  cuisine: "Your Cuisine",
  image: "https://your-image-url.com",
  ingredients: ["ingredient1", "ingredient2"],
  calories: 300
}
```

Happy meal planning! 🍽️✨
