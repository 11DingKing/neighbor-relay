package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Users struct{ DB *storage.DB }

func (r Users) Create(ctx context.Context, u domain.User, passwordHash string) error {
	active := 0
	if u.Active {
		active = 1
	}
	_, err := r.DB.Exec(ctx, `INSERT INTO users(id,name,email,role,active,created_at) VALUES(?,?,?,?,?,?)`, u.ID, u.Name, u.Email, u.Role, active, u.CreatedAt.UTC().Format(timeLayout))
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(ctx, `INSERT INTO user_credentials(user_id,password_hash) VALUES(?,?)`, u.ID, passwordHash)
	return err
}
func (r Users) ByEmail(ctx context.Context, email string) (domain.User, string, error) {
	var u domain.User
	var role string
	var active int
	var created, hash string
	err := r.DB.QueryRow(ctx, `SELECT u.id,u.name,u.email,u.role,u.active,u.created_at,c.password_hash FROM users u JOIN user_credentials c ON c.user_id=u.id WHERE u.email=?`, email).Scan(&u.ID, &u.Name, &u.Email, &role, &active, &created, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, "", domain.ErrNotFound
		}
		return u, "", err
	}
	u.Role = domain.Role(role)
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(timeLayout, created)
	return u, hash, nil
}
func (r Users) ByID(ctx context.Context, id string) (domain.User, error) {
	var u domain.User
	var role, created string
	var active int
	err := r.DB.QueryRow(ctx, `SELECT id,name,email,role,active,created_at FROM users WHERE id=?`, id).Scan(&u.ID, &u.Name, &u.Email, &role, &active, &created)
	if err == sql.ErrNoRows {
		return u, domain.ErrNotFound
	}
	u.Role = domain.Role(role)
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(timeLayout, created)
	return u, err
}

const timeLayout = "2006-01-02T15:04:05.999999999Z07:00"
