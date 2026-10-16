package msg_style_lowercased

import "log/slog"

func _() {
	slog.Info("")
	slog.Info("msg")
	slog.Info("Msg")         // want `message should be lowercased`
	slog.Info(" Msg")        // want `message should be lowercased`
	slog.Info("\tMsg")       // want `message should be lowercased`
	slog.Info("\nMsg")       // want `message should be lowercased`
	slog.Info("\u004dsg")    // want `message should be lowercased`
	slog.Info(`Msg`)         // want `message should be lowercased`
	slog.Info("Ünicode msg") // want `message should be lowercased`

	// Special cases:
	slog.Info("200 OK")
	slog.Info("U.S.")
	slog.Info("HTTP")
}
