# Nourish - Indian Meal Planner Development Instructions

**CRITICAL**: Always follow these instructions first and only fallback to additional search and context gathering if the information here is incomplete or found to be in error.

## Repository Overview

Nourish is a full-stack meal planning application with:
- **Backend**: Go + Gin framework + MongoDB (located in `/backend`)
- **Frontend**: React + Vite + Tailwind CSS (located in `/frontend`) 
- **Database**: MongoDB (local or MongoDB Atlas)
- **Architecture**: Clean architecture with repository pattern

### Important Documentation Files

- `README.md` - Main project overview and quick start
- `backend/SETUP.md` - Detailed backend setup instructions  
- `frontend/SETUP-INSTRUCTIONS.md` - Frontend setup guide
- `SETUP.md` - Full application setup with user features
- `DEPLOYMENT.md` - Production deployment guide
- `backend/test-meals-api.sh` - API testing script

## Prerequisites Installation

### Required Tools

**Go 1.21+ (Backend)**:
```bash
# Check current version
go version

# Install if needed:
# Ubuntu/Debian: sudo apt update && sudo apt install golang-go
# macOS: brew install go  
# Windows: Download from https://golang.org/dl/
```

**Node.js 16+ and npm (Frontend)**:
```bash
# Check current versions
node --version && npm --version

# Install if needed:
# Ubuntu/Debian: curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - && sudo apt-get install -y nodejs
# macOS: brew install node
# Windows: Download from https://nodejs.org/
```

## Bootstrap, Build, and Test Commands

### Complete Setup (RECOMMENDED)

**NEVER CANCEL builds or long commands - use adequate timeouts:**

```bash
# First make scripts executable (Linux/Mac only)
chmod +x start-dev.sh

# Option 1: Automated setup (Linux/Mac)
./start-dev.sh
# TIMEOUT: 120 seconds - sets up both backend and frontend

# Option 2: Automated setup (Windows PowerShell) 
./start-dev.ps1
# TIMEOUT: 120 seconds - sets up both backend and frontend
```

### Manual Setup

**Backend Setup**:
```bash
cd backend
go mod tidy                    # Install dependencies - TIMEOUT: 30 seconds
cp .env.example .env          # Create environment file
# Edit .env with your MongoDB URI
make build                    # Build application - TIMEOUT: 60 seconds  
```

**Frontend Setup**:
```bash
cd frontend  
npm install                   # Install dependencies - TIMEOUT: 60 seconds
npm run build                 # Build application - TIMEOUT: 30 seconds
```

### Development Tools Installation

```bash
cd backend
./scripts/setup.sh            # Install Go dev tools - TIMEOUT: 180 seconds
# OR manually:
make install                  # Installs golangci-lint, air, godoc - TIMEOUT: 180 seconds
```

## Running the Application

### Development Mode (with hot reload)

**Start Both Services**:
```bash
# From repository root
./start-dev.sh               # Linux/Mac - TIMEOUT: 60 seconds
# OR
./start-dev.ps1              # Windows - TIMEOUT: 60 seconds
```

**Manual Start**:
```bash
# Terminal 1 - Backend
cd backend
make dev                     # Starts with hot reload - TIMEOUT: 30 seconds

# Terminal 2 - Frontend  
cd frontend
npm start                    # Starts dev server - TIMEOUT: 30 seconds
```

### Production Mode

```bash
# Backend
cd backend
make build-prod              # Production build - TIMEOUT: 60 seconds
make run                     # Run built binary

# Frontend
cd frontend
npm run build               # Production build - TIMEOUT: 30 seconds
npm run preview             # Preview built app
```

## Build Times and Timeout Requirements

**CRITICAL - NEVER CANCEL these operations:**

- **Go dependencies**: ~9 seconds - TIMEOUT: 30+ seconds
- **Go build**: ~20 seconds - TIMEOUT: 60+ seconds  
- **Go tests**: ~3 seconds - TIMEOUT: 30+ seconds
- **Go dev tools install**: ~72 seconds - TIMEOUT: 180+ seconds
- **Go linting**: ~7 seconds - TIMEOUT: 30+ seconds
- **npm install**: ~11 seconds - TIMEOUT: 60+ seconds
- **Frontend build**: ~4 seconds - TIMEOUT: 30+ seconds
- **Setup scripts**: ~45-60 seconds - TIMEOUT: 120+ seconds

## Testing and Validation

### Manual Validation Scenarios

**ALWAYS perform these validation steps after making changes:**

1. **Backend Health Check** (requires MongoDB):
   ```bash
   # Start backend (requires database connection)
   cd backend && make dev
   
   # In another terminal, test health endpoint
   curl http://localhost:5000/api/health
   # Expected: {"status": "healthy", "timestamp": "..."}
   ```

2. **Frontend Development Server**:
   ```bash
   cd frontend && npm start
   # Should start on http://localhost:3000 in ~1 second
   # Check console for any build errors
   ```

3. **Database Connection Validation**:
   ```bash
   # Backend will log connection status on startup
   # Success: "Connected to MongoDB" or similar
   # Failure: "Failed to connect to database" error
   ```

4. **Frontend-Backend Integration** (when database available):
   - Navigate to http://localhost:3000
   - Verify the app loads without console errors
   - Test adding a meal plan item
   - Test viewing monthly calendar
   - Verify dishes load from backend API

