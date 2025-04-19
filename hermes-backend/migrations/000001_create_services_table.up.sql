-- Migration: Create services table

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE service_status AS ENUM ('HEALTHY', 'UNHEALTHY', 'UNKNOWN', 'WARNING');

CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    status service_status NOT NULL DEFAULT 'UNKNOWN',
    type VARCHAR(100),
    endpoint VARCHAR(255) NOT NULL,
    metadata JSONB DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add indexes to improve query performance
CREATE INDEX idx_services_status ON services(status);
CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_services_type ON services(type);
CREATE INDEX idx_services_tags ON services USING GIN(tags);