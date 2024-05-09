package args_on_sep_lines

import "log/slog"

func tests() {
	slog.Info("msg")
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

	slog.With(
		"foo", 1,
		"bar", 2,
	).Info("msg")

	slog.With(
		slog.Int("foo", 1),
		slog.Int("bar", 2),
	).Info("msg")

	slog.Info("msg", "foo", 1, "bar", 2)                     // want `arguments should be put on separate lines`
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2)) // want `arguments should be put on separate lines`

	slog.Info("msg", "foo", 1, // want `arguments should be put on separate lines`
		"bar", 2,
	)
	slog.Info("msg", // want `arguments should be put on separate lines`
		"foo", 1, "bar", 2,
	)
	slog.Info("msg", slog.Int("foo", 1), // want `arguments should be put on separate lines`
		slog.Int("bar", 2),
	)
	slog.Info("msg", // want `arguments should be put on separate lines`
		slog.Int("foo", 1), slog.Int("bar", 2),
	)

	slog.With("msg", "foo", 1, "bar", 2).Info("msg")                     // want `arguments should be put on separate lines`
	slog.With("msg", slog.Int("foo", 1), slog.Int("bar", 2)).Info("msg") // want `arguments should be put on separate lines`
}
