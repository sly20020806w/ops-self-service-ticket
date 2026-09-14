package config

import (
        "os"
        "strconv"
)

type Config struct {
        Addr          string
        DBPath        string
        JWTSecret     string
        NotifyWebhook string
        SeedDemo      bool
}

func Load() Config {
        return Config{
                Addr:          getenv("ADDR", ":8080"),
                DBPath:        getenv("DB_PATH", "./data/ticket.db"),
                JWTSecret:     getenv("JWT_SECRET", "ops-ticket-demo-secret-change-me"),
                NotifyWebhook: getenv("NOTIFY_WEBHOOK", ""),
                SeedDemo:      getenvBool("SEED_DEMO", true),
        }
}

func getenv(k, def string) string {
        if v := os.Getenv(k); v != "" {
                return v
        }
        return def
}

func getenvBool(k string, def bool) bool {
        v := os.Getenv(k)
        if v == "" {
                return def
        }
        b, err := strconv.ParseBool(v)
        if err != nil {
                return def
        }
        return b
}
