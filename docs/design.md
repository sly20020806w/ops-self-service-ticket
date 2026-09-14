# 设计说明：自主工单状态机

## 状态

draft → pending_approval → approved → done
                 ↘ rejected → closed
draft/done/rejected → close → closed

## 动作权限（简化版，可按角色再收紧）

| 状态 | 可用动作 |
|------|----------|
| draft | submit, close |
| pending_approval | approve, reject, comment |
| approved | execute, comment |
| done / rejected | close, comment |
| closed | comment |

## 排他网关

审批通过后进入 `gateway`：

- `urgent==true` → 仍进入 `exec`（演示可扩展为加急队列）
- `default` → `exec`

生产可把 condition 做成表达式引擎（cel / expr）。
