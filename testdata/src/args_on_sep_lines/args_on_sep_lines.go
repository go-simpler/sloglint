package args_on_sep_lines

import "log/slog"

func _() {
	slog.Info("msg", "foo", 1)
	slog.Info("msg", slog.Int("foo", 1))

	slog.Info("msg",
		"foo", 1,
		"bar", 2,
	)
	slog.Info("msg",
		slog.Int("foo", 1),
		slog.Int("bar", 2),
	)

	slog.Info("msg", "foo", 1, "bar", 2)                     // want `arguments should be put on separate lines`
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2)) // want `arguments should be put on separate lines`

	// Arguments are checked in source order, even if key-value pairs and
	// attributes are mixed within a single call.
	slog.Info("msg",
		"foo", 1, slog.Int("bar", 2), // want `arguments should be put on separate lines`
		"baz", 3,
	)
	slog.Info("msg",
		slog.Int("foo", 1), "bar", 2, // want `arguments should be put on separate lines`
		"baz", 3,
	)
	slog.Info("msg",
		"foo", 1,
		slog.Int("bar", 2),
		"baz", 3,
	)
}
