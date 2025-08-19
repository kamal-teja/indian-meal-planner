# 🚀 Meal Planner App - Free Deployment Guide

This guide will help you deploy your Indian Meal Planner app for **FREE** using MongoDB Atlas and various hosting providers.

## 📋 Prerequisites

Before deploying, ensure you have:
- A GitHub account
- Basic knowledge of Git commands

## 🗄️ Step 1: Setup MongoDB Atlas (Free Database)

### 1.1 Create MongoDB Atlas Account
1. Go to [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Click "Try Free" and create an account
3. Choose "Build a Database" → "Free" (M0 Sandbox)
4. Choose your preferred cloud provider and region
5. Create cluster (this takes 2-3 minutes)

### 1.2 Configure Database Access
1. **Database Access**: Go to "Database Access" in the sidebar
   - Click "Add New Database User"
   - Choose "Password" authentication
   - Create username and strong password
   - Database User Privileges: "Read and write to any database"
   - Click "Add User"

2. **Network Access**: Go to "Network Access" in the sidebar
   - Click "Add IP Address"
   - Choose "Allow Access from Anywhere" (0.0.0.0/0)
   - Click "Confirm"

### 1.3 Get Connection String
1. Go to "Clusters" → Click "Connect" on your cluster
2. Choose "Connect your application"
3. Copy the connection string (looks like: `mongodb+srv://username:password@cluster0.xxxxx.mongodb.net/`)
4. Replace `<password>` with your actual password
5. Add database name at the end: `mongodb+srv://username:password@cluster0.xxxxx.mongodb.net/meal-planner`

## 🚀 Step 2: Deploy to Render (Recommended - Easy & Free)

### 2.1 Prepare Your Repository
1. Push your code to GitHub:
```bash
git add .
git commit -m "Add MongoDB persistence and deployment config"
git push origin main
```

### 2.2 Deploy Backend (API)
1. Go to [Render](https://render.com) and sign up/login
2. Click "New" → "Web Service"
3. Connect your GitHub repository
4. Configure:
   - **Name**: `meal-planner-api`
   - **Root Directory**: `backend`
   - **Environment**: `Node`
   - **Build Command**: `npm install`
   - **Start Command**: `npm start`
   - **Instance Type**: Free

5. **Environment Variables** (click "Advanced"):
   - `NODE_ENV`: `production`
   - `MONGODB_URI`: `[Your MongoDB connection string from Step 1.3]`
   - `PORT`: `5000`

6. Click "Create Web Service"
7. Wait for deployment (5-10 minutes)
8. Copy the service URL (e.g., `https://meal-planner-api.onrender.com`)

### 2.3 Deploy Frontend
1. In Render dashboard, click "New" → "Static Site"
2. Connect the same GitHub repository
3. Configure:
   - **Name**: `meal-planner-frontend`
   - **Root Directory**: `frontend`
   - **Build Command**: `npm install && npm run build`
   - **Publish Directory**: `build`

4. **Environment Variables**:
   - `REACT_APP_API_URL`: `https://meal-planner-api.onrender.com/api`

5. Click "Create Static Site"
6. Wait for deployment (5-10 minutes)

### 2.4 Update CORS (Backend)
1. Go to your backend service in Render
2. Add environment variable:
   - `FRONTEND_URL`: `https://your-frontend-url.onrender.com`

## 🔄 Alternative Deployment Options

### Option 2: Vercel (Frontend) + Railway (Backend)

#### Deploy Backend to Railway:
1. Go to [Railway](https://railway.app)
2. Connect GitHub repo
3. Deploy backend folder
4. Add MongoDB URI in environment variables
5. Copy the backend URL

#### Deploy Frontend to Vercel:
1. Go to [Vercel](https://vercel.com)
2. Import GitHub repository
3. Set root directory to `frontend`
4. Add environment variable: `REACT_APP_API_URL` with Railway backend URL

### Option 3: Netlify (Frontend) + Heroku (Backend)

#### Deploy Backend to Heroku:
1. Install Heroku CLI
2. Create Heroku app:
```bash
cd backend
heroku create meal-planner-api
heroku config:set MONGODB_URI="your-mongodb-uri"
heroku config:set NODE_ENV=production
git add .
git commit -m "Deploy to Heroku"
git push heroku main
```

#### Deploy Frontend to Netlify:
1. Go to [Netlify](https://netlify.com)
2. Drag and drop your `frontend/build` folder (after running `npm run build`)
3. Set environment variables in Netlify dashboard

## 🔧 Step 3: Test Your Deployment

1. **Backend Health Check**: Visit `https://your-backend-url/api/health`
2. **Frontend**: Visit your frontend URL
3. **Test Features**:
   - View default dishes
   - Add new dishes
   - Plan meals
   - View shopping list

## 🛠️ Troubleshooting

### Common Issues:

1. **MongoDB Connection Fails**:
   - Check your connection string
   - Verify network access allows all IPs
   - Ensure database user has correct permissions

2. **Frontend Can't Connect to Backend**:
   - Verify `REACT_APP_API_URL` environment variable
   - Check CORS settings in backend
   - Ensure backend is deployed and running

3. **Build Failures**:
   - Check Node.js version compatibility
   - Verify all dependencies are in package.json
   - Check build logs for specific errors

### Environment Variables Checklist:

**Backend**:
- ✅ `MONGODB_URI`
- ✅ `NODE_ENV=production`
- ✅ `PORT` (usually auto-set by hosting provider)

**Frontend**:
- ✅ `REACT_APP_API_URL`

## 📱 Step 4: Optional Enhancements

### Custom Domain (Free with most providers):
1. Register a free domain at [Freenom](https://freenom.com)
2. Configure DNS in your hosting provider
3. Set up HTTPS (usually automatic)

### Performance Monitoring:
1. Set up basic monitoring in your hosting dashboard
2. Add error tracking with [Sentry](https://sentry.io) (free tier)

## 🎉 Success!

Your meal planner app is now live and accessible worldwide! The free tier limitations:

- **MongoDB Atlas**: 512 MB storage, shared cluster
- **Render**: 500 build minutes/month, apps sleep after 15 min inactivity
- **Vercel**: 100 GB bandwidth, unlimited static sites

These limits are more than sufficient for personal use and small-scale applications.

## 🔄 Making Updates

To update your deployed app:
1. Make changes locally
2. Commit and push to GitHub
3. Deployments will trigger automatically

## 💡 Tips for Production

1. **Regular Backups**: MongoDB Atlas provides point-in-time backups
2. **Monitoring**: Set up uptime monitoring with [UptimeRobot](https://uptimerobot.com)
3. **Security**: Use environment variables for all sensitive data
4. **Performance**: Consider upgrading to paid tiers as usage grows

Happy cooking and meal planning! 🍽️✨
