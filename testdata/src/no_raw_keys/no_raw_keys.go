package no_raw_keys

import "log/slog"

const foo = "foo"

func Foo(value int) slog.Attr {
	return slog.Int("foo", value)
}

func tests() {
	slog.Info("msg")
	slog.Info("msg", foo, 1)
	slog.Info("msg", Foo(1))
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

	slog.With("foo", 1).Info("msg")                                       // want `raw keys should not be used`
	slog.With(slog.Int("foo", 1)).Info("msg")                             // want `raw keys should not be used`
	slog.With(slog.Attr{"foo", slog.IntValue(1)}).Info("msg")             // want `raw keys should not be used`
	slog.With(slog.Attr{Key: "foo"}).Info("msg")                          // want `raw keys should not be used`
	slog.With(slog.Attr{Key: "foo", Value: slog.IntValue(1)}).Info("msg") // want `raw keys should not be used`
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: "foo"}).Info("msg") // want `raw keys should not be used`
}
