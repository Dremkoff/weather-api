package logger

import (
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/logger"))

		log.Info("Setting up logger middleware")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			t1 := time.Now()

			defer func() {
				entry.Info(
					"request completed",
					slog.Duration("duration", time.Since(t1)),
				)
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
