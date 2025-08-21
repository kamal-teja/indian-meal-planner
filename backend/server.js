require('dotenv').config();
const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
const mongoose = require('mongoose');
const connectDB = require('./config/database');
const Dish = require('./models/Dish');
const Meal = require('./models/Meal');
const User = require('./models/User');
const MealPlan = require('./models/MealPlan');
const seedDishes = require('./scripts/seedDishes');
const authRoutes = require('./routes/auth');
const { auth, optionalAuth } = require('./middleware/auth');

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
  origin: true, // Allow all origins for development
  credentials: true,
  optionsSuccessStatus: 200,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With'],
  preflightContinue: false
};

// Middleware
app.use(cors(corsOptions));

// Handle preflight requests explicitly
app.options('*', cors(corsOptions));

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
  user: meal.user,
  notes: meal.notes || '',
  rating: meal.rating,
  createdAt: meal.createdAt
});

// Routes

// Auth routes
app.use('/api/auth', authRoutes);

// Get all dishes with pagination (with optional user context for favorites)
app.get('/api/dishes', optionalAuth, async (req, res) => {
  try {
    const page = parseInt(req.query.page) || 1;
    const limit = parseInt(req.query.limit) || 20;
    const skip = (page - 1) * limit;
    
    // Get search and filter parameters
    const { search, type, cuisine } = req.query;
    
    // Build query object
    let query = {};
    
    if (search) {
      query.$or = [
        { name: { $regex: search, $options: 'i' } },
        { cuisine: { $regex: search, $options: 'i' } },
        { ingredients: { $in: [new RegExp(search, 'i')] } }
      ];
    }
    
    if (type && type !== 'all') {
      query.type = type;
    }
    
    if (cuisine && cuisine !== 'all') {
      query.cuisine = cuisine;
    }
    
    const [dishes, totalCount] = await Promise.all([
      Dish.find(query)
        .sort({ createdAt: -1 })
        .skip(skip)
        .limit(limit),
      Dish.countDocuments(query)
    ]);
    
    const transformedDishes = dishes.map(dish => {
      const dishData = transformDish(dish);
      if (req.user) {
        dishData.isFavorite = req.user.favoriteDishes.includes(dish._id);
      }
      return dishData;
    });
    
    res.json({
      dishes: transformedDishes,
      pagination: {
        currentPage: page,
        totalPages: Math.ceil(totalCount / limit),
        totalCount,
        hasMore: page < Math.ceil(totalCount / limit)
      }
    });
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

// Get meals for a specific date (requires authentication)
app.get('/api/meals/:date', auth, async (req, res) => {
  try {
    const { date } = req.params;
    const meals = await Meal.find({ date, user: req.user._id }).populate('dish').sort({ createdAt: 1 });
    res.json(meals.map(transformMeal));
  } catch (error) {
    console.error('Error fetching meals:', error);
    res.status(500).json({ error: 'Failed to fetch meals' });
  }
});

// Get meals for a month (requires authentication)
app.get('/api/meals/month/:year/:month', auth, async (req, res) => {
  try {
    const { year, month } = req.params;
    const startDate = `${year}-${month.padStart(2, '0')}-01`;
    const endDate = `${year}-${month.padStart(2, '0')}-31`;
    
    const meals = await Meal.find({
      date: { $gte: startDate, $lte: endDate },
      user: req.user._id
    }).populate('dish').sort({ date: 1, createdAt: 1 });
    
    res.json(meals.map(transformMeal));
  } catch (error) {
    console.error('Error fetching monthly meals:', error);
    res.status(500).json({ error: 'Failed to fetch monthly meals' });
  }
});

// Add a meal (requires authentication)
app.post('/api/meals', auth, async (req, res) => {
  try {
    const { date, mealType, dishId, notes, rating } = req.body;
    
    const dish = await Dish.findById(dishId);
    if (!dish) {
      return res.status(404).json({ error: 'Dish not found' });
    }

    const newMeal = new Meal({
      date,
      mealType,
      dish: dishId,
      user: req.user._id,
      notes: notes || '',
      rating: rating || null
    });

    const savedMeal = await newMeal.save();
    const populatedMeal = await Meal.findById(savedMeal._id).populate('dish');
    res.status(201).json(transformMeal(populatedMeal));
  } catch (error) {
    console.error('Error creating meal:', error);
    res.status(500).json({ error: 'Failed to create meal' });
  }
});

// Update a meal (requires authentication)
app.put('/api/meals/:id', auth, async (req, res) => {
  try {
    const { id } = req.params;
    const { dishId, notes, rating } = req.body;
    
    // Find the meal and check ownership
    const meal = await Meal.findOne({ _id: id, user: req.user._id });
    if (!meal) {
      return res.status(404).json({ error: 'Meal not found or access denied' });
    }

    const updates = {};
    if (dishId) {
      const dish = await Dish.findById(dishId);
      if (!dish) {
        return res.status(404).json({ error: 'Dish not found' });
      }
      updates.dish = dishId;
    }
    if (notes !== undefined) updates.notes = notes;
    if (rating !== undefined) updates.rating = rating;

    const updatedMeal = await Meal.findByIdAndUpdate(
      id,
      updates,
      { new: true }
    ).populate('dish');

    res.json(transformMeal(updatedMeal));
  } catch (error) {
    console.error('Error updating meal:', error);
    res.status(500).json({ error: 'Failed to update meal' });
  }
});

// Delete a meal (requires authentication)
app.delete('/api/meals/:id', auth, async (req, res) => {
  try {
    const { id } = req.params;
    const deletedMeal = await Meal.findOneAndDelete({ _id: id, user: req.user._id });
    
    if (!deletedMeal) {
      return res.status(404).json({ error: 'Meal not found or access denied' });
    }

    res.json({ message: 'Meal deleted successfully' });
  } catch (error) {
    console.error('Error deleting meal:', error);
    res.status(500).json({ error: 'Failed to delete meal' });
  }
});

// Search dishes with advanced filters
app.get('/api/dishes/search', optionalAuth, async (req, res) => {
  try {
    const { q, cuisine, type, maxCalories, dietaryTags, spiceLevel, maxPrepTime, difficulty } = req.query;
    let query = {};

    if (q) {
      query.$or = [
        { name: { $regex: q, $options: 'i' } },
        { ingredients: { $in: [new RegExp(q, 'i')] } }
      ];
    }

    if (cuisine) query.cuisine = cuisine;
    if (type) query.type = type;
    if (maxCalories) query.calories = { $lte: parseInt(maxCalories) };

    if (dietaryTags) {
      const tags = Array.isArray(dietaryTags) ? dietaryTags : [dietaryTags];
      query.dietaryTags = { $in: tags };
    }

    if (spiceLevel) query.spiceLevel = spiceLevel;
    if (maxPrepTime) query.prepTime = { $lte: parseInt(maxPrepTime) };
    if (difficulty) query.difficulty = difficulty;

    const dishes = await Dish.find(query).sort({ createdAt: -1 });
    const transformedDishes = dishes.map(dish => {
      const dishData = transformDish(dish);
      if (req.user) {
        dishData.isFavorite = req.user.favoriteDishes.includes(dish._id);
      }
      return dishData;
    });
    
    res.json(transformedDishes);
  } catch (error) {
    console.error('Error searching dishes:', error);
    res.status(500).json({ error: 'Failed to search dishes' });
  }
});

// Get user's favorite dishes
app.get('/api/dishes/favorites', auth, async (req, res) => {
  try {
    const user = await User.findById(req.user._id).populate('favoriteDishes');
    const favoriteDishes = user.favoriteDishes.map(dish => {
      const dishData = transformDish(dish);
      dishData.isFavorite = true;
      return dishData;
    });
    res.json(favoriteDishes);
  } catch (error) {
    console.error('Error fetching favorite dishes:', error);
    res.status(500).json({ error: 'Failed to fetch favorite dishes' });
  }
});

// Get meal analytics for user
app.get('/api/analytics/meals', auth, async (req, res) => {
  try {
    const { period = '30' } = req.query; // days
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - parseInt(period));
    
    const meals = await Meal.find({
      user: req.user._id,
      createdAt: { $gte: startDate }
    }).populate('dish');

    const analytics = {
      totalMeals: meals.length,
      averageCaloriesPerDay: 0,
      mealTypeDistribution: {
        breakfast: 0,
        lunch: 0,
        dinner: 0,
        snack: 0
      },
      topCuisines: {},
      averageRating: 0,
      totalCalories: 0
    };

    let totalRatings = 0;
    let ratingCount = 0;

    meals.forEach(meal => {
      analytics.totalCalories += meal.dish.calories;
      analytics.mealTypeDistribution[meal.mealType]++;
      
      if (!analytics.topCuisines[meal.dish.cuisine]) {
        analytics.topCuisines[meal.dish.cuisine] = 0;
      }
      analytics.topCuisines[meal.dish.cuisine]++;

      if (meal.rating) {
        totalRatings += meal.rating;
        ratingCount++;
      }
    });

    analytics.averageCaloriesPerDay = Math.round(analytics.totalCalories / parseInt(period));
    analytics.averageRating = ratingCount > 0 ? (totalRatings / ratingCount).toFixed(1) : 0;

    res.json(analytics);
  } catch (error) {
    console.error('Error fetching meal analytics:', error);
    res.status(500).json({ error: 'Failed to fetch meal analytics' });
  }
});

// Generate shopping list for date range
app.get('/api/shopping-list', auth, async (req, res) => {
  try {
    const { startDate, endDate } = req.query;
    
    if (!startDate || !endDate) {
      return res.status(400).json({ error: 'Start date and end date are required' });
    }

    const meals = await Meal.find({
      user: req.user._id,
      date: { $gte: startDate, $lte: endDate }
    }).populate('dish');

    const ingredientMap = {};
    
    meals.forEach(meal => {
      meal.dish.ingredients.forEach(ingredient => {
        const key = ingredient.toLowerCase();
        if (!ingredientMap[key]) {
          ingredientMap[key] = {
            name: ingredient,
            count: 0,
            dishes: []
          };
        }
        ingredientMap[key].count++;
        if (!ingredientMap[key].dishes.includes(meal.dish.name)) {
          ingredientMap[key].dishes.push(meal.dish.name);
        }
      });
    });

    const shoppingList = Object.values(ingredientMap)
      .sort((a, b) => b.count - a.count);

    res.json({
      period: { startDate, endDate },
      totalItems: shoppingList.length,
      items: shoppingList
    });
  } catch (error) {
    console.error('Error generating shopping list:', error);
    res.status(500).json({ error: 'Failed to generate shopping list' });
  }
});

// Get meal recommendations based on user preferences
app.get('/api/recommendations', auth, async (req, res) => {
  try {
    const { mealType = 'lunch', date } = req.query;
    const user = await User.findById(req.user._id);
    
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }

    // Build base query - start with all dishes and filter progressively
    let query = {};
    
    // If user has favorite regions, prefer those cuisines but don't make it exclusive
    const favoriteRegions = user.profile?.favoriteRegions || [];
    
    // Get user's recent meals to avoid repetition
    const recentMeals = await Meal.find({ 
      user: req.user._id,
      date: { $gte: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0] }
    }).populate('dish');
    
    const recentDishIds = recentMeals.map(meal => meal.dish?._id?.toString()).filter(Boolean);
    if (recentDishIds.length > 0) {
      query._id = { $nin: recentDishIds };
    }

    // Get dishes and score them
    const dishes = await Dish.find(query).limit(100);
    
    if (dishes.length === 0) {
      // If no dishes found, get some basic dishes without the recent exclusion
      const fallbackDishes = await Dish.find({}).limit(20);
      return res.json({
        mealType,
        date,
        recommendations: fallbackDishes.map(dish => ({
          ...transformDish(dish),
          recommendationScore: 1,
          isFavorite: user.favoriteDishes?.includes(dish._id) || false
        })),
        totalFound: fallbackDishes.length
      });
    }

    // Enhanced scoring algorithm
    const scoredDishes = dishes.map(dish => {
      let score = 1; // Base score for all dishes
      
      // Favorite cuisine bonus (higher weight)
      if (favoriteRegions.length > 0 && favoriteRegions.includes(dish.cuisine)) {
        score += 5;
      }
      
      // Dietary preference match (only if dish has dietary tags)
      if (dish.dietaryTags && Array.isArray(dish.dietaryTags) && dish.dietaryTags.length > 0) {
        const userPreferences = user.profile?.dietaryPreferences || [];
        const matchingTags = dish.dietaryTags.filter(tag => 
          userPreferences.includes(tag)
        );
        score += matchingTags.length * 3;
      }
      
      // Vegetarian preference handling
      const userPreferences = user.profile?.dietaryPreferences || [];
      if (userPreferences.includes('vegetarian') && dish.type === 'Veg') {
        score += 3;
      }
      
      // Favorite dish bonus
      if (user.favoriteDishes && user.favoriteDishes.includes(dish._id)) {
        score += 10;
      }
      
      // Meal type appropriateness based on calories
      const calories = dish.calories || 0;
      if (mealType === 'breakfast' && calories >= 200 && calories <= 500) score += 2;
      else if (mealType === 'lunch' && calories >= 300 && calories <= 700) score += 2;
      else if (mealType === 'dinner' && calories >= 400 && calories <= 800) score += 2;
      else if (mealType === 'snack' && calories >= 100 && calories <= 400) score += 2;
      
      // Spice level match (only if both user and dish have spice level)
      if (user.profile?.spiceLevel && dish.spiceLevel) {
        const spiceLevels = ['mild', 'medium', 'hot', 'extra-hot'];
        const userSpiceIndex = spiceLevels.indexOf(user.profile.spiceLevel);
        const dishSpiceIndex = spiceLevels.indexOf(dish.spiceLevel);
        
        if (dishSpiceIndex <= userSpiceIndex) {
          score += 2;
        }
      }
      
      // Random factor to add variety
      score += Math.random() * 0.5;
      
      return { dish: transformDish(dish), score };
    });

    // Sort by score and return top recommendations
    const recommendations = scoredDishes
      .sort((a, b) => b.score - a.score)
      .slice(0, 12)
      .map(item => ({
        ...item.dish,
        recommendationScore: Math.round(item.score * 10) / 10,
        isFavorite: user.favoriteDishes?.includes(item.dish.id) || false
      }));

    res.json({
      mealType,
      date,
      recommendations,
      totalFound: recommendations.length,
      userPreferences: {
        dietaryPreferences: user.profile?.dietaryPreferences || [],
        favoriteRegions: user.profile?.favoriteRegions || [],
        spiceLevel: user.profile?.spiceLevel || 'medium'
      }
    });
  } catch (error) {
    console.error('Error getting recommendations:', error);
    res.status(500).json({ error: 'Failed to get recommendations' });
  }
});

