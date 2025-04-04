package discard_handler

import (
	"io"
	"log/slog"
)

func tests() {
	slog.NewTextHandler(io.Discard, nil) // want `use slog.DiscardHandler instead`
	slog.NewJSONHandler(io.Discard, nil) // want `use slog.DiscardHandler instead`
}
