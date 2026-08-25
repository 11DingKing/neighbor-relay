package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/repository"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"strings"
	"time"
)

type Service struct {
	DB          *storage.DB
	Users       repository.Users
	Sessions    repository.Sessions
	Households  repository.Households
	Cases       repository.Cases
	Visits      repository.Visits
	Audit       repository.Audit
	Outbox      repository.Outbox
	Idempotency repository.Idempotency
	Config      config.Config
	Logger      *slog.Logger
}

const timeLayout = "2006-01-02T15:04:05.999999999Z07:00"

func New(db *storage.DB, c config.Config, l *slog.Logger) *Service {
	return &Service{DB: db, Users: repository.Users{DB: db}, Sessions: repository.Sessions{DB: db}, Households: repository.Households{DB: db}, Cases: repository.Cases{DB: db}, Visits: repository.Visits{DB: db}, Audit: repository.Audit{DB: db}, Outbox: repository.Outbox{DB: db}, Idempotency: repository.Idempotency{DB: db}, Config: c, Logger: l}
}
func ID() string {
	b := make([]byte, 16)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func Hash(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}
func Encode(v any) string { b, _ := json.Marshal(v); return string(b) }
func (s *Service) Authenticate(ctx context.Context, token string) (domain.User, error) {
	if strings.TrimSpace(token) == "" {
		return domain.User{}, domain.ErrUnauthorized
	}
	session, e := s.Sessions.Active(ctx, Hash(token), time.Now().UTC())
	if e != nil {
		return domain.User{}, e
	}
	u, e := s.Users.ByID(ctx, session.UserID)
	if e != nil {
		return domain.User{}, e
	}
	if !u.Active {
		return domain.User{}, domain.ErrUnauthorized
	}
	return u, nil
}
func (s *Service) RequireRole(u domain.User, roles ...domain.Role) error {
	for _, r := range roles {
		if u.Role == r {
			return nil
		}
	}
	return domain.ErrForbidden
}
func mapSQLError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint") {
		return domain.ErrConflict
	}
	return err
}
func isNotFound(err error) bool { return errors.Is(err, domain.ErrNotFound) }

var _ = fmt.Sprintf
