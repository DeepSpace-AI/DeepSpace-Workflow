-- 0004_models.up.sql
-- Models for pricing and capabilities

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS models (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  provider TEXT NOT NULL,
  price_input NUMERIC(20, 6) NOT NULL DEFAULT 0,
  price_output NUMERIC(20, 6) NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'USD',
  capabilities JSONB NOT NULL DEFAULT '[]',
  status TEXT NOT NULL DEFAULT 'active',
  metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_models_provider ON models (provider);
CREATE INDEX IF NOT EXISTS idx_models_status ON models (status);

COMMIT;
