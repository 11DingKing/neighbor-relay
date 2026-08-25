package httpapi

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestPrivateForbiddenRequestKeeps403(t *testing.T) {
	ctx := context.Background()
	db, e := storage.Open(ctx, filepath.Join(t.TempDir(), "h.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	a := service.New(db, config.Load(), slog.Default())
	u := domain.User{ID: "w7", Name: "W", Email: "w7@e", Role: domain.RoleFieldWorker, Active: true}
	if e = a.EnsureUser(ctx, u, "p"); e != nil {
		t.Fatal(e)
	}
	_, tok, e := a.Login(ctx, u.Email, "p")
	if e != nil {
		t.Fatal(e)
	}
	r := httptest.NewRequest("POST", "/v1/households", nil)
	r.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	New(a, slog.Default()).ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}
