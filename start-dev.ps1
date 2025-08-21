# Development startup script for Meal Planner App (Windows PowerShell)

Write-Host "🍽️  Starting Meal Planner Development Environment..." -ForegroundColor Green

# Check if .env file exists in backend
if (-not (Test-Path "backend\.env")) {
    Write-Host "⚠️  Creating backend\.env file..." -ForegroundColor Yellow
    Copy-Item "backend\.env.example" "backend\.env"
    Write-Host "📝 Please update backend\.env with your MongoDB URI and other settings" -ForegroundColor Cyan
}

# Install backend dependencies
Write-Host "📦 Installing backend dependencies..." -ForegroundColor Blue
Set-Location backend
go mod tidy

# Install frontend dependencies  
Write-Host "📦 Installing frontend dependencies..." -ForegroundColor Blue
Set-Location ..\frontend
npm install

# Start backend in background
Write-Host "🚀 Starting backend server..." -ForegroundColor Green
Set-Location ..\backend

# Start backend as background job
$backendJob = Start-Job -ScriptBlock {
    Set-Location $using:PWD\backend
    go run cmd/server/main.go
}

# Wait a moment for backend to start
Start-Sleep -Seconds 5

# Start frontend
Write-Host "🎨 Starting frontend development server..." -ForegroundColor Green
Set-Location ..\frontend

# Start frontend as background job
$frontendJob = Start-Job -ScriptBlock {
    Set-Location $using:PWD\frontend
    npm start
}

Write-Host "✅ Development servers started!" -ForegroundColor Green
Write-Host "🌐 Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "🔧 Backend: http://localhost:5000" -ForegroundColor Cyan
Write-Host "📚 API Health: http://localhost:5000/api/health" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press Ctrl+C to stop all services" -ForegroundColor Yellow

# Function to cleanup on exit
function Cleanup {
    Write-Host "🛑 Stopping development servers..." -ForegroundColor Red
    Stop-Job $backendJob, $frontendJob -Force
    Remove-Job $backendJob, $frontendJob -Force
}

# Register cleanup function for Ctrl+C
Register-EngineEvent PowerShell.Exiting -Action { Cleanup }

try {
    # Wait for jobs to complete (they won't unless stopped)
    Wait-Job $backendJob, $frontendJob
}
finally {
    Cleanup
}
