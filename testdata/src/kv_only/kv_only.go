package kv_only

import (
	"context"
	"log/slog"
)

func _(ctx context.Context, logger *slog.Logger) {
	slog.Info("msg", "foo", 1, "bar", 2)
	slog.Info("msg", "foo", 1, slog.Group("group", "bar", 2))
	slog.Log(ctx, slog.LevelInfo, "msg", "foo", 1)
	logger.Log(ctx, slog.LevelInfo, "msg", "foo", 1)

	slog.Info("msg", "foo", 1, slog.Int("bar", 2))                           // want `attributes should not be used`
	slog.Info("msg", "foo", 1, slog.GroupAttrs("group", slog.Int("bar", 2))) // want `use slog.Group with key-value pairs instead`
	slog.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("foo", 1))            // want `use slog.Log with key-value pairs instead`
	logger.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("foo", 1))          // want `use slog.Logger.Log with key-value pairs instead`
}
