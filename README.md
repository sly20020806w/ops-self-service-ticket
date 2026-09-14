# OpsTicket · 自主工单平台

> Go + Vue3 工单子系统：动态表单、流程审批、排他网关、服务树绑定、状态动作、IM/Webhook 通知，附 SLA、审计与运营看板。

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](#)
[![Vue3](https://img.shields.io/badge/Vue-3-42b883?logo=vue.js)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 能力一览

| 能力 | 说明 |
|------|------|
| 动态表单 | Schema 驱动创建页字段渲染 |
| 流程引擎 | 开始 / 审批 / 执行 / 结束节点 |
| 排他网关 | 按表单条件分流 |
| 状态动作 | 按状态返回可操作按钮 |
| 服务树绑定 | 创建工单绑定叶子节点 |
| 审批执行闭环 | 提交 → 审批 → 执行 → 关闭 |
| IM 通知 | 控制台 + Webhook，通知落库 |
| SLA | 按优先级截止与超时标记 |
| 运营看板 | 总量 / 状态分布 / SLA 风险 |
| 审计日志 | 操作全链路留痕 |

## 快速开始

### 后端

```bash
cd backend
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
go run ./cmd/server
```

默认：`http://127.0.0.1:8080`

### 前端

```bash
cd frontend
npm install
npm run dev
```

打开：`http://127.0.0.1:5173`

### 账号（密码 `admin123`）

| 账号 | 角色 |
|------|------|
| alice | 申请人 |
| bob | 审批人 |
| carol | 执行人 |
| admin | 管理员 |

## Docker

```bash
docker compose up -d --build
```

## License

MIT
