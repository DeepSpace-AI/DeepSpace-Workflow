-- 0006_plan_billing.down.sql
-- 回滚套餐计费相关表

BEGIN;

DROP TABLE IF EXISTS plan_subscriptions;
DROP TABLE IF EXISTS plan_model_prices;
DROP TABLE IF EXISTS plans;

COMMIT;
