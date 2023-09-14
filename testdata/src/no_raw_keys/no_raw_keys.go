package no_raw_keys

import (
	"log/slog"
)

const (
	foo = "foo"
	bar = "bar"
)

func tests() {
	slog.Info("msg", foo, 1)
	slog.Info("msg", slog.Int(foo, 1))
	slog.Info("msg", foo, 1, slog.Int(bar, 2)) // want `key-value pairs and attributes should not be mixed`

	slog.Info("msg", "foo", 1)                     // want `raw keys should not be used`
	slog.Info("msg", slog.Int("foo", 1))           // want `raw keys should not be used`
	slog.Info("msg", "foo", 1, slog.Int("bar", 2)) // want `raw keys should not be used` `key-value pairs and attributes should not be mixed`
}
