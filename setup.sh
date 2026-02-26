#!/bin/bash

echo "🏋️  Workout Tracker Setup"
echo "========================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.23 or later."
    exit 1
fi
echo "✅ Go $(go version | awk '{print $3}') found"

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or later."
    exit 1
fi
echo "✅ Node.js $(node --version) found"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker."
    exit 1
fi
echo "✅ Docker found"

echo ""
echo "📦 Installing backend dependencies..."
go mod download

echo ""
echo "📦 Installing frontend dependencies..."
cd frontend
npm install
cd ..

echo ""
echo "✅ Setup complete!"
echo ""
echo "🚀 Quick Start:"
echo ""
echo "1. Start PostgreSQL:"
echo "   docker-compose up postgres -d"
echo ""
echo "2. Run backend (in one terminal):"
echo "   go run main.go"
echo ""
echo "3. Run frontend (in another terminal):"
echo "   cd frontend && npm run dev"
echo ""
echo "Or use the Makefile:"
echo "   make docker-up       # Start PostgreSQL"
echo "   make dev             # Run backend"
echo "   make dev-frontend    # Run frontend"
echo ""
echo "📚 Visit http://localhost:3000 to see your app!"
