![DeepSpace Logo](./docs/deepspace-logo.png)
# DeepSpace Workflows

> **DeepSpace Workflows** 是一个面向科研与研发场景的 **AI Workflow 平台**，以 Chat 为入口，以 Workflow / Pipeline 为核心，以 **写作、协作、技能** 为基础能力，构建可长期演进的 AI Research Infrastructure。

---

## 1. 项目定位

DeepSpace Workflows 并不是一个简单的 ChatGPT 替代品，而是一个 **AI 工作流与科研协作平台（AI Research OS）**，核心目标包括：

* **多步骤 AI 任务编排**：支持复杂的多步骤任务（如文献综述、实验设计、数据分析等），而非单一问答。
* **可插拔能力链路**：通过 Pipeline 机制，灵活串联
  RAG、技能（Skills）、工具（Tools）等能力，满足多样化需求。
* **项目与知识库管理**：支持多项目隔离与知识库构建

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
Admin Web (Nuxt + Nuxt UI)
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

## 7. Worker 邮件队列

Worker 用于消费邮件发送队列（Redis List），通过 SMTP 发送邮件，并支持失败重试与死信队列。

运行方式：

```
cd services/worker
go run ./cmd/worker
```

关键行为：

* 使用 Redis `BRPOP` 阻塞拉取队列
* 发送失败重试（默认 5 次）
* 超过重试次数写入死信队列（默认 `email:dead`）

主要配置项见 `.env.example` 的 Worker/邮件部分。

## 8. Docker 运行

使用 Docker Compose 启动（Web/Admin 对外暴露，Gateway 仅内网访问）：

```
docker compose --env-file .env.docker up -d --build
```

首次运行初始化管理员用户（确保数据库已启动）：

```
docker compose --env-file .env.docker run --rm \
  -e ADMIN_EMAIL=admin@example.com \
  -e ADMIN_PASSWORD=change-me \
  -e ADMIN_DISPLAY_NAME=管理员 \
  gateway /app/admin-init
```

默认端口：

* Web: http://localhost:8080
* Admin: http://localhost:8081

关键环境变量说明（见 `.env.docker`）：

* `NEWAPI_BASE_URL`：NewAPI 服务地址（生产环境需要修改为真实地址）
* `JWT_SECRET`：JWT 密钥（生产环境必须替换）
* 邮件相关变量：`EMAIL_FROM_ADDRESS`、`SMTP_HOST`、`SMTP_USER`、`SMTP_PASSWORD` 必须填写真实值

停止与清理：

```
docker compose down
```


## 9. Roadmap（简化）

1. Chat + Gateway + Billing MVP
2. Project + RAG（科研最小闭环）
3. Workflow Engine（JSON 定义）
4. MCP / Skills 扩展


## 10.目录结构

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
│  ├─ web/                         # Nuxt 4 + Nuxt UI（用户端：Chat / Workflow / KB）
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
│  └─ admin/                       # Nuxt 4 + Nuxt UI（管理端：定价/对账/审计/风控）
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