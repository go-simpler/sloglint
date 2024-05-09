package key_naming_case

import (
	"context"
	"log/slog"
)

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

	// With snake_case key
	slog.Info("msg")
	slog.With("foo_bar", 1).Info("msg")
	slog.With(snakeKey, 1).Info("msg")
	slog.With(slog.Int("foo_bar", 1)).Info("msg")
	slog.With(slog.Int(snakeKey, 1)).Info("msg")
	slog.With(slog.Attr{}).Info("msg")
	slog.With(slog.Attr{"foo_bar", slog.IntValue(1)}).Info("msg")
	slog.With(slog.Attr{snakeKey, slog.IntValue(1)}).Info("msg")
	slog.With(slog.Attr{Key: "foo_bar"}).Info("msg")
	slog.With(slog.Attr{Key: snakeKey}).Info("msg")
	slog.With(slog.Attr{Key: "foo_bar", Value: slog.IntValue(1)}).Info("msg")
	slog.With(slog.Attr{Key: snakeKey, Value: slog.IntValue(1)}).Info("msg")
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: "foo_bar"}).Info("msg")
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: snakeKey}).Info("msg")

	slog.With("foo-bar", 1).Info("msg")                                       // want `keys should be written in snake_case`
	slog.With(kebabKey, 1).Info("msg")                                        // want `keys should be written in snake_case`
	slog.With(slog.Int("foo-bar", 1)).Info("msg")                             // want `keys should be written in snake_case`
	slog.With(slog.Int(kebabKey, 1)).Info("msg")                              // want `keys should be written in snake_case`
	slog.With(slog.Attr{"foo-bar", slog.IntValue(1)}).Info("msg")             // want `keys should be written in snake_case`
	slog.With(slog.Attr{kebabKey, slog.IntValue(1)}).Info("msg")              // want `keys should be written in snake_case`
	slog.With(slog.Attr{Key: "foo-bar"}).Info("msg")                          // want `keys should be written in snake_case`
	slog.With(slog.Attr{Key: kebabKey}).Info("msg")                           // want `keys should be written in snake_case`
	slog.With(slog.Attr{Key: "foo-bar", Value: slog.IntValue(1)}).Info("msg") // want `keys should be written in snake_case`
	slog.With(slog.Attr{Key: kebabKey, Value: slog.IntValue(1)}).Info("msg")  // want `keys should be written in snake_case`
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: "foo-bar"}).Info("msg") // want `keys should be written in snake_case`
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: kebabKey}).Info("msg")  // want `keys should be written in snake_case`

	slog.LogAttrs(context.TODO(), slog.LevelInfo, "msg", slog.Attr{Value: slog.IntValue(1), Key: kebabKey}) // want `keys should be written in snake_case`
}

func issue35() {
	intAttr := slog.Int
	slog.Info("msg", intAttr("foo_bar", 1))
}
