-- 0006_plan_billing.up.sql
-- 套餐计费相关表

BEGIN;

CREATE TABLE IF NOT EXISTS plans (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  currency TEXT NOT NULL DEFAULT 'USD',
  billing_mode TEXT NOT NULL CHECK (billing_mode IN ('token','request')),
  price_input NUMERIC(20, 6) NOT NULL DEFAULT 0,
  price_output NUMERIC(20, 6) NOT NULL DEFAULT 0,
  price_request NUMERIC(20, 6) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS plan_model_prices (
  id BIGSERIAL PRIMARY KEY,
  plan_id BIGINT NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
  model_id UUID NOT NULL REFERENCES models(id) ON DELETE CASCADE,
  currency TEXT NOT NULL DEFAULT 'USD',
  price_input NUMERIC(20, 6) NOT NULL DEFAULT 0,
  price_output NUMERIC(20, 6) NOT NULL DEFAULT 0,
  price_request NUMERIC(20, 6) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS plan_subscriptions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  plan_id BIGINT NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
  status TEXT NOT NULL DEFAULT 'active',
  start_at TIMESTAMPTZ NOT NULL,
  end_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_plan_model_prices_plan_model
  ON plan_model_prices (plan_id, model_id);
CREATE INDEX IF NOT EXISTS idx_plan_subscriptions_user_status
  ON plan_subscriptions (user_id, status);
CREATE INDEX IF NOT EXISTS idx_plan_subscriptions_user_start_end
  ON plan_subscriptions (user_id, start_at, end_at);
CREATE INDEX IF NOT EXISTS idx_plans_status
  ON plans (status);

COMMIT;
