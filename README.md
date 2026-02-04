# DeepSpace Workflows

> **DeepSpace Workflows** 是一个面向科研与研发场景的 **可商用 AI Workflow 平台**，以 Chat 为入口，以 Workflow / Pipeline 为核心，以 **计费、审计、可追溯** 为基础能力，构建可长期演进的 AI Research Infrastructure。

---

## 1. 项目定位

DeepSpace Workflows 并不是一个简单的 ChatGPT 替代品，而是一个 **AI 工作流与科研协作平台（AI Research OS）**，核心目标包括：

* 🧪 **科研可复现**：每一次 AI 调用都可追溯（trace_id）
* 💰 **商用可计费**：支持冻结 / 扣费 / 解冻的完整计费闭环
* 🧩 **能力可编排**：RAG、Tools、Skills、Workflow 统一通过 Pipeline 组织
* 🏢 **组织级治理**：用户 / 组织 / 项目 / 权限 / 审计

适用对象：

* 高校科研团队
* 企业 R&D 团队
* AI 应用与平台型产品
* 教育与科研 AI 平台

---

## 2. 核心产品形态

```
Organization
 └─ Project
     ├─ Chat Sessions
     ├─ Knowledge Base (RAG)
     ├─ Workflows
     ├─ API Keys (Scoped)
     └─ Usage & Cost
```

* **Chat**：默认交互入口（DeepSpace Chat）
* **Workflow**：多步 AI 任务编排（DeepSpace Workflows 核心）
* **Pipeline**：可插拔能力链路（RAG / MCP / Skills）
* **Gateway**：唯一 AI 出口，负责鉴权、计费、审计

---

## 3. 总体架构

```
Web (Nuxt + PrimeVue)
  → Gateway (Go: Auth / Billing / SSE / Pipeline)
    → NewAPI → Upstream Models
Admin Web (React + Ant Design)
  → Admin API (Go)
```

**架构原则**：

* Gateway 是唯一 AI 出口
* 前端禁止直连模型或 NewAPI
* 所有调用必须可计费、可审计

---

## 4. 核心能力

### 4.1 Chat（入口能力）

* SSE 流式响应
* 会话 / 消息管理
* Project 隔离

### 4.2 Workflow（核心能力）

* 多步骤 AI 编排（Plan → Execute → Verify）
* 串联 RAG / Skills / Tools
* Workflow 定义支持 JSON / YAML

### 4.3 Pipeline（扩展能力）

```
Auth → RateLimit → Budget Hold → Policy
→ Context Build (RAG)
→ Tool / Skill Dispatch
→ Workflow Orchestration
→ NewAPI Call
→ Usage Capture → Stream Response
```

---

## 5. 计费与审计（商用核心）

* Wallet（余额 / 冻结）
* Transactions（hold / capture / release）
* Usage Records（token / cost / model）
* Audit Logs（trace_id 全链路追踪）

---

## 6. 技术栈

| 层级                  | 技术                                            |
| ------------------- | --------------------------------------------- |
| Web                 | Nuxt, NuxtUI |
| Gateway / Admin API | Go                                            |
| 数据库                 | PostgreSQL                                    |
| 缓存 / 队列             | Redis                                         |
| 对象存储                | LocalFiles/MinIO/S3                           |
| AI 接入               | NewAPI                                        |

---

## 7. Roadmap（简化）

1. Chat + Gateway + Billing MVP
2. Project + RAG（科研最小闭环）
3. Workflow Engine（JSON 定义）
4. MCP / Skills 扩展


## 8.目录结构

