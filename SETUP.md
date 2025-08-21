# Nourish - Complete Setup Guide

## Features Added

✅ **Authentication System**
- User registration and login with JWT tokens
- Secure password hashing with bcryptjs
- Protected routes and middleware
- User profile management

✅ **Enhanced User Experience**
- Beautiful login/register forms
- User dropdown menu with quick navigation
- Protected routes for authenticated users
- Responsive design

✅ **Search & Filter**
- Advanced dish search by name, cuisine, and ingredients
- Filter by meal type, cuisine, and calorie count
- Real-time search with debouncing
- Integrated with favorites system

✅ **Favorites System**
- Mark dishes as favorites with heart icon
- Dedicated favorites page
- Quick add to meal plan from favorites
- Persistent across sessions

✅ **Analytics Dashboard**
- Meal consumption analytics
- Calorie tracking and insights
- Cuisine diversity metrics
- Meal frequency analysis
- Health recommendations
- Beautiful charts and visualizations

✅ **Shopping List Generator**
- Generate ingredient lists from meal plans
- Customizable date ranges
- Interactive checklist with progress tracking
- Export functionality
- Ingredient count by dish usage

✅ **User Profile**
- Dietary preferences management
- Spice level preferences
- Favorite regional cuisines
- Account settings

## Installation & Setup

### Prerequisites
- Go 1.21+ ([Download from golang.org](https://golang.org/dl/))
- MongoDB database
- npm or yarn (for frontend)

### Backend Setup

1. **Navigate to backend directory:**
   ```bash
   cd backend
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Environment Configuration:**
   ```bash
   cp .env.example .env
   ```
   
   Update `.env` with your values:
   ```env
   MONGODB_URI=mongodb+srv://username:password@cluster0.example.mongodb.net/meal-planner
   PORT=5000
   ENVIRONMENT=development
   JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters-long
   ```

4. **Start the backend server:**
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup

1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Environment Configuration:**
   Create `.env` file:
   ```env
   REACT_APP_API_URL=http://localhost:5000/api
   ```

4. **Start the frontend server:**
   ```bash
   npm start
   ```

### Quick Development Start

**Windows (PowerShell):**
```powershell
cd backend
.\setup.ps1
```

**Linux/macOS:**
```bash
cd backend
chmod +x setup.sh
./setup.sh
```

Or use the provided development script:
```bash
chmod +x start-dev.sh
./start-dev.sh
```

## New API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user profile
- `PUT /api/auth/profile` - Update user profile

### Dishes
- `GET /api/dishes` - Get all dishes with optional search/filter
- `POST /api/dishes` - Create new dish
- `GET /api/dishes/favorites` - Get user's favorite dishes
- `POST /api/dishes/:id/favorite` - Toggle favorite dish

### Meals
- `GET /api/meals` - Get meals with optional date filter
- `POST /api/meals` - Create meal plan
- `PUT /api/meals/:id` - Update meal
- `DELETE /api/meals/:id` - Delete meal

### Analytics
- `GET /api/analytics/nutrition-summary` - Get nutrition analytics for user

## User Flow

### New User Experience
1. **Registration**: User creates account with name, email, password
2. **Profile Setup**: Set dietary preferences, spice level, favorite regions
3. **Meal Planning**: Browse dishes, add to meal plan, mark favorites
4. **Analytics**: View eating patterns and health insights
5. **Shopping**: Generate shopping lists from meal plans

### Key Features for Users

#### Authentication
- Secure login/logout
- Password protection
- Persistent sessions

#### Meal Planning
- Day and month views
- Search and filter dishes
- Add custom dishes
- Rate meals (1-5 stars)
- Add notes to meals

#### Favorites Management
- Heart icon to favorite dishes
- Dedicated favorites page
- Quick add to meal plan

#### Analytics & Insights
- Calorie tracking
- Meal frequency analysis
- Cuisine diversity
- Health recommendations
- Progress over time

#### Shopping Lists
- Auto-generate from meal plans
- Interactive checklist
- Export to text file
- Ingredient usage statistics

## Database Schema

### Users Collection
```javascript
{
  name: String,
  email: String (unique),
  password: String (hashed),
  profile: {
    dietaryPreferences: [String],
    spiceLevel: String,
    favoriteRegions: [String],
    avatar: String
  },
  favoriteDishes: [ObjectId],
  isEmailVerified: Boolean,
  lastLoginAt: Date,
  createdAt: Date,
  updatedAt: Date
}
```

### Meals Collection (Updated)
```javascript
{
  date: String,
  mealType: String,
  dish: ObjectId,
  user: ObjectId, // NEW
  notes: String,  // NEW
  rating: Number, // NEW
  createdAt: Date,
  updatedAt: Date
}
```

## Security Features

- JWT token authentication
- Password hashing with bcryptjs
- Protected routes middleware
- Input validation with express-validator
- CORS configuration
- Request rate limiting (recommended for production)

## Production Deployment

### Environment Variables
Make sure to set strong JWT secrets and proper MongoDB URIs in production.

### Security Recommendations
1. Use HTTPS in production
2. Set strong JWT secrets (32+ characters)
3. Configure proper CORS origins
4. Add rate limiting
5. Use environment-specific configurations

## Troubleshooting

### Common Issues
1. **Authentication not working**: Check JWT_SECRET is set
2. **CORS errors**: Verify FRONTEND_URL in backend .env
3. **Database connection**: Ensure MongoDB URI is correct
4. **Port conflicts**: Change PORT in .env files

### Development Tips
- Use MongoDB Compass to inspect database
- Check browser Network tab for API errors
- Enable detailed error logging in development
- Use React Developer Tools for frontend debugging

## Next Steps & Enhancements

Potential future improvements:
- Email verification system
- Password reset functionality
- Meal plan templates
- Social features (sharing meal plans)
- Nutritional information integration
- Recipe instructions and cooking times
- Push notifications for meal reminders
- Dark mode theme
- Mobile app development
- Integration with grocery delivery services
