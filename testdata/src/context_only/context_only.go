package context_only

import (
	"context"
	"io"
	"log/slog"
)

func tests() {
	ctx := context.Background()

	slog.Debug("msg") // want `methods that do not take a context should not be used`
	slog.Info("msg")  // want `methods that do not take a context should not be used`
	slog.Warn("msg")  // want `methods that do not take a context should not be used`
	slog.Error("msg") // want `methods that do not take a context should not be used`

	slog.Log(context.Background(), slog.LevelInfo, "msg")
	slog.DebugContext(context.TODO(), "msg")
	slog.InfoContext(context.WithoutCancel(ctx), "msg")
	slog.WarnContext(ctx, "msg")
	slog.ErrorContext(ctx, "msg")

	logger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))

	logger.Debug("msg") // want `methods that do not take a context should not be used`
	logger.Info("msg")  // want `methods that do not take a context should not be used`
	logger.Warn("msg")  // want `methods that do not take a context should not be used`
	logger.Error("msg") // want `methods that do not take a context should not be used`

	logger.Log(ctx, slog.LevelInfo, "msg")
	logger.DebugContext(ctx, "msg")
	logger.InfoContext(ctx, "msg")
	logger.WarnContext(ctx, "msg")
	logger.ErrorContext(ctx, "msg")
}
