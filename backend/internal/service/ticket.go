package service

import (
        "encoding/json"
        "errors"
        "fmt"
        "time"

        "github.com/sly20020806w/ops-self-service-ticket/internal/model"
        "github.com/sly20020806w/ops-self-service-ticket/internal/notify"
        "gorm.io/gorm"
)

type TicketService struct {
        DB       *gorm.DB
        Notifier *notify.Notifier
}

type flowNode struct {
        Key  string `json:"key"`
        Name string `json:"name"`
        Type string `json:"type"`
        Role string `json:"role"`
}

type flowEdge struct {
        From      string `json:"from"`
        To        string `json:"to"`
        On        string `json:"on"`
        Condition string `json:"condition"`
}

var allowedActions = map[string][]string{
        model.StatusDraft:     {"submit", "close"},
        model.StatusPending:   {"approve", "reject", "comment"},
        model.StatusApproved:  {"execute", "comment"},
        model.StatusExecuting: {"execute", "comment"},
        model.StatusDone:      {"close", "comment"},
        model.StatusRejected:  {"close", "comment"},
        model.StatusClosed:    {"comment"},
}

func (s *TicketService) nextTicketNo() string {
        var count int64
        today := time.Now().Format("20060102")
        s.DB.Model(&model.Ticket{}).Where("ticket_no LIKE ?", "WO-"+today+"-%").Count(&count)
        return fmt.Sprintf("WO-%s-%04d", today, count+1)
}

func (s *TicketService) Create(creatorID uint, title string, flowID uint, treeNodeID *uint, priority string, formData map[string]any) (*model.Ticket, error) {
        var flow model.FlowDefinition
        if err := s.DB.First(&flow, flowID).Error; err != nil {
                return nil, errors.New("流程不存在")
        }
        raw, _ := json.Marshal(formData)
        deadline := time.Now().Add(slaDuration(priority))
        t := &model.Ticket{
                TicketNo:     s.nextTicketNo(),
                Title:        title,
                FlowID:       flow.ID,
                FormID:       flow.FormID,
                TreeNodeID:   treeNodeID,
                CreatorID:    creatorID,
                Status:       model.StatusDraft,
                CurrentNode:  "start",
                Priority:     priority,
                FormDataJSON: string(raw),
                SLADeadline:  &deadline,
        }
        if err := s.DB.Create(t).Error; err != nil {
                return nil, err
        }
        s.audit(creatorID, "ticket.create", t.ID, map[string]any{"title": title})
        return t, nil
}

func slaDuration(priority string) time.Duration {
        switch priority {
        case "P0":
                return 1 * time.Hour
        case "P1":
                return 4 * time.Hour
        case "P2":
                return 24 * time.Hour
        default:
                return 72 * time.Hour
        }
}

func (s *TicketService) Action(ticketID, actorID uint, action, comment string) (*model.Ticket, error) {
        var t model.Ticket
        if err := s.DB.First(&t, ticketID).Error; err != nil {
                return nil, errors.New("工单不存在")
        }
        if !contains(allowedActions[t.Status], action) {
                return nil, fmt.Errorf("当前状态 %s 不允许动作 %s", t.Status, action)
        }

        from := t.CurrentNode
        to := from
        newStatus := t.Status

        switch action {
        case "submit":
                to = "mgr_approve"
                newStatus = model.StatusPending
                s.Notifier.Send(t.ID, "im", "approver", fmt.Sprintf("【待审批】%s %s", t.TicketNo, t.Title))
        case "approve":
                to = s.routeExclusive(&t, "approve")
                newStatus = model.StatusApproved
                if to == "exec" {
                        newStatus = model.StatusApproved
                }
                s.Notifier.Send(t.ID, "im", "executor", fmt.Sprintf("【已审批】%s 请执行", t.TicketNo))
                s.Notifier.Send(t.ID, "im", "creator", fmt.Sprintf("【审批通过】%s", t.TicketNo))
        case "reject":
                to = "end"
                newStatus = model.StatusRejected
                s.Notifier.Send(t.ID, "im", "creator", fmt.Sprintf("【已驳回】%s %s", t.TicketNo, comment))
        case "execute":
                to = "end"
                newStatus = model.StatusDone
                now := time.Now()
                t.ClosedAt = &now
                s.Notifier.Send(t.ID, "im", "creator", fmt.Sprintf("【已完成】%s", t.TicketNo))
        case "close":
                to = "end"
                newStatus = model.StatusClosed
                now := time.Now()
                t.ClosedAt = &now
        case "comment":
                // 仅记录
        }

        if t.SLADeadline != nil && time.Now().After(*t.SLADeadline) && newStatus != model.StatusDone && newStatus != model.StatusClosed {
                t.SLABreached = true
        }

        t.CurrentNode = to
        t.Status = newStatus
        if err := s.DB.Save(&t).Error; err != nil {
                return nil, err
        }
        act := model.TicketAction{
                TicketID: t.ID,
                ActorID:  actorID,
                Action:   action,
                FromNode: from,
                ToNode:   to,
                Comment:  comment,
        }
        _ = s.DB.Create(&act).Error
        s.audit(actorID, "ticket."+action, t.ID, map[string]any{"from": from, "to": to, "comment": comment})
        if s.Notifier.Webhook != "" {
                s.Notifier.Send(t.ID, "webhook", "default", fmt.Sprintf("%s %s by user#%d", t.TicketNo, action, actorID))
        }
        return &t, nil
}

