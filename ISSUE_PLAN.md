# DeepSpace Workflows 可执行 Issue 列表（优先级 + 依赖）

说明：
- 优先级：P0（阻塞主流程）/ P1（重要）/ P2（可延后）
- 依赖以 Issue ID 表示

## P0 基础与主链路

### DS-001 建立数据库迁移基线（P0）
- 目标：引入迁移工具并建立基础表结构
- 交付：`infra/migrations` 完整初始版本
- 依赖：无

### DS-002 完成配置体系（P0）
- 目标：Gateway/Worker 配置统一加载（DB/Redis/NewAPI/Storage）
- 交付：`services/gateway/internal/config` + 示例 env
- 依赖：无

### DS-003 Gateway 基础中间件（P0）
- 目标：Auth、trace_id、统一错误格式
- 交付：中间件与全局注入
- 依赖：DS-002

### DS-004 API Key 鉴权（P0）
- 目标：API Key 生成/校验/组织绑定
- 交付：API Key 管理与验证逻辑
- 依赖：DS-001, DS-003

### DS-005 SSE 透传与 chat completion 代理（P0）
- 目标：/v1/chat/completions SSE 透传
- 交付：NewAPI 代理 + trace_id 全链路
- 依赖：DS-003

### DS-006 计费闭环（P0）
- 目标：Wallet + Transaction + Usage 记录闭环
- 交付：hold/capture/release
- 依赖：DS-001, DS-003

### DS-007 Pipeline 最小框架（P0）
- 目标：构建 step 链式执行
- 交付：auth→policy→budget_hold→newapi_call→usage_capture
- 依赖：DS-003, DS-005, DS-006

### DS-008 Web Chat MVP（P0）
- 目标：基本聊天 UI + SSE 渲染
- 交付：消息列表、输入框、SSE hooks
- 依赖：DS-005

### DS-009 Web 计费视图（P0）
- 目标：展示钱包余额与 usage
- 交付：Wallet / Usage 页面
- 依赖：DS-006

---

## P1 运营能力与扩展

### DS-010 Audit Log 落库与查询（P1）
- 目标：请求日志可追溯
- 交付：audit_logs 表 + 查询接口
- 依赖：DS-001, DS-003

### DS-011 Admin 基础骨架（P1）
- 目标：管理端基础路由与布局
- 交付：Dashboard/Users/Models/Billing/Audit/Risk
- 依赖：无

### DS-012 Admin 用户与组织（P1）
- 目标：用户、组织、角色与权限管理
- 交付：Users/Orgs 页面
- 依赖：DS-011, DS-001

### DS-013 Admin 计费与对账（P1）
- 目标：Wallets/Transactions/Usage Records
- 交付：Billing 管理界面
- 依赖：DS-011, DS-006

### DS-014 Admin 模型与定价（P1）
- 目标：模型上架与定价规则
- 交付：Models/Pricing 页面
- 依赖：DS-011

### DS-015 Web 项目管理（P1）
- 目标：Project 列表与切换
- 交付：Projects 页面
- 依赖：DS-004

### DS-016 API 文档完善（P1）
- 目标：OpenAPI + 示例
- 交付：`api/openapi.yaml` + examples
- 依赖：DS-005, DS-006

### DS-017 Worker 基础框架（P1）
- 目标：异步任务调度
- 交付：job 基类 + 队列
- 依赖：DS-002

### DS-018 Usage 聚合任务（P1）
- 目标：批量聚合 usage
- 交付：定时任务与表结构
- 依赖：DS-017, DS-006

---

## P2 后续扩展

### DS-019 RAG 索引与检索（P2）
- 目标：文档解析、索引构建
- 交付：worker 任务 + gateway pipeline hook
- 依赖：DS-017

### DS-020 Workflow DSL 引擎（P2）
- 目标：Workflow JSON/YAML 执行
- 交付：workflow 定义 + run 管道
- 依赖：DS-007

### DS-021 Admin 风控策略（P2）
- 目标：限速/预算/IP 规则
- 交付：Risk & Policy 页面
- 依赖：DS-011

### DS-022 导出任务（P2）
- 目标：CSV/JSON/Markdown 导出
- 交付：Worker 任务
- 依赖：DS-017

