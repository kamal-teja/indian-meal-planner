require('dotenv').config();
const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
const connectDB = require('./config/database');
const Dish = require('./models/Dish');
const Meal = require('./models/Meal');
const seedDishes = require('./scripts/seedDishes');

const app = express();
const PORT = process.env.PORT || 5000;

// Connect to MongoDB and start server only after successful connection
const startServer = async () => {
  try {
    // Connect to database first
    await connectDB();
    console.log('✅ Database connected successfully');
    
    // Seed default dishes if database is empty
    await seedDishes();
    console.log('✅ Database seeded with default dishes');
    
    // Start the server only after database is ready
    app.listen(PORT, () => {
      console.log(`🚀 Server running on port ${PORT}`);
      console.log(`🌐 API Health Check: http://localhost:${PORT}/api/health`);
    });
    
  } catch (error) {
    console.error('❌ Failed to start server:', error.message);
    process.exit(1);
  }
};

// CORS configuration for production
const corsOptions = {
  origin: process.env.NODE_ENV === 'production' 
    ? [process.env.FRONTEND_URL, /\.onrender\.com$/, /\.vercel\.app$/, /\.netlify\.app$/]
    : 'http://localhost:3000',
  credentials: true,
  optionsSuccessStatus: 200
};

// Middleware
app.use(cors(corsOptions));
app.use(bodyParser.json());

// Helper function to transform MongoDB document to include id field
const transformDish = (dish) => ({
  id: dish._id.toString(),
  name: dish.name,
  type: dish.type,
  cuisine: dish.cuisine,
  image: dish.image,
  ingredients: dish.ingredients,
  calories: dish.calories
});

const transformMeal = (meal) => ({
  id: meal._id.toString(),
  date: meal.date,
  mealType: meal.mealType,
  dish: transformDish(meal.dish),
  createdAt: meal.createdAt
});

// Routes

// Get all dishes
app.get('/api/dishes', async (req, res) => {
  try {
    const dishes = await Dish.find().sort({ createdAt: -1 });
    res.json(dishes.map(transformDish));
  } catch (error) {
    console.error('Error fetching dishes:', error);
    res.status(500).json({ error: 'Failed to fetch dishes' });
  }
});

// Add a new dish
app.post('/api/dishes', async (req, res) => {
  try {
    const { name, type, cuisine, ingredients, calories, image } = req.body;
    
    // Validation
    if (!name || !type || !cuisine || !ingredients || !calories) {
      return res.status(400).json({ error: 'Missing required fields: name, type, cuisine, ingredients, calories' });
    }
    
    if (!Array.isArray(ingredients) || ingredients.length === 0) {
      return res.status(400).json({ error: 'Ingredients must be a non-empty array' });
    }
    
    if (typeof calories !== 'number' || calories <= 0) {
      return res.status(400).json({ error: 'Calories must be a positive number' });
    }
    
    const newDish = new Dish({
      name: name.trim(),
      type,
      cuisine: cuisine.trim(),
      image: image || `https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=500&h=300&fit=crop`,
      ingredients: ingredients.map(ing => ing.trim()).filter(ing => ing.length > 0),
      calories: parseInt(calories)
    });
    
    const savedDish = await newDish.save();
    res.status(201).json(transformDish(savedDish));
  } catch (error) {
    console.error('Error creating dish:', error);
    res.status(500).json({ error: 'Failed to create dish' });
  }
});

// Get meals for a specific date
app.get('/api/meals/:date', async (req, res) => {
  try {
    const { date } = req.params;
    const meals = await Meal.find({ date }).populate('dish').sort({ createdAt: 1 });
    res.json(meals.map(transformMeal));
  } catch (error) {
    console.error('Error fetching meals:', error);
    res.status(500).json({ error: 'Failed to fetch meals' });
  }
});

// Get meals for a month
app.get('/api/meals/month/:year/:month', async (req, res) => {
  try {
    const { year, month } = req.params;
    const startDate = `${year}-${month.padStart(2, '0')}-01`;
    const endDate = `${year}-${month.padStart(2, '0')}-31`;
    
    const meals = await Meal.find({
      date: { $gte: startDate, $lte: endDate }
    }).populate('dish').sort({ date: 1, createdAt: 1 });
    
    res.json(meals.map(transformMeal));
  } catch (error) {
    console.error('Error fetching monthly meals:', error);
    res.status(500).json({ error: 'Failed to fetch monthly meals' });
  }
});

// Add a meal
app.post('/api/meals', async (req, res) => {
  try {
    const { date, mealType, dishId } = req.body;
    
    const dish = await Dish.findById(dishId);
    if (!dish) {
      return res.status(404).json({ error: 'Dish not found' });
    }

    const newMeal = new Meal({
      date,
      mealType,
      dish: dishId
    });

    const savedMeal = await newMeal.save();
    const populatedMeal = await Meal.findById(savedMeal._id).populate('dish');
    res.status(201).json(transformMeal(populatedMeal));
  } catch (error) {
    console.error('Error creating meal:', error);
    res.status(500).json({ error: 'Failed to create meal' });
  }
});

// Update a meal
app.put('/api/meals/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const { dishId } = req.body;
    
    const dish = await Dish.findById(dishId);
    if (!dish) {
      return res.status(404).json({ error: 'Dish not found' });
    }

    const updatedMeal = await Meal.findByIdAndUpdate(
      id,
      { dish: dishId },
      { new: true }
    ).populate('dish');

    if (!updatedMeal) {
      return res.status(404).json({ error: 'Meal not found' });
    }

    res.json(transformMeal(updatedMeal));
  } catch (error) {
    console.error('Error updating meal:', error);
    res.status(500).json({ error: 'Failed to update meal' });
  }
});

// Delete a meal
app.delete('/api/meals/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const deletedMeal = await Meal.findByIdAndDelete(id);
    
    if (!deletedMeal) {
      return res.status(404).json({ error: 'Meal not found' });
    }

    res.json({ message: 'Meal deleted successfully' });
  } catch (error) {
    console.error('Error deleting meal:', error);
    res.status(500).json({ error: 'Failed to delete meal' });
  }
});

// Health check
app.get('/api/health', (req, res) => {
  res.json({ 
    status: 'OK', 
    message: 'Meal Planner API is running',
    database: 'Connected',
    timestamp: new Date().toISOString()
  });
});

// Start the server
startServer();
