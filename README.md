# 🍽️ Indian Meal Planner App

A beautiful, modern meal planning application focused on Indian cuisine with shopping list functionality and MongoDB persistence.

![Meal Planner Demo](https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=800&h=400&fit=crop)

## ✨ Features

### 📱 Core Functionality
- **Day View**: Plan meals for specific dates with beautiful meal cards
- **Month View**: Overview of meal plans across the month
- **Dish Management**: Browse, search, and filter from 8+ default Indian dishes
- **Custom Dishes**: Add your own recipes with ingredients and nutritional info
- **Shopping List**: Auto-generated ingredients list with progress tracking

### 🎨 User Experience
- **Modern UI**: Tailwind CSS with gradient backgrounds and smooth animations
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile
- **Intuitive Navigation**: Easy meal planning with drag-and-drop feel
- **Visual Feedback**: Loading states, hover effects, and micro-interactions

### 🗄️ Data Persistence
- **MongoDB Atlas**: Cloud database with free tier
- **Real-time Updates**: Changes saved instantly
- **Default Dishes**: Pre-loaded with popular Indian dishes
- **Custom Recipes**: Your dishes persist across sessions

## 🚀 Quick Start

### Option 1: Local Development

1. **Clone the repository**:
```bash
git clone <your-repo-url>
cd meal-planner-app
```

2. **Setup MongoDB** (Choose one):
   - **Option A**: Use MongoDB Atlas (recommended)
     - Create account at [MongoDB Atlas](https://www.mongodb.com/atlas)
     - Create a free cluster
     - Get connection string
   - **Option B**: Local MongoDB
     - Install MongoDB locally
     - Use connection string: `mongodb://localhost:27017/meal-planner`

3. **Backend Setup**:
```bash
cd backend
npm install
cp env.example .env
# Edit .env and add your MongoDB URI
npm run dev
```

4. **Frontend Setup** (new terminal):
```bash
cd frontend
npm install
npm start
```

5. **Access the app**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:5000
   - Health Check: http://localhost:5000/api/health

### Option 2: One-Command Development (Linux/Mac)
```bash
./start-dev.sh
```

## 🌐 Deploy for FREE

Follow our comprehensive [Deployment Guide](DEPLOYMENT.md) to deploy to:
- **Render** (Recommended - Backend + Frontend)
- **Vercel** (Frontend) + **Railway** (Backend)  
- **Netlify** (Frontend) + **Heroku** (Backend)

All options include free tiers perfect for personal use!

## 📁 Project Structure

```
meal-planner-app/
├── frontend/                 # React.js application
│   ├── src/
│   │   ├── components/      # React components
│   │   │   ├── DayView.js   # Daily meal planning
│   │   │   ├── MonthView.js # Monthly overview
│   │   │   ├── DishSelector.js # Dish selection modal
│   │   │   ├── AddDishForm.js  # Add custom dishes
│   │   │   ├── IngredientsList.js # Shopping list
│   │   │   └── ...
│   │   ├── services/        # API services
│   │   └── styles/          # Tailwind CSS
│   └── package.json
├── backend/                  # Node.js + Express API
│   ├── models/              # MongoDB models
│   │   ├── Dish.js         # Dish schema
│   │   └── Meal.js         # Meal schema
│   ├── config/             # Database configuration
│   ├── scripts/            # Utility scripts
│   └── server.js           # Main server file
├── DEPLOYMENT.md           # Deployment guide
└── README.md              # This file
```

## 🛠️ Technology Stack

### Frontend
- **React 18** - Modern React with hooks
- **React Router** - Client-side routing
- **Tailwind CSS** - Utility-first styling
- **Lucide React** - Beautiful icons
- **Date-fns** - Date manipulation
- **Axios** - HTTP client

### Backend
- **Node.js** - Runtime environment
- **Express.js** - Web framework
- **MongoDB** - Database
- **Mongoose** - MongoDB ODM
- **CORS** - Cross-origin resource sharing
- **dotenv** - Environment variables

### DevOps & Deployment
- **MongoDB Atlas** - Cloud database
- **Render** - Hosting platform
- **Vercel** - Static site hosting
- **Docker** - Containerization
- **GitHub Actions** - CI/CD (optional)

## 📊 API Endpoints

### Dishes
- `GET /api/dishes` - Get all dishes
- `POST /api/dishes` - Create new dish

### Meals
- `GET /api/meals/:date` - Get meals for specific date
- `GET /api/meals/month/:year/:month` - Get monthly meals
- `POST /api/meals` - Create meal plan
- `PUT /api/meals/:id` - Update meal
- `DELETE /api/meals/:id` - Delete meal

### Health
- `GET /api/health` - API health check

## 🎯 Key Features Explained

### 🛒 Smart Shopping List
- **Auto-consolidation**: Combines ingredients from multiple meals
- **Progress tracking**: Check off items as you shop
- **Meal context**: See which dishes need each ingredient
- **Responsive design**: Works great on mobile while shopping

### 🍛 Dish Management
- **Default collection**: 8 popular Indian dishes included
- **Custom recipes**: Add unlimited personal dishes
- **Rich metadata**: Cuisine type, ingredients, calories, images
- **Smart search**: Find dishes by name, cuisine, or ingredient

### 📅 Meal Planning
- **Intuitive interface**: Click to add meals to any time slot
- **Visual calendar**: Month view for planning ahead
- **Calorie tracking**: See total daily calories
- **Quick editing**: Edit or delete meals with hover actions

## 🔧 Environment Variables

### Backend (.env)
```bash
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/meal-planner
NODE_ENV=development
PORT=5000
FRONTEND_URL=http://localhost:3000
```

### Frontend (.env)
```bash
REACT_APP_API_URL=http://localhost:5000/api
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Having issues? Check our [Deployment Guide](DEPLOYMENT.md) or open an issue.

## 🎉 Acknowledgments

- Default dish images from [Unsplash](https://unsplash.com)
- Icons by [Lucide](https://lucide.dev)
- Inspired by the rich culinary traditions of India

---

**Built with ❤️ for food lovers and meal planning enthusiasts!**

*Happy cooking! 🍛✨*
