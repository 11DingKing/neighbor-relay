package config

import (
	"os"
	"time"
)

type Config struct {
	Addr, DBPath               string
	SessionTTL, WorkerInterval time.Duration
	WorkerMaxAttempts          int
}

func Load() Config {
	return Config{Addr: env("ADDR", ":8080"), DBPath: env("DB_PATH", "./neighbor-relay.db"), SessionTTL: duration("SESSION_TTL", 12*time.Hour), WorkerInterval: duration("WORKER_INTERVAL", 2*time.Second), WorkerMaxAttempts: 5}
}
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func duration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return fallback
}
