# Test script for the new /meals/month endpoint (Windows PowerShell)

$BASE_URL = "http://localhost:5000/api"
$TOKEN = ""

Write-Host "🗓️  Testing Meals Month Endpoint" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green

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

if ([string]::IsNullOrEmpty($TOKEN)) {
    Write-Host ""
    Write-Host "⚠️  No JWT token provided. Please set TOKEN variable." -ForegroundColor Yellow
    Write-Host "Example: `$TOKEN = 'your_jwt_token_here'; .\test-month-endpoint.ps1" -ForegroundColor Gray
    Write-Host ""
    Write-Host "To get a token, login first:" -ForegroundColor Cyan
    Write-Host "Invoke-RestMethod -Uri '$BASE_URL/auth/login' -Method POST -ContentType 'application/json' -Body '{\"email\":\"test@example.com\",\"password\":\"password123\"}'" -ForegroundColor Gray
    exit 1
}

$headers = @{ Authorization = "Bearer $TOKEN"; "Content-Type" = "application/json" }

# Test the new month endpoint
Write-Host "📅 Testing GET /api/meals/month/2025/8" -ForegroundColor Cyan
Write-Host ""

try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/meals/month/2025/8" -Method GET -Headers $headers
    Write-Host "✅ Success (200):" -ForegroundColor Green
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Error:" -ForegroundColor Red
    Write-Host "Status: $($_.Exception.Response.StatusCode)" -ForegroundColor Red
    Write-Host "Message: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host ""

# Test with invalid parameters
Write-Host "❌ Testing with invalid year" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/meals/month/invalid/8" -Method GET -Headers $headers
    Write-Host "Unexpected success:" -ForegroundColor Yellow
    $response | ConvertTo-Json
} catch {
    Write-Host "✅ Correctly returned error:" -ForegroundColor Green
    Write-Host "Status: $($_.Exception.Response.StatusCode)" -ForegroundColor White
}

Write-Host ""

Write-Host "❌ Testing with invalid month" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/meals/month/2025/13" -Method GET -Headers $headers
    Write-Host "Unexpected success:" -ForegroundColor Yellow
    $response | ConvertTo-Json
} catch {
    Write-Host "✅ Correctly returned error:" -ForegroundColor Green
    Write-Host "Status: $($_.Exception.Response.StatusCode)" -ForegroundColor White
}

Write-Host ""
Write-Host "✅ Month endpoint testing complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Expected behavior:" -ForegroundColor Magenta
Write-Host "- Valid month request should return 200 with meals data" -ForegroundColor White
Write-Host "- Invalid year should return 400 with error message" -ForegroundColor White
Write-Host "- Invalid month should return 400 with error message" -ForegroundColor White
