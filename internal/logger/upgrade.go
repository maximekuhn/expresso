package logger

import (
	"context"
	"log/slog"
)

func UpgradeWithRequestId(ctx context.Context, key interface{}, l *slog.Logger) *slog.Logger {
	requestId, ok := ctx.Value(key).(string)
	if !ok || requestId == "" {
		requestId = "unknown"
	}
	return l.With(slog.String("requestId", requestId))
}
