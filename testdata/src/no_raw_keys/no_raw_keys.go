package no_raw_keys

import (
	"log/slog"
	"no_raw_keys/keys"
)

const foo = "foo"

func _() {
	slog.Info("msg", foo, 1)
	slog.Info("msg", keys.Foo, 1)
	slog.Info("msg", slog.Int(foo, 1))
	slog.Info("msg", slog.Attr{})
	slog.Info("msg", slog.Attr{foo, slog.IntValue(1)})
	slog.Info("msg", slog.Attr{Key: foo})
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1)})
	slog.Info("msg", slog.Attr{Key: foo, Value: slog.IntValue(1)})
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: foo})

	slog.Info("msg", "foo", 1)                                       // want `raw keys should not be used`
	slog.Info("msg", slog.Int("foo", 1))                             // want `raw keys should not be used`
	slog.Info("msg", slog.Attr{"foo", slog.IntValue(1)})             // want `raw keys should not be used`
	slog.Info("msg", slog.Attr{Key: "foo"})                          // want `raw keys should not be used`
	slog.Info("msg", slog.Attr{Key: "foo", Value: slog.IntValue(1)}) // want `raw keys should not be used`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: "foo"}) // want `raw keys should not be used`
}
