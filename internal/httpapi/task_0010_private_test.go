package httpapi

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestPrivateLogoutUnknownTokenReportsNotFound(t *testing.T) {
	db, e := storage.Open(context.Background(), filepath.Join(t.TempDir(), "l.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	a := service.New(db, config.Load(), slog.Default())
	r := httptest.NewRequest("POST", "/v1/auth/logout", nil)
	r.Header.Set("Authorization", "Bearer unknown-token")
	w := httptest.NewRecorder()
	New(a, slog.Default()).ServeHTTP(w, r)
	if w.Code != 404 {
		t.Fatalf("status=%d", w.Code)
	}
}
