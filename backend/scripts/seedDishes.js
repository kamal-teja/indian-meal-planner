require('dotenv').config();
const mongoose = require('mongoose');
const Dish = require('../models/Dish');

const defaultDishes = [
  {
    name: "Butter Chicken",
    type: "Non-Veg",
    cuisine: "North Indian",
    image: "https://images.unsplash.com/photo-1603894584373-5ac82b2ae398?w=500&h=300&fit=crop",
    ingredients: ["Chicken", "Butter", "Tomatoes", "Cream", "Spices"],
    calories: 450,
    isDefault: true
  },
  {
    name: "Dal Tadka",
    type: "Veg",
    cuisine: "North Indian", 
    image: "https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=500&h=300&fit=crop",
    ingredients: ["Yellow Lentils", "Onions", "Tomatoes", "Spices"],
    calories: 200,
    isDefault: true
  },
  {
    name: "Biryani",
    type: "Non-Veg",
    cuisine: "Hyderabadi",
    image: "https://images.unsplash.com/photo-1563379091339-03246963d7d3?w=500&h=300&fit=crop",
    ingredients: ["Basmati Rice", "Chicken", "Yogurt", "Saffron", "Spices"],
    calories: 550,
    isDefault: true
  },
  {
    name: "Masala Dosa",
    type: "Veg",
    cuisine: "South Indian",
    image: "https://images.unsplash.com/photo-1567188040759-fb8a883dc6d8?w=500&h=300&fit=crop",
    ingredients: ["Rice Batter", "Potato", "Onions", "Curry Leaves"],
    calories: 300,
    isDefault: true
  },
  {
    name: "Palak Paneer",
    type: "Veg",
    cuisine: "North Indian",
    image: "https://images.unsplash.com/photo-1585937421612-70a008356fbe?w=500&h=300&fit=crop",
    ingredients: ["Spinach", "Paneer", "Cream", "Spices"],
    calories: 280,
    isDefault: true
  },
  {
    name: "Tandoori Chicken",
    type: "Non-Veg",
    cuisine: "Punjabi",
    image: "https://images.unsplash.com/photo-1599487488170-d11ec9c172f0?w=500&h=300&fit=crop",
    ingredients: ["Chicken", "Yogurt", "Tandoori Masala", "Lemon"],
    calories: 350,
    isDefault: true
  },
  {
    name: "Samosa",
    type: "Veg",
    cuisine: "North Indian",
    image: "https://images.unsplash.com/photo-1601050690597-df0568f70950?w=500&h=300&fit=crop",
    ingredients: ["Flour", "Potato", "Peas", "Spices"],
    calories: 150,
    isDefault: true
  },
  {
    name: "Fish Curry",
    type: "Non-Veg",
    cuisine: "South Indian",
    image: "https://images.unsplash.com/photo-1585937421612-70a008356fbe?w=500&h=300&fit=crop",
    ingredients: ["Fish", "Coconut", "Curry Leaves", "Spices"],
    calories: 400,
    isDefault: true
  }
];

const seedDishes = async () => {
  try {
    await mongoose.connect(process.env.MONGODB_URI || 'mongodb://localhost:27017/meal-planner');
    console.log('Connected to MongoDB');

    // Check if dishes already exist
    const existingDishes = await Dish.countDocuments({ isDefault: true });
    
    if (existingDishes === 0) {
      await Dish.insertMany(defaultDishes);
      console.log('Default dishes seeded successfully');
    } else {
      console.log('Default dishes already exist, skipping seed');
    }

    await mongoose.connection.close();
    console.log('Database connection closed');
  } catch (error) {
    console.error('Error seeding dishes:', error);
    process.exit(1);
  }
};

// Run if called directly
if (require.main === module) {
  seedDishes();
}

module.exports = seedDishes;
