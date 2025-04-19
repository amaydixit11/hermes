-- Migration: Create users table

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('ADMIN', 'USER', 'GUEST');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role user_role NOT NULL DEFAULT 'USER',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add indexes to improve query performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Create an initial admin user with password 'admin' (in production, use a secure password and change immediately)
INSERT INTO users (
    username, 
    email, 
    password_hash, 
    first_name, 
    last_name, 
    role
) VALUES (
    'admin',
    'admin@hermes.local',
    '$2a$10$EYQyxu4t5ntk3RIUnUQ1A.V5YkBUNHF3Yk8mA9/tpHSgY7tBrchX2', -- bcrypt hash for 'admin'
    'Admin',
    'User',
    'ADMIN'
);