const mongoose = require('mongoose');

const mealSchema = new mongoose.Schema({
  date: {
    type: String,
    required: true,
    match: /^\d{4}-\d{2}-\d{2}$/ // YYYY-MM-DD format
  },
  mealType: {
    type: String,
    required: true,
    enum: ['breakfast', 'lunch', 'dinner', 'snack']
  },
  dish: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Dish',
    required: true
  },
  user: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
    required: true
  },
  notes: {
    type: String,
    maxlength: 500,
    default: ''
  },
  rating: {
    type: Number,
    min: 1,
    max: 5,
    default: null
  }
}, {
  timestamps: true
});

// Create compound indexes for faster queries
mealSchema.index({ date: 1, mealType: 1 });
mealSchema.index({ user: 1, date: 1 });
mealSchema.index({ user: 1, createdAt: -1 });

module.exports = mongoose.model('Meal', mealSchema);
