# Go Backend for Indian Meal Planner

This is a professional Go backend implementation for the Indian Meal Planner application, converted from the original Node.js backend.

## Architecture

This Go backend follows clean architecture principles with the following structure:

```
backend-go/
├── cmd/server/           # Application entry point
├── internal/             # Private application code
│   ├── api/             # HTTP handlers and routing
│   │   ├── handlers/    # HTTP request handlers
│   │   └── middleware/  # HTTP middleware
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and seeding
│   ├── models/          # Data models and DTOs
│   ├── repository/      # Data access layer
│   └── service/         # Business logic layer
└── pkg/                 # Public packages
    └── logger/          # Logging utilities
```

## Technologies Used

- **Framework**: Gin (HTTP web framework)
- **Database**: MongoDB with official Go driver
- **Authentication**: JWT tokens
- **Validation**: go-playground/validator
- **Configuration**: Environment variables with godotenv
- **Logging**: Built-in slog package
- **Password Hashing**: bcrypt

## Features

- ✅ User authentication (register/login) with JWT
- ✅ User profile management with dietary preferences
- ✅ Dish management with search and filtering
- ✅ Favorites system
- ✅ Meal tracking and planning
- ✅ Nutrition summary and analytics
- ✅ CORS support for frontend integration
- ✅ Comprehensive logging
- ✅ Input validation
- ✅ Error handling
- ✅ Database seeding with default dishes

## Getting Started

### Prerequisites

- Go 1.21 or higher ([Download from golang.org](https://golang.org/dl/))
- MongoDB (local or Atlas)

### Option 1: Automated Setup (Recommended)

**Windows (PowerShell):**
```powershell
.\setup.ps1
```

**Linux/macOS:**
```bash
chmod +x setup.sh
./setup.sh
```

### Option 2: Manual Installation

1. Navigate to the Go backend directory:
   ```bash
   cd backend-go
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

4. Update the `.env` file with your configuration:
   ```
   PORT=5000
   MONGODB_URI=mongodb://localhost:27017/meal-planner
   JWT_SECRET=your-super-secret-jwt-key
   ```

### Troubleshooting

If you encounter dependency issues, see [DEPENDENCY_FIX.md](./DEPENDENCY_FIX.md) for solutions.

5. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

### Building for Production

```bash
go build -o meal-planner-server cmd/server/main.go
./meal-planner-server
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Dishes
- `GET /api/dishes` - Get dishes with pagination and filtering
- `GET /api/dishes/:id` - Get specific dish
- `GET /api/dishes/favorites` - Get user's favorite dishes (auth required)
- `POST /api/dishes/:id/favorite` - Add dish to favorites (auth required)
- `DELETE /api/dishes/:id/favorite` - Remove dish from favorites (auth required)

### User
- `GET /api/user/profile` - Get user profile (auth required)
- `PUT /api/user/profile` - Update user profile (auth required)
- `DELETE /api/user/account` - Delete user account (auth required)

### Meals
- `POST /api/meals` - Create meal entry (auth required)
- `GET /api/meals` - Get user's meals with pagination (auth required)
- `GET /api/meals?startDate=2024-01-01&endDate=2024-01-31` - Get meals by date range
- `GET /api/meals/nutrition-summary` - Get nutrition summary (auth required)
- `GET /api/meals/:id` - Get specific meal (auth required)
- `PUT /api/meals/:id` - Update meal (auth required)
- `DELETE /api/meals/:id` - Delete meal (auth required)

### Health
- `GET /api/health` - Health check endpoint

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `5000` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017/meal-planner` |
| `JWT_SECRET` | JWT signing secret | `fallback-secret-key` |
| `JWT_EXPIRES_IN` | JWT expiration time | `7d` |
| `NODE_ENV` | Environment mode | `development` |
| `ALLOWED_ORIGINS` | CORS allowed origins | `http://localhost:3000,http://localhost:5173` |
| `LOG_LEVEL` | Logging level | `info` |
| `LOG_FORMAT` | Log format (json/text) | `json` |

## Architecture Patterns

### Repository Pattern
- Clean separation between data access and business logic
- Easy to test and mock
- Database-agnostic interfaces

### Service Layer
- Contains all business logic
- Orchestrates between repositories
- Handles validation and error handling

### Dependency Injection
- All dependencies are injected through constructors
- Makes testing easier
- Promotes loose coupling

### Clean Architecture
- Domain models in the center
- Infrastructure details at the edges
- Dependency inversion principle

## Error Handling

All API responses follow a consistent format:

```json
{
  "success": true/false,
  "message": "Success message",
  "data": "Response data",
  "error": "Error message",
  "details": "Additional error details"
}
```

## Logging

The application uses structured logging with configurable levels:
- `debug`: Detailed debugging information
- `info`: General information
- `warn`: Warning messages
- `error`: Error messages

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- CORS protection
- Input validation
- SQL injection prevention (NoSQL)
- Error message sanitization

## Performance Optimizations

- MongoDB indexes for efficient queries
- Connection pooling
- Pagination for large datasets
- Structured logging for minimal overhead
- Graceful shutdown handling

## Comparison with Node.js Backend

| Feature | Node.js | Go |
|---------|---------|-----|
| Framework | Express | Gin |
| Language | JavaScript | Go |
| Performance | Good | Excellent |
| Type Safety | No (unless TypeScript) | Yes |
| Memory Usage | Higher | Lower |
| Concurrency | Event Loop | Goroutines |
| Ecosystem | Massive | Growing |
| Learning Curve | Easy | Moderate |

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
golangci-lint run
```

### Adding New Features

1. Add models in `internal/models/`
2. Add repository interfaces and implementations in `internal/repository/`
3. Add business logic in `internal/service/`
4. Add HTTP handlers in `internal/api/handlers/`
5. Update routes in `internal/api/router.go`

This Go backend maintains API compatibility with the original Node.js version while providing better performance, type safety, and more robust error handling.
