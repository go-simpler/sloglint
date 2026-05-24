package kv_only

import (
	"context"
	"log/slog"
)

func _(ctx context.Context, logger *slog.Logger) {
	slog.Info("msg", "foo", 1)
	slog.Info("msg", slog.Group("group", "foo", 1))

	slog.Info("msg", slog.Int("foo", 1))                            // want `attributes should not be used`
	slog.GroupAttrs("group", slog.Int("foo", 1))                    // want `use slog.Group with key-value pairs instead`
	slog.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("foo", 1))   // want `use slog.Log with key-value pairs instead`
	logger.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("foo", 1)) // want `use slog.Logger.Log with key-value pairs instead`
}
