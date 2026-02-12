package msg_style_capitalized

import "log/slog"

func _() {
	slog.Info("")
	slog.Info("msg") // want `message should be capitalized`
	slog.Info("Msg")

	// Special cases:
	slog.Info("iPhone")
}
