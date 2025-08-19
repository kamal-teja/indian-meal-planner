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
  }
}, {
  timestamps: true
});

// Create compound index for date and mealType for faster queries
mealSchema.index({ date: 1, mealType: 1 });

module.exports = mongoose.model('Meal', mealSchema);
