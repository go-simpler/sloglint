package typed_attributes_disabled

import "log/slog"

func _(n int) {
	slog.Any("number", n)
}