```
deepspace-workflows/
├─ README.md
├─ LICENSE
├─ .gitignore
├─ .editorconfig
├─ .env.example
├─ docker-compose.yml
├─ Makefile
│
├─ apps/
│  ├─ web/                         # Nuxt 3 + Nuxt UI（用户端：Chat / Workflow / KB）
│  │  ├─ nuxt.config.ts
│  │  ├─ package.json
│  │  ├─ app.vue
│  │  ├─ pages/
│  │  │  ├─ index.vue
│  │  │  ├─ login.vue
│  │  │  ├─ projects/
│  │  │  │  └─ [projectId].vue
│  │  │  ├─ chat/
│  │  │  │  └─ [projectId].vue
│  │  │  ├─ workflows/
│  │  │  │  ├─ [projectId].vue
│  │  │  │  └─ run-[runId].vue
│  │  │  └─ knowledge/
│  │  │     └─ [projectId].vue
│  │  ├─ components/
│  │  │  ├─ chat/
│  │  │  ├─ workflow/
│  │  │  ├─ kb/
│  │  │  └─ common/
│  │  ├─ composables/              # Nuxt composables（替代 hooks）
│  │  │  ├─ useGatewayClient.ts
│  │  │  ├─ useAuth.ts
│  │  │  ├─ useChatStream.ts
│  │  │  ├─ useProjects.ts
│  │  │  ├─ useWorkflows.ts
│  │  │  └─ useBilling.ts
│  │  ├─ stores/                   # Pinia（建议，用于 auth/project/chat 等）
│  │  │  ├─ auth.ts
│  │  │  ├─ project.ts
│  │  │  ├─ chat.ts
│  │  │  ├─ workflow.ts
│  │  │  └─ billing.ts
│  │  ├─ middleware/
│  │  │  └─ auth.global.ts
│  │  ├─ plugins/
│  │  │  └─ gateway.client.ts      # 可选：注入 $gateway
│  │  ├─ utils/
│  │  └─ assets/
│  │
│  └─ admin/                       # Nuxt 3 + Nuxt UI（管理端：定价/对账/审计/风控）
│     ├─ nuxt.config.ts
│     ├─ package.json
│     ├─ pages/
│     │  ├─ index.vue
│     │  ├─ users.vue
│     │  ├─ orgs.vue
│     │  ├─ projects.vue
│     │  ├─ models.vue
│     │  ├─ pricing.vue
│     │  ├─ billing/
│     │  │  ├─ wallets.vue
│     │  │  ├─ transactions.vue
│     │  │  └─ usage.vue
│     │  ├─ audit.vue
│     │  └─ policy.vue
│     ├─ stores/
│     ├─ composables/
│     └─ components/
│
├─ services/
│  ├─ gateway/                     # Go + Gin（唯一 AI 出口：Auth/Billing/SSE/Pipeline）
│  │  ├─ cmd/
│  │  │  └─ gateway/
│  │  │     └─ main.go
│  │  ├─ internal/
│  │  │  ├─ api/                   # Gin handlers + routes + middleware
│  │  │  │  ├─ routes.go
│  │  │  │  ├─ middleware/
│  │  │  │  │  ├─ auth.go
│  │  │  │  │  ├─ trace.go
│  │  │  │  │  ├─ rate_limit.go
│  │  │  │  │  └─ recover.go
│  │  │  │  ├─ health.go
│  │  │  │  ├─ auth_handlers.go
│  │  │  │  ├─ project_handlers.go
│  │  │  │  ├─ chat_handlers.go
│  │  │  │  ├─ workflow_handlers.go
│  │  │  │  ├─ file_handlers.go
│  │  │  │  └─ admin_handlers.go   # 也可拆到独立 admin 服务
│  │  │  ├─ config/
│  │  │  ├─ domain/                # 领域模型
│  │  │  │  ├─ auth/
│  │  │  │  ├─ projects/
│  │  │  │  ├─ billing/
│  │  │  │  ├─ chat/
│  │  │  │  ├─ workflow/
│  │  │  │  ├─ rag/
│  │  │  │  └─ audit/
│  │  │  ├─ service/               # 用例编排（Chat/RunWorkflow/Billing）
│  │  │  ├─ repo/                  # Postgres repositories
│  │  │  ├─ pipeline/              # Pipeline + steps
│  │  │  │  ├─ pipeline.go
│  │  │  │  └─ steps/
│  │  │  │     ├─ auth.go
│  │  │  │     ├─ policy.go
│  │  │  │     ├─ budget_hold.go
│  │  │  │     ├─ context_rag.go
│  │  │  │     ├─ tool_skill.go
│  │  │  │     ├─ workflow.go
│  │  │  │     ├─ newapi_call.go
│  │  │  │     └─ usage_capture.go
│  │  │  ├─ integrations/
│  │  │  │  ├─ newapi/
│  │  │  │  ├─ storage/            # LocalFS/MinIO/S3-compatible
│  │  │  │  │  ├─ storage.go
│  │  │  │  │  ├─ localfs.go
│  │  │  │  │  └─ minio.go
│  │  │  │  └─ redis/
│  │  │  └─ observability/
│  │  │     ├─ logger.go
│  │  │     └─ trace.go
│  │  └─ test/
│  │
│  └─ worker/                      # Go（异步：文档解析/索引/usage聚合/workflow异步）
│     ├─ cmd/
│     │  └─ worker/
│     │     └─ main.go
│     ├─ internal/
│     │  ├─ config/
│     │  ├─ jobs/
│     │  │  ├─ ingest_document.go
│     │  │  ├─ build_index.go
│     │  │  ├─ usage_aggregate.go
│     │  │  └─ workflow_async.go
│     │  ├─ repo/
│     │  ├─ integrations/
│     │  │  ├─ storage/
│     │  │  └─ redis/
│     │  └─ observability/
│     └─ test/
│
├─ infra/
│  ├─ migrations/                  # SQL migrations（强烈建议）
│  │  ├─ 0001_init.up.sql
│  │  ├─ 0001_init.down.sql
│  │  ├─ 0002_projects.up.sql
│  │  ├─ 0002_projects.down.sql
│  │  └─ ...
│  ├─ scripts/
│  │  ├─ dev-up.sh
│  │  ├─ dev-down.sh
│  │  ├─ migrate-up.sh
│  │  └─ seed.sh
│  └─ docker/
│     ├─ postgres/
│     │  └─ init.sql
│     └─ minio/                    # 可选
│        └─ init.sh
│
├─ api/                            # 协议文档（推荐 OpenAPI）
│  ├─ openapi.yaml
│  └─ examples/
│     ├─ chat_stream.json
│     └─ workflow_run.json
│
├─ docs/
│  ├─ plan.md
│  ├─ architecture.md
│  ├─ billing.md
│  ├─ workflow_dsl.md
│  └─ rag.md
│
├─ go.mod
└─ pnpm-workspace.yaml
```