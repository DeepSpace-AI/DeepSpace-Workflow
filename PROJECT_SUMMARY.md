# DeepSpace Workflows 项目总结

## 项目定位
DeepSpace Workflows 是面向科研与研发场景的可商用 AI Workflow/Research OS。它以 Chat 为入口、Workflow/Pipeline 为核心，以计费、审计、可追溯为基础能力，强调科研可复现、商用计费闭环、能力可编排与组织级治理。

## 总体架构与链路
固定链路为：
User Web（Chat/Workflow UI） → Gateway（Auth/Billing/Policy/Pipeline/SSE） → NewAPI → 上游模型。管理端 Admin Web 通过 Admin API 进行运营、计费、审计与风控管理；Worker 负责异步任务（索引、聚合、导出等）。

架构原则：
- Gateway 是唯一 AI 出口，前端与管理端不得直连模型或 NewAPI。
- 所有调用必须可计费、可审计且可追溯（trace_id）。
- Pipeline 以可插拔步骤组织能力（RAG/MCP/Skills/Workflow）。

## 核心能力
- Chat：SSE 流式响应、会话与消息管理、按 Project 隔离。
- Workflow：多步骤编排（Plan → Execute → Verify），可串联 RAG/Skills/Tools，支持 JSON/YAML 定义。
- Pipeline：Auth → RateLimit → Budget Hold → Policy → Context Build(RAG) → Tool/Skill Dispatch → Workflow Orchestration → NewAPI Call → Usage Capture → Stream Response。
- 计费与审计：钱包余额/冻结、交易流水（hold/capture/release）、Usage 记录、Audit 日志与 trace_id 全链路追踪。

## 目录与模块概览
- `apps/web/`：用户端 Web（Nuxt 3 + Nuxt UI），包含 Chat、Workflow、Knowledge Base 等页面与组件。
- `apps/admin/`：管理端 Web（Nuxt 3 + Nuxt UI），包含用户/组织、模型定价、计费对账、审计与风控等页面。
- `services/gateway/`：Go + Gin 的 Gateway，负责鉴权、计费、SSE 透传与 Pipeline 编排。
- `services/worker/`：Go 异步任务（文档解析、索引构建、usage 聚合、workflow 异步）。
- `docs/system-architecture-plan.md`：系统分层、实现清单、Admin IA 与 Gateway 数据表设计细化。
- `package.json`：开发脚本入口。

## 关键数据与领域模型（设计方向）
Gateway 侧核心领域模型包含：
- 组织与用户：users、orgs、org_members、api_keys。
- 计费闭环：wallets、transactions、usage records。
- 会话与消息、审计与追踪（trace_id）。

## Admin 信息架构（设计方向）
管理端规划包含：Dashboard、Users/Orgs、Models/Pricing、Billing（Wallets/Transactions/Usage）、Audit（trace_id 与请求日志）、Risk & Policy（限速/预算/IP）。

## 技术栈
- Web：Nuxt / Nuxt UI
- Gateway & Worker：Go
- 数据库：PostgreSQL
- 缓存/队列：Redis
- 对象存储：LocalFS / MinIO / S3
- AI 接入：NewAPI

## 本地开发脚本
- `pnpm web:dev`：启动用户端 Web
- `pnpm admin:dev`：启动管理端 Web
- `pnpm gateway:dev`：启动 Gateway（使用 air）

## 现状与下一步（来自规划文档）
- Web 侧重点是 Chat MVP 与计费/项目视图，SSE 逻辑集中在 hooks。
- Admin 侧重点是运营能力、计费对账与审计追溯。
- Gateway 需完善鉴权、计费闭环、审计追踪、Pipeline 可插拔步骤。
- Worker 承担 usage 聚合、导出与 RAG 索引构建等异步任务。

