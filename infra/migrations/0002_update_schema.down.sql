-- 0002_update_schema.down.sql

BEGIN;

DROP INDEX IF EXISTS idx_conversations_org_project_updated;
DROP INDEX IF EXISTS idx_usage_records_org_project;
DROP INDEX IF EXISTS idx_usage_records_org_api_key;
DROP INDEX IF EXISTS idx_api_keys_last_used;
DROP INDEX IF EXISTS idx_api_keys_org_status;
DROP INDEX IF EXISTS idx_api_keys_key_hash;
DROP INDEX IF EXISTS idx_transactions_org_ref;

-- restore global ref_id uniqueness
ALTER TABLE transactions ADD CONSTRAINT transactions_ref_id_key UNIQUE (ref_id);

COMMIT;
