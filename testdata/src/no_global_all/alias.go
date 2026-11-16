package no_global_all

import sl "log/slog"

func _() {
	sl.Info("msg") // want `default logger should not be used`
}
