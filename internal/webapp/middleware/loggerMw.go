package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/expresso/internal/logger"
)

type LoggerMiddleware struct {
	l *slog.Logger
}

func NewLoggerMiddleware(l *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		l: l.With(logger.LoggerNameField, "LoggerMiddleware"),
	}
}

func (mw *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.UpgradeWithRequestId(r.Context(), RequestIdKey{}, mw.l)
		logline := fmt.Sprintf("%s %s", r.Method, r.URL)
		l.Info(logline)
		next.ServeHTTP(w, r)
	})
}
