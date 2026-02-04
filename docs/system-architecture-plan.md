# DeepSpace Workflows — 系统分层与实现清单（模块级）

## 0. 平台定位

DeepSpace Workflows 是一个 面向科研与研发场景的 AI Workflow / Research OS，是：
以 AI 为执行引擎、以 Workflow 为组织形式、以计费与审计为底座的科研基础设施。
核心关键词：
Workflow First
Project Oriented
Audit & Billing Native
AI Infrastructure

> 目标：在固定架构（Web → Gateway → NewAPI）不变的前提下，补全可商用部署所需的系统分层、实现清单、Admin 信息架构，以及 Gateway 数据表细化设计。
>
> 技术栈固定：Web（React + Ant Design + Ant Design X + Zustand + SSE）、Gateway/Worker（Go）、数据库（PostgreSQL）、缓存/队列（Redis）。

---

## 1. 总体分层（固定链路）

```
User Web (Chat / Workflow UI)
   → Gateway (Auth / Billing / Policy / Pipeline / SSE)
   → NewAPI → Upstream Models
Admin Web (Ops / Billing / Policy)
   → Admin API (Go, 与Gateway不同路由)
Worker (Async / RAG / Workflow Tasks)
Storage (LocalFS / MinIO / S3 Compatible)
```

**关键原则**
- Gateway 是唯一 AI 出口，前端/管理端禁止直连 NewAPI。
- 商用闭环优先：计费幂等、冻结/扣费/解冻、审计追踪、风控限额。
- Pipeline 可插拔（RAG/MCP/Skills/Workflow/多模态）不可写死逻辑。


**组织与项目架构**
```
Organization
└─ Project
├─ Chats
├─ Workflows
├─ Knowledge Base
├─ Files / Artifacts
├─ API Keys
└─ Usage / Billing
```

---

## 2. 系统分层与实现清单（模块级任务）

### 2.1 Web（用户端）

**目标：** 商用 Chat MVP + 计费与项目视图。

**模块任务清单：**
1. Chat UI（Ant Design X）
   - MessageList / MessageItem / PromptInput / AgentThought 组件
   - 保持纯展示、无副作用
2. Chat Hooks（SSE + 业务）
   - `useChatStream.ts`：SSE 连接、流式处理
   - `useChat.ts`：会话/消息状态与 Gateway API 调用
3. Gateway Client
   - 统一 API Client（仅允许调用 Gateway）
4. 状态管理（Zustand）
   - chatStore / userStore / billingStore
5. 路由与页面
   - Chat / Billing / Projects / Profile
6. 商用基础能力
   - 登录、API Key 管理入口
   - Wallet 与流水视图

**约束：**
- SSE 逻辑仅存在于 hooks（`useChat.ts`/`useChatStream.ts`）。
- Ant Design X 组件仅放在 `apps/web/src/components/ai/`。
- 前端只允许调用 Gateway API。

---

### 2.2 Admin Web（管理端）

**目标：** 商用运营必备能力与审计可追溯。

**模块任务清单：**
1. 用户与组织
   - 用户列表、组织列表、角色与权限
2. 模型与定价
   - 模型上架、价格配置
3. 计费与对账
   - Wallets / Transactions / Usage Records
4. 风控策略
   - 速率限制、预算上限、IP 白名单/黑名单
5. 审计与追踪
   - trace_id 追溯，系统请求日志检索

---

### 2.3 Gateway（Go）

**目标：** 统一鉴权、计费、审计、SSE 透传与 Pipeline 编排。

**模块任务清单：**
1. HTTP 层
   - 鉴权中间件（JWT / API Key / Org）
   - SSE 路由
2. Service 层
   - Chat Completion 业务
   - Billing（hold/capture/release）
3. Domain 层
   - Wallet / Transaction / Usage / Conversation / Message
4. Repository 层
   - Postgres 读写、事务与幂等支持
5. AI 层
   - NewAPI 适配与 SSE 透传
6. Pipeline
   - 可插拔步骤：RAG / MCP / Skills / Workflow / 多模态
7. Audit
   - trace_id 全链路贯穿（请求、计费、日志）

---

### 2.4 Worker（异步任务）

**目标：** 与在线请求解耦的异步能力。

**模块任务清单：**
- usage 聚合与对账
- 导出（CSV/JSON/Markdown）
- 索引构建（为 RAG 预留）

---

## 3. Admin 端信息架构（IA）

### 3.1 导航结构（推荐）

