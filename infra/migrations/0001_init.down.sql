-- 0001_init.down.sql

BEGIN;

DROP INDEX IF EXISTS idx_messages_conversation_created;
DROP INDEX IF EXISTS idx_audit_logs_trace;
DROP INDEX IF EXISTS idx_audit_logs_user_created;
DROP INDEX IF EXISTS idx_transactions_user_created;
DROP INDEX IF EXISTS idx_usage_records_trace;
DROP INDEX IF EXISTS idx_usage_records_user_created;

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS usage_records;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;

COMMIT;
