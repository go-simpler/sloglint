package msg_format_upper

import (
	"context"
	"log/slog"
)

func tests() {
	ctx := context.Background()

	slog.Info("")
	slog.Info("Msg")
	slog.InfoContext(ctx, "Msg")
	slog.Log(ctx, slog.LevelInfo, "Msg")
	slog.With("key", "value").Info("Msg")

	slog.Info("msg")                      // want `message should start with uppercase character`
	slog.InfoContext(ctx, "msg")          // want `message should start with uppercase character`
	slog.Log(ctx, slog.LevelInfo, "msg")  // want `message should start with uppercase character`
	slog.With("key", "value").Info("msg") // want `message should start with uppercase character`

	// special cases:
	slog.Info("U.S. dollar")
	slog.Info("HTTP request")
	slog.Info("iPhone 18")
}
