package model

import "time"

const (
        StatusDraft     = "draft"
        StatusPending   = "pending_approval"
        StatusApproved  = "approved"
        StatusExecuting = "executing"
        StatusDone      = "done"
        StatusRejected  = "rejected"
        StatusClosed    = "closed"
)

type User struct {
        ID        uint      `gorm:"primaryKey" json:"id"`
        Username  string    `gorm:"uniqueIndex;size:64" json:"username"`
        Display   string    `gorm:"size:64" json:"display"`
        Role      string    `gorm:"size:32" json:"role"`
        Password  string    `json:"-"`
        CreatedAt time.Time `json:"createdAt"`
}

type ServiceTreeNode struct {
        ID       uint   `gorm:"primaryKey" json:"id"`
        ParentID *uint  `json:"parentId"`
        Name     string `gorm:"size:128" json:"name"`
        Path     string `gorm:"size:512;index" json:"path"`
        Leaf     bool   `json:"leaf"`
}

type FormTemplate struct {
        ID          uint      `gorm:"primaryKey" json:"id"`
        Name        string    `gorm:"size:128" json:"name"`
        Description string    `gorm:"size:512" json:"description"`
        FieldsJSON  string    `gorm:"type:text" json:"fields"`
        Enabled     bool      `json:"enabled"`
        CreatedAt   time.Time `json:"createdAt"`
}

type FlowDefinition struct {
        ID        uint      `gorm:"primaryKey" json:"id"`
        Name      string    `gorm:"size:128" json:"name"`
        FormID    uint      `json:"formId"`
        NodesJSON string    `gorm:"type:text" json:"nodes"`
        EdgesJSON string    `gorm:"type:text" json:"edges"`
        Enabled   bool      `json:"enabled"`
        CreatedAt time.Time `json:"createdAt"`
}

type Ticket struct {
        ID           uint       `gorm:"primaryKey" json:"id"`
        TicketNo     string     `gorm:"uniqueIndex;size:32" json:"ticketNo"`
        Title        string     `gorm:"size:256" json:"title"`
        FlowID       uint       `json:"flowId"`
        FormID       uint       `json:"formId"`
        TreeNodeID   *uint      `json:"treeNodeId"`
        CreatorID    uint       `json:"creatorId"`
        AssigneeID   *uint      `json:"assigneeId"`
        Status       string     `gorm:"size:32;index" json:"status"`
        CurrentNode  string     `gorm:"size:64" json:"currentNode"`
        Priority     string     `gorm:"size:16" json:"priority"`
        FormDataJSON string     `gorm:"type:text" json:"formData"`
        SLADeadline  *time.Time `json:"slaDeadline"`
        SLABreached  bool       `json:"slaBreached"`
        CreatedAt    time.Time  `json:"createdAt"`
        UpdatedAt    time.Time  `json:"updatedAt"`
        ClosedAt     *time.Time `json:"closedAt"`

        Creator  *User            `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
        Assignee *User            `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
        TreeNode *ServiceTreeNode `gorm:"foreignKey:TreeNodeID" json:"treeNode,omitempty"`
}

type TicketAction struct {
        ID        uint      `gorm:"primaryKey" json:"id"`
        TicketID  uint      `gorm:"index" json:"ticketId"`
        ActorID   uint      `json:"actorId"`
        Action    string    `gorm:"size:32" json:"action"`
        FromNode  string    `gorm:"size:64" json:"fromNode"`
        ToNode    string    `gorm:"size:64" json:"toNode"`
        Comment   string    `gorm:"size:1024" json:"comment"`
        CreatedAt time.Time `json:"createdAt"`
        Actor     *User     `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

type AuditLog struct {
        ID        uint      `gorm:"primaryKey" json:"id"`
        ActorID   uint      `json:"actorId"`
        Action    string    `gorm:"size:64" json:"action"`
        Resource  string    `gorm:"size:64" json:"resource"`
        Detail    string    `gorm:"type:text" json:"detail"`
        CreatedAt time.Time `json:"createdAt"`
}

type NotifyRecord struct {
        ID        uint      `gorm:"primaryKey" json:"id"`
        TicketID  uint      `gorm:"index" json:"ticketId"`
        Channel   string    `gorm:"size:32" json:"channel"`
        Target    string    `gorm:"size:128" json:"target"`
        Content   string    `gorm:"size:1024" json:"content"`
        Status    string    `gorm:"size:16" json:"status"`
        CreatedAt time.Time `json:"createdAt"`
}
