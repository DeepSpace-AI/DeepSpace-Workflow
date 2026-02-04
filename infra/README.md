# Infra - Database (GORM)

This project uses GORM AutoMigrate for database schema management.

## Required env
Set DB parts (the app composes DSN from these):
- `DB_HOST`
- `DB_PORT` (default: 5432)
- `DB_USER`
- `DB_PASSWORD` (optional)
- `DB_NAME`
- `DB_SSLMODE` (default: disable)

Example:
```
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="deepspace"
export DB_SSLMODE="disable"
```

## Migration controls
- `DB_AUTO_MIGRATE` (default: true)
- `DB_RESET_ON_START` (default: false)

When `DB_RESET_ON_START=true`, Gateway will drop all tables and recreate them on startup.
