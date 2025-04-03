package msg_style_lowercased

import (
	"context"
	"log/slog"
)

func tests() {
	ctx := context.Background()

	slog.Info("")
	slog.Info("msg")
	slog.InfoContext(ctx, "msg")
	slog.Log(ctx, slog.LevelInfo, "msg")
	slog.With("key", "value").Info("msg")

	slog.Info("Msg")                      // want `message should be lowercased`
	slog.InfoContext(ctx, "Msg")          // want `message should be lowercased`
	slog.Log(ctx, slog.LevelInfo, "Msg")  // want `message should be lowercased`
	slog.With("key", "value").Info("Msg") // want `message should be lowercased`

	// special cases:
	slog.Info("U.S. dollar")
	slog.Info("HTTP request")
	slog.Info("iPhone 18")
}
