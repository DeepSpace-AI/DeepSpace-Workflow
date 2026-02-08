# DeepSpace Workflows 管理端 (Admin) 开发计划

## 1. 项目现状评估

### 1.1 代码与架构现状
- **技术栈**: 项目采用了最前沿的技术栈，包括 **Nuxt 4**、**Nuxt UI v4** (基于 Vue 3.5+) 和 **Tailwind CSS 4**。
- **目录结构**: 采用了 Nuxt 4 的 `app/` 目录模式 (`apps/admin/app/`)，符合现代化 Nuxt 项目结构。
- **当前状态**: 项目处于 **初始化阶段** (Greenfield)。
    - **配置文件**: `package.json` 和 `nuxt.config.ts` 已配置基础依赖和模块。
    - **入口文件**: `app.vue` 已包含基础的 `<UApp>` 和路由出口。
    - **缺失模块**: 核心的 `pages/` 目录缺失，导致目前没有任何可访问页面。
    - **低级错误**: `layouts/` 下存在拼写错误 `dashboar.vue` (应为 `dashboard.vue`) 且文件为空。
    - **功能完整性**: 0%。仅有空壳，没有任何业务功能实现。

### 1.2 技术债务与风险
- **Nuxt 4 & UI v4 稳定性**: 使用了非常新的版本，可能会遇到文档不全或 Breaking Changes 的风险。
- **空文件**: 存在无效的布局文件，需要立即清理。

## 2. 需求优先级排序

根据 `ISSUE_PLAN.md`，管理端的需求优先级划分如下：

| 优先级 | 需求模块 | 对应 Issue | 业务价值 |
| :--- | :--- | :--- | :--- |
| **P0 (Critical)** | **Admin 基础骨架** | **DS-011** | 搭建管理端运行环境，确立导航与布局结构，是所有功能的容器。 |
| **P1 (High)** | **用户与组织管理** | **DS-012** | 基础数据管理，支撑业务运行。 |
| **P1 (High)** | **模型与定价管理** | **DS-014** | 核心商业化配置能力，决定平台售卖内容。 |
| **P1 (High)** | **计费与对账** | **DS-013** | 财务核心，展示钱包、流水与用量，支撑营收核算。 |
| **P1 (High)** | **审计日志** | **DS-010** | 满足合规性与可追溯性要求。 |
| **P2 (Medium)** | **风控策略** | **DS-021** | 增强平台安全性与稳定性 (限流/黑白名单)。 |

## 3. 技术实施方案

### 3.1 基础架构修正与搭建 (DS-011)
- **修正布局**: 重命名并实现 `apps/admin/app/layouts/dashboard.vue`，构建侧边栏导航 (Sidebar) 和顶部栏 (Header)。
- **路由结构**: 创建 `apps/admin/app/pages/` 目录，并按照以下结构初始化页面：
    - `index.vue` (Dashboard 概览)
    - `users.vue` & `orgs.vue` (用户/组织)
    - `models.vue` & `pricing.vue` (模型/定价)
    - `billing/` (子目录: `wallets.vue`, `transactions.vue`, `usage.vue`)
    - `audit.vue` (审计)
    - `policy.vue` (风控)
- **状态管理**: 初始化 Pinia Store，建立 `useAdminAuth` (模拟/对接 Gateway) 和 `useGlobalState`。

### 3.2 核心模块实现
- **UI 组件规范**: 基于 `@nuxt/ui` v4，封装统一的 `DataTable` (带分页/筛选)、`PageHeader`、`StatusBadge` 等业务组件。
- **数据交互**: 封装 `useAdminFetch` composable，统一处理 API 请求、错误提示 (Toast) 和 Token 注入。

## 4. 开发里程碑

| 阶段 | 目标 | 交付物 | 预计周期 |
| :--- | :--- | :--- | :--- |
| **M1: 骨架构建** | 完成项目结构修复，实现 Dashboard 布局与所有菜单路由跳转。 | 修复的 `dashboard.vue`，完整的 `pages` 目录结构，可导航的空页面。 | 1 天 |
| **M2: 基础数据** | 完成用户、组织、模型、定价页面的 CRUD 界面 (Mock 数据)。 | Users/Orgs/Models/Pricing 列表与表单页。 | 2-3 天 |
| **M3: 计费与审计** | 完成计费相关报表视图与审计日志查询视图。 | Billing (3个子页面) 与 Audit 页面。 | 2-3 天 |
| **M4: 对接联调** | 对接 Gateway API (需配合后端进度)，实现真实数据交互。 | 移除 Mock，全功能联调通过。 | 待定 (依赖后端) |

## 5. 测试策略

- **单元测试 (Unit Test)**: 针对核心工具函数 (`utils/`) 和复杂 Composables (`useBilling`) 编写 Vitest 测试用例。
- **组件测试**: 确保通用业务组件 (如 `PricingCard`, `UserTable`) 在不同 Props 下渲染正确。
- **手动验证**: 每个 Milestone 交付前，进行全链路点击测试，确保路由跳转无死链，UI 响应式正常。

