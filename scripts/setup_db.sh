#!/bin/bash

# PostgreSQL Database Setup Script for Inventory System

echo "Setting up PostgreSQL database for Inventory System..."

# Database configuration from environment or defaults
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-inventory_system}

echo "Database Configuration:"
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "User: $DB_USER"
echo "Database: $DB_NAME"
echo ""

# Check if PostgreSQL is running
echo "Checking PostgreSQL connection..."
pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER

if [ $? -ne 0 ]; then
    echo "❌ PostgreSQL is not running or not accessible."
    echo "Please ensure PostgreSQL is installed and running."
    echo ""
    echo "Installation instructions:"
    echo "macOS: brew install postgresql && brew services start postgresql"
    echo "Ubuntu: sudo apt-get install postgresql postgresql-contrib"
    echo "CentOS/RHEL: sudo yum install postgresql-server postgresql-contrib"
    echo "Docker: docker run --name postgres -e POSTGRES_PASSWORD=$DB_PASSWORD -p 5432:5432 -d postgres"
    exit 1
fi

echo "✅ PostgreSQL is running!"

# Create database if it doesn't exist
echo "Creating database '$DB_NAME' if it doesn't exist..."

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME"

if [ $? -eq 0 ]; then
    echo "✅ Database '$DB_NAME' is ready!"
else
    echo "❌ Failed to create database '$DB_NAME'"
    exit 1
fi

# Test connection to the new database
echo "Testing connection to database..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version();" > /dev/null

if [ $? -eq 0 ]; then
    echo "✅ Successfully connected to database '$DB_NAME'!"
    echo ""
    echo "Database setup complete! You can now run:"
    echo "  go run main.go"
    echo ""
    echo "The application will automatically create the required tables on first run."
else
    echo "❌ Failed to connect to database '$DB_NAME'"
    exit 1
fi
