require('dotenv').config();
const fs = require('fs');
const path = require('path');
const mongoose = require('mongoose');
const Dish = require('../models/Dish');

const log = (msg) => console.log(msg);
const error = (msg) => console.error(msg);

const getJsonPath = () => {
  const argPath = process.argv[2];
  const envPath = process.env.DISHES_JSON;
  const defaultPath = path.join(__dirname, '..', 'data', '550-dishes-complete.json');
  return path.resolve(argPath || envPath || defaultPath);
};

const connectDB = async () => {
  const uri = process.env.MONGODB_URI;
  if (!uri) {
    throw new Error('MONGODB_URI env var is not set');
  }
  await mongoose.connect(uri);
  log('✅ Connected to MongoDB');
};

const importDishes = async () => {
  const jsonPath = getJsonPath();
  if (!fs.existsSync(jsonPath)) {
    throw new Error(`JSON file not found at: ${jsonPath}`);
  }

  const raw = fs.readFileSync(jsonPath, 'utf8');
  const dishes = JSON.parse(raw);
  log(`📦 Loaded ${dishes.length} dishes from ${jsonPath}`);

  const clear = (process.env.CLEAR_EXISTING || 'false').toLowerCase() === 'true';
  const upsert = (process.env.UPSERT || 'true').toLowerCase() === 'true';

  if (clear) {
    const del = await Dish.deleteMany({});
    log(`🧹 Cleared existing dishes: ${del.deletedCount}`);
  }

  if (upsert) {
    let success = 0;
    for (const d of dishes) {
      await Dish.updateOne(
        { name: d.name, cuisine: d.cuisine },
        {
          $set: {
            type: d.type,
            cuisine: d.cuisine,
            image: d.image,
            ingredients: d.ingredients,
            calories: d.calories,
            name: d.name
          }
        },
        { upsert: true }
      );
      success++;
      if (success % 100 === 0) log(`   ...upserted ${success}`);
    }
    log(`✅ Upserted ${success} dishes`);
  } else {
    const result = await Dish.insertMany(dishes, { ordered: false });
    log(`✅ Inserted ${result.length} dishes`);
  }
};

(async () => {
  try {
    await connectDB();
    await importDishes();
  } catch (e) {
    error('❌ Import failed: ' + e.message);
    process.exitCode = 1;
  } finally {
    await mongoose.connection.close();
    log('🔌 Disconnected from MongoDB');
  }
})();


