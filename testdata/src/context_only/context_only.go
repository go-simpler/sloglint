package context_only

import (
	"context"
	"log/slog"
)

func tests() {
	logger := slog.New(nil)
	ctx := context.Background()

	slog.Log(ctx, slog.LevelInfo, "msg")
	slog.DebugContext(ctx, "msg")
	slog.InfoContext(ctx, "msg")
	slog.WarnContext(ctx, "msg")
	slog.ErrorContext(ctx, "msg")

	logger.Log(ctx, slog.LevelInfo, "msg")
	logger.DebugContext(ctx, "msg")
	logger.InfoContext(ctx, "msg")
	logger.WarnContext(ctx, "msg")
	logger.ErrorContext(ctx, "msg")

	slog.Debug("msg") // want `methods without a context should not be used`
	slog.Info("msg")  // want `methods without a context should not be used`
	slog.Warn("msg")  // want `methods without a context should not be used`
	slog.Error("msg") // want `methods without a context should not be used`

	logger.Debug("msg") // want `methods without a context should not be used`
	logger.Info("msg")  // want `methods without a context should not be used`
	logger.Warn("msg")  // want `methods without a context should not be used`
	logger.Error("msg") // want `methods without a context should not be used`
}
