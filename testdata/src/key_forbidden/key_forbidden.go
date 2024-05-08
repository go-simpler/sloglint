package key_forbidden

import "log/slog"

const (
	snakeKey = "foo_bar"
)

func tests() {
	slog.Info("msg")
	slog.Info("msg", "foo-bar", 1)
	slog.Info("msg", "foo_bar", 1)           // want `keys include forbidden values`
	slog.Info("msg", snakeKey, 1)            // want `keys include forbidden values`
	slog.Info("msg", slog.Int("foo_bar", 1)) // want `keys include forbidden values`
	slog.Info("msg", slog.Int(snakeKey, 1))  // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{})
	slog.Info("msg", slog.Attr{"foo_bar", slog.IntValue(1)})             // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{snakeKey, slog.IntValue(1)})              // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{Key: "foo_bar"})                          // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{Key: snakeKey})                           // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{Key: "foo_bar", Value: slog.IntValue(1)}) // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{Key: snakeKey, Value: slog.IntValue(1)})  // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: "foo_bar"}) // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: snakeKey})  // want `keys include forbidden values`
	slog.With(slog.Attr{"foo_bar", slog.IntValue(1)}).Info("msg")        // want `keys include forbidden values`
	slog.With(slog.Attr{snakeKey, slog.IntValue(1)})                     // want `keys include forbidden values`
	slog.With(snakeKey, slog.IntValue(1))                                // want `keys include forbidden values`
	slog.With("foo_bar", slog.IntValue(1))                               // want `keys include forbidden values`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: `foo_bar`}) // want `keys include forbidden values`
}
