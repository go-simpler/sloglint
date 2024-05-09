package no_global_default

import "log/slog"

var logger = slog.New(nil)

func tests() {
	slog.Info("msg")          // want `default logger should not be used`
	slog.With("key", "value") // want `default logger should not be used`
	logger.Info("msg")
	logger.With("key", "value")
}
