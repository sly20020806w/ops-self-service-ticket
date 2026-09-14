# 自主工单平台 · Ops Self-Service Ticket

> 面向生产演示与落地复用的 **运维自主工单** 板块。对齐「大运维平台」课程中 **工单模块** 核心能力，并在此基础上做了可向老板展示的增强。

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![Vue3](https://img.shields.io/badge/Vue-3-42b883?logo=vue.js)
![License](https://img.shields.io/badge/license-MIT-green)

## 为什么做这个项目

运维团队最常见的诉求是：**业务自己提单 → 审批 → 自动/半自动执行 → 全程可追溯**。  
本仓库把「自主工单」做成可独立运行的完整工程，方便：

- 直接本地/容器演示给领导看价值
- 作为企业内部工单底座二次开发
- 作为简历/作品集的 Go + Vue 全栈运维开发案例

## 课程对齐能力（有）

| 能力 | 说明 |
|------|------|
| 表单设计器产物 | 动态字段 Schema（input/textarea/select/number/switch）驱动创建页 |
| 流程管理 | 开始 / 审批 / 执行 / 结束节点可配置 |
| 排他网关 | 按表单字段表达式分流（如 `priority==high`） |
| 状态动作按钮 | 按工单状态 + 角色动态返回可操作 action |
| 服务树叶子绑定 | 创建工单时强制绑定叶子节点 |
| 审批/执行链路 | 提交 → 审批通过/驳回 → 执行 → 关闭/撤销 |
| IM 通知 | 控制台 + Webhook（企微/飞书/钉钉兼容文本消息）通知日志落库 |

## 加分增强（课程外优化）

| 增强 | 说明 |
|------|------|
| SLA 截止与超时标记 | 流程级 SLA，后台定时扫描超时 |
| 工单模板 | 一键预填标题/优先级/表单，加快自主申请 |
| 老板看板 | 总量/进行中/今日/超时 + 状态与优先级分布 |
| 全量审计日志 | 谁在何时对什么资源做了什么 |
| 多角色演示账号 | admin / 申请人 / 审批人 / 执行人一键切换 |
| SQLite 零依赖启动 | 无需先装 MySQL，30 秒出 Demo |
| Docker Compose | 一键前后端演示环境 |

## 快速开始（本地）

### 1) 启动后端

```bash
cd backend
go mod tidy
go run ./cmd/server
```

后端默认：`http://127.0.0.1:8080`

### 2) 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端默认：`http://127.0.0.1:5173`

### 演示账号（密码均为 `Passw0rd!`）

| 账号 | 角色 | 用途 |
|------|------|------|
| `zhangsan` | applicant | 申请人：建单、提交 |
| `lisi` | approver | 审批人：通过/驳回 |
| `wangwu` | executor | 执行人：执行完成 |
| `admin` | admin | 管理员：看全局看板/审计 |

### 推荐演示路径（给老板 3 分钟）

1. 用 `zhangsan` 登录 → 看板 → **新建工单**
2. 选「云主机自主申请流程」，绑定服务树叶子，`priority=high`，提交
3. 切换 `lisi` 审批通过（触发排他网关走「加急执行」）
4. 切换 `wangwu` 执行完成
5. 回到看板看状态分布，进「审计」看通知与操作轨迹

## Docker 一键演示

```bash
docker compose up -d --build
```

- 前端：http://127.0.0.1:5173  
- 后端：http://127.0.0.1:8080  

## 目录结构

```text
ops-self-service-ticket/
├── backend/                 # Go + Gin + Gorm
│   ├── cmd/server/          # 入口
│   ├── configs/             # 配置
│   └── internal/            # api/auth/model/service/notify/seed
├── frontend/                # Vue3 + Vite + Element Plus
├── docker-compose.yml
└── docs/ARCHITECTURE.md
```

## 核心 API（节选）

- `POST /api/auth/login` 登录
- `GET  /api/dashboard` 看板指标
- `GET  /api/tickets` 工单列表
- `POST /api/tickets` 创建（可直接 submit）
- `GET  /api/tickets/:id` 详情 + `availableActions`
- `POST /api/tickets/:id/{submit|approve|reject|execute|close|cancel|comment}`
- `GET  /api/workflows` / `GET /api/forms` / `GET /api/tree/nodes`
- `GET  /api/audit` / `GET /api/notify-logs`

## 配置

`backend/configs/config.yaml`：

- `notify.webhook_url`：填入企业微信/飞书/钉钉机器人地址即可推送
- `sla.default_hours`：默认 SLA
- `auth.jwt_secret`：生产务必修改

## 生产落地建议

1. 将 SQLite 切换为 MySQL/PostgreSQL（Gorm Dialector 可替换）
2. 对接真实服务树与 CMDB、SSO/LDAP
3. 执行节点对接任务中心 / Ansible / 云 API
4. 通知通道按公司 IM 卡片协议定制
5. 补齐细粒度 Casbin 权限（可与平台权限模块共用）

## License

MIT
