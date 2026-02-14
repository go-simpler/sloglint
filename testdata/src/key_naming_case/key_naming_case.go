package key_naming_case

import "log/slog"

const (
	snakeKey = "foo_bar"
	kebabKey = "foo-bar"
)

func _() {
	slog.Info("msg", "foo_bar", 1)
	slog.Info("msg", snakeKey, 1)
	slog.Info("msg", slog.Int("foo_bar", 1))
	slog.Info("msg", slog.Int(snakeKey, 1))

	slog.Info("msg", "foo-bar", 1)           // want `keys should be written in snake case`
	slog.Info("msg", kebabKey, 1)            // want `keys should be written in snake case`
	slog.Info("msg", slog.Int("foo-bar", 1)) // want `keys should be written in snake case`
	slog.Info("msg", slog.Int(kebabKey, 1))  // want `keys should be written in snake case`
}
