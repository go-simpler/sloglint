package typed_attributes

import . "log/slog"

func _(Int func(string, int) Attr, n int) {
	Any("dot import with shadowed replacement", n) // want `use slog.Int instead of slog.Any`
	_ = Int
}
