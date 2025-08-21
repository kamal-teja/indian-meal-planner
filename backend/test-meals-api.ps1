# Meal API Test Script for Windows
# This script tests the meal endpoints to help debug issues

$BASE_URL = "http://localhost:5000/api"
$TOKEN = ""

Write-Host "🧪 Meal API Test Script" -ForegroundColor Green
Write-Host "=======================" -ForegroundColor Green

# Check if server is running
Write-Host "🔍 Checking if server is running..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$BASE_URL/health" -Method GET -TimeoutSec 5 -ErrorAction Stop
    Write-Host "✅ Server is running" -ForegroundColor Green
} catch {
    Write-Host "❌ Server is not running at $BASE_URL" -ForegroundColor Red
    Write-Host "Please start the server with: go run cmd/server/main.go" -ForegroundColor Yellow
    exit 1
}

# Get JWT token info
Write-Host ""
Write-Host "🔑 To test authenticated endpoints, you need a JWT token." -ForegroundColor Cyan
Write-Host "First, register/login to get a token:" -ForegroundColor White
Write-Host ""
Write-Host "Register:" -ForegroundColor Yellow
Write-Host "Invoke-RestMethod -Uri '$BASE_URL/auth/register' -Method POST -ContentType 'application/json' -Body '{\"email\":\"test@example.com\",\"password\":\"password123\",\"name\":\"Test User\"}'" -ForegroundColor Gray
Write-Host ""
Write-Host "Login:" -ForegroundColor Yellow
Write-Host "Invoke-RestMethod -Uri '$BASE_URL/auth/login' -Method POST -ContentType 'application/json' -Body '{\"email\":\"test@example.com\",\"password\":\"password123\"}'" -ForegroundColor Gray
Write-Host ""

if ([string]::IsNullOrEmpty($TOKEN)) {
    Write-Host "⚠️  No JWT token provided. Please set TOKEN variable or run login commands above." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Example: `$TOKEN = 'your_jwt_token_here'; .\test-meals-api.ps1" -ForegroundColor Gray
    Write-Host ""
}

# Test GET meals by date (this was failing)
Write-Host "📅 Testing GET /api/meals/2025-08-21" -ForegroundColor Cyan
if (-not [string]::IsNullOrEmpty($TOKEN)) {
    try {
        $headers = @{ Authorization = "Bearer $TOKEN"; "Content-Type" = "application/json" }
        $response = Invoke-RestMethod -Uri "$BASE_URL/meals/2025-08-21" -Method GET -Headers $headers
        Write-Host "✅ Success:" -ForegroundColor Green
        $response | ConvertTo-Json -Depth 3
    } catch {
        Write-Host "❌ Error:" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
} else {
    Write-Host "Skipping - no token provided" -ForegroundColor Yellow
}
Write-Host ""

# Test POST create meal (this was failing)
Write-Host "📝 Testing POST /api/meals" -ForegroundColor Cyan
if (-not [string]::IsNullOrEmpty($TOKEN)) {
    try {
        $headers = @{ Authorization = "Bearer $TOKEN"; "Content-Type" = "application/json" }
        $body = @{
            date = "2025-08-21T10:30:00Z"
            mealType = "breakfast"
            dishId = "64f1a2b3c4d5e6f789012345"
            notes = "Test meal from PowerShell script"
        } | ConvertTo-Json
        
        $response = Invoke-RestMethod -Uri "$BASE_URL/meals" -Method POST -Headers $headers -Body $body
        Write-Host "✅ Success:" -ForegroundColor Green
        $response | ConvertTo-Json -Depth 3
    } catch {
        Write-Host "❌ Error:" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
} else {
    Write-Host "Skipping - no token provided" -ForegroundColor Yellow
}
Write-Host ""

# Test GET dishes (this was working)
Write-Host "🍽️  Testing GET /api/dishes (should work)" -ForegroundColor Cyan
if (-not [string]::IsNullOrEmpty($TOKEN)) {
    try {
        $headers = @{ Authorization = "Bearer $TOKEN" }
        $response = Invoke-RestMethod -Uri "$BASE_URL/dishes?page=1&limit=5" -Method GET -Headers $headers
        Write-Host "✅ Success:" -ForegroundColor Green
        Write-Host "Found $($response.data.length) dishes" -ForegroundColor White
    } catch {
        Write-Host "❌ Error:" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
} else {
    Write-Host "Skipping - no token provided" -ForegroundColor Yellow
}
Write-Host ""

# Test GET meals by month (new endpoint)
Write-Host "📅 Testing GET /api/meals/month/2025/8 (new endpoint)" -ForegroundColor Cyan
if (-not [string]::IsNullOrEmpty($TOKEN)) {
    try {
        $headers = @{ Authorization = "Bearer $TOKEN"; "Content-Type" = "application/json" }
        $response = Invoke-RestMethod -Uri "$BASE_URL/meals/month/2025/8" -Method GET -Headers $headers
        Write-Host "✅ Success:" -ForegroundColor Green
        Write-Host "Found $($response.data.length) meals for August 2025" -ForegroundColor White
    } catch {
        Write-Host "❌ Error:" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
} else {
    Write-Host "Skipping - no token provided" -ForegroundColor Yellow
}
Write-Host ""

Write-Host "🎯 Common Issues and Solutions:" -ForegroundColor Magenta
Write-Host ""
Write-Host "1. 400 Bad Request for GET /api/meals/2025-08-21:" -ForegroundColor Yellow
Write-Host "   - Fixed: Now accepts date format YYYY-MM-DD" -ForegroundColor White
Write-Host "   - Also accepts ObjectID for specific meals" -ForegroundColor White
Write-Host ""
Write-Host "2. 400 Bad Request for POST /api/meals:" -ForegroundColor Yellow
Write-Host "   - Check date format (use ISO8601): '2025-08-21T10:30:00Z'" -ForegroundColor White
Write-Host "   - Ensure mealType is one of: breakfast, lunch, dinner, snack" -ForegroundColor White
Write-Host "   - Ensure dishId is a valid ObjectID" -ForegroundColor White
Write-Host "   - Rating is optional (0-5, where 0 = no rating)" -ForegroundColor White
Write-Host "   - Check server logs for detailed error messages" -ForegroundColor White
Write-Host ""
Write-Host "3. 401 Unauthorized:" -ForegroundColor Yellow
Write-Host "   - Make sure you're sending JWT token in Authorization header" -ForegroundColor White
Write-Host "   - Token format: 'Bearer your_jwt_token_here'" -ForegroundColor White
Write-Host ""
Write-Host "💡 Check server logs for detailed error messages!" -ForegroundColor Cyan
