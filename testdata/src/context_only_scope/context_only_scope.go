package context_only_scope

import (
	"context"
	"log/slog"
	"net/http"
)

func withContext(ctx context.Context) {
	slog.Info("msg") // want `InfoContext should be used instead`
	slog.InfoContext(ctx, "msg")

	_ = func() {
		slog.Info("msg") // want `InfoContext should be used instead`
		slog.InfoContext(ctx, "msg")
	}
}

func withRequest(r *http.Request) {
	slog.Info("msg") // want `InfoContext should be used instead`
	slog.InfoContext(r.Context(), "msg")

	_ = func() {
		slog.Info("msg") // want `InfoContext should be used instead`
		slog.InfoContext(r.Context(), "msg")
	}
}

func withoutContext() {
	slog.Info("msg")
	slog.InfoContext(context.Background(), "msg")

	_ = func(ctx context.Context) {
		slog.Info("msg") // want `InfoContext should be used instead`
		slog.InfoContext(ctx, "msg")
	}
}