func (s *TicketService) routeExclusive(t *model.Ticket, on string) string {
        var flow model.FlowDefinition
        if err := s.DB.First(&flow, t.FlowID).Error; err != nil {
                return "exec"
        }
        var edges []flowEdge
        _ = json.Unmarshal([]byte(flow.EdgesJSON), &edges)
        formData := map[string]any{}
        _ = json.Unmarshal([]byte(t.FormDataJSON), &formData)

        // 课程「排他网关」：按表单 urgent 字段分流（增强可读条件）
        urgent, _ := formData["urgent"].(bool)
        for _, e := range edges {
                if e.From == "gateway" || (e.From == t.CurrentNode && e.On == on) {
                        if e.Condition == "urgent==true" && urgent {
                                return e.To
                        }
                        if e.Condition == "default" && !urgent {
                                return e.To
                        }
                        if e.On == on && e.Condition == "" {
                                return e.To
                        }
                }
        }
        // approve 后进入 gateway 再落到 exec
        for _, e := range edges {
                if e.From == t.CurrentNode && e.On == on {
                        next := e.To
                        if next == "gateway" {
                                for _, g := range edges {
                                        if g.From == "gateway" {
                                                if g.Condition == "urgent==true" && urgent {
                                                        return g.To
                                                }
                                                if g.Condition == "default" {
                                                        return g.To
                                                }
                                        }
                                }
                        }
                        return next
                }
        }
        return "exec"
}

func (s *TicketService) List(status string) ([]model.Ticket, error) {
        var list []model.Ticket
        q := s.DB.Preload("Creator").Preload("Assignee").Preload("TreeNode").Order("id desc")
        if status != "" {
                q = q.Where("status = ?", status)
        }
        // 刷新 SLA 超时标记
        _ = s.DB.Model(&model.Ticket{}).
                Where("sla_deadline < ? AND status NOT IN ? AND sla_breached = ?", time.Now(), []string{model.StatusDone, model.StatusClosed}, false).
                Update("sla_breached", true)
        err := q.Find(&list).Error
        return list, err
}

func (s *TicketService) Get(id uint) (*model.Ticket, []model.TicketAction, error) {
        var t model.Ticket
        if err := s.DB.Preload("Creator").Preload("Assignee").Preload("TreeNode").First(&t, id).Error; err != nil {
                return nil, nil, err
        }
        var acts []model.TicketAction
        _ = s.DB.Preload("Actor").Where("ticket_id = ?", id).Order("id asc").Find(&acts).Error
        return &t, acts, nil
}

func (s *TicketService) Dashboard() map[string]any {
        type kv struct {
                Status string
                Cnt    int64
        }
        var rows []kv
        s.DB.Model(&model.Ticket{}).Select("status, count(*) as cnt").Group("status").Scan(&rows)
        var breached int64
        s.DB.Model(&model.Ticket{}).Where("sla_breached = ?", true).Count(&breached)
        var total int64
        s.DB.Model(&model.Ticket{}).Count(&total)
        return map[string]any{
                "total":        total,
                "slaBreached":  breached,
                "byStatus":     rows,
                "generatedAt":  time.Now().Format(time.RFC3339),
        }
}

func (s *TicketService) audit(actorID uint, action string, ticketID uint, detail map[string]any) {
        raw, _ := json.Marshal(detail)
        _ = s.DB.Create(&model.AuditLog{
                ActorID:  actorID,
                Action:   action,
                Resource: fmt.Sprintf("ticket:%d", ticketID),
                Detail: string(raw),
        }).Error
}

func contains(ss []string, x string) bool {
        for _, s := range ss {
                if s == x {
                        return true
                }
        }
        return false
}

func AvailableActions(status string) []string {
        return append([]string{}, allowedActions[status]...)
}
