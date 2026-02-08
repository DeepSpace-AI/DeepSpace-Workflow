# 套餐计费对接说明（Admin）

> 目标：新增套餐管理能力，支持按 Token 或按请求计费；当组织有生效套餐时，优先使用套餐资费进行扣费；无套餐时回退模型默认价格。

## 1. 术语与原则

- 套餐（Plan）：定义计费模式与费率
- 订阅（Subscription）：组织与套餐的生效关系
- 计费优先级：生效套餐 > 模型默认价
- 同一组织：仅允许一个 active 订阅

## 2. 计费模式

支持两种模式，按套餐配置：

1) Token 计费
- 使用上游 usage 中的 token 数据
- 费用计算：
  - cost = prompt_tokens/1,000,000 * price_input + completion_tokens/1,000,000 * price_output

2) 请求计费
- 每次请求固定价格（按模型区分）
- 成功响应即扣费（失败按现有策略 release）

## 3. 数据结构（建议）

### 3.1 plans（套餐定义）
- id
- name
- status（active/disabled）
- currency
- billing_mode（token/request）
- price_input（token 模式）
- price_output（token 模式）
- price_request（request 模式）
- created_at / updated_at

### 3.2 plan_model_prices（套餐-模型价格覆盖）
- plan_id
- model_id
- currency
- price_input（token 模式）
- price_output（token 模式）
- price_request（request 模式）

规则：
- 若某模型存在覆盖价，则使用覆盖价
- 未配置覆盖价时，使用 plans 默认价格

### 3.3 plan_subscriptions（组织订阅）
- id
- org_id
- plan_id
- status（active/expired/canceled）
- start_at
- end_at
- created_at / updated_at

约束：
- 同一 org 仅允许一个 active 订阅

## 4. 计费流程（Gateway）

1) 解析上游 usage（已接入）
2) 查询 org 当前生效订阅
3) 若存在订阅：
   - 按 billing_mode 选择 token/request 计费
   - 优先使用 plan_model_prices 覆盖价
4) 若不存在订阅：
   - 回退模型默认价
5) 记录 usage 与扣费

## 5. Admin 接口（建议）

### 5.1 套餐管理
- GET /api/admin/plans
- POST /api/admin/plans
- PATCH /api/admin/plans/{id}

### 5.2 套餐订阅
- GET /api/admin/subscriptions
- POST /api/admin/subscriptions
- PATCH /api/admin/subscriptions/{id}

### 5.3 组织当前订阅
- GET /api/admin/orgs/{id}/subscription

## 6. 前端对接要点

- 套餐支持两种计费模式（token/request）
- request 模式下需展示每次请求价格
- token 模式下需展示 input/output 费率
- 提供模型覆盖价的配置入口（可选）
- 订阅有效期、状态切换与冲突提示（同 org 仅 1 个 active）

## 7. 兼容与回退

- 未绑定套餐时：使用模型默认价
- 订阅过期或取消：立即回退模型默认价

## 8. 下一步落地顺序

1) 数据迁移（plans / plan_model_prices / plan_subscriptions）
2) Repo/Service（套餐、订阅、价格覆盖）
3) 计费逻辑接入（UsageCapture）
4) Admin API + OpenAPI 文档
