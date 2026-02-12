package no_global_default

import "log/slog"

var logger *slog.Logger

func _() {
	slog.Info("msg") // want `default logger should not be used`
	logger.Info("msg")
}
