package msg_format_lower

import (
	"context"
	"log/slog"
)

func tests() {
	ctx := context.Background()

	slog.Info("msg")
	slog.Info("žluťoučký kůň úpěl ďábelské ódy")

	slog.Info("U.S. dollars")
	slog.Info("HTTP request failed")

	slog.Info("Msg")                     // want `message literal should start with lower character`
	slog.Info("Ångström is very small")  // want `message literal should start with lower character`
	slog.InfoContext(ctx, "Msg")         // want `message literal should start with lower character`
	slog.Log(ctx, slog.LevelInfo, "Msg") // want `message literal should start with lower character`

	slog.With("key", "value").Info("Msg") // want `message literal should start with lower character`
}