```
Admin
├─ Dashboard
├─ Users
│  ├─ Users List
│  ├─ Organizations
│  └─ Roles & Permissions
├─ Models
│  ├─ Models Catalog
│  └─ Pricing Rules
├─ Billing
│  ├─ Wallets
│  ├─ Transactions
│  └─ Usage Records
├─ Audit
│  ├─ Trace Lookup
│  └─ Request Logs
├─ Risk & Policy
│  ├─ Rate Limits
│  ├─ Budget Caps
│  └─ IP Allow/Deny
└─ System
   ├─ API Keys
   └─ Webhooks (optional)
```

### 3.2 页面与功能说明
- Dashboard：整体使用量、成本趋势、异常提示
- Users/Orgs：用户状态、组织归属、角色与权限
- Models/Pricing：模型上架、输入/输出单价
- Billing：钱包余额、冻结、流水、usage 记录
- Audit：trace_id 检索、请求全链路追溯
- Risk：限速、预算上限、IP 访问控制

---

## 4. Gateway 数据表设计（细化）

### 4.1 用户与组织

**users**
- id (PK)
- email
- password_hash
- status
- created_at

**orgs**
- id (PK)
- name
- owner_user_id (FK)
- created_at

**org_members**
- org_id (FK)
- user_id (FK)
- role

**api_keys**
- id (PK)
- org_id (FK)
- key_hash
- scopes
- created_at

---

### 4.2 计费核心（商用闭环）

**wallets**
- org_id (PK/FK)
- balance
- frozen_balance
- updated_at

**transactions**
- id (PK)
- org_id (FK)
- type (hold/capture/release/refund)
- amount
- ref_id (billing_ref_id, unique)
- created_at

**prices**
- model (PK)
- input_price
- output_price
- effective_at

**usage_records**
- id (PK)
- org_id (FK)
- model
- input_tokens
- output_tokens
- cost
- trace_id
- created_at

---

### 4.3 会话与消息

**conversations**
- id (PK)
- org_id (FK)
- user_id (FK)
- title
- created_at

**messages**
- id (PK)
- conversation_id (FK)
- role (user/assistant/system)
- content_json
- created_at

---

### 4.4 审计与追踪（建议）

**audit_logs**
- id (PK)
- org_id (FK)
- trace_id
- event_type
- payload_json
- created_at

---

## 5. 实施顺序建议（MVP）

1. Gateway 基础链路 + SSE
2. 计费闭环（Wallet / Transactions / Usage）
3. Web Chat MVP（Ant Design X）
4. Admin 基础功能（Models / Pricing / Billing）
5. Pipeline 扩展能力（RAG / MCP / Skills / Workflow / 多模态）

---

## 6. RAG + MCP + Skills + Workflow 扩展方案（细化）

### 6.1 Pipeline 扩展接口（概念）

```
Auth → RateLimit → Budget Hold → Policy
→ Context Build (RAG)
→ Tool/Skill Dispatch (MCP/Skills)
→ Workflow Orchestration
→ NewAPI Call
→ Usage Capture → Stream Response
```

**要点：**
- RAG、MCP、Skills、Workflow 均通过 Pipeline 插件步骤接入。
- Gateway 统一管理 trace_id 与 billing_ref_id，确保审计与计费不丢失。
- Workflow 负责多步骤编排，Skills 负责单点工具能力。

---

### 6.2 RAG（检索增强生成）

**能力目标：**
- 支持项目级/组织级知识库。
- 引用来源可追溯（citation），保存到消息元数据。

**模块任务清单：**
1. 文档接入与解析（Worker 异步）
2. 向量化与索引（Worker）
3. 检索与重排（Gateway Pipeline）
4. Context Build（Gateway Pipeline）
5. 引用落库（usage_records/消息元数据）

---

### 6.3 MCP（工具接入）

**能力目标：**
- Gateway 作为工具路由入口，统一审计和计费。
- 工具能力以“可挂载适配器”形式扩展。

**模块任务清单：**
1. MCP Client 适配层（Gateway）
2. 工具注册与版本管理（Admin）
3. 调用记录与审计（usage_records + audit_logs）
4. 失败重试与幂等策略

---

### 6.4 Skills（原子能力）

**能力目标：**
- 技能是可复用的原子能力（如：摘要、翻译、SQL 生成）。
- Skills 可供 Workflow/Agent 编排。

**模块任务清单：**
1. Skills 目录与元数据（Admin）
2. Skills 路由（Gateway Pipeline）
3. 技能计费与调用审计

---

### 6.5 Workflow（多步编排）

**能力目标：**
- 面向复杂任务的多步流程编排（计划 → 执行 → 验证）。
- 允许串联 RAG、MCP、Skills。

**模块任务清单：**
1. Workflow 定义（JSON/YAML 形态存储）
2. Workflow 执行引擎（Gateway/Worker）
3. 运行状态与可追溯记录（audit_logs）
4. 失败处理与可恢复机制
