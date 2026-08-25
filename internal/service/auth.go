package service

import (
	"context"
	"crypto/subtle"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"time"
)

func (s *Service) Login(ctx context.Context, email, password string) (domain.User, string, error) {
	u, stored, e := s.Users.ByEmail(ctx, email)
	if e != nil {
		return u, "", domain.ErrUnauthorized
	}
	if !u.Active {
		return domain.User{}, "", domain.ErrUnauthorized
	}
	given := Hash(password)
	if subtle.ConstantTimeCompare([]byte(given), []byte(stored)) != 1 {
		return domain.User{}, "", domain.ErrUnauthorized
	}
	token := ID() + ID()
	now := time.Now().UTC()
	exp := now.Add(s.Config.SessionTTL)
	ses := domain.Session{ID: ID(), UserID: u.ID, TokenHash: Hash(token), ExpiresAt: &exp, CreatedAt: &now}
	if e = s.Sessions.Create(ctx, ses); e != nil {
		return domain.User{}, "", e
	}
	return u, token, nil
}
func (s *Service) Logout(ctx context.Context, token string) error {
	err := s.Sessions.Revoke(ctx, Hash(token), time.Now().UTC())
	if err != nil {
		return nil
	}
	return nil
}
func (s *Service) EnsureUser(ctx context.Context, u domain.User, password string) error {
	return mapSQLError(s.Users.Create(ctx, u, Hash(password)))
}
