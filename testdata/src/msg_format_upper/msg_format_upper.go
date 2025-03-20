package msg_format_upper

import (
	"context"
	"log/slog"
)

func tests() {
	ctx := context.Background()

	slog.Info("Msg")
	slog.Info("Žluťoučký kůň úpěl ďábelské ódy")

	slog.Info("U.S. dollars")
	slog.Info("HTTP request failed")
	slog.Info("iPhone 18 Plus Ultra Max")

	slog.Info("msg")                     // want `message literal should start with upper character`
	slog.Info("ångström is very small")  // want `message literal should start with upper character`
	slog.InfoContext(ctx, "msg")         // want `message literal should start with upper character`
	slog.Log(ctx, slog.LevelInfo, "msg") // want `message literal should start with upper character`

	slog.With("key", "value").Info("msg") // want `message literal should start with upper character`
}
