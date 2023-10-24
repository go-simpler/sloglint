package mixed_args

import (
	"context"
	"log/slog"
)

func tests() {
	logger := slog.New(nil)
	ctx := context.Background()

	slog.Info("msg")
	slog.Info("msg", "foo", 1)
	slog.Info("msg", "foo", 1, "bar", 2)
	slog.Info("msg", slog.Int("foo", 1))
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2))

	slog.Log(ctx, slog.LevelInfo, "msg", "foo", 1, slog.Int("bar", 2)) // want `key-value pairs and attributes should not be mixed`
	slog.Debug("msg", "foo", 1, slog.Int("bar", 2))                    // want `key-value pairs and attributes should not be mixed`
	slog.Info("msg", "foo", 1, slog.Int("bar", 2))                     // want `key-value pairs and attributes should not be mixed`
	slog.Warn("msg", "foo", 1, slog.Int("bar", 2))                     // want `key-value pairs and attributes should not be mixed`
	slog.Error("msg", "foo", 1, slog.Int("bar", 2))                    // want `key-value pairs and attributes should not be mixed`
	slog.DebugContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `key-value pairs and attributes should not be mixed`
	slog.InfoContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `key-value pairs and attributes should not be mixed`
	slog.WarnContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `key-value pairs and attributes should not be mixed`
	slog.ErrorContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `key-value pairs and attributes should not be mixed`

	logger.Log(ctx, slog.LevelInfo, "msg", "foo", 1, slog.Int("bar", 2)) // want `key-value pairs and attributes should not be mixed`
	logger.Debug("msg", "foo", 1, slog.Int("bar", 2))                    // want `key-value pairs and attributes should not be mixed`
	logger.Info("msg", "foo", 1, slog.Int("bar", 2))                     // want `key-value pairs and attributes should not be mixed`
	logger.Warn("msg", "foo", 1, slog.Int("bar", 2))                     // want `key-value pairs and attributes should not be mixed`
	logger.Error("msg", "foo", 1, slog.Int("bar", 2))                    // want `key-value pairs and attributes should not be mixed`
	logger.DebugContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `key-value pairs and attributes should not be mixed`
	logger.InfoContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `key-value pairs and attributes should not be mixed`
	logger.WarnContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))         // want `key-value pairs and attributes should not be mixed`
	logger.ErrorContext(ctx, "msg", "foo", 1, slog.Int("bar", 2))        // want `key-value pairs and attributes should not be mixed`
}
