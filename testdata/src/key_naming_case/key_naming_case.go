package key_naming_case

import "log/slog"

const (
	snakeKey       = "foo_bar"
	kebabKey       = "foo-bar"
	dottedSnakeKey = "http.status_code"
	dottedKebabKey = "http.status-code"
)

func _() {
	slog.Info("msg", "foo_bar", 1)
	slog.Info("msg", snakeKey, 1)
	slog.Info("msg", slog.Int("foo_bar", 1))
	slog.Info("msg", slog.Int(snakeKey, 1))

	slog.Info("msg", "http.status_code", 1)
	slog.Info("msg", "service.instance.id", 1)
	slog.Info("msg", dottedSnakeKey, 1)
	slog.Info("msg", slog.Int("http.status_code", 1))

	slog.Info("msg", "foo-bar", 1)           // want `keys should be written in snake_case`
	slog.Info("msg", kebabKey, 1)            // want `keys should be written in snake_case`
	slog.Info("msg", slog.Int("foo-bar", 1)) // want `keys should be written in snake_case`
	slog.Info("msg", slog.Int(kebabKey, 1))  // want `keys should be written in snake_case`

	slog.Info("msg", "http.status-code", 1)           // want `keys should be written in snake_case`
	slog.Info("msg", "http.statusCode", 1)            // want `keys should be written in snake_case`
	slog.Info("msg", "Http.Status_Code", 1)           // want `keys should be written in snake_case`
	slog.Info("msg", dottedKebabKey, 1)               // want `keys should be written in snake_case`
	slog.Info("msg", slog.Int("http.status-code", 1)) // want `keys should be written in snake_case`
	slog.Info("msg", slog.Int("http.statusCode", 1))  // want `keys should be written in snake_case`
	slog.Info("msg", slog.Int(dottedKebabKey, 1))     // want `keys should be written in snake_case`
}
