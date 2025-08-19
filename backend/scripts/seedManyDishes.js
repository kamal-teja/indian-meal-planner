require('dotenv').config();
const mongoose = require('mongoose');
const Dish = require('../models/Dish');

const connectDB = async () => {
  try {
    await mongoose.connect(process.env.MONGODB_URI);
    console.log('MongoDB connected for seeding');
  } catch (error) {
    console.error('MongoDB connection error:', error);
    process.exit(1);
  }
};

const getRandomCalories = (min = 150, max = 800) => {
  return Math.floor(Math.random() * (max - min + 1)) + min;
};

const getRandomImage = (cuisine) => {
  const imageMap = {
    'Indian': [
      'https://images.unsplash.com/photo-1565557623262-b51c2513a641?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1574085733277-851d9d856a3a?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1596797038530-2c107229654b?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1585937421612-70a008356fbe?w=500&h=300&fit=crop'
    ],
    'Chinese': [
      'https://images.unsplash.com/photo-1563379091339-03246963d96a?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1568096889942-6eedde686635?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1496116218417-1a781b1c416c?w=500&h=300&fit=crop'
    ],
    'Italian': [
      'https://images.unsplash.com/photo-1551782450-17144efb9c50?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1540189549336-e6e99c3679fe?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1565299624946-b28f40a0ca4b?w=500&h=300&fit=crop'
    ],
    'Thai': [
      'https://images.unsplash.com/photo-1559847844-5315695dadae?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1573821663912-6df460f9c684?w=500&h=300&fit=crop'
    ],
    'Korean': [
      'https://images.unsplash.com/photo-1590301157890-4810ed352733?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1598511757337-fe2cafc31ba1?w=500&h=300&fit=crop'
    ],
    'French': [
      'https://images.unsplash.com/photo-1600891964092-4316c288032e?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1547592180-85f173990554?w=500&h=300&fit=crop'
    ],
    'Continental': [
      'https://images.unsplash.com/photo-1567620905732-2d1ec7ab7445?w=500&h=300&fit=crop',
      'https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=500&h=300&fit=crop'
    ]
  };
  
  const images = imageMap[cuisine] || imageMap['Continental'];
  return images[Math.floor(Math.random() * images.length)];
};

