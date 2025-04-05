package discard_handler

import (
	"io"
	"log/slog"
)

func _() {
	_ = slog.NewTextHandler(io.Discard, nil) // want `use slog.DiscardHandler instead`
	_ = slog.NewJSONHandler(io.Discard, nil) // want `use slog.DiscardHandler instead`
}
