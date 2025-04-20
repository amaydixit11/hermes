#!/bin/bash

# Check if command is provided
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <command> [args]"
    echo "Commands:"
    echo "  up        - Run all migrations"
    echo "  down      - Rollback last migration"
    echo "  reset     - Rollback all migrations and apply them again"
    echo "  create    - Create a new migration (requires name argument)"
    echo "Example: $0 create add_status_column"
    exit 1
fi

COMMAND=$1
PROJECT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
MIGRATIONS_DIR="$PROJECT_ROOT/migrations"

# Ensure migrations directory exists
mkdir -p "$MIGRATIONS_DIR"

case $COMMAND in
    up)
        echo "Running migrations..."
        go run "$PROJECT_ROOT/cmd/api/main.go" -migrate -migrations "$MIGRATIONS_DIR"
        ;;
    down)
        echo "Rolling back last migration..."
        go run "$PROJECT_ROOT/cmd/api/main.go" -rollback -migrations "$MIGRATIONS_DIR"
        ;;
    reset)
        echo "Resetting database (rollback all and migrate up)..."
        # You might want to implement this specific command in your Go code
        go run "$PROJECT_ROOT/cmd/api/main.go" -rollback-all -migrations "$MIGRATIONS_DIR"
        go run "$PROJECT_ROOT/cmd/api/main.go" -migrate -migrations "$MIGRATIONS_DIR"
        ;;
    create)
        if [ "$#" -ne 2 ]; then
            echo "Migration name required."
            echo "Example: $0 create add_status_column"
            exit 1
        fi
        
        MIGRATION_NAME=$2
        VERSION=$(date +%Y%m%d%H%M%S)
        FILE_PREFIX="${VERSION}_${MIGRATION_NAME}"
        
        echo "Creating migration: $FILE_PREFIX"
        
        # Create up migration
        cat > "$MIGRATIONS_DIR/${FILE_PREFIX}.up.sql" << EOF
-- Migration: ${MIGRATION_NAME}
-- Created at: $(date)
-- Up migration SQL

EOF
        
        # Create down migration
        cat > "$MIGRATIONS_DIR/${FILE_PREFIX}.down.sql" << EOF
-- Migration: ${MIGRATION_NAME}
-- Created at: $(date)
-- Down migration SQL

EOF
        
        echo "Created migration files:"
        echo "  - $MIGRATIONS_DIR/${FILE_PREFIX}.up.sql"
        echo "  - $MIGRATIONS_DIR/${FILE_PREFIX}.down.sql"
        ;;
    *)
        echo "Unknown command: $COMMAND"
        echo "Valid commands: up, down, reset, create"
        exit 1
        ;;
esac