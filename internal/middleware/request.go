package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"net/http"
)

type key string

const RequestIDKey key = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 8)
			_, _ = rand.Read(b)
			id = hex.EncodeToString(b)
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func ID(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return v
	}
	return "unknown"
}
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				WriteError(w, r, domain.ErrInvalid)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func WriteError(w http.ResponseWriter, r *http.Request, e error) {
	status := http.StatusInternalServerError
	code := domain.CodeInternal
	switch e {
	case domain.ErrUnauthorized, domain.ErrExpired:
		status = http.StatusUnauthorized
		code = domain.CodeUnauthorized
	case domain.ErrForbidden:
		status = http.StatusForbidden
		code = domain.CodeForbidden
	case domain.ErrNotFound:
		status = http.StatusNotFound
		code = domain.CodeNotFound
	case domain.ErrConflict:
		status = http.StatusConflict
		code = domain.CodeConflict
	case domain.ErrInvalid:
		status = http.StatusBadRequest
		code = domain.CodeInvalid
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":{"code":"` + string(code) + `","message":"request could not be completed","request_id":"` + ID(r.Context()) + `"}}`))
}
