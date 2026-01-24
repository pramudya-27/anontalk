#!/bin/bash

echo "Setting up AnonTalk Backend..."

# Check if MySQL is running
if ! command -v mysql &> /dev/null; then
    echo "MySQL not found. Please install MySQL first."
    exit 1
fi

# Create .env file from template
if [ ! -f .env ]; then
    cp .env.example .env
    echo "Created .env file. Please update DATABASE_URL and JWT_SECRET"
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download

echo "Setup complete!"
echo ""
echo "To start the backend:"
echo "1. Update .env with your database credentials"
echo "2. Run: go run main.go"
echo ""
echo "Or use Docker:"
echo "docker-compose up -d"
