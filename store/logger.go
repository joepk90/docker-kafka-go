package store

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct{}

func (Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	kvs := []any{
		slog.String("level", level.String()),
	}

	for k, v := range data {
		kvs = append(kvs, slog.Any(k, v))
	}
	slog.DebugContext(ctx, msg, kvs...)
}
