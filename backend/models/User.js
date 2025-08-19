const mongoose = require('mongoose');
const bcrypt = require('bcryptjs');

const userSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
    trim: true,
    maxlength: 50
  },
  email: {
    type: String,
    required: true,
    unique: true,
    trim: true,
    lowercase: true,
    match: [/^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/, 'Please enter a valid email']
  },
  password: {
    type: String,
    required: true,
    minlength: 6
  },
  profile: {
    dietaryPreferences: {
      type: [String],
      enum: ['vegetarian', 'vegan', 'gluten-free', 'dairy-free', 'nut-free', 'keto', 'paleo', 'low-carb', 'high-protein', 'low-sodium', 'sugar-free'],
      default: []
    },
    spiceLevel: {
      type: String,
      enum: ['mild', 'medium', 'hot', 'extra-hot'],
      default: 'medium'
    },
    favoriteRegions: {
      type: [String],
      enum: ['North Indian', 'South Indian', 'Bengali', 'Gujarati', 'Punjabi', 'Rajasthani', 'Maharashtrian', 'Italian', 'Chinese', 'Thai', 'French', 'Korean', 'Continental'],
      default: []
    },
    avatar: {
      type: String,
      default: null
    },
    nutritionGoals: {
      dailyCalories: { type: Number, default: 2000, min: 1000, max: 5000 },
      protein: { type: Number, default: 150, min: 0 }, // grams
      carbs: { type: Number, default: 250, min: 0 }, // grams
      fat: { type: Number, default: 65, min: 0 }, // grams
      fiber: { type: Number, default: 25, min: 0 }, // grams
      sodium: { type: Number, default: 2300, min: 0 } // milligrams
    },
    activityLevel: {
      type: String,
      enum: ['sedentary', 'lightly-active', 'moderately-active', 'very-active', 'extremely-active'],
      default: 'moderately-active'
    },
    healthGoals: {
      type: [String],
      enum: ['weight-loss', 'weight-gain', 'muscle-building', 'maintenance', 'heart-healthy', 'diabetes-friendly'],
      default: ['maintenance']
    }
  },
  favoriteDishes: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Dish'
  }],
  isEmailVerified: {
    type: Boolean,
    default: false
  },
  lastLoginAt: {
    type: Date,
    default: null
  }
}, {
  timestamps: true
});

// Hash password before saving
userSchema.pre('save', async function(next) {
  if (!this.isModified('password')) return next();
  
  try {
    const salt = await bcrypt.genSalt(12);
    this.password = await bcrypt.hash(this.password, salt);
    next();
  } catch (error) {
    next(error);
  }
});

// Compare password method
userSchema.methods.comparePassword = async function(candidatePassword) {
  return bcrypt.compare(candidatePassword, this.password);
};

// Remove password from JSON output
userSchema.methods.toJSON = function() {
  const userObject = this.toObject();
  delete userObject.password;
  return userObject;
};

module.exports = mongoose.model('User', userSchema);
