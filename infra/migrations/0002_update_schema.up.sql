-- 0002_update_schema.up.sql
-- Adjust constraints and add indexes for multi-tenant and query performance

BEGIN;

-- transactions: make ref_id unique per org, not globally
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_ref_id_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_org_ref ON transactions (org_id, ref_id);

-- api_keys: ensure hashed keys are unique and add lookup indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys (key_hash);
CREATE INDEX IF NOT EXISTS idx_api_keys_org_status ON api_keys (org_id, status);
CREATE INDEX IF NOT EXISTS idx_api_keys_last_used ON api_keys (last_used_at DESC);

-- usage_records: add indexes for common filters
CREATE INDEX IF NOT EXISTS idx_usage_records_org_api_key ON usage_records (org_id, api_key_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_org_project ON usage_records (org_id, project_id);

-- conversations: list per project by updated time
CREATE INDEX IF NOT EXISTS idx_conversations_org_project_updated ON conversations (org_id, project_id, updated_at DESC);

COMMIT;
