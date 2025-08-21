# Go Backend Setup and Build Script for Windows
# This script automates the setup and build process for the Go backend

Write-Host "🚀 Go Backend Setup Script" -ForegroundColor Green
Write-Host "==========================" -ForegroundColor Green

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    Write-Host "After installation, restart your terminal and try again." -ForegroundColor Yellow
    exit 1
}

$goVersion = go version
Write-Host "✅ Go version: $goVersion" -ForegroundColor Green

# Navigate to backend-go directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptPath
Write-Host "📁 Working directory: $(Get-Location)" -ForegroundColor Cyan

# Clean previous builds
Write-Host "🧹 Cleaning previous builds..." -ForegroundColor Yellow
if (Test-Path "bin") {
    Remove-Item -Recurse -Force bin
}

# Clean module cache if requested
if ($args[0] -eq "--clean-cache") {
    Write-Host "🧹 Cleaning Go module cache..." -ForegroundColor Yellow
    go clean -modcache
}

# Remove go.sum to force fresh dependency resolution
if (Test-Path "go.sum") {
    Write-Host "🗑️  Removing go.sum for fresh dependency resolution..." -ForegroundColor Yellow
    Remove-Item go.sum
}

# Download and tidy dependencies
Write-Host "📦 Downloading dependencies..." -ForegroundColor Cyan
go mod download

Write-Host "🔧 Tidying dependencies..." -ForegroundColor Cyan
go mod tidy

# Verify dependencies
Write-Host "✅ Verifying dependencies..." -ForegroundColor Cyan
go mod verify

# Build the application
Write-Host "🔨 Building application..." -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path bin | Out-Null
go build -o bin/server.exe cmd/server/main.go

Write-Host "✅ Build successful!" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 To run the server:" -ForegroundColor Yellow
Write-Host "   .\bin\server.exe" -ForegroundColor White
Write-Host ""
Write-Host "🔧 Or run in development mode:" -ForegroundColor Yellow
Write-Host "   go run cmd/server/main.go" -ForegroundColor White
Write-Host ""
Write-Host "📋 Environment setup:" -ForegroundColor Yellow
Write-Host "   1. Copy .env.example to .env" -ForegroundColor White
Write-Host "   2. Update database connection settings" -ForegroundColor White
Write-Host "   3. Set JWT_SECRET" -ForegroundColor White
Write-Host ""
Write-Host "🎉 Setup complete!" -ForegroundColor Green
