package kv_only

import (
	"context"
	"log/slog"
)

func tests() {
	ctx := context.Background()

	slog.Info("msg")
	slog.Info("msg", "foo", 1, "bar", 2)

	slog.Log(ctx, slog.LevelInfo, "msg", "foo", 1, slog.Int("bar", 2)) // want `attributes should not be used`
	slog.Debug("msg", "foo", 1, slog.Int("bar", 2))                    // want `attributes should not be used`
	slog.Info("msg", "foo", 1, slog.Int("bar", 2))                     // want `attributes should not be used`
	slog.Warn("msg", "foo", 1, slog.Int("bar", 2))                     // want `attributes should not be used`
	slog.Error("msg", "foo", 1, slog.Int("bar", 2))                    // want `attributes should not be used`
	slog.DebugContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `attributes should not be used`
	slog.InfoContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `attributes should not be used`
	slog.WarnContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `attributes should not be used`
	slog.ErrorContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `attributes should not be used`

	logger := slog.New(nil)
	logger.Log(ctx, slog.LevelInfo, "msg", "foo", 1, slog.Int("bar", 2)) // want `attributes should not be used`
	logger.Debug("msg", "foo", 1, slog.Int("bar", 2))                    // want `attributes should not be used`
	logger.Info("msg", "foo", 1, slog.Int("bar", 2))                     // want `attributes should not be used`
	logger.Warn("msg", "foo", 1, slog.Int("bar", 2))                     // want `attributes should not be used`
	logger.Error("msg", "foo", 1, slog.Int("bar", 2))                    // want `attributes should not be used`
	logger.DebugContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `attributes should not be used`
	logger.InfoContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `attributes should not be used`
	logger.WarnContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `attributes should not be used`
	logger.ErrorContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `attributes should not be used`
}
