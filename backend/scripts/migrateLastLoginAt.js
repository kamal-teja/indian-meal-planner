require('dotenv').config();
const mongoose = require('mongoose');
const User = require('../models/User');

// Connect to MongoDB
const connectDB = async () => {
  try {
    const conn = await mongoose.connect(process.env.MONGODB_URI || 'mongodb://localhost:27017/meal-planner', {
      maxPoolSize: 10,
      serverSelectionTimeoutMS: 5000,
      socketTimeoutMS: 45000,
      bufferCommands: false
    });
    console.log('✅ MongoDB Connected for migration:', conn.connection.host);
    console.log('📊 Database:', conn.connection.name);
  } catch (error) {
    console.error('❌ Database connection failed:', error);
    process.exit(1);
  }
};

const updateLastLoginAt = async () => {
  try {
    await connectDB();
    
    // Update all users where lastLoginAt is null to set it to their createdAt date
    // This gives a reasonable default for existing users
    const result = await User.updateMany(
      { lastLoginAt: { $in: [null, undefined] } },
      { $set: { lastLoginAt: null } } // Keep as null, will be updated on next login
    );
    
    console.log(`✅ Migration completed. ${result.modifiedCount} users processed.`);
    console.log('ℹ️  Note: lastLoginAt will be set to the actual login time when users log in next.');
    
    // Alternatively, if you want to set it to createdAt for existing users:
    // const result = await User.updateMany(
    //   { lastLoginAt: { $in: [null, undefined] } },
    //   [{ $set: { lastLoginAt: "$createdAt" } }]
    // );
    
    process.exit(0);
  } catch (error) {
    console.error('❌ Migration failed:', error);
    process.exit(1);
  }
};

// Run migration if this file is executed directly
if (require.main === module) {
  updateLastLoginAt();
}

module.exports = { updateLastLoginAt };
