const mongoose = require('mongoose');

const connectDB = async () => {
  try {
    // Enhanced connection options for better reliability
    const conn = await mongoose.connect(process.env.MONGODB_URI || 'mongodb://localhost:27017/meal-planner', {
      maxPoolSize: 10, // Maintain up to 10 socket connections
      serverSelectionTimeoutMS: 5000, // Keep trying to send operations for 5 seconds
      socketTimeoutMS: 45000, // Close sockets after 45 seconds of inactivity
      bufferCommands: false // Disable mongoose buffering
    });

    console.log(`✅ MongoDB Connected: ${conn.connection.host}`);
    console.log(`📊 Database: ${conn.connection.name}`);
    
    // Handle connection events
    mongoose.connection.on('error', (err) => {
      console.error('❌ MongoDB connection error:', err);
    });
    
    mongoose.connection.on('disconnected', () => {
      console.warn('⚠️ MongoDB disconnected');
    });
    
    mongoose.connection.on('reconnected', () => {
      console.log('✅ MongoDB reconnected');
    });

    return conn;
  } catch (error) {
    console.error('❌ Error connecting to MongoDB:', error.message);
    
    // Provide helpful error messages
    if (error.message.includes('ENOTFOUND')) {
      console.error('🔍 DNS resolution failed. Check your MongoDB URI.');
    } else if (error.message.includes('authentication')) {
      console.error('🔑 Authentication failed. Check your username and password.');
    } else if (error.message.includes('timeout')) {
      console.error('⏰ Connection timeout. Check your network and MongoDB Atlas whitelist.');
    }
    
    throw error; // Re-throw to be handled by startServer
  }
};

module.exports = connectDB;
