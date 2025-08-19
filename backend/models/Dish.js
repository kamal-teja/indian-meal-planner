const mongoose = require('mongoose');

const dishSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
    trim: true
  },
  type: {
    type: String,
    required: true,
    enum: ['Veg', 'Non-Veg']
  },
  cuisine: {
    type: String,
    required: true,
    trim: true
  },
  image: {
    type: String,
    default: 'https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=500&h=300&fit=crop'
  },
  ingredients: [{
    type: String,
    required: true,
    trim: true
  }],
  calories: {
    type: Number,
    required: true,
    min: 0
  },
  nutrition: {
    protein: { type: Number, default: 0, min: 0 }, // grams
    carbs: { type: Number, default: 0, min: 0 }, // grams
    fat: { type: Number, default: 0, min: 0 }, // grams
    fiber: { type: Number, default: 0, min: 0 }, // grams
    sugar: { type: Number, default: 0, min: 0 }, // grams
    sodium: { type: Number, default: 0, min: 0 } // milligrams
  },
  dietaryTags: [{
    type: String,
    enum: ['vegetarian', 'vegan', 'gluten-free', 'dairy-free', 'nut-free', 'keto', 'paleo', 'low-carb', 'high-protein', 'low-sodium', 'sugar-free']
  }],
  spiceLevel: {
    type: String,
    enum: ['mild', 'medium', 'hot', 'extra-hot'],
    default: 'medium'
  },
  prepTime: {
    type: Number, // minutes
    default: 30,
    min: 0
  },
  cookTime: {
    type: Number, // minutes
    default: 30,
    min: 0
  },
  servings: {
    type: Number,
    default: 2,
    min: 1
  },
  difficulty: {
    type: String,
    enum: ['easy', 'medium', 'hard'],
    default: 'medium'
  },
  instructions: [{
    step: { type: Number, required: true },
    description: { type: String, required: true }
  }],
  isDefault: {
    type: Boolean,
    default: false
  }
}, {
  timestamps: true
});

module.exports = mongoose.model('Dish', dishSchema);
