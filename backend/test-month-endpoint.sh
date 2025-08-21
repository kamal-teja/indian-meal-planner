#!/bin/bash

# Test script for the new /meals/month endpoint

BASE_URL="http://localhost:5000/api"
TOKEN=""

echo "🗓️  Testing Meals Month Endpoint"
echo "================================"

# Check if server is running
echo "🔍 Checking if server is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ Server is not running at $BASE_URL"
    echo "Please start the server with: go run cmd/server/main.go"
    exit 1
fi
echo "✅ Server is running"

if [ -z "$TOKEN" ]; then
    echo ""
    echo "⚠️  No JWT token provided. Please set TOKEN variable."
    echo "Example: TOKEN=\"your_jwt_token_here\" $0"
    echo ""
    echo "To get a token, login first:"
    echo "curl -X POST $BASE_URL/auth/login \\"
    echo "  -H 'Content-Type: application/json' \\"
    echo "  -d '{\"email\":\"test@example.com\",\"password\":\"password123\"}'"
    exit 1
fi

# Test the new month endpoint
echo "📅 Testing GET /api/meals/month/2025/8"
echo ""

curl -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$BASE_URL/meals/month/2025/8" \
     -w "\nHTTP Status: %{http_code}\n" \
     -s

echo ""
echo ""

# Test with invalid parameters
echo "❌ Testing with invalid year"
curl -H "Authorization: Bearer $TOKEN" \
     "$BASE_URL/meals/month/invalid/8" \
     -w "\nHTTP Status: %{http_code}\n" \
     -s

echo ""
echo ""

echo "❌ Testing with invalid month"
curl -H "Authorization: Bearer $TOKEN" \
     "$BASE_URL/meals/month/2025/13" \
     -w "\nHTTP Status: %{http_code}\n" \
     -s

echo ""
echo ""
echo "✅ Month endpoint testing complete!"
echo ""
echo "Expected behavior:"
echo "- Valid month request should return 200 with meals data"
echo "- Invalid year should return 400 with error message"
echo "- Invalid month should return 400 with error message"
