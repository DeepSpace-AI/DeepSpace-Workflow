-- 0002_update_schema.down.sql

BEGIN;

DROP INDEX IF EXISTS idx_conversations_user_project_updated;
DROP INDEX IF EXISTS idx_usage_records_user_project;
DROP INDEX IF EXISTS idx_api_keys_last_used;
DROP INDEX IF EXISTS idx_api_keys_user_status;
DROP INDEX IF EXISTS idx_api_keys_key_hash;
DROP INDEX IF EXISTS idx_transactions_user_ref;

-- restore global ref_id uniqueness
ALTER TABLE transactions ADD CONSTRAINT transactions_ref_id_key UNIQUE (ref_id);

COMMIT;
