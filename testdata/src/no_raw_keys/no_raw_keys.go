package no_raw_keys

import (
	"log/slog"
)

const foo = "foo"

func Foo(value int) slog.Attr {
	return slog.Int("foo", value)
}

func tests() {
	slog.Info("msg", foo, 1)
	slog.Info("msg", Foo(1))
	slog.Info("msg", slog.Int(foo, 1))

	slog.Info("msg", "foo", 1)           // want `raw keys should not be used`
	slog.Info("msg", slog.Int("foo", 1)) // want `raw keys should not be used`
}
