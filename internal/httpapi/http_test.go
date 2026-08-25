package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func httpServer(t *testing.T) (http.Handler, *service.Service) {
	d, e := storage.Open(context.Background(), filepath.Join(t.TempDir(), "h.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { d.Close() })
	c := config.Load()
	s := service.New(d, c, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	u := domain.User{ID: "a", Name: "Admin", Email: "a@x", Role: domain.RoleAdmin, Active: true}
	if e = s.EnsureUser(context.Background(), u, "pw"); e != nil {
		t.Fatal(e)
	}
	return New(s, slog.Default()), s
}
func req(t *testing.T, h http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	q := httptest.NewRequest(method, path, r)
	if token != "" {
		q.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, q)
	return w
}
func TestHealthAndReady(t *testing.T) {
	h, _ := httpServer(t)
	if w := req(t, h, "GET", "/healthz", nil, ""); w.Code != 200 {
		t.Fatal(w.Code)
	}
	if w := req(t, h, "GET", "/readyz", nil, ""); w.Code != 200 {
		t.Fatal(w.Code)
	}
}
func TestAuthRequired(t *testing.T) {
	h, _ := httpServer(t)
	if w := req(t, h, "GET", "/v1/households", nil, ""); w.Code != 401 {
		t.Fatal(w.Code)
	}
	if w := req(t, h, "POST", "/v1/auth/login", map[string]string{"email": "a@x", "password": "bad"}, ""); w.Code != 401 {
		t.Fatal(w.Code)
	}
}
func TestLoginAndHousehold(t *testing.T) {
	h, s := httpServer(t)
	w := req(t, h, "POST", "/v1/auth/login", map[string]string{"email": "a@x", "password": "pw"}, "")
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var out struct{ Token string }
	if e := json.Unmarshal(w.Body.Bytes(), &out); e != nil || out.Token == "" {
		t.Fatal(out, e)
	}
	w = req(t, h, "POST", "/v1/households", map[string]any{"address": "1 Main", "contact_name": "A", "phone": "p", "priority": 1}, out.Token)
	if w.Code != 201 {
		t.Fatal(w.Code, w.Body.String())
	}
	items, e := s.ListHouseholds(context.Background(), adminForHTTP(), 10)
	if e != nil || len(items) != 1 {
		t.Fatalf("items=%d err=%v", len(items), e)
	}
}
func TestRequestIDEcho(t *testing.T) {
	h, _ := httpServer(t)
	q := httptest.NewRequest("GET", "/healthz", nil)
	q.Header.Set("X-Request-ID", "abc")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, q)
	if w.Header().Get("X-Request-ID") != "abc" {
		t.Fatal("request id missing")
	}
}
func adminForHTTP() domain.User {
	return domain.User{ID: "a", Name: "Admin", Email: "a@x", Role: domain.RoleAdmin, Active: true}
}
func TestMalformedJSON(t *testing.T) {
	h, _ := httpServer(t)
	q := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString("{"))
	q.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, q)
	if w.Code != 400 {
		t.Fatal(w.Code)
	}
}

var _ = time.Second
