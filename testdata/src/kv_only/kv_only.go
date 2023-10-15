package kv_only

import "log/slog"

func tests() {
	slog.Info("msg")
	slog.Info("msg", "foo", 1)
	slog.Info("msg", "foo", 1, "bar", 2)

	slog.Info("msg", slog.Int("foo", 1))                     // want `attributes should not be used`
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2)) // want `attributes should not be used`
	slog.Info("msg", "foo", 1, slog.Int("bar", 2))           // want `attributes should not be used`
}