const dishes = [
  // Indian Dishes (150 dishes)
  { name: 'Butter Chicken', type: 'Non-Veg', cuisine: 'North Indian', ingredients: ['Chicken', 'Tomatoes', 'Cream', 'Butter', 'Garam Masala', 'Ginger', 'Garlic'] },
  { name: 'Chicken Biryani', type: 'Non-Veg', cuisine: 'North Indian', ingredients: ['Basmati Rice', 'Chicken', 'Saffron', 'Yogurt', 'Onions', 'Spices'] },
  { name: 'Mutton Biryani', type: 'Non-Veg', cuisine: 'North Indian', ingredients: ['Basmati Rice', 'Mutton', 'Saffron', 'Yogurt', 'Onions', 'Spices'] },
  { name: 'Paneer Makhani', type: 'Veg', cuisine: 'North Indian', ingredients: ['Paneer', 'Tomatoes', 'Cream', 'Cashews', 'Spices'] },
  { name: 'Dal Tadka', type: 'Veg', cuisine: 'North Indian', ingredients: ['Yellow Lentils', 'Onions', 'Tomatoes', 'Cumin', 'Turmeric', 'Cilantro'] },
  { name: 'Dal Makhani', type: 'Veg', cuisine: 'North Indian', ingredients: ['Black Lentils', 'Cream', 'Butter', 'Tomatoes', 'Spices'] },
  { name: 'Rogan Josh', type: 'Non-Veg', cuisine: 'North Indian', ingredients: ['Mutton', 'Yogurt', 'Fennel', 'Ginger', 'Garlic', 'Red Chili'] },
  { name: 'Chole Bhature', type: 'Veg', cuisine: 'North Indian', ingredients: ['Chickpeas', 'Flour', 'Yogurt', 'Spices', 'Oil'] },
  { name: 'Rajma', type: 'Veg', cuisine: 'North Indian', ingredients: ['Kidney Beans', 'Tomatoes', 'Onions', 'Ginger', 'Garlic', 'Spices'] },
  { name: 'Palak Paneer', type: 'Veg', cuisine: 'North Indian', ingredients: ['Spinach', 'Paneer', 'Onions', 'Tomatoes', 'Cream', 'Spices'] },
  { name: 'Aloo Gobi', type: 'Veg', cuisine: 'North Indian', ingredients: ['Potatoes', 'Cauliflower', 'Turmeric', 'Cumin', 'Coriander'] },
  { name: 'Bhindi Masala', type: 'Veg', cuisine: 'North Indian', ingredients: ['Okra', 'Onions', 'Tomatoes', 'Spices'] },
  { name: 'Kadai Chicken', type: 'Non-Veg', cuisine: 'North Indian', ingredients: ['Chicken', 'Bell Peppers', 'Onions', 'Tomatoes', 'Kadai Masala'] },
  { name: 'Tandoori Chicken', type: 'Non-Veg', cuisine: 'North Indian', ingredients: ['Chicken', 'Yogurt', 'Tandoori Masala', 'Ginger', 'Garlic'] },
  { name: 'Chicken Tikka Masala', type: 'Non-Veg', cuisine: 'North Indian', ingredients: ['Chicken Tikka', 'Tomatoes', 'Cream', 'Onions', 'Spices'] },
  { name: 'Naan', type: 'Veg', cuisine: 'North Indian', ingredients: ['Flour', 'Yogurt', 'Yeast', 'Salt', 'Oil'] },
  { name: 'Garlic Naan', type: 'Veg', cuisine: 'North Indian', ingredients: ['Flour', 'Yogurt', 'Yeast', 'Garlic', 'Butter'] },
  { name: 'Amritsari Kulcha', type: 'Veg', cuisine: 'Punjabi', ingredients: ['Flour', 'Potatoes', 'Onions', 'Spices', 'Yogurt'] },
  { name: 'Sarson Ka Saag', type: 'Veg', cuisine: 'Punjabi', ingredients: ['Mustard Greens', 'Spinach', 'Cornmeal', 'Ginger', 'Garlic'] },
  { name: 'Makki Ki Roti', type: 'Veg', cuisine: 'Punjabi', ingredients: ['Corn Flour', 'Water', 'Salt'] },
  { name: 'Masala Dosa', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Urad Dal', 'Potatoes', 'Mustard Seeds', 'Curry Leaves'] },
  { name: 'Plain Dosa', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Urad Dal', 'Salt'] },
  { name: 'Rava Dosa', type: 'Veg', cuisine: 'South Indian', ingredients: ['Semolina', 'Rice Flour', 'Cumin', 'Ginger', 'Curry Leaves'] },
  { name: 'Sambar', type: 'Veg', cuisine: 'South Indian', ingredients: ['Toor Dal', 'Tamarind', 'Vegetables', 'Sambar Powder', 'Curry Leaves'] },
  { name: 'Idli', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Urad Dal', 'Salt'] },
  { name: 'Vada', type: 'Veg', cuisine: 'South Indian', ingredients: ['Urad Dal', 'Green Chili', 'Ginger', 'Curry Leaves'] },
  { name: 'Rasam', type: 'Veg', cuisine: 'South Indian', ingredients: ['Tomatoes', 'Tamarind', 'Toor Dal', 'Black Pepper', 'Cumin'] },
  { name: 'Uttapam', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Urad Dal', 'Vegetables', 'Spices'] },
  { name: 'Fish Curry', type: 'Non-Veg', cuisine: 'South Indian', ingredients: ['Fish', 'Coconut', 'Tamarind', 'Curry Leaves', 'Spices'] },
  { name: 'Upma', type: 'Veg', cuisine: 'South Indian', ingredients: ['Semolina', 'Vegetables', 'Mustard Seeds', 'Curry Leaves', 'Ginger'] },
  { name: 'Pongal', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Moong Dal', 'Black Pepper', 'Cumin', 'Ginger'] },
  { name: 'Appam', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Coconut', 'Yeast', 'Sugar'] },
  { name: 'Puttu', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice Flour', 'Coconut', 'Salt'] },
  { name: 'Payasam', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Milk', 'Sugar', 'Cardamom', 'Nuts'] },
  { name: 'Curd Rice', type: 'Veg', cuisine: 'South Indian', ingredients: ['Rice', 'Yogurt', 'Mustard Seeds', 'Curry Leaves', 'Ginger'] },
  { name: 'Vada Pav', type: 'Veg', cuisine: 'Maharashtrian', ingredients: ['Potatoes', 'Bread', 'Gram Flour', 'Green Chili', 'Ginger'] },
  { name: 'Poha', type: 'Veg', cuisine: 'Maharashtrian', ingredients: ['Flattened Rice', 'Onions', 'Turmeric', 'Mustard Seeds', 'Curry Leaves'] },
  { name: 'Misal Pav', type: 'Veg', cuisine: 'Maharashtrian', ingredients: ['Mixed Lentils', 'Spices', 'Bread', 'Onions', 'Coriander'] },
  { name: 'Pav Bhaji', type: 'Veg', cuisine: 'Maharashtrian', ingredients: ['Mixed Vegetables', 'Bread', 'Butter', 'Onions', 'Spices'] },
  { name: 'Puran Poli', type: 'Veg', cuisine: 'Maharashtrian', ingredients: ['Wheat Flour', 'Chana Dal', 'Jaggery', 'Cardamom'] },
  { name: 'Dhokla', type: 'Veg', cuisine: 'Gujarati', ingredients: ['Gram Flour', 'Yogurt', 'Ginger', 'Green Chili', 'Mustard Seeds'] },
  { name: 'Thepla', type: 'Veg', cuisine: 'Gujarati', ingredients: ['Wheat Flour', 'Fenugreek Leaves', 'Spices', 'Yogurt'] },
  { name: 'Undhiyu', type: 'Veg', cuisine: 'Gujarati', ingredients: ['Mixed Vegetables', 'Green Beans', 'Sweet Potato', 'Spices'] },
  { name: 'Khandvi', type: 'Veg', cuisine: 'Gujarati', ingredients: ['Gram Flour', 'Yogurt', 'Ginger', 'Green Chili', 'Mustard Seeds'] },
  { name: 'Fafda', type: 'Veg', cuisine: 'Gujarati', ingredients: ['Gram Flour', 'Oil', 'Turmeric', 'Salt', 'Baking Soda'] },
  { name: 'Khaman', type: 'Veg', cuisine: 'Gujarati', ingredients: ['Gram Flour', 'Ginger', 'Green Chili', 'Lemon', 'Sugar'] },
  { name: 'Gujarati Dal', type: 'Veg', cuisine: 'Gujarati', ingredients: ['Toor Dal', 'Jaggery', 'Tamarind', 'Turmeric', 'Mustard Seeds'] },
  { name: 'Litti Chokha', type: 'Veg', cuisine: 'Bihari', ingredients: ['Wheat Flour', 'Roasted Gram', 'Eggplant', 'Tomatoes', 'Spices'] },
  { name: 'Machher Jhol', type: 'Non-Veg', cuisine: 'Bengali', ingredients: ['Fish', 'Potatoes', 'Tomatoes', 'Ginger', 'Turmeric', 'Mustard Oil'] },
  { name: 'Mishti Doi', type: 'Veg', cuisine: 'Bengali', ingredients: ['Milk', 'Sugar', 'Yogurt Culture'] },
  { name: 'Rasgulla', type: 'Veg', cuisine: 'Bengali', ingredients: ['Cottage Cheese', 'Sugar', 'Cardamom'] },
  { name: 'Sandesh', type: 'Veg', cuisine: 'Bengali', ingredients: ['Cottage Cheese', 'Sugar', 'Cardamom'] },
  { name: 'Aloo Posto', type: 'Veg', cuisine: 'Bengali', ingredients: ['Potatoes', 'Poppy Seeds', 'Green Chili', 'Mustard Oil'] },
  { name: 'Kosha Mangsho', type: 'Non-Veg', cuisine: 'Bengali', ingredients: ['Mutton', 'Onions', 'Ginger', 'Garlic', 'Spices'] },
  { name: 'Shorshe Ilish', type: 'Non-Veg', cuisine: 'Bengali', ingredients: ['Hilsa Fish', 'Mustard Seeds', 'Green Chili', 'Turmeric'] },
  { name: 'Chingri Malai Curry', type: 'Non-Veg', cuisine: 'Bengali', ingredients: ['Prawns', 'Coconut Milk', 'Ginger', 'Garlic', 'Spices'] },
  { name: 'Beguni', type: 'Veg', cuisine: 'Bengali', ingredients: ['Eggplant', 'Gram Flour', 'Turmeric', 'Salt', 'Oil'] },
  { name: 'Puchka', type: 'Veg', cuisine: 'Bengali', ingredients: ['Semolina', 'Potatoes', 'Tamarind Water', 'Spices'] },

  // Chinese Dishes (90 dishes)
  { name: 'Sweet and Sour Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Bell Peppers', 'Pineapple', 'Vinegar', 'Sugar', 'Soy Sauce'] },
  { name: 'Kung Pao Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Peanuts', 'Dried Chilies', 'Soy Sauce', 'Garlic'] },
  { name: 'Vegetable Fried Rice', type: 'Veg', cuisine: 'Chinese', ingredients: ['Rice', 'Vegetables', 'Soy Sauce', 'Garlic', 'Ginger'] },
  { name: 'Chicken Fried Rice', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Rice', 'Chicken', 'Vegetables', 'Soy Sauce', 'Eggs'] },
  { name: 'Shrimp Fried Rice', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Rice', 'Shrimp', 'Vegetables', 'Soy Sauce', 'Eggs'] },
  { name: 'Beef Fried Rice', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Rice', 'Beef', 'Vegetables', 'Soy Sauce', 'Oyster Sauce'] },
  { name: 'Vegetable Chow Mein', type: 'Veg', cuisine: 'Chinese', ingredients: ['Noodles', 'Vegetables', 'Soy Sauce', 'Sesame Oil'] },
  { name: 'Chicken Chow Mein', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Noodles', 'Chicken', 'Vegetables', 'Soy Sauce'] },
  { name: 'Hot and Sour Soup', type: 'Veg', cuisine: 'Chinese', ingredients: ['Mushrooms', 'Tofu', 'Vinegar', 'White Pepper', 'Soy Sauce'] },
  { name: 'Wonton Soup', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Wontons', 'Chicken Broth', 'Green Onions', 'Soy Sauce'] },
  { name: 'Peking Duck', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Duck', 'Hoisin Sauce', 'Cucumber', 'Spring Onions', 'Pancakes'] },
  { name: 'Mapo Tofu', type: 'Veg', cuisine: 'Chinese', ingredients: ['Tofu', 'Doubanjiang', 'Sichuan Peppercorns', 'Ground Pork'] },
  { name: 'Dim Sum', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Flour', 'Shrimp', 'Pork', 'Soy Sauce', 'Ginger'] },
  { name: 'General Tso Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Soy Sauce', 'Sugar', 'Vinegar', 'Garlic'] },
  { name: 'Ma Po Eggplant', type: 'Veg', cuisine: 'Chinese', ingredients: ['Eggplant', 'Doubanjiang', 'Garlic', 'Ginger', 'Soy Sauce'] },
  { name: 'Orange Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Orange Juice', 'Soy Sauce', 'Sugar', 'Ginger'] },
  { name: 'Beef and Broccoli', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Beef', 'Broccoli', 'Soy Sauce', 'Oyster Sauce', 'Garlic'] },
  { name: 'Spring Rolls', type: 'Veg', cuisine: 'Chinese', ingredients: ['Spring Roll Wrappers', 'Cabbage', 'Carrots', 'Bean Sprouts'] },
  { name: 'Szechuan Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Szechuan Peppercorns', 'Dried Chilies', 'Soy Sauce'] },
  { name: 'Black Bean Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Black Bean Sauce', 'Bell Peppers', 'Onions'] },
  { name: 'Cashew Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Cashews', 'Soy Sauce', 'Hoisin Sauce'] },
  { name: 'Salt and Pepper Shrimp', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Shrimp', 'Salt', 'White Pepper', 'Jalapeños'] },
  { name: 'Dan Dan Noodles', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Noodles', 'Ground Pork', 'Sesame Paste', 'Chili Oil'] },
  { name: 'Tea Smoked Duck', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Duck', 'Tea Leaves', 'Sugar', 'Rice', 'Spices'] },
  { name: 'Mongolian Beef', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Beef', 'Scallions', 'Ginger', 'Soy Sauce', 'Hoisin Sauce'] },
  { name: 'Honey Garlic Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Honey', 'Garlic', 'Soy Sauce', 'Ginger'] },
  { name: 'Sesame Chicken', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Chicken', 'Sesame Seeds', 'Soy Sauce', 'Sugar', 'Vinegar'] },
  { name: 'Chinese Broccoli with Oyster Sauce', type: 'Veg', cuisine: 'Chinese', ingredients: ['Chinese Broccoli', 'Oyster Sauce', 'Garlic', 'Ginger'] },
  { name: 'Fish with Black Bean Sauce', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Fish', 'Black Bean Sauce', 'Bell Peppers', 'Onions'] },
  { name: 'Pork with Sweet and Sour Sauce', type: 'Non-Veg', cuisine: 'Chinese', ingredients: ['Pork', 'Pineapple', 'Bell Peppers', 'Sweet and Sour Sauce'] },

  // Italian Dishes (100 dishes)
  { name: 'Spaghetti Carbonara', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Spaghetti', 'Eggs', 'Pancetta', 'Parmesan', 'Black Pepper'] },
  { name: 'Margherita Pizza', type: 'Veg', cuisine: 'Italian', ingredients: ['Pizza Dough', 'Tomato Sauce', 'Mozzarella', 'Basil', 'Olive Oil'] },
  { name: 'Pepperoni Pizza', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Pizza Dough', 'Tomato Sauce', 'Mozzarella', 'Pepperoni'] },
  { name: 'Meat Lovers Pizza', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Pizza Dough', 'Tomato Sauce', 'Mozzarella', 'Pepperoni', 'Sausage', 'Ham'] },
  { name: 'Veggie Pizza', type: 'Veg', cuisine: 'Italian', ingredients: ['Pizza Dough', 'Tomato Sauce', 'Mozzarella', 'Bell Peppers', 'Mushrooms', 'Olives'] },
  { name: 'Lasagna', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Pasta Sheets', 'Ground Beef', 'Tomato Sauce', 'Bechamel', 'Cheese'] },
  { name: 'Vegetable Lasagna', type: 'Veg', cuisine: 'Italian', ingredients: ['Pasta Sheets', 'Mixed Vegetables', 'Tomato Sauce', 'Bechamel', 'Cheese'] },
  { name: 'Mushroom Risotto', type: 'Veg', cuisine: 'Italian', ingredients: ['Arborio Rice', 'Mushrooms', 'Vegetable Stock', 'White Wine', 'Parmesan'] },
  { name: 'Seafood Risotto', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Arborio Rice', 'Mixed Seafood', 'Fish Stock', 'White Wine', 'Saffron'] },
  { name: 'Fettuccine Alfredo', type: 'Veg', cuisine: 'Italian', ingredients: ['Fettuccine', 'Butter', 'Heavy Cream', 'Parmesan', 'Garlic'] },
  { name: 'Chicken Parmigiana', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Chicken Breast', 'Breadcrumbs', 'Marinara Sauce', 'Mozzarella'] },
  { name: 'Eggplant Parmigiana', type: 'Veg', cuisine: 'Italian', ingredients: ['Eggplant', 'Breadcrumbs', 'Marinara Sauce', 'Mozzarella'] },
  { name: 'Osso Buco', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Veal Shanks', 'Tomatoes', 'White Wine', 'Carrots', 'Celery'] },
  { name: 'Tiramisu', type: 'Veg', cuisine: 'Italian', ingredients: ['Ladyfingers', 'Coffee', 'Mascarpone', 'Eggs', 'Cocoa Powder'] },
  { name: 'Gelato', type: 'Veg', cuisine: 'Italian', ingredients: ['Milk', 'Sugar', 'Eggs', 'Vanilla', 'Various Flavors'] },
  { name: 'Panna Cotta', type: 'Veg', cuisine: 'Italian', ingredients: ['Heavy Cream', 'Sugar', 'Gelatin', 'Vanilla'] },
  { name: 'Potato Gnocchi', type: 'Veg', cuisine: 'Italian', ingredients: ['Potatoes', 'Flour', 'Eggs', 'Salt'] },
  { name: 'Ricotta Gnocchi', type: 'Veg', cuisine: 'Italian', ingredients: ['Ricotta', 'Flour', 'Eggs', 'Parmesan'] },
  { name: 'Minestrone Soup', type: 'Veg', cuisine: 'Italian', ingredients: ['Mixed Vegetables', 'Beans', 'Tomatoes', 'Pasta', 'Herbs'] },
  { name: 'Penne Arrabbiata', type: 'Veg', cuisine: 'Italian', ingredients: ['Penne Pasta', 'Tomatoes', 'Red Chilies', 'Garlic', 'Olive Oil'] },
  { name: 'Penne alla Vodka', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Penne Pasta', 'Vodka', 'Tomatoes', 'Cream', 'Pancetta'] },
  { name: 'Cheese Ravioli', type: 'Veg', cuisine: 'Italian', ingredients: ['Pasta Dough', 'Ricotta', 'Parmesan', 'Spinach', 'Nutmeg'] },
  { name: 'Meat Ravioli', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Pasta Dough', 'Ground Beef', 'Ricotta', 'Herbs'] },
  { name: 'Prosciutto Pizza', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Pizza Dough', 'Prosciutto', 'Arugula', 'Mozzarella', 'Olive Oil'] },
  { name: 'Cannoli', type: 'Veg', cuisine: 'Italian', ingredients: ['Pastry Shells', 'Ricotta', 'Sugar', 'Vanilla', 'Chocolate Chips'] },
  { name: 'Bruschetta', type: 'Veg', cuisine: 'Italian', ingredients: ['Bread', 'Tomatoes', 'Basil', 'Garlic', 'Olive Oil'] },
  { name: 'Caprese Salad', type: 'Veg', cuisine: 'Italian', ingredients: ['Mozzarella', 'Tomatoes', 'Basil', 'Olive Oil', 'Balsamic'] },
  { name: 'Antipasto Platter', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Cured Meats', 'Cheese', 'Olives', 'Vegetables'] },
  { name: 'Tortellini in Brodo', type: 'Non-Veg', cuisine: 'Italian', ingredients: ['Tortellini', 'Chicken Broth', 'Parmesan'] },
  { name: 'Pesto Pasta', type: 'Veg', cuisine: 'Italian', ingredients: ['Pasta', 'Basil', 'Pine Nuts', 'Parmesan', 'Olive Oil'] },

  // Thai Dishes (70 dishes)
  { name: 'Pad Thai', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Rice Noodles', 'Shrimp', 'Bean Sprouts', 'Tamarind', 'Fish Sauce'] },
  { name: 'Vegetable Pad Thai', type: 'Veg', cuisine: 'Thai', ingredients: ['Rice Noodles', 'Tofu', 'Bean Sprouts', 'Tamarind', 'Soy Sauce'] },
  { name: 'Green Curry Chicken', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Green Curry Paste', 'Coconut Milk', 'Chicken', 'Thai Basil', 'Eggplant'] },
  { name: 'Green Curry Vegetables', type: 'Veg', cuisine: 'Thai', ingredients: ['Green Curry Paste', 'Coconut Milk', 'Mixed Vegetables', 'Thai Basil'] },
  { name: 'Red Curry Chicken', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Red Curry Paste', 'Coconut Milk', 'Chicken', 'Thai Basil'] },
  { name: 'Red Curry Beef', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Red Curry Paste', 'Coconut Milk', 'Beef', 'Thai Basil'] },
  { name: 'Tom Yum Soup', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Shrimp', 'Lemongrass', 'Lime Leaves', 'Galangal', 'Chili'] },
  { name: 'Tom Kha Gai', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Chicken', 'Coconut Milk', 'Galangal', 'Lime Leaves', 'Mushrooms'] },
  { name: 'Massaman Curry', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Massaman Paste', 'Beef', 'Potatoes', 'Coconut Milk', 'Peanuts'] },
  { name: 'Som Tam', type: 'Veg', cuisine: 'Thai', ingredients: ['Green Papaya', 'Tomatoes', 'Green Beans', 'Lime', 'Fish Sauce'] },
  { name: 'Mango Sticky Rice', type: 'Veg', cuisine: 'Thai', ingredients: ['Sticky Rice', 'Mango', 'Coconut Milk', 'Sugar'] },
  { name: 'Thai Fried Rice', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Rice', 'Chicken', 'Pineapple', 'Cashews', 'Raisins'] },
  { name: 'Larb', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Ground Pork', 'Mint', 'Cilantro', 'Lime', 'Fish Sauce'] },
  { name: 'Thai Basil Chicken', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Chicken', 'Thai Basil', 'Chilies', 'Garlic', 'Soy Sauce'] },
  { name: 'Pad See Ew', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Wide Rice Noodles', 'Pork', 'Chinese Broccoli', 'Dark Soy Sauce'] },
  { name: 'Thai Spring Rolls', type: 'Veg', cuisine: 'Thai', ingredients: ['Rice Paper', 'Vegetables', 'Herbs', 'Peanut Sauce'] },
  { name: 'Thai Satay', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Chicken Skewers', 'Turmeric', 'Coconut Milk', 'Peanut Sauce'] },
  { name: 'Thai Fish Cakes', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Fish Paste', 'Red Curry Paste', 'Green Beans', 'Lime Leaves'] },
  { name: 'Pad Kra Pao', type: 'Non-Veg', cuisine: 'Thai', ingredients: ['Ground Pork', 'Holy Basil', 'Chilies', 'Fish Sauce'] },
  { name: 'Thai Coconut Soup', type: 'Veg', cuisine: 'Thai', ingredients: ['Coconut Milk', 'Vegetables', 'Lemongrass', 'Galangal'] },

  // Korean Dishes (60 dishes)
  { name: 'Kimchi', type: 'Veg', cuisine: 'Korean', ingredients: ['Napa Cabbage', 'Korean Chili Flakes', 'Garlic', 'Ginger', 'Fish Sauce'] },
  { name: 'Bulgogi', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Beef', 'Soy Sauce', 'Sugar', 'Sesame Oil', 'Garlic'] },
  { name: 'Bibimbap', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Rice', 'Mixed Vegetables', 'Beef', 'Egg', 'Gochujang'] },
  { name: 'Vegetable Bibimbap', type: 'Veg', cuisine: 'Korean', ingredients: ['Rice', 'Mixed Vegetables', 'Egg', 'Gochujang'] },
  { name: 'Korean Fried Chicken', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Chicken', 'Korean Chili Sauce', 'Garlic', 'Ginger'] },
  { name: 'Japchae', type: 'Veg', cuisine: 'Korean', ingredients: ['Sweet Potato Noodles', 'Vegetables', 'Sesame Oil', 'Soy Sauce'] },
  { name: 'Korean BBQ Pork', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Pork Belly', 'Gochujang', 'Garlic', 'Soy Sauce'] },
  { name: 'Korean BBQ Beef', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Beef Short Ribs', 'Soy Sauce', 'Pear', 'Garlic'] },
  { name: 'Kimchi Fried Rice', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Rice', 'Kimchi', 'Pork', 'Egg', 'Sesame Oil'] },
  { name: 'Korean Corn Dogs', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Hot Dogs', 'Mozzarella', 'Potato Cubes', 'Batter'] },
  { name: 'Sundubu Jjigae', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Soft Tofu', 'Kimchi', 'Pork', 'Gochugaru', 'Scallions'] },
  { name: 'Korean Pancakes', type: 'Veg', cuisine: 'Korean', ingredients: ['Flour', 'Scallions', 'Kimchi', 'Vegetables'] },
  { name: 'Tteokbokki', type: 'Veg', cuisine: 'Korean', ingredients: ['Rice Cakes', 'Gochujang', 'Fish Cake', 'Scallions'] },
  { name: 'Korean Beef Soup', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Beef', 'Radish', 'Scallions', 'Garlic'] },
  { name: 'Galbi', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Beef Short Ribs', 'Soy Sauce', 'Pear', 'Garlic'] },
  { name: 'Korean Chicken Soup', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Chicken', 'Ginseng', 'Rice', 'Garlic'] },
  { name: 'Korean Tofu Soup', type: 'Veg', cuisine: 'Korean', ingredients: ['Tofu', 'Vegetables', 'Soy Sauce', 'Sesame Oil'] },
  { name: 'Korean Fish Stew', type: 'Non-Veg', cuisine: 'Korean', ingredients: ['Fish', 'Vegetables', 'Gochujang', 'Tofu'] },
  { name: 'Korean Spinach Salad', type: 'Veg', cuisine: 'Korean', ingredients: ['Spinach', 'Sesame Oil', 'Garlic', 'Soy Sauce'] },
  { name: 'Korean Bean Sprout Soup', type: 'Veg', cuisine: 'Korean', ingredients: ['Bean Sprouts', 'Scallions', 'Garlic', 'Soy Sauce'] },

  // French Dishes (60 dishes)
  { name: 'Coq au Vin', type: 'Non-Veg', cuisine: 'French', ingredients: ['Chicken', 'Red Wine', 'Mushrooms', 'Bacon', 'Pearl Onions'] },
  { name: 'Bouillabaisse', type: 'Non-Veg', cuisine: 'French', ingredients: ['Mixed Fish', 'Tomatoes', 'Saffron', 'Fennel', 'Olive Oil'] },
  { name: 'Ratatouille', type: 'Veg', cuisine: 'French', ingredients: ['Eggplant', 'Zucchini', 'Tomatoes', 'Bell Peppers', 'Herbs'] },
  { name: 'French Onion Soup', type: 'Veg', cuisine: 'French', ingredients: ['Onions', 'Beef Broth', 'Gruyere Cheese', 'Bread', 'Thyme'] },
  { name: 'Beef Bourguignon', type: 'Non-Veg', cuisine: 'French', ingredients: ['Beef', 'Red Wine', 'Mushrooms', 'Carrots', 'Bacon'] },
  { name: 'Croque Monsieur', type: 'Non-Veg', cuisine: 'French', ingredients: ['Ham', 'Gruyere Cheese', 'Bread', 'Bechamel Sauce'] },
  { name: 'Croque Madame', type: 'Non-Veg', cuisine: 'French', ingredients: ['Ham', 'Gruyere Cheese', 'Bread', 'Bechamel Sauce', 'Fried Egg'] },
  { name: 'Escargot', type: 'Non-Veg', cuisine: 'French', ingredients: ['Snails', 'Garlic', 'Parsley', 'Butter'] },
  { name: 'Duck Confit', type: 'Non-Veg', cuisine: 'French', ingredients: ['Duck Legs', 'Duck Fat', 'Garlic', 'Thyme'] },
  { name: 'Quiche Lorraine', type: 'Non-Veg', cuisine: 'French', ingredients: ['Pastry', 'Bacon', 'Eggs', 'Cream', 'Gruyere'] },
  { name: 'Vegetable Quiche', type: 'Veg', cuisine: 'French', ingredients: ['Pastry', 'Mixed Vegetables', 'Eggs', 'Cream', 'Cheese'] },
  { name: 'Cassoulet', type: 'Non-Veg', cuisine: 'French', ingredients: ['White Beans', 'Duck', 'Sausage', 'Tomatoes'] },
  { name: 'Crepes', type: 'Veg', cuisine: 'French', ingredients: ['Flour', 'Eggs', 'Milk', 'Butter', 'Sugar'] },
  { name: 'Savory Crepes', type: 'Non-Veg', cuisine: 'French', ingredients: ['Flour', 'Eggs', 'Milk', 'Ham', 'Cheese'] },
  { name: 'Macarons', type: 'Veg', cuisine: 'French', ingredients: ['Almond Flour', 'Egg Whites', 'Sugar', 'Food Coloring'] },
  { name: 'Croissant', type: 'Veg', cuisine: 'French', ingredients: ['Flour', 'Butter', 'Yeast', 'Sugar', 'Salt'] },
  { name: 'Pain au Chocolat', type: 'Veg', cuisine: 'French', ingredients: ['Pastry', 'Dark Chocolate', 'Butter'] },
  { name: 'Nicoise Salad', type: 'Non-Veg', cuisine: 'French', ingredients: ['Tuna', 'Olives', 'Tomatoes', 'Eggs', 'Anchovies'] },
  { name: 'French Toast', type: 'Veg', cuisine: 'French', ingredients: ['Bread', 'Eggs', 'Milk', 'Vanilla', 'Cinnamon'] },
  { name: 'Tarte Tatin', type: 'Veg', cuisine: 'French', ingredients: ['Apples', 'Pastry', 'Butter', 'Sugar'] },

  // Continental Dishes (80 dishes)
  { name: 'Caesar Salad', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Romaine Lettuce', 'Parmesan', 'Croutons', 'Anchovies', 'Caesar Dressing'] },
  { name: 'Greek Salad', type: 'Veg', cuisine: 'Continental', ingredients: ['Tomatoes', 'Cucumbers', 'Olives', 'Feta Cheese', 'Olive Oil'] },
  { name: 'Club Sandwich', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Bread', 'Turkey', 'Bacon', 'Lettuce', 'Tomato', 'Mayo'] },
  { name: 'BLT Sandwich', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Bread', 'Bacon', 'Lettuce', 'Tomato', 'Mayo'] },
  { name: 'Grilled Cheese Sandwich', type: 'Veg', cuisine: 'Continental', ingredients: ['Bread', 'Cheese', 'Butter'] },
  { name: 'Fish and Chips', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Fish Fillet', 'Potatoes', 'Flour', 'Beer Batter'] },
  { name: 'Beef Steak', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Beef Sirloin', 'Salt', 'Pepper', 'Garlic', 'Herbs'] },
  { name: 'Ribeye Steak', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Ribeye Beef', 'Salt', 'Pepper', 'Butter', 'Thyme'] },
  { name: 'Grilled Chicken Breast', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Chicken Breast', 'Herbs', 'Olive Oil', 'Lemon'] },
  { name: 'Roasted Chicken', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Whole Chicken', 'Herbs', 'Vegetables', 'Butter'] },
  { name: 'Shepherd\'s Pie', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Ground Lamb', 'Potatoes', 'Vegetables', 'Gravy'] },
  { name: 'Cottage Pie', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Ground Beef', 'Potatoes', 'Vegetables', 'Gravy'] },
  { name: 'Beef Wellington', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Beef Tenderloin', 'Puff Pastry', 'Mushroom Duxelles', 'Prosciutto'] },
  { name: 'Chicken Cordon Bleu', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Chicken Breast', 'Ham', 'Swiss Cheese', 'Breadcrumbs'] },
  { name: 'Lamb Chops', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Lamb Rack', 'Herbs', 'Garlic', 'Olive Oil'] },
  { name: 'Stuffed Bell Peppers', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Bell Peppers', 'Ground Beef', 'Rice', 'Tomatoes'] },
  { name: 'Chicken Kiev', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Chicken Breast', 'Garlic Butter', 'Breadcrumbs', 'Herbs'] },
  { name: 'Beef Stroganoff', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Beef Strips', 'Mushrooms', 'Sour Cream', 'Onions'] },
  { name: 'Pork Tenderloin', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Pork Tenderloin', 'Herbs', 'Garlic', 'Apple Sauce'] },
  { name: 'Grilled Salmon', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Salmon Fillet', 'Lemon', 'Dill', 'Olive Oil'] },
  { name: 'Baked Cod', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Cod Fillet', 'Breadcrumbs', 'Lemon', 'Herbs'] },
  { name: 'Chicken Wings', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Chicken Wings', 'Buffalo Sauce', 'Celery', 'Blue Cheese'] },
  { name: 'Meatloaf', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Ground Beef', 'Breadcrumbs', 'Eggs', 'Onions', 'Ketchup'] },
  { name: 'Pot Roast', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Beef Chuck', 'Vegetables', 'Beef Broth', 'Herbs'] },
  { name: 'Bangers and Mash', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Sausages', 'Mashed Potatoes', 'Onion Gravy'] },
  { name: 'Full English Breakfast', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['Eggs', 'Bacon', 'Sausages', 'Baked Beans', 'Toast'] },
  { name: 'Eggs Benedict', type: 'Non-Veg', cuisine: 'Continental', ingredients: ['English Muffins', 'Poached Eggs', 'Ham', 'Hollandaise Sauce'] },
  { name: 'Pancakes', type: 'Veg', cuisine: 'Continental', ingredients: ['Flour', 'Eggs', 'Milk', 'Sugar', 'Baking Powder'] },
  { name: 'Waffles', type: 'Veg', cuisine: 'Continental', ingredients: ['Flour', 'Eggs', 'Milk', 'Sugar', 'Butter'] },
  { name: 'French Fries', type: 'Veg', cuisine: 'Continental', ingredients: ['Potatoes', 'Oil', 'Salt'] }
];

const seedManyDishes = async () => {
  try {
    // Clear existing dishes
    await Dish.deleteMany({});
    console.log('Cleared existing dishes');

    // Add images and calories to dishes
    const dishesWithDetails = dishes.map(dish => ({
      ...dish,
      image: getRandomImage(dish.cuisine),
      calories: getRandomCalories()
    }));

    // Insert all dishes
    const result = await Dish.insertMany(dishesWithDetails);
    console.log(`✅ Successfully seeded ${result.length} dishes!`);
    
    // Log cuisine breakdown
    const cuisineCounts = {};
    dishesWithDetails.forEach(dish => {
      cuisineCounts[dish.cuisine] = (cuisineCounts[dish.cuisine] || 0) + 1;
    });
    
    console.log('\n📊 Cuisine breakdown:');
    Object.entries(cuisineCounts).forEach(([cuisine, count]) => {
      console.log(`${cuisine}: ${count} dishes`);
    });

    console.log(`\n🎉 Total dishes in database: ${result.length}`);
    
  } catch (error) {
    console.error('❌ Error seeding dishes:', error);
  }
};

const runSeed = async () => {
  await connectDB();
  await seedManyDishes();
  await mongoose.connection.close();
  console.log('Database connection closed');
};

// Run the seeding if this file is executed directly
if (require.main === module) {
  runSeed();
}

module.exports = seedManyDishes;
