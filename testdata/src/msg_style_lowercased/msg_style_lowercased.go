package msg_style_lowercased

import "log/slog"

func _() {
	slog.Info("")
	slog.Info("msg")
	slog.Info("Msg") // want `message should be lowercased`

	// Special cases:
	slog.Info("U.S.")
	slog.Info("HTTP")
}
