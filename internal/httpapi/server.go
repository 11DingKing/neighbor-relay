package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/middleware"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	App    *service.Service
	Logger *slog.Logger
}

func New(app *service.Service, l *slog.Logger) http.Handler {
	srv := &Server{App: app, Logger: l}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", srv.health)
	mux.HandleFunc("GET /readyz", srv.ready)
	mux.HandleFunc("POST /v1/auth/login", srv.login)
	mux.HandleFunc("POST /v1/auth/logout", srv.logout)
	mux.HandleFunc("GET /v1/households", srv.households)
	mux.HandleFunc("POST /v1/households", srv.createHousehold)
	mux.HandleFunc("POST /v1/cases", srv.createCase)
	mux.HandleFunc("POST /v1/visits", srv.planVisit)
	mux.HandleFunc("POST /v1/visits/{id}/start", srv.startVisit)
	mux.HandleFunc("POST /v1/visits/{id}/complete", srv.completeVisit)
	return middleware.RequestID(middleware.Recover(mux))
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if err := s.App.DB.SQL.PingContext(r.Context()); err != nil {
		middleware.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}
func decode(r *http.Request, v any) error {
	if r.Body == nil {
		return domain.ErrInvalid
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return domain.ErrInvalid
	}
	return nil
}
func token(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
func (s *Server) user(r *http.Request) (domain.User, error) {
	return s.App.Authenticate(r.Context(), token(r))
}
func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if e := decode(r, &in); e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	u, t, e := s.App.Login(r.Context(), in.Email, in.Password)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	write(w, map[string]any{"user": u, "token": t, "expires_in": int(s.App.Config.SessionTTL.Seconds())})
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if e := s.App.Logout(r.Context(), token(r)); e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) households(w http.ResponseWriter, r *http.Request) {
	u, e := s.user(r)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, e := s.App.ListHouseholds(r.Context(), u, limit)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	write(w, map[string]any{"items": items})
}
func (s *Server) createHousehold(w http.ResponseWriter, r *http.Request) {
	u, e := s.user(r)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	var in struct {
		Address     string `json:"address"`
		ContactName string `json:"contact_name"`
		Phone       string `json:"phone"`
		Priority    int    `json:"priority"`
	}
	if e = decode(r, &in); e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	h, e := s.App.CreateHousehold(r.Context(), u, in.Address, in.ContactName, in.Phone, in.Priority)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	w.WriteHeader(http.StatusCreated)
	write(w, h)
}
func (s *Server) createCase(w http.ResponseWriter, r *http.Request) {
	u, e := s.user(r)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	var in struct {
		HouseholdID string `json:"household_id"`
		Title       string `json:"title"`
		Summary     string `json:"summary"`
	}
	if e = decode(r, &in); e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	c, e := s.App.CreateCase(r.Context(), u, in.HouseholdID, in.Title, in.Summary)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	w.WriteHeader(http.StatusCreated)
	write(w, c)
}
func (s *Server) planVisit(w http.ResponseWriter, r *http.Request) {
	u, e := s.user(r)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	var in struct {
		CaseID         string    `json:"case_id"`
		AssigneeID     string    `json:"assignee_id"`
		IdempotencyKey string    `json:"idempotency_key"`
		ScheduledFor   time.Time `json:"scheduled_for"`
	}
	if e = decode(r, &in); e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	if in.IdempotencyKey == "" {
		in.IdempotencyKey = r.Header.Get("Idempotency-Key")
	}
	v, e := s.App.PlanVisit(r.Context(), u, in.CaseID, in.AssigneeID, in.ScheduledFor, in.IdempotencyKey, middleware.ID(r.Context()))
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	w.WriteHeader(http.StatusCreated)
	write(w, v)
}
func (s *Server) startVisit(w http.ResponseWriter, r *http.Request) {
	u, e := s.user(r)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	e = s.App.StartVisit(r.Context(), u, r.PathValue("id"))
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) completeVisit(w http.ResponseWriter, r *http.Request) {
	u, e := s.user(r)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	var in struct {
		Note string `json:"note"`
	}
	if e = decode(r, &in); e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	e = s.App.CompleteVisit(r.Context(), u, r.PathValue("id"), in.Note)
	if e != nil {
		middleware.WriteError(w, r, e)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

var _ = errors.Is
