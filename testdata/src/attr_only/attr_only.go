package attr_only

import "log/slog"

func tests() {
	slog.Info("msg")
	slog.Info("msg", slog.Int("foo", 1))
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2))
	slog.With(slog.Int("foo", 1), slog.Int("bar", 2)).Info("msg")

	slog.Info("msg", "foo", 1)                          // want `key-value pairs should not be used`
	slog.Info("msg", "foo", 1, "bar", 2)                // want `key-value pairs should not be used`
	slog.Info("msg", "foo", 1, slog.Int("bar", 2))      // want `key-value pairs should not be used`
	slog.With("foo", 1, slog.Int("bar", 2)).Info("msg") // want `key-value pairs should not be used`

	args := []slog.Attr{slog.Int("foo", 1), slog.Int("bar", 2)}
	slog.LogAttrs(nil, slog.LevelInfo, "msg", args...)
}
