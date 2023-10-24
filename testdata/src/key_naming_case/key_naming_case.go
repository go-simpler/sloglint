package key_naming_case

import "log/slog"

const (
	snakeKey = "foo_bar"
	kebabKey = "foo-bar"
)

func tests() {
	slog.Info("msg")
	slog.Info("msg", "foo_bar", 1)
	slog.Info("msg", snakeKey, 1)
	slog.Info("msg", slog.Int("foo_bar", 1))
	slog.Info("msg", slog.Int(snakeKey, 1))
	slog.Info("msg", slog.Attr{})
	slog.Info("msg", slog.Attr{"foo_bar", slog.IntValue(1)})
	slog.Info("msg", slog.Attr{snakeKey, slog.IntValue(1)})
	slog.Info("msg", slog.Attr{Key: "foo_bar"})
	slog.Info("msg", slog.Attr{Key: snakeKey})
	slog.Info("msg", slog.Attr{Key: "foo_bar", Value: slog.IntValue(1)})
	slog.Info("msg", slog.Attr{Key: snakeKey, Value: slog.IntValue(1)})
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: "foo_bar"})
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: snakeKey})

	slog.Info("msg", "foo-bar", 1)                                       // want `keys should be written in snake_case`
	slog.Info("msg", kebabKey, 1)                                        // want `keys should be written in snake_case`
	slog.Info("msg", slog.Int("foo-bar", 1))                             // want `keys should be written in snake_case`
	slog.Info("msg", slog.Int(kebabKey, 1))                              // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{"foo-bar", slog.IntValue(1)})             // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{kebabKey, slog.IntValue(1)})              // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{Key: "foo-bar"})                          // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{Key: kebabKey})                           // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{Key: "foo-bar", Value: slog.IntValue(1)}) // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{Key: kebabKey, Value: slog.IntValue(1)})  // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: "foo-bar"}) // want `keys should be written in snake_case`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: kebabKey})  // want `keys should be written in snake_case`
}
