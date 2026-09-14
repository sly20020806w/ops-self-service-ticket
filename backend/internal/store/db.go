package store

import (
        "os"
        "path/filepath"

        "github.com/sly20020806w/ops-self-service-ticket/internal/model"
        "github.com/glebarez/sqlite"
        "gorm.io/gorm"
        "gorm.io/gorm/logger"
)

func Open(dbPath string) (*gorm.DB, error) {
        if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
                return nil, err
        }
        db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
                Logger: logger.Default.LogMode(logger.Warn),
        })
        if err != nil {
                return nil, err
        }
        if err := db.AutoMigrate(
                &model.User{},
                &model.ServiceTreeNode{},
                &model.FormTemplate{},
                &model.FlowDefinition{},
                &model.Ticket{},
                &model.TicketAction{},
                &model.AuditLog{},
                &model.NotifyRecord{},
        ); err != nil {
                return nil, err
        }
        return db, nil
}
