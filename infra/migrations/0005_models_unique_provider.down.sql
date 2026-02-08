-- 0005_models_unique_provider.down.sql
-- Revert model uniqueness to name only

BEGIN;

DROP INDEX IF EXISTS idx_models_name_provider;
CREATE UNIQUE INDEX IF NOT EXISTS models_name_key ON models (name);

COMMIT;
