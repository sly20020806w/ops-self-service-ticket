package main

import (
        "log"

        "github.com/gin-gonic/gin"
        "github.com/sly20020806w/ops-self-service-ticket/internal/config"
        "github.com/sly20020806w/ops-self-service-ticket/internal/handler"
        "github.com/sly20020806w/ops-self-service-ticket/internal/middleware"
        "github.com/sly20020806w/ops-self-service-ticket/internal/notify"
        "github.com/sly20020806w/ops-self-service-ticket/internal/service"
        "github.com/sly20020806w/ops-self-service-ticket/internal/store"
)

func main() {
        cfg := config.Load()
        db, err := store.Open(cfg.DBPath)
        if err != nil {
                log.Fatal(err)
        }
        if cfg.SeedDemo {
                if err := store.SeedDemo(db); err != nil {
                        log.Fatal("seed:", err)
                }
        }

        n := notify.New(db, cfg.NotifyWebhook)
        ts := &service.TicketService{DB: db, Notifier: n}
        api := &handler.API{DB: db, Secret: cfg.JWTSecret, Ticket: ts}

        r := gin.Default()
        r.Use(middleware.CORS())
        api.Register(r)

        log.Println("ops-self-service-ticket listening on", cfg.Addr)
        if err := r.Run(cfg.Addr); err != nil {
                log.Fatal(err)
        }
}
