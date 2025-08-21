#!/bin/bash

# Development setup script for Go backend

set -e

echo "🚀 Setting up Go backend development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [[ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or higher."
    exit 1
fi

echo "✅ Go version $GO_VERSION is compatible"

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy

# Install development tools
echo "🔧 Installing development tools..."

# Air for hot reload
if ! command -v air &> /dev/null; then
    echo "Installing Air for hot reload..."
    go install github.com/air-verse/air@latest
else
    echo "✅ Air is already installed"
fi

# golangci-lint for linting
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
else
    echo "✅ golangci-lint is already installed"
fi

# godoc for documentation
if ! command -v godoc &> /dev/null; then
    echo "Installing godoc..."
    go install golang.org/x/tools/cmd/godoc@latest
else
    echo "✅ godoc is already installed"
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ Created .env file. Please update it with your configuration."
else
    echo "✅ .env file already exists"
fi

# Create necessary directories
mkdir -p tmp
mkdir -p build
mkdir -p logs

echo ""
echo "🎉 Development environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env file with your configuration"
echo "2. Start MongoDB (local or use Atlas URI)"
echo "3. Run 'make dev' to start development server with hot reload"
echo "4. Or run 'make run' to build and run normally"
echo ""
echo "Available commands:"
echo "  make dev      - Start with hot reload"
echo "  make build    - Build the application"
echo "  make test     - Run tests"
echo "  make lint     - Run linter"
echo "  make format   - Format code"
echo "  make help     - Show all available commands"
