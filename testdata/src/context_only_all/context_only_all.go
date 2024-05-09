package context_only_all

import (
	"context"
	"log/slog"
)

func tests(ctx context.Context) {
	slog.Log(ctx, slog.LevelInfo, "msg")
	slog.DebugContext(ctx, "msg")
	slog.InfoContext(ctx, "msg")
	slog.WarnContext(ctx, "msg")
	slog.ErrorContext(ctx, "msg")
	slog.With("key", "value").ErrorContext(ctx, "msg")

	slog.Debug("msg") // want `DebugContext should be used instead`
	slog.Info("msg")  // want `InfoContext should be used instead`
	slog.Warn("msg")  // want `WarnContext should be used instead`
	slog.Error("msg") // want `ErrorContext should be used instead`

	logger := slog.New(nil)
	logger.Log(ctx, slog.LevelInfo, "msg")
	logger.DebugContext(ctx, "msg")
	logger.InfoContext(ctx, "msg")
	logger.WarnContext(ctx, "msg")
	logger.ErrorContext(ctx, "msg")

	logger.Debug("msg")                      // want `DebugContext should be used instead`
	logger.Info("msg")                       // want `InfoContext should be used instead`
	logger.Warn("msg")                       // want `WarnContext should be used instead`
	logger.Error("msg")                      // want `ErrorContext should be used instead`
	logger.With("key", "value").Error("msg") // want `ErrorContext should be used instead`
}
