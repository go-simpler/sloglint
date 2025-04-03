package msg_format_lower

import (
	"context"
	"log/slog"
)

func tests() {
	ctx := context.Background()

	slog.Info("msg")
	slog.InfoContext(ctx, "msg")
	slog.Log(ctx, slog.LevelInfo, "msg")
	slog.With("key", "value").Info("msg")

	slog.Info("Msg")                      // want `message should start with lowercase character`
	slog.InfoContext(ctx, "Msg")          // want `message should start with lowercase character`
	slog.Log(ctx, slog.LevelInfo, "Msg")  // want `message should start with lowercase character`
	slog.With("key", "value").Info("Msg") // want `message should start with lowercase character`

	// special cases:
	slog.Info("U.S. dollar")
	slog.Info("HTTP request")
	slog.Info("iPhone 18")
}