// Get nutrition progress for a date range
app.get('/api/nutrition/progress', auth, async (req, res) => {
  try {
    const { startDate, endDate } = req.query;
    const user = await User.findById(req.user._id);
    
    const meals = await Meal.find({
      user: req.user._id,
      date: { $gte: startDate, $lte: endDate }
    }).populate('dish');

    // Calculate daily nutrition
    const dailyNutrition = {};
    
    meals.forEach(meal => {
      const date = meal.date;
      if (!dailyNutrition[date]) {
        dailyNutrition[date] = {
          calories: 0,
          protein: 0,
          carbs: 0,
          fat: 0,
          fiber: 0,
          sodium: 0,
          meals: []
        };
      }
      
      const dish = meal.dish;
      dailyNutrition[date].calories += dish.calories;
      dailyNutrition[date].protein += dish.nutrition?.protein || 0;
      dailyNutrition[date].carbs += dish.nutrition?.carbs || 0;
      dailyNutrition[date].fat += dish.nutrition?.fat || 0;
      dailyNutrition[date].fiber += dish.nutrition?.fiber || 0;
      dailyNutrition[date].sodium += dish.nutrition?.sodium || 0;
      dailyNutrition[date].meals.push({
        mealType: meal.mealType,
        dish: transformDish(dish),
        rating: meal.rating
      });
    });

    // Calculate progress against goals
    const goals = user.profile.nutritionGoals;
    const progressData = Object.entries(dailyNutrition).map(([date, nutrition]) => ({
      date,
      nutrition,
      progress: {
        calories: (nutrition.calories / goals.dailyCalories) * 100,
        protein: (nutrition.protein / goals.protein) * 100,
        carbs: (nutrition.carbs / goals.carbs) * 100,
        fat: (nutrition.fat / goals.fat) * 100,
        fiber: (nutrition.fiber / goals.fiber) * 100,
        sodium: Math.min((nutrition.sodium / goals.sodium) * 100, 100) // Cap at 100% for sodium
      },
      goalsStatus: {
        caloriesOnTrack: Math.abs(nutrition.calories - goals.dailyCalories) <= goals.dailyCalories * 0.1,
        proteinMet: nutrition.protein >= goals.protein * 0.9,
        sodiumOk: nutrition.sodium <= goals.sodium
      }
    }));

    res.json({
      period: { startDate, endDate },
      goals,
      dailyData: progressData,
      summary: {
        avgCalories: progressData.reduce((sum, day) => sum + day.nutrition.calories, 0) / progressData.length,
        avgProtein: progressData.reduce((sum, day) => sum + day.nutrition.protein, 0) / progressData.length,
        daysOnTrack: progressData.filter(day => day.goalsStatus.caloriesOnTrack).length
      }
    });
  } catch (error) {
    console.error('Error getting nutrition progress:', error);
    res.status(500).json({ error: 'Failed to get nutrition progress' });
  }
});

// Update user nutrition goals
app.put('/api/profile/nutrition-goals', auth, async (req, res) => {
  try {
    const { dailyCalories, protein, carbs, fat, fiber, sodium } = req.body;
    
    const updateData = {};
    if (dailyCalories) updateData['profile.nutritionGoals.dailyCalories'] = dailyCalories;
    if (protein) updateData['profile.nutritionGoals.protein'] = protein;
    if (carbs) updateData['profile.nutritionGoals.carbs'] = carbs;
    if (fat) updateData['profile.nutritionGoals.fat'] = fat;
    if (fiber) updateData['profile.nutritionGoals.fiber'] = fiber;
    if (sodium) updateData['profile.nutritionGoals.sodium'] = sodium;

    const user = await User.findByIdAndUpdate(
      req.user._id,
      { $set: updateData },
      { new: true }
    ).select('-password');

    res.json({
      message: 'Nutrition goals updated successfully',
      nutritionGoals: user.profile.nutritionGoals
    });
  } catch (error) {
    console.error('Error updating nutrition goals:', error);
    res.status(500).json({ error: 'Failed to update nutrition goals' });
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
