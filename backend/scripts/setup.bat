@echo off
REM Development setup script for Go backend (Windows)

echo 🚀 Setting up Go backend development environment...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go 1.21 or higher.
    exit /b 1
)

echo ✅ Go is installed

REM Install dependencies
echo 📦 Installing Go dependencies...
go mod download
go mod tidy

REM Install development tools
echo 🔧 Installing development tools...

REM Air for hot reload
where air >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing Air for hot reload...
    go install github.com/air-verse/air@latest
) else (
    echo ✅ Air is already installed
)

REM golangci-lint for linting
where golangci-lint >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing golangci-lint...
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
) else (
    echo ✅ golangci-lint is already installed
)

REM godoc for documentation
where godoc >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing godoc...
    go install golang.org/x/tools/cmd/godoc@latest
) else (
    echo ✅ godoc is already installed
)

REM Create .env file if it doesn't exist
if not exist .env (
    echo 📝 Creating .env file from template...
    copy .env.example .env
    echo ✅ Created .env file. Please update it with your configuration.
) else (
    echo ✅ .env file already exists
)

REM Create necessary directories
if not exist tmp mkdir tmp
if not exist build mkdir build
if not exist logs mkdir logs

echo.
echo 🎉 Development environment setup complete!
echo.
echo Next steps:
echo 1. Update .env file with your configuration
echo 2. Start MongoDB (local or use Atlas URI)
echo 3. Run 'make dev' to start development server with hot reload
echo 4. Or run 'make run' to build and run normally
echo.
echo Available commands:
echo   make dev      - Start with hot reload
echo   make build    - Build the application
echo   make test     - Run tests
echo   make lint     - Run linter
echo   make format   - Format code
echo   make help     - Show all available commands
