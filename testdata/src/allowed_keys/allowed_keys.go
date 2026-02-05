package allowed_keys

import "log/slog"

const (
	snakeKey = "foo-bar"
)

func tests() {
	slog.Info("msg")
	slog.Info("msg", "foo_bar", 1)
	slog.With("foo-bar", 1).Info("msg")      // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", "foo-bar", 1)           // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", snakeKey, 1)            // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Int("foo-bar", 1)) // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Int(snakeKey, 1))  // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{})
	slog.Info("msg", slog.Attr{"foo-bar", slog.IntValue(1)})                  // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{snakeKey, slog.IntValue(1)})                   // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{Key: "foo-bar"})                               // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{Key: snakeKey})                                // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{Key: "foo-bar", Value: slog.IntValue(1)})      // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{Key: snakeKey, Value: slog.IntValue(1)})       // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: "foo-bar"})      // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: snakeKey})       // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.Info("msg", slog.Attr{Value: slog.IntValue(1), Key: `foo-bar`})      // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{"foo-bar", slog.IntValue(1)}).Info("msg")             // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{snakeKey, slog.IntValue(1)}).Info("msg")              // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{Key: "foo-bar"}).Info("msg")                          // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{Key: snakeKey}).Info("msg")                           // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{Key: "foo-bar", Value: slog.IntValue(1)}).Info("msg") // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{Key: snakeKey, Value: slog.IntValue(1)}).Info("msg")  // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: "foo-bar"}).Info("msg") // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: snakeKey}).Info("msg")  // want `"foo-bar" key is not in the allowed keys list and should not be used`
	slog.With(slog.Attr{Value: slog.IntValue(1), Key: `foo-bar`}).Info("msg") // want `"foo-bar" key is not in the allowed keys list and should not be used`
}
