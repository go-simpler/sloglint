package typed_attributes

import sl "log/slog"

func _(n int) {
	sl.Any("alias", n) // want `use slog.Int instead of slog.Any`
}
