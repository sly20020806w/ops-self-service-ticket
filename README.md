# 自主工单平台 · Ops Self-Service Ticket

> 可直接演示、可二次开发的 **运维自主工单** 完整工程。  
> 对齐「大运维平台」课程 **工单模块** 核心能力，并补齐老板演示与生产落地增强项。

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](#)
[![Vue3](https://img.shields.io/badge/Vue-3-42b883?logo=vue.js)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 一句话价值

业务同学自助提单 → 审批人审批（可走排他网关）→ 执行人落地 → 全链路审计与 IM 通知；带 SLA 与看板，适合汇报演示。

## 课程对齐（有）

| 能力 | 实现位置 |
|------|----------|
| 动态表单设计器 | `FormTemplate.fields` JSON 驱动前端渲染 |
| 流程管理 | `FlowDefinition.nodes/edges` |
| 排他网关 | 审批后按 `urgent==true` / `default` 分流 |
| 状态动作按钮 | `AvailableActions(status)` 驱动前端按钮 |
| 服务树叶子绑定 | 创建工单可选 `treeNodeId` |
| 审批/执行/驳回/关闭 | `POST /api/tickets/:id/actions` |
| IM 通知整合 | 通知记录 + Webhook 推送 |

## 加分增强（课程外）

- SLA：按优先级 P0/P1/P2/P3 截止时间，列表自动标记超时
- 老板看板：总量 / 状态分布 / SLA 风险
- 审计日志：谁对哪张工单做了什么
- Prometheus `/metrics` 便于监控接入
- SQLite 零依赖 + Docker Compose 一键演示

## 快速开始

### 后端

```bash
cd backend
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
go run ./cmd/server
```

默认监听：`http://127.0.0.1:8080`

### 前端

```bash
cd frontend
npm install
npm run dev
```

打开：`http://127.0.0.1:5173`

### 演示账号（密码均为 `admin123`）

| 账号 | 角色 | 建议操作 |
|------|------|----------|
| alice | requester | 新建草稿并 submit |
| bob | approver | approve / reject |
| carol | executor | execute |
| admin | admin | 看看板与审计 |

### 3 分钟老板演示脚本

1. `alice` 登录 → 工单 → 新建（勾选紧急开关）→ 进入详情 → `submit`
2. 退出，换 `bob` → 打开该单 → `approve`（走排他网关到执行）
3. 换 `carol` → `execute`
4. 看板看状态分布；审计页看通知与操作轨迹

## Docker

```bash
docker compose up -d --build
```

- Web: http://127.0.0.1:5173  
- API: http://127.0.0.1:8080  

可选环境变量：

- `NOTIFY_WEBHOOK`：企业微信/钉钉/飞书机器人地址
- `JWT_SECRET`：生产务必修改
- `DB_PATH`：数据库文件路径

## 目录

```text
backend/   Go + Gin + Gorm + SQLite
frontend/  Vue3 + Vite + Element Plus
docs/      架构与课程对照
```

## License

MIT
