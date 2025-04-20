-- Migration: create_health_checks_table
-- Up migration SQL

-- First check if the health_checks table doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'health_checks') THEN
        CREATE TABLE health_checks (
            id SERIAL PRIMARY KEY,
            service_id VARCHAR(255) NOT NULL,
            name VARCHAR(255) NOT NULL,
            type VARCHAR(50) NOT NULL,
            endpoint VARCHAR(255),
            interval INTEGER DEFAULT 60,
            timeout INTEGER DEFAULT 5,
            method VARCHAR(10) DEFAULT 'GET',
            expected_status INTEGER,
            expected_body TEXT,
            headers JSONB,
            retries INTEGER DEFAULT 1,
            threshold_count INTEGER DEFAULT 3,
            timeout_count INTEGER DEFAULT 0,
            enabled BOOLEAN DEFAULT TRUE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        );

        -- Create indexes
        CREATE INDEX idx_health_checks_service_id ON health_checks(service_id);
    END IF;
END
$$;

-- Migration: create_health_history_table
-- Up migration SQL

-- First check if the health_history table doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE tablename = 'health_history') THEN
        CREATE TABLE health_history (
            id SERIAL PRIMARY KEY,
            service_id VARCHAR(255) NOT NULL,
            check_id INTEGER,
            status VARCHAR(50) NOT NULL,
            message TEXT,
            response_time_ms INTEGER,
            status_code INTEGER,
            details JSONB,
            timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            
            CONSTRAINT fk_health_history_check_id FOREIGN KEY (check_id)
                REFERENCES health_checks(id) ON DELETE SET NULL
        );

        -- Create indexes
        CREATE INDEX idx_health_history_service_id ON health_history(service_id);
        CREATE INDEX idx_health_history_check_id ON health_history(check_id);
        CREATE INDEX idx_health_history_timestamp ON health_history(timestamp);
        CREATE INDEX idx_health_history_status ON health_history(status);
    END IF;
END
$$;