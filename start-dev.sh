#!/bin/bash

# Development startup script for Nourish App
echo "🍽️  Starting Nourish Development Environment..."

# Check if .env file exists in backend
if [ ! -f "backend/.env" ]; then
    echo "⚠️  Creating backend/.env file..."
    cp backend/.env.example backend/.env
    echo "📝 Please update backend/.env with your MongoDB URI and other settings"
fi

# Install backend dependencies
echo "📦 Installing backend dependencies..."
cd backend && go mod tidy

# Install frontend dependencies  
echo "📦 Installing frontend dependencies..."
cd ../frontend && npm install

# Start backend in background
echo "🚀 Starting backend server..."
cd ../backend && go run cmd/server/main.go &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 5

# Start frontend
echo "🎨 Starting frontend development server..."
cd ../frontend && npm start &
FRONTEND_PID=$!

echo "✅ Development servers started!"
echo "🌐 Frontend: http://localhost:3000"
echo "🔧 Backend: http://localhost:5000"
echo "📚 API Health: http://localhost:5000/api/health"
echo ""
echo "Press Ctrl+C to stop all services"

# Function to cleanup on exit
cleanup() {
    echo "🛑 Stopping development servers..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit 0
}

# Trap Ctrl+C and call cleanup
trap cleanup INT

# Wait for both processes
wait $BACKEND_PID $FRONTEND_PID
