package no_global_all

import "log/slog"

var logger *slog.Logger

func _() {
	slog.Info("msg")   // want `global logger should not be used`
	logger.Info("msg") // want `global logger should not be used`
}
