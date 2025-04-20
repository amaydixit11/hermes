-- Migration: create_health_checks_table
-- Down migration SQL

DROP TABLE IF EXISTS health_checks;

-- Migration: create_health_history_table
-- Down migration SQL

DROP TABLE IF EXISTS health_history;