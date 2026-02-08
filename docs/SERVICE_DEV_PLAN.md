# DeepSpace Workflows 服务级开发计划

本计划按服务/模块拆解到具体任务列表，基于当前仓库结构与 `docs/system-architecture-plan.md` 定义的目标。

## 总体原则与里程碑
- 先落地可商用 Chat MVP（Gateway 鉴权 + 计费闭环 + SSE 透传），再扩展 Workflow/RAG。
- Gateway 是唯一 AI 出口；所有调用可计费、可审计、可追溯（trace_id）。
- Pipeline 可插拔，避免业务写死。

里程碑建议：
1. Gateway MVP + DB 迁移 + 基础日志/trace
2. Web Chat MVP + 计费视图
3. Admin MVP（运营/计费/审计）
4. Worker 异步任务 + RAG/Workflow 扩展

---

## 服务与模块任务清单

### 1) Gateway（`services/gateway`）
目标：统一鉴权、计费、审计、SSE 透传与 Pipeline 编排。

**A. 基础能力（P0）**
- 配置体系
  - 完善 `internal/config`：DB/Redis/NewAPI/Storage/日志/环境变量
- 中间件
  - Auth（API Key / JWT 基础）
  - trace_id 贯穿请求与日志
  - Rate Limit（基于 Redis 或本地实现）
  - Recovery 与统一错误格式
- 路由与健康检查
  - `/health`
  - `/v1/*` 透传（代理 NewAPI）

**B. 计费闭环（P0）**
- Wallet 余额与冻结模型
- Transaction 事务语义（hold / capture / release / refund）
- Usage 记录落库
- 幂等与对账字段（billing_ref_id）

**C. Pipeline（P1）**
- Pipeline 框架
  - step 接口与执行链
  - 统一上下文传递（用户/组织/费用/trace）
- 核心步骤实现
  - auth / policy / budget_hold / newapi_call / usage_capture
- 可插拔扩展接口
  - RAG / MCP / Skills / Workflow

**D. 审计与追踪（P1）**
- Audit Log 落库
- trace_id 查询接口

**E. API 规范（P1）**
- OpenAPI 文档补全与 examples

---

### 2) Worker（`services/worker`）
目标：解耦在线请求的异步能力。

**A. 基础框架（P1）**
- Job 管理
  - 统一任务接口
  - 任务调度（Redis/队列）
- 日志/trace 接入

**B. 任务实现（P1）**
- Usage 聚合
- 导出任务（CSV/JSON/Markdown）
- 文档解析与索引构建（为 RAG 预留）
- Workflow 异步执行（长时任务）

---

### 3) Web 用户端（`apps/web`）
目标：可商用 Chat MVP + 计费与项目视图。

**A. 基础功能（P0）**
- 登录/会话管理
- 项目选择与切换
- Chat UI（消息列表、输入框、流式渲染）
- SSE Hooks（`useChatStream.ts`）

**B. 计费视图（P0）**
- Wallet 余额与冻结
- Usage 列表

**C. 项目管理（P1）**
- Project 列表与详情
- API Key 管理入口

**D. 统一 API Client（P0）**
- 只允许调用 Gateway
- 统一错误与 trace_id 展示

---

### 4) Admin 管理端（`apps/admin`）
目标：运营必备能力与审计可追溯。

**A. 用户与组织（P1）**
- 用户列表、组织列表
- 角色与权限

**B. 模型与定价（P1）**
- 模型上架
- 价格配置

**C. 计费与对账（P1）**
- Wallets / Transactions / Usage Records

**D. 审计与追踪（P1）**
- trace_id 查询
- 请求日志检索

**E. 风控策略（P2）**
- 速率限制
- 预算上限
- IP 白名单/黑名单

---

### 5) 数据库与迁移（`infra/migrations`）
目标：为 Gateway/Worker 提供核心数据模型。

**P0 表结构**
- users / orgs / org_members
- api_keys
- wallets / transactions / usage_records
- conversations / messages
- audit_logs

**P1 扩展**
- workflows / workflow_runs
- rag_indexes / rag_chunks
- files / artifacts

---

### 6) 文档与规范（`docs/` + 根目录）
- 补充 API 文档与示例
- 更新系统实施计划
- 维护变更日志（建议新增 `CHANGELOG.md`）

---

## 建议的下一步执行顺序
1. 数据库迁移（P0 表结构）
2. Gateway 中间件（auth/trace/error）
3. Chat Completion SSE 透传
4. 计费闭环（hold/capture/release）
5. Web Chat MVP（只接 Gateway）

