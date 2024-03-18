package context_only_scope

import (
	"context"
	"log/slog"
)

func withContext(ctx context.Context) {
	slog.Info("msg") // want `InfoContext should be used instead`
	slog.InfoContext(ctx, "msg")
}

func withoutContext() {
	slog.Info("msg")
}
