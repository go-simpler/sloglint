package msg_style_capitalized

import "log/slog"

func _() {
	slog.Info("")
	slog.Info("msg")         // want `message should be capitalized`
	slog.Info(" msg")        // want `message should be capitalized`
	slog.Info("\tmsg")       // want `message should be capitalized`
	slog.Info("\nmsg")       // want `message should be capitalized`
	slog.Info("\u006dsg")    // want `message should be capitalized`
	slog.Info(`msg`)         // want `message should be capitalized`
	slog.Info("ünicode msg") // want `message should be capitalized`
	slog.Info("Msg")

	// Special cases:
	slog.Info("200 OK")
	slog.Info("iPhone")
}
