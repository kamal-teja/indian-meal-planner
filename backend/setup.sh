#!/bin/bash

# Go Backend Setup and Build Script
# This script automates the setup and build process for the Go backend

set -e

echo "🚀 Go Backend Setup Script"
echo "=========================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed or not in PATH"
    echo "Please install Go from https://golang.org/dl/"
    echo "Or use your package manager:"
    echo "  - Windows: Download and install from golang.org"
    echo "  - macOS: brew install go"
    echo "  - Linux: sudo apt install golang-go"
    exit 1
fi

echo "✅ Go version: $(go version)"

# Navigate to backend-go directory
cd "$(dirname "$0")"
echo "📁 Working directory: $(pwd)"

# Clean previous builds
echo "🧹 Cleaning previous builds..."
if [ -d "bin" ]; then
    rm -rf bin
fi

# Clean module cache if requested
if [ "$1" = "--clean-cache" ]; then
    echo "🧹 Cleaning Go module cache..."
    go clean -modcache
fi

# Remove go.sum to force fresh dependency resolution
if [ -f "go.sum" ]; then
    echo "🗑️  Removing go.sum for fresh dependency resolution..."
    rm go.sum
fi

# Download and tidy dependencies
echo "📦 Downloading dependencies..."
go mod download

echo "🔧 Tidying dependencies..."
go mod tidy

# Verify dependencies
echo "✅ Verifying dependencies..."
go mod verify

# Build the application
echo "🔨 Building application..."
mkdir -p bin
go build -o bin/server cmd/server/main.go

echo "✅ Build successful!"
echo ""
echo "🚀 To run the server:"
echo "   ./bin/server"
echo ""
echo "🔧 Or run in development mode:"
echo "   go run cmd/server/main.go"
echo ""
echo "📋 Environment setup:"
echo "   1. Copy .env.example to .env"
echo "   2. Update database connection settings"
echo "   3. Set JWT_SECRET"
echo ""
echo "🎉 Setup complete!"
