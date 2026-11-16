package no_global_all

import "log/slog"

var logger *slog.Logger

func _() {
	slog.Info("msg")   // want `default logger should not be used`
	logger.Info("msg") // want `global logger should not be used`

	// A local logger named "slog" is neither the default nor a global logger.
	slog := slog.New(nil)
	slog.Info("msg")
}
