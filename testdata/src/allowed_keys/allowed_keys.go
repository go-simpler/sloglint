package allowed_keys

import "log/slog"

const (
	fooKey = "foo"
	barKey = "bar"
)

func _() {
	slog.Info("msg", "foo", 1)
	slog.Info("msg", fooKey, 1)
	slog.Info("msg", slog.Int("foo", 1))
	slog.Info("msg", slog.Int(fooKey, 1))

	slog.Info("msg", "bar", 1)            // want `"bar" key is not allowed and should not be used`
	slog.Info("msg", barKey, 1)           // want `"bar" key is not allowed and should not be used`
	slog.Info("msg", slog.Int("bar", 1))  // want `"bar" key is not allowed and should not be used`
	slog.Info("msg", slog.Int(barKey, 1)) // want `"bar" key is not allowed and should not be used`
}
