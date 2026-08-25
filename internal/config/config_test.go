package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaults(t *testing.T) {
	for _, k := range []string{"ADDR", "DB_PATH", "SESSION_TTL", "WORKER_INTERVAL"} {
		_ = os.Unsetenv(k)
	}
	c := Load()
	if c.Addr == "" || c.DBPath == "" || c.SessionTTL <= 0 || c.WorkerInterval <= 0 {
		t.Fatalf("bad defaults %#v", c)
	}
}
func TestOverrides(t *testing.T) {
	os.Setenv("ADDR", "127.0.0.1:9")
	os.Setenv("SESSION_TTL", "2h")
	os.Setenv("WORKER_INTERVAL", "50ms")
	defer func() { os.Unsetenv("ADDR"); os.Unsetenv("SESSION_TTL"); os.Unsetenv("WORKER_INTERVAL") }()
	c := Load()
	if c.Addr != "127.0.0.1:9" || c.SessionTTL != 2*time.Hour || c.WorkerInterval != 50*time.Millisecond {
		t.Fatal(c)
	}
}
func TestInvalidDurationFallback(t *testing.T) {
	os.Setenv("SESSION_TTL", "bad")
	defer os.Unsetenv("SESSION_TTL")
	if Load().SessionTTL <= 0 {
		t.Fatal("duration not recovered")
	}
}
