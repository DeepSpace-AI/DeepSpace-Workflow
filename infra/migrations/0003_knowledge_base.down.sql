-- 0003_knowledge_base.down.sql
-- Rollback knowledge base and document management

BEGIN;

DROP TABLE IF EXISTS knowledge_documents;
DROP TABLE IF EXISTS knowledge_bases;

COMMIT;