5. **Core User Scenarios** (manual testing):
   - **Meal Planning**: Add breakfast/lunch/dinner for current date
   - **Dish Management**: Browse dishes, add custom dish with ingredients
   - **Shopping List**: Generate shopping list from meal plans
   - **Calendar View**: Switch between day and month views

6. **Build Validation** (always possible):
   ```bash
   # Backend build test
   cd backend && make build
   # Should complete in ~20 seconds without errors
   
   # Frontend build test
   cd frontend && npm run build
   # Should complete in ~4 seconds without errors
   ```

### API Testing Scripts

```bash
cd backend
# Test meal API endpoints
./test-meals-api.sh           # Manual API testing script
./test-month-endpoint.sh      # Test monthly meal retrieval
```

### Code Quality Checks

**Backend Validation**:
```bash
cd backend

# NOTE: Must add Go tools to PATH for linting
export PATH=$PATH:$(go env GOPATH)/bin

make format                   # Format code - EXPECT go vet warnings (non-breaking)
make lint                     # Run golangci-lint - TIMEOUT: 30 seconds - EXPECT errors
make test                     # Run tests (currently no tests exist)
```

**Frontend Validation**:
```bash
cd frontend
npm run build                 # Validates code during build
# No separate lint command - relies on build validation
```

## Known Issues and Workarounds

### Backend Issues

1. **`make format` shows go vet warnings**: 
   - Warnings about unkeyed BSON struct literals and self-assignment - not build-breaking
   - Continue with development, these are code style issues

2. **`make lint` shows multiple errors**:
   - Unchecked error returns, ineffectual assignments, struct literal issues
   - Code builds and runs despite these issues - they are code quality improvements
   - NEVER let linting failures block development work

3. **Go tools PATH issue**:
   - Run `export PATH=$PATH:$(go env GOPATH)/bin` before using make lint
   - Tools install to Go's bin directory which may not be in PATH

4. **No tests exist**: 
   - Focus on manual validation scenarios above
   - Test changes by running the application and exercising features

5. **MongoDB connection required**:
   - Use local MongoDB: `mongodb://localhost:27017/meal-planner`
   - OR use MongoDB Atlas cloud database
   - App will fail to start without database connection

6. **Scripts need execute permission**:
   - Run `chmod +x start-dev.sh` on Linux/Mac before using scripts

### Frontend Issues

1. **Security vulnerabilities in npm audit**:
   - Known issue with 2 moderate vulnerabilities
   - Run `npm audit fix` if needed, but not critical for development

## Repository Navigation

### Key Backend Files

```
backend/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── api/handlers/           # HTTP request handlers
│   ├── config/config.go        # Environment configuration
│   ├── models/                 # Data models and DTOs
│   ├── repository/             # Database access layer
│   └── service/                # Business logic layer  
├── Makefile                    # Build automation commands
└── .env                        # Environment variables (create from .env.example)
```

### Key Frontend Files

```
frontend/
├── src/
│   ├── components/             # React components
│   │   ├── DayView.js         # Daily meal planning interface
│   │   ├── MonthView.js       # Monthly calendar view
│   │   ├── DishSelector.js    # Modal for selecting dishes
│   │   ├── ShoppingList.js    # Shopping list generation
│   │   └── auth/              # Authentication components
│   ├── services/              # API service calls
│   └── App.js                 # Main application component
└── package.json               # Dependencies and scripts
```

### Configuration Files

- `backend/.env` - Backend environment variables (create from .env.example)
- `backend/go.mod` - Go module dependencies  
- `frontend/package.json` - Node.js dependencies and scripts
- `frontend/vite.config.js` - Vite build configuration
- `frontend/tailwind.config.js` - Tailwind CSS configuration

## Environment Variables

### Required Backend Variables (.env)

```bash
MONGODB_URI=mongodb://localhost:27017/meal-planner  # OR Atlas URI
PORT=5000
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ENVIRONMENT=development
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Optional Frontend Variables

```bash
REACT_APP_API_URL=http://localhost:5000/api  # Backend API URL
```

## Development Workflow

### Making Changes

1. **Edit Code**: Focus on `backend/internal/` or `frontend/src/`
2. **Format Backend**: `cd backend && make format` (ignore warnings)
3. **Build**: `make build` (backend) or `npm run build` (frontend)  
4. **Test Manually**: Follow validation scenarios above
5. **ALWAYS validate end-to-end functionality before completing changes**

### Common Tasks

- **Add new dish**: Edit default dishes in backend seeding code
- **Modify UI**: Update React components in `frontend/src/components/`
- **Change API**: Update handlers in `backend/internal/api/handlers/`
- **Database schema**: Update models in `backend/internal/models/`

## Deployment

- **Development**: Use local setup with MongoDB connection
- **Production**: See `DEPLOYMENT.md` for free hosting options (Render, Vercel, etc.)

## Troubleshooting

1. **Port conflicts**: Backend runs on :5000, frontend on :3000
2. **Database issues**: Check MongoDB connection and URI format
3. **Build failures**: Ensure Go 1.21+ and Node.js 16+ installed
4. **API errors**: Check backend logs and validate request format

## Critical Reminders

- **NEVER CANCEL** long-running builds or installs - they complete within documented timeouts
- **ALWAYS** test manually after making changes - no automated test suite exists
- **Database connection required** - ensure MongoDB is accessible before starting backend
- **Build artifacts excluded** - build/, dist/, node_modules/ are in .gitignore