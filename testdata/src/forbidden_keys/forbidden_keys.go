package forbidden_keys

import (
	"log/slog"

	"forbidden_keys/keys"
)

const (
	fooKey = "foo"
	barKey = "bar"

	foo, bar = "foo", "bar" // a multi-name declaration
)

func _() {
	slog.Info("msg", "foo", 1)
	slog.Info("msg", fooKey, 1)
	slog.Info("msg", slog.Int("foo", 1))
	slog.Info("msg", slog.Int(fooKey, 1))
	slog.Info("msg", foo, 1)

	slog.Info("msg", "bar", 1)              // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", barKey, 1)             // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int("bar", 1))    // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int(barKey, 1))   // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", bar, 1)                // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", fileKey, 1)            // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int(fileKey, 1))  // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", keys.Bar, 1)           // want `the "bar" key is forbidden and should not be used`
	slog.Info("msg", slog.Int(keys.Bar, 1)) // want `the "bar" key is forbidden and should not be used`
}
