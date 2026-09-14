package function_names

import "log/slog"

type Logger struct{}

func (*Logger) Info(msg string, args ...any) {}

func custom(msg string, args ...any) {}

func unconfigured(msg string, args ...any) {}

func calls(logger *Logger, standard *slog.Logger, msg string, dynamic func(string, ...any)) {
	const key = "key"

	slog.Info(msg, key, 1)     // want `message should be a string literal or a constant`
	standard.Info(msg, key, 1) // want `message should be a string literal or a constant`
	custom(msg, key, 1)        // want `message should be a string literal or a constant`
	logger.Info(msg, key, 1)   // want `message should be a string literal or a constant`

	custom("message", "raw", 1)   // want `the "raw" key should be a constant`
	slog.Group("group", "raw", 1) // want `the "group" key should be a constant` `the "raw" key should be a constant`
	slog.Int("raw", 1)            // want `the "raw" key should be a constant`

	unconfigured(msg, "raw", 1)
	dynamic(msg, "raw", 1)
	custom("message", key, 1)
}
