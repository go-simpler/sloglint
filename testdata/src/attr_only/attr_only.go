package attr_only

import (
	"context"
	"log/slog"
)

func _(ctx context.Context, logger *slog.Logger) {
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2))
	slog.Info("msg", slog.Int("foo", 1), slog.GroupAttrs("group", slog.Int("bar", 2)))
	slog.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("foo", 1))
	logger.LogAttrs(ctx, slog.LevelInfo, "msg", slog.Int("foo", 1))

	slog.Info("msg", slog.Int("foo", 1), "bar", 2)                      // want `key-value pairs should not be used`
	slog.Info("msg", slog.Int("foo", 1), slog.Group("group", "bar", 2)) // want `use slog.GroupAttrs with attributes instead`
	slog.Log(ctx, slog.LevelInfo, "msg", "foo", 1)                      // want `use slog.LogAttrs with attributes instead`
	logger.Log(ctx, slog.LevelInfo, "msg", "foo", 1)                    // want `use slog.Logger.LogAttrs with attributes instead`
}