## 6. 部署计划

- **构建**: 使用 `pnpm build` 生成 `.output` 产物。
- **环境**: 依赖 `ENV` 环境变量区分 `development` (Mock/Local API) 和 `production` (Live API)。
- **验证**: 部署后访问 `/admin` 根路径，检查静态资源加载与 API 连通性。

## 7. 风险评估与应对措施

- **风险**: Nuxt UI v4 变动频繁。
    - **应对措施**: 锁定 `package.json` 版本，遇到 Bug 及时查阅 GitHub Issues 或回退到稳定小版本。
- **风险**: 后端 API 接口未定。
    - **应对措施**: 前端优先定义 TypeScript Interface (DTO)，使用 Mock 数据开发，确保 UI 逻辑独立于后端进度。

## 8. M4 对接联调清单

### 8.1 接口清单（Admin 端）

- **已具备可对接**
    - `GET /api/admin/users` 用户列表（分页/搜索）
    - `POST /api/admin/users` 用户创建
    - `GET /api/admin/users/:id` 用户详情
    - `PATCH /api/admin/users/:id` 用户更新
    - `DELETE /api/admin/users/:id` 用户删除
    - `GET /api/admin/orgs` 组织列表（分页/搜索）
    - `POST /api/admin/orgs` 组织创建
    - `GET /api/admin/orgs/:id` 组织详情
    - `PATCH /api/admin/orgs/:id` 组织更新
    - `DELETE /api/admin/orgs/:id` 组织删除

- **已具备可对接（组织维度）**
    - `GET /api/billing/wallet` 钱包摘要
    - `GET /api/billing/usage` 用量列表（分页/时间区间）
    - `POST /api/billing/hold` 预算冻结
    - `POST /api/billing/capture` 预算扣减
    - `POST /api/billing/release` 预算释放

- **需补齐（Gateway 侧新增）**
    - `GET /api/admin/models` 模型列表
    - `POST /api/admin/models` 模型新增
    - `PATCH /api/admin/models/:id` 模型更新
    - `DELETE /api/admin/models/:id` 模型下架
    - `GET /api/admin/pricing` 定价规则列表
    - `POST /api/admin/pricing` 定价规则新增
    - `PATCH /api/admin/pricing/:id` 定价规则更新
    - `DELETE /api/admin/pricing/:id` 定价规则删除
    - `GET /api/admin/wallets` 组织钱包列表（跨组织）
    - `GET /api/admin/transactions` 交易流水列表（跨组织）
    - `GET /api/admin/usage` 用量列表（跨组织）
    - `GET /api/admin/audit` 审计日志列表
    - `GET /api/admin/policies` 风控策略列表
    - `POST /api/admin/policies` 风控策略新增
    - `PATCH /api/admin/policies/:id` 风控策略更新
    - `DELETE /api/admin/policies/:id` 风控策略删除

### 8.2 鉴权与权限约定

- **认证**: Cookie 会话为主，Authorization Bearer 为兜底。
- **权限**: `/api/admin/*` 需要管理员角色。
- **错误码**: 401 未登录，403 无权限，500 服务端异常。

### 8.3 分页与字段约定

- **分页**: `page`、`page_size` 输入；响应统一返回 `items`、`total`、`page`、`page_size`。
- **时间**: 使用 RFC3339 格式。
- **枚举**: 统一以英文枚举值传输，前端映射中文展示。

### 8.4 前端对接清单（页面映射）

- **用户管理**: `users.vue` ↔ `/api/admin/users`
- **组织管理**: `orgs.vue` ↔ `/api/admin/orgs`
- **模型管理**: `models.vue` ↔ `/api/admin/models`
- **定价管理**: `pricing.vue` ↔ `/api/admin/pricing`
- **钱包摘要**: `billing/wallets.vue` ↔ `/api/billing/wallet`
- **交易流水**: `billing/transactions.vue` ↔ `/api/admin/transactions`
- **用量明细**: `billing/usage.vue` ↔ `/api/admin/usage`（联调前期可先接 `/api/billing/usage`）
- **审计日志**: `audit.vue` ↔ `/api/admin/audit`
- **风控策略**: `policy.vue` ↔ `/api/admin/policies`

### 8.5 联调步骤

- 先对接已具备接口（用户/组织/钱包/用量）。
- 补齐 Gateway 缺口接口并完成接口契约对齐。
- 全量替换 Mock 数据，并统一错误处理与空态展示。
- 逐页验证筛选、分页、操作按钮逻辑与权限拦截。

### 8.6 验收清单

- **鉴权**: 未登录跳转/提示正确，管理员权限校验生效。
- **分页**: 每页条数、总数、翻页逻辑正确。
- **筛选**: 关键词与下拉筛选准确联动。
- **错误**: 401/403/500 有明确提示且不影响页面稳定。
