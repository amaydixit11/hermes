CREATE TABLE IF NOT EXISTS service_dependencies (
    id SERIAL PRIMARY KEY,
    service_id VARCHAR NOT NULL,
    dependency_id VARCHAR NOT NULL,
    dependency_type VARCHAR NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT idx_service_dependency_service_id_dependency_id UNIQUE (service_id, dependency_id)
);

CREATE INDEX idx_service_dependencies_service_id ON service_dependencies(service_id);
CREATE INDEX idx_service_dependencies_dependency_id ON service_dependencies(dependency_id);

CREATE TABLE IF NOT EXISTS service_versions (
    id SERIAL PRIMARY KEY,
    service_id VARCHAR,
    version VARCHAR NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    endpoint VARCHAR NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_service_versions_service_id ON service_versions(service_id);
