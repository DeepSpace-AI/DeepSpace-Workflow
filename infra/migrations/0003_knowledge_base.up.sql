-- 0003_knowledge_base.up.sql
-- Knowledge base and document management

BEGIN;

CREATE TABLE IF NOT EXISTS knowledge_bases (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  project_id BIGINT REFERENCES projects(id) ON DELETE SET NULL,
  scope TEXT NOT NULL CHECK (scope IN ('org', 'project')),
  name TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS knowledge_documents (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  project_id BIGINT REFERENCES projects(id) ON DELETE SET NULL,
  knowledge_base_id BIGINT NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
  file_name TEXT NOT NULL,
  content_type TEXT,
  size_bytes BIGINT,
  storage_path TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('uploaded','failed')),
  metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_knowledge_bases_user_scope_project_created
  ON knowledge_bases (user_id, scope, project_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_knowledge_documents_kb_created
  ON knowledge_documents (knowledge_base_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_knowledge_documents_user_project
  ON knowledge_documents (user_id, project_id);

COMMIT;
