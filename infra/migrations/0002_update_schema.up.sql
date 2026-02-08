-- 0002_update_schema.up.sql
-- Adjust constraints and add indexes for multi-tenant and query performance

BEGIN;

-- transactions: make ref_id unique per user, not globally
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_ref_id_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_user_ref ON transactions (user_id, ref_id);

-- api_keys: ensure hashed keys are unique and add lookup indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys (key_hash);
CREATE INDEX IF NOT EXISTS idx_api_keys_user_status ON api_keys (user_id, status);
CREATE INDEX IF NOT EXISTS idx_api_keys_last_used ON api_keys (last_used_at DESC);

-- usage_records: add indexes for common filters
CREATE INDEX IF NOT EXISTS idx_usage_records_user_project ON usage_records (user_id, project_id);

-- conversations: list per project by updated time
CREATE INDEX IF NOT EXISTS idx_conversations_user_project_updated ON conversations (user_id, project_id, updated_at DESC);

COMMIT;
