-- Initial database setup for Inventory System
-- This file is automatically executed when the PostgreSQL container starts

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
SET timezone = 'UTC';

-- Create indexes for better performance (will be created after tables are migrated by GORM)
-- These are just placeholders and will be created by the application

-- Database ready message
SELECT 'Inventory System PostgreSQL database initialized successfully!' as message;
