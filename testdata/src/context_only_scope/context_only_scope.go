package context_only_scope

import (
	"context"
	"log/slog"
)

func tests(ctx context.Context) {
	slog.Info("msg") // want `InfoContext should be used instead`
	slog.InfoContext(ctx, "msg")

	slog.With("key", "value").Info("msg") // want `InfoContext should be used instead`
	slog.With("key", "value").InfoContext(ctx, "msg")

	if true {
		slog.Info("msg") // want `InfoContext should be used instead`
		slog.InfoContext(ctx, "msg")
	}

	_ = func() {
		slog.Info("msg") // want `InfoContext should be used instead`
		slog.InfoContext(ctx, "msg")
	}
}

func noctx() {
	slog.Info("msg")
}
