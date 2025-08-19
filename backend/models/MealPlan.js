const mongoose = require('mongoose');

const mealPlanSchema = new mongoose.Schema({
  user: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
    required: true
  },
  name: {
    type: String,
    required: true,
    trim: true,
    maxlength: 100
  },
  description: {
    type: String,
    maxlength: 500,
    default: ''
  },
  startDate: {
    type: Date,
    required: true
  },
  endDate: {
    type: Date,
    required: true
  },
  meals: [{
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
    notes: {
      type: String,
      maxlength: 200,
      default: ''
    }
  }],
  isTemplate: {
    type: Boolean,
    default: false
  },
  isPublic: {
    type: Boolean,
    default: false
  },
  tags: [{
    type: String,
    trim: true
  }],
  nutritionSummary: {
    totalCalories: { type: Number, default: 0 },
    avgDailyCalories: { type: Number, default: 0 },
    totalProtein: { type: Number, default: 0 },
    totalCarbs: { type: Number, default: 0 },
    totalFat: { type: Number, default: 0 }
  },
  likes: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User'
  }],
  createdBy: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
    required: true
  }
}, {
  timestamps: true
});

// Indexes for better query performance
mealPlanSchema.index({ user: 1, startDate: 1 });
mealPlanSchema.index({ isPublic: 1, createdAt: -1 });
mealPlanSchema.index({ isTemplate: 1, isPublic: 1 });
mealPlanSchema.index({ tags: 1 });

module.exports = mongoose.model('MealPlan', mealPlanSchema);
