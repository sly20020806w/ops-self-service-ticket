# 架构说明

## 业务状态机

```text
draft → pending → (gateway?) → executing → done → closed
                 ↘ rejected ─────────────────────↗
任意未完成状态可 cancel（创建人或管理员）
```

## 排他网关

流程节点 `type=gateway`，表达式示例：

- `priority==high` → `nextOnTrue`
- 否则 → `nextOnFalse`

表单字段来自工单 `formDataJson`。

## 动作按钮策略

后端根据 `status + role + assignee` 计算 `availableActions`，前端只渲染返回的按钮，避免前后端状态漂移。

## 模块边界

本仓库聚焦「自主工单」独立板块，可嵌入更大运维平台：

- 上游：用户权限、服务树
- 下游：任务执行 Agent、云资源开通、通知中心
