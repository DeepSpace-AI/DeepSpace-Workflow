# DeepSpace Workflows - Copilot 指引

## 架构与边界（必须遵守）
- 固定链路：User Web -> Gateway -> NewAPI；Admin Web -> Admin API；Worker 处理异步任务（索引、聚合等）。参见 [README.md](README.md)。
- Gateway 是唯一 AI 出口，前端与管理端禁止直连模型或 NewAPI。参见 [README.md](README.md)。
- Pipeline 以可插拔步骤组织能力，结构参考 [services/gateway/internal/pipeline](services/gateway/internal/pipeline)。

## 关键目录与职责
- 用户端：Nuxt 应用在 [apps/web](apps/web)（Chat/Workflow/Knowledge Base）。
- 管理端：Nuxt 应用在 [apps/admin](apps/admin)（计费、审计、风控、定价）。
- 服务端：Gateway 在 [services/gateway](services/gateway)（Auth/Billing/SSE/Pipeline）。
- 异步任务：Worker 在 [services/worker](services/worker)。
- 架构细化与数据设计：见 [docs/system-architecture-plan.md](docs/system-architecture-plan.md) 与 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)。

## 约定与特殊规则（项目特有）
- 项目首选中文：展示文案与代码注释必须是中文；变量命名必须使用英文。
- 侧边栏用户信息按钮必须有弹出菜单，包含“个人信息”和“退出登录”。
- Nuxt 项目中不要重复引用 Vue 相关库（如 vue, vue-router, vuex）。
- SSE 交互与 Chat 流式处理必须走 Gateway，不在前端直连模型。

## 开发与运行
- 用户端开发：`pnpm web:dev`（参见 [package.json](package.json)）。
- 管理端开发：`pnpm admin:dev`（参见 [package.json](package.json)）。
- Gateway 开发：`pnpm gateway:dev`（使用 air，参见 [package.json](package.json)）。

## 数据库与迁移
- 使用 GORM AutoMigrate，环境变量与迁移开关见 [infra/README.md](infra/README.md)。
- `DB_RESET_ON_START=true` 会在 Gateway 启动时清空并重建表。

## 设计与实现提示（仅限已存在模式）
- 计费/审计/追踪（trace_id）贯穿调用链路，设计参考 [README.md](README.md) 与 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)。
- Pipeline 步骤组织能力（RAG/Skills/Workflow/Usage Capture）参考 [README.md](README.md) 与 [services/gateway/internal/pipeline](services/gateway/internal/pipeline)。
