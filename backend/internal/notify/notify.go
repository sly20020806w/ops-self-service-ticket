package notify

import (
        "bytes"
        "encoding/json"
        "net/http"
        "time"

        "github.com/sly20020806w/ops-self-service-ticket/internal/model"
        "gorm.io/gorm"
)

// Notifier 对齐课程「工单和 IM 通知整合」，增强 webhook 卡片推送。
type Notifier struct {
        DB      *gorm.DB
        Webhook string
        Client  *http.Client
}

func New(db *gorm.DB, webhook string) *Notifier {
        return &Notifier{
                DB:      db,
                Webhook: webhook,
                Client:  &http.Client{Timeout: 5 * time.Second},
        }
}

func (n *Notifier) Send(ticketID uint, channel, target, content string) {
        rec := model.NotifyRecord{
                TicketID: ticketID,
                Channel:  channel,
                Target:   target,
                Content:  content,
                Status:   "sent",
        }
        if channel == "webhook" && n.Webhook != "" {
                body, _ := json.Marshal(map[string]any{
                        "msg_type": "text",
                        "content":  map[string]string{"text": content},
                })
                resp, err := n.Client.Post(n.Webhook, "application/json", bytes.NewReader(body))
                if err != nil {
                        rec.Status = "failed"
                } else {
                        _ = resp.Body.Close()
                        if resp.StatusCode >= 300 {
                                rec.Status = "failed"
                        }
                }
        }
        _ = n.DB.Create(&rec).Error
}
