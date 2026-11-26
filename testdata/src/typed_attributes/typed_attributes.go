package typed_attributes

import (
	"log/slog"
	"time"
)

type integerAlias = int
type durationAlias = time.Duration
type timeAlias = time.Time

type namedInt int
type namedString string

func (n namedInt) LogValue() slog.Value { return slog.StringValue("custom") }

type otherLogger struct{}

func (otherLogger) Any(_ string, _ any) slog.Attr { return slog.Attr{} }

func _(
	s string, i int, i64 int64, u64 uint64, f64 float64, b bool,
	d time.Duration, tm time.Time, ai integerAlias, ad durationAlias, at timeAlias,
	unknown any, n namedInt,
) {
	slog.Any("string", s)                   // want `use slog.String instead of slog.Any`
	slog.Any("int", i)                      // want `use slog.Int instead of slog.Any`
	slog.Any("int64", i64)                  // want `use slog.Int64 instead of slog.Any`
	slog.Any("uint64", u64)                 // want `use slog.Uint64 instead of slog.Any`
	slog.Any("float64", f64)                // want `use slog.Float64 instead of slog.Any`
	slog.Any("bool", b)                     // want `use slog.Bool instead of slog.Any`
	slog.Any("duration", d)                 // want `use slog.Duration instead of slog.Any`
	slog.Any("time", tm)                    // want `use slog.Time instead of slog.Any`
	slog.Any("int alias", ai)               // want `use slog.Int instead of slog.Any`
	slog.Any("duration alias", ad)          // want `use slog.Duration instead of slog.Any`
	slog.Any("time alias", at)              // want `use slog.Time instead of slog.Any`
	slog.Any("literal string", "value")     // want `use slog.String instead of slog.Any`
	slog.Any("literal int", 42)             // want `use slog.Int instead of slog.Any`
	slog.Any("literal float", 1.5)          // want `use slog.Float64 instead of slog.Any`
	slog.Any("literal bool", true)          // want `use slog.Bool instead of slog.Any`
	slog.Info("nested", slog.Any("key", i)) // want `use slog.Int instead of slog.Any`

	// These values do not have an identical type-specific replacement.
	slog.Any("dynamic", unknown)
	slog.Any("named LogValuer", n)
	slog.Any("named string", namedString(s))
	slog.Any("nil", nil)
	slog.Any("pointer", &i)
	slog.Any("int32", int32(1))
	slog.Any("rune literal", 'a')
	slog.Any("float32", float32(1))
	slog.Any("uint", uint(1))
	slog.Any("value", slog.IntValue(1))
	slog.Any("attrs", []slog.Attr{slog.Int("key", 1)})
	otherLogger{}.Any("custom function", 1)
	anyFunc := slog.Any
	anyFunc("function value", 1)
}

func _[T ~int](value T) {
	slog.Any("type parameter", value)
}
