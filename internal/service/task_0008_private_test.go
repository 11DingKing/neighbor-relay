package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"path/filepath"
	"sync"
	"testing"
)

func TestPrivateConcurrentCaseCloseConflicts(t *testing.T) {
	ctx := context.Background()
	db, e := storage.Open(ctx, filepath.Join(t.TempDir(), "cc.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	a := New(db, config.Load(), slog.Default())
	u := domain.User{ID: "a8", Name: "A", Email: "a8@e", Role: domain.RoleAdmin, Active: true}
	if e = a.EnsureUser(ctx, u, "p"); e != nil {
		t.Fatal(e)
	}
	h, e := a.CreateHousehold(ctx, u, "1", "R", "1", 1)
	if e != nil {
		t.Fatal(e)
	}
	c, e := a.CreateCase(ctx, u, h.ID, "T", "S")
	if e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); errs <- a.TransitionCase(ctx, u, c.ID, domain.CaseClosed) }()
	}
	wg.Wait()
	close(errs)
	ok := 0
	conf := 0
	for e := range errs {
		if e == nil {
			ok++
		}
		if e == domain.ErrConflict {
			conf++
		}
	}
	if ok != 1 || conf != 1 {
		t.Fatalf("ok=%d conflict=%d", ok, conf)
	}
}
