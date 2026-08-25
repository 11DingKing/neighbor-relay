package observability

import (
	"log/slog"
	"net/http"
	"time"
)

func Request(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("http request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
func FieldError(logger *slog.Logger, op string, err error) {
	logger.Error("operation failed", "operation", op, "error", err)
}
func Job(logger *slog.Logger, id, status string) {
	logger.Info("outbox job", "job_id", id, "status", status)
}
