#!/bin/bash

# Meal API Test Script
# This script tests the meal endpoints to help debug issues

BASE_URL="http://localhost:5000/api"
TOKEN=""

echo "🧪 Meal API Test Script"
echo "======================="

# Check if server is running
echo "🔍 Checking if server is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ Server is not running at $BASE_URL"
    echo "Please start the server with: go run cmd/server/main.go"
    exit 1
fi
echo "✅ Server is running"

# Get JWT token (you'll need to replace this with your actual login)
echo ""
echo "🔑 To test authenticated endpoints, you need a JWT token."
echo "First, register/login to get a token:"
echo ""
echo "Register:"
echo "curl -X POST $BASE_URL/auth/register \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"email\":\"test@example.com\",\"password\":\"password123\",\"name\":\"Test User\"}'"
echo ""
echo "Login:"
echo "curl -X POST $BASE_URL/auth/login \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"email\":\"test@example.com\",\"password\":\"password123\"}'"
echo ""

if [ -z "$TOKEN" ]; then
    echo "⚠️  No JWT token provided. Please set TOKEN variable or run login commands above."
    echo ""
    echo "Example: TOKEN=\"your_jwt_token_here\" $0"
    echo ""
fi

# Test GET meals by date (this was failing)
echo "📅 Testing GET /api/meals/2025-08-21"
if [ -n "$TOKEN" ]; then
    curl -H "Authorization: Bearer $TOKEN" \
         -H "Content-Type: application/json" \
         "$BASE_URL/meals/2025-08-21"
else
    echo "Skipping - no token provided"
fi
echo ""
echo ""

# Test POST create meal (this was failing)
echo "📝 Testing POST /api/meals"
if [ -n "$TOKEN" ]; then
    curl -X POST \
         -H "Authorization: Bearer $TOKEN" \
         -H "Content-Type: application/json" \
         -d '{
           "date": "2025-08-21T10:30:00Z",
           "mealType": "breakfast",
           "dishId": "64f1a2b3c4d5e6f789012345",
           "notes": "Test meal from script",
           "rating": 4
         }' \
         "$BASE_URL/meals"
else
    echo "Skipping - no token provided"
fi
echo ""
echo ""

# Test GET dishes (this was working)
echo "🍽️  Testing GET /api/dishes (should work)"
if [ -n "$TOKEN" ]; then
    curl -H "Authorization: Bearer $TOKEN" \
         "$BASE_URL/dishes?page=1&limit=5"
else
    echo "Skipping - no token provided"
fi
echo ""
echo ""

echo "🎯 Common Issues and Solutions:"
echo ""
echo "1. 400 Bad Request for GET /api/meals/2025-08-21:"
echo "   - Fixed: Now accepts date format YYYY-MM-DD"
echo "   - Also accepts ObjectID for specific meals"
echo ""
echo "2. 400 Bad Request for POST /api/meals:"
echo "   - Check date format (use ISO8601): '2025-08-21T10:30:00Z'"
echo "   - Ensure mealType is one of: breakfast, lunch, dinner, snack"
echo "   - Ensure dishId is a valid ObjectID"
echo "   - Check server logs for detailed error messages"
echo ""
echo "3. 401 Unauthorized:"
echo "   - Make sure you're sending JWT token in Authorization header"
echo "   - Token format: 'Bearer your_jwt_token_here'"
echo ""
echo "💡 Check server logs for detailed error messages!"
