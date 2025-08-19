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
  isDefault: {
    type: Boolean,
    default: false
  }
}, {
  timestamps: true
});

module.exports = mongoose.model('Dish', dishSchema);
