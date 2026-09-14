package store

import (
        "encoding/json"
        "time"

        "github.com/sly20020806w/ops-self-service-ticket/internal/model"
        "golang.org/x/crypto/bcrypt"
        "gorm.io/gorm"
)

func SeedDemo(db *gorm.DB) error {
        var n int64
        db.Model(&model.User{}).Count(&n)
        if n > 0 {
                return nil
        }

        hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
        users := []model.User{
                {Username: "admin", Display: "平台管理员", Role: "admin", Password: string(hash)},
                {Username: "alice", Display: "Alice-申请人", Role: "requester", Password: string(hash)},
                {Username: "bob", Display: "Bob-审批人", Role: "approver", Password: string(hash)},
                {Username: "carol", Display: "Carol-执行人", Role: "executor", Password: string(hash)},
        }
        if err := db.Create(&users).Error; err != nil {
                return err
        }

        root := model.ServiceTreeNode{Name: "公司", Path: "/公司", Leaf: false}
        if err := db.Create(&root).Error; err != nil {
                return err
        }
        biz := model.ServiceTreeNode{ParentID: &root.ID, Name: "电商业务", Path: "/公司/电商业务", Leaf: false}
        if err := db.Create(&biz).Error; err != nil {
                return err
        }
        leaves := []model.ServiceTreeNode{
                {ParentID: &biz.ID, Name: "订单服务", Path: "/公司/电商业务/订单服务", Leaf: true},
                {ParentID: &biz.ID, Name: "支付服务", Path: "/公司/电商业务/支付服务", Leaf: true},
                {ParentID: &biz.ID, Name: "库存服务", Path: "/公司/电商业务/库存服务", Leaf: true},
        }
        if err := db.Create(&leaves).Error; err != nil {
                return err
        }

        fields, _ := json.Marshal([]map[string]any{
                {"key": "reason", "label": "申请原因", "type": "textarea", "required": true},
                {"key": "env", "label": "环境", "type": "select", "required": true, "options": []string{"dev", "test", "prod"}},
                {"key": "cpu", "label": "CPU核数", "type": "number", "required": true},
                {"key": "mem", "label": "内存GB", "type": "number", "required": true},
                {"key": "urgent", "label": "是否紧急", "type": "switch", "required": false},
        })
        form := model.FormTemplate{
                Name:        "虚拟机资源申请",
                Description: "对齐课程：动态表单设计器字段",
                FieldsJSON:  string(fields),
                Enabled:     true,
        }
        if err := db.Create(&form).Error; err != nil {
                return err
        }

        nodes, _ := json.Marshal([]map[string]any{
                {"key": "start", "name": "开始", "type": "start"},
                {"key": "mgr_approve", "name": "主管审批", "type": "approve", "role": "approver"},
                {"key": "gateway", "name": "紧急排他网关", "type": "exclusive"},
                {"key": "exec", "name": "运维执行", "type": "execute", "role": "executor"},
                {"key": "end", "name": "结束", "type": "end"},
        })
        edges, _ := json.Marshal([]map[string]any{
                {"from": "start", "to": "mgr_approve"},
                {"from": "mgr_approve", "to": "gateway", "on": "approve"},
                {"from": "gateway", "to": "exec", "condition": "urgent==true"},
                {"from": "gateway", "to": "exec", "condition": "default"},
                {"from": "exec", "to": "end", "on": "execute"},
        })
        flow := model.FlowDefinition{
                Name:      "VM申请标准流程",
                FormID:    form.ID,
                NodesJSON: string(nodes),
                EdgesJSON: string(edges),
                Enabled:   true,
        }
        if err := db.Create(&flow).Error; err != nil {
                return err
        }

        formData, _ := json.Marshal(map[string]any{
                "reason": "大促扩容", "env": "prod", "cpu": 8, "mem": 32, "urgent": true,
        })
        deadline := time.Now().Add(4 * time.Hour)
        leafID := leaves[0].ID
        assignee := users[1].ID
        ticket := model.Ticket{
                TicketNo:     "WO-" + time.Now().Format("20060102") + "-0001",
                Title:        "【演示】订单服务扩容 8C32G",
                FlowID:       flow.ID,
                FormID:       form.ID,
                TreeNodeID:   &leafID,
                CreatorID:    users[1].ID,
                AssigneeID:   &assignee,
                Status:       model.StatusPending,
                CurrentNode:  "mgr_approve",
                Priority:     "P1",
                FormDataJSON: string(formData),
                SLADeadline:  &deadline,
        }
        return db.Create(&ticket).Error
}
