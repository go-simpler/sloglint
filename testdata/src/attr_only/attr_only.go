package attr_only

import "log/slog"

func _() {
	slog.Info("msg", "foo", 1, "bar", 2)            // want `key-value pairs should not be used`
	slog.Info("msg", "foo", 1, slog.Int("bar", 2))  // want `key-value pairs should not be used`
	slog.Info("msg", "foo", 1, slog.Group("group")) // want `key-value pairs should not be used`
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2))
}
