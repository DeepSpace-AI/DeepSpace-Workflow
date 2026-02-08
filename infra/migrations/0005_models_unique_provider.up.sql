-- 0005_models_unique_provider.up.sql
-- Ensure model uniqueness on (name, provider)

BEGIN;

DROP INDEX IF EXISTS models_name_key;
DROP INDEX IF EXISTS idx_models_name;
DROP INDEX IF EXISTS idx_models_unique_name;

CREATE UNIQUE INDEX IF NOT EXISTS idx_models_name_provider ON models (name, provider);

COMMIT;
