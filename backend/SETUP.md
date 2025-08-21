# Setup Guide for Go Backend

## Prerequisites Installation

### 1. Install Go

#### Windows:
1. Download Go from https://golang.org/dl/
2. Download the Windows installer (.msi file)
3. Run the installer and follow the instructions
4. Add Go to your PATH if not automatically added:
   - Add `C:\Program Files\Go\bin` to your system PATH
5. Verify installation: Open PowerShell and run `go version`

#### macOS:
```bash
# Using Homebrew
brew install go

# Or download from https://golang.org/dl/
```

#### Linux:
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Or download from https://golang.org/dl/
```

### 2. Install MongoDB

#### Option 1: Local MongoDB
- Download from https://www.mongodb.com/try/download/community
- Follow installation instructions for your OS
- Start MongoDB service

#### Option 2: MongoDB Atlas (Recommended)
1. Go to https://www.mongodb.com/atlas
2. Create a free account
3. Create a new cluster
4. Get the connection string
5. Update your .env file with the connection string

## Quick Start

### 1. Setup Environment

#### Windows PowerShell:
```powershell
cd backend-go
.\scripts\setup.bat
```

#### Linux/macOS:
```bash
cd backend-go
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### 2. Configure Environment Variables

Copy and edit the environment file:
```bash
cp .env.example .env
```

Update `.env` with your configuration:
```env
PORT=5000
MONGODB_URI=mongodb://localhost:27017/meal-planner
# Or use MongoDB Atlas:
# MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/meal-planner

JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=7d
NODE_ENV=development
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3. Run the Application

#### Development mode (with hot reload):
```bash
make dev
```

#### Production mode:
```bash
make build
make run
```

#### Using Go directly:
```bash
go run cmd/server/main.go
```

## API Testing

Once the server is running, you can test the API:

### Health Check
```bash
curl http://localhost:5000/api/health
```

### Register a User
```bash
curl -X POST http://localhost:5000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Dishes
```bash
curl http://localhost:5000/api/dishes
```

## Development Workflow

### 1. Make Changes
Edit files in the `internal/` directory following the established patterns.

### 2. Format Code
```bash
make format
```

### 3. Run Linter
```bash
make lint
```

### 4. Run Tests
```bash
make test
```

### 5. Build for Production
```bash
make build-prod
```

## Project Structure Explanation

```
backend-go/
├── cmd/server/main.go          # Application entry point
├── internal/                   # Private application code
│   ├── api/                   # HTTP layer
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── middleware/        # HTTP middleware (auth, CORS, etc.)
│   │   └── router.go          # Route definitions
│   ├── config/                # Configuration management
│   ├── database/              # Database connection and seeding
│   ├── models/                # Data models and DTOs
│   ├── repository/            # Data access layer (MongoDB operations)
│   └── service/               # Business logic layer
├── pkg/                       # Public packages (can be imported by other projects)
│   └── logger/                # Logging utilities
├── scripts/                   # Setup and utility scripts
├── .env.example              # Environment variables template
├── go.mod                    # Go modules file
├── go.sum                    # Go modules checksums
├── Makefile                  # Build automation
├── Dockerfile                # Docker configuration
├── .air.toml                 # Hot reload configuration
└── README.md                 # Documentation
```

## Key Features Implemented

### 🔐 Authentication & Authorization
- JWT-based authentication
- Password hashing with bcrypt
- Protected and optional auth routes
- User registration and login

### 📊 User Management
- User profiles with dietary preferences
- Nutrition goals tracking
- Account management

### 🍽️ Dish Management
- Comprehensive dish database
- Search and filtering capabilities
- Favorites system
- Nutritional information

### 📅 Meal Planning
- Meal entry tracking
- Date range queries
- Nutrition summary and analytics
- Meal history

### 🏗️ Architecture
- Clean architecture with clear separation of concerns
- Repository pattern for data access
- Service layer for business logic
- Dependency injection
- Comprehensive error handling
- Structured logging

### 🔧 Developer Experience
- Hot reload for development
- Comprehensive Makefile
- Docker support
- Linting and formatting tools
- Environment-based configuration

## Production Deployment

### Using Docker
```bash
# Build image
docker build -t meal-planner-backend-go .

# Run container
docker run -p 5000:5000 --env-file .env meal-planner-backend-go
```

### Direct Deployment
```bash
# Build for production
make build-prod

# Copy binary and .env to server
scp build/meal-planner-server user@server:/app/
scp .env user@server:/app/

# Run on server
./meal-planner-server
```

## Migration from Node.js

This Go backend is a complete rewrite of the Node.js backend with:

### ✅ Full API Compatibility
- All endpoints from the Node.js version are implemented
- Same request/response formats
- Same authentication flow
- Same database schema

### 🚀 Performance Improvements
- Better memory usage
- Faster response times
- More efficient database queries
- Better concurrency handling

### 🛡️ Enhanced Reliability
- Strong type safety
- Better error handling
- Comprehensive input validation
- Structured logging

### 📈 Scalability
- Better resource utilization
- Efficient goroutines for concurrency
- Connection pooling
- Graceful shutdown

You can switch from the Node.js backend to this Go backend without any changes to your frontend application.

## Troubleshooting

### Common Issues

#### 1. Go not found
- Ensure Go is installed and in your PATH
- Verify with `go version`

#### 2. MongoDB connection issues
- Check if MongoDB is running
- Verify connection string in .env
- Check firewall settings for MongoDB Atlas

#### 3. Port already in use
- Change PORT in .env file
- Kill process using the port: `lsof -ti:5000 | xargs kill -9` (macOS/Linux)

#### 4. Module not found errors
- Run `go mod tidy`
- Ensure you're in the backend-go directory

### Getting Help

1. Check the logs for error messages
2. Ensure all environment variables are set correctly
3. Verify MongoDB connection
4. Check if all dependencies are installed

For more help, refer to the detailed README.md in the project root.
