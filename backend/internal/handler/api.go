package handler

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/sly20020806w/ops-self-service-ticket/internal/middleware"
        "github.com/sly20020806w/ops-self-service-ticket/internal/model"
        "github.com/sly20020806w/ops-self-service-ticket/internal/service"
        "golang.org/x/crypto/bcrypt"
        "gorm.io/gorm"
)

type API struct {
        DB     *gorm.DB
        Secret string
        Ticket *service.TicketService
}

func (a *API) Register(r *gin.Engine) {
        r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
        r.GET("/metrics", func(c *gin.Context) { c.String(200, "ops_ticket_up 1\n") })

        api := r.Group("/api")
        api.POST("/login", a.login)

        auth := api.Group("")
        auth.Use(middleware.Auth(a.Secret))
        {
                auth.GET("/me", a.me)
                auth.GET("/dashboard", a.dashboard)

                auth.GET("/forms", a.listForms)
                auth.POST("/forms", a.createForm)
                auth.GET("/flows", a.listFlows)
                auth.POST("/flows", a.createFlow)
                auth.GET("/tree", a.listTree)

                auth.GET("/tickets", a.listTickets)
                auth.POST("/tickets", a.createTicket)
                auth.GET("/tickets/:id", a.getTicket)
                auth.POST("/tickets/:id/actions", a.ticketAction)

                auth.GET("/audits", a.listAudits)
                auth.GET("/notifies", a.listNotifies)
        }
}

func (a *API) login(c *gin.Context) {
        var req struct {
                Username string `json:"username"`
                Password string `json:"password"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
                return
        }
        var u model.User
        if err := a.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
                return
        }
        if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
                return
        }
        token, err := middleware.Sign(a.Secret, u.ID, u.Username, u.Role)
        if err != nil {
                c.JSON(500, gin.H{"error": "签发失败"})
                return
        }
        c.JSON(200, gin.H{"token": token, "user": u})
}

func (a *API) me(c *gin.Context) {
        uid := c.GetUint("uid")
        var u model.User
        _ = a.DB.First(&u, uid)
        c.JSON(200, u)
}

func (a *API) dashboard(c *gin.Context) {
        c.JSON(200, a.Ticket.Dashboard())
}

func (a *API) listForms(c *gin.Context) {
        var list []model.FormTemplate
        _ = a.DB.Order("id desc").Find(&list)
        c.JSON(200, list)
}

func (a *API) createForm(c *gin.Context) {
        var f model.FormTemplate
        if err := c.ShouldBindJSON(&f); err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
        }
        f.Enabled = true
        if err := a.DB.Create(&f).Error; err != nil {
                c.JSON(500, gin.H{"error": err.Error()})
                return
        }
        c.JSON(200, f)
}

func (a *API) listFlows(c *gin.Context) {
        var list []model.FlowDefinition
        _ = a.DB.Order("id desc").Find(&list)
        c.JSON(200, list)
}

func (a *API) createFlow(c *gin.Context) {
        var f model.FlowDefinition
        if err := c.ShouldBindJSON(&f); err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
        }
        f.Enabled = true
        if err := a.DB.Create(&f).Error; err != nil {
                c.JSON(500, gin.H{"error": err.Error()})
                return
        }
        c.JSON(200, f)
}

func (a *API) listTree(c *gin.Context) {
        var list []model.ServiceTreeNode
        _ = a.DB.Order("id asc").Find(&list)
        c.JSON(200, list)
}

func (a *API) listTickets(c *gin.Context) {
        list, err := a.Ticket.List(c.Query("status"))
        if err != nil {
                c.JSON(500, gin.H{"error": err.Error()})
                return
        }
        c.JSON(200, list)
}

func (a *API) createTicket(c *gin.Context) {
        var req struct {
                Title      string         `json:"title"`
                FlowID     uint           `json:"flowId"`
                TreeNodeID *uint          `json:"treeNodeId"`
                Priority   string         `json:"priority"`
                FormData   map[string]any `json:"formData"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
        }
        if req.Priority == "" {
                req.Priority = "P2"
        }
        t, err := a.Ticket.Create(c.GetUint("uid"), req.Title, req.FlowID, req.TreeNodeID, req.Priority, req.FormData)
        if err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
        }
        c.JSON(200, t)
}

func (a *API) getTicket(c *gin.Context) {
        id, _ := strconv.Atoi(c.Param("id"))
        t, acts, err := a.Ticket.Get(uint(id))
        if err != nil {
                c.JSON(404, gin.H{"error": "不存在"})
                return
        }
        c.JSON(200, gin.H{
                "ticket":           t,
                "actions":          acts,
                "availableActions": service.AvailableActions(t.Status),
        })
}

func (a *API) ticketAction(c *gin.Context) {
        id, _ := strconv.Atoi(c.Param("id"))
        var req struct {
                Action  string `json:"action"`
                Comment string `json:"comment"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
        }
        t, err := a.Ticket.Action(uint(id), c.GetUint("uid"), req.Action, req.Comment)
        if err != nil {
                c.JSON(400, gin.H{"error": err.Error()})
                return
        }
        c.JSON(200, t)
}

func (a *API) listAudits(c *gin.Context) {
        var list []model.AuditLog
        _ = a.DB.Order("id desc").Limit(100).Find(&list)
        c.JSON(200, list)
}

func (a *API) listNotifies(c *gin.Context) {
        var list []model.NotifyRecord
        _ = a.DB.Order("id desc").Limit(100).Find(&list)
        c.JSON(200, list)
}
