package sloglint

import (
	"errors"
	"testing"
)

func TestOptions(t *testing.T) {
	tests := map[string]struct {
		opts Options
		err  error
	}{
		"KVOnly+AttrOnly: incompatible": {
			opts: Options{KVOnly: true, AttrOnly: true},
			err:  errIncompatible,
		},
		"NoGlobal: invalid value": {
			opts: Options{NoGlobal: "-"},
			err:  errInvalidValue,
		},
		"ContextOnly: invalid value": {
			opts: Options{ContextOnly: "-"},
			err:  errInvalidValue,
		},
		"KeyNamingCase: invalid value": {
			opts: Options{KeyNamingCase: "-"},
			err:  errInvalidValue,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			analyzer := New(&test.opts)
			if _, err := analyzer.Run(nil); !errors.Is(err, test.err) {
				t.Errorf("errors.Is() mismatch\ngot:  %v\nwant: %v", err, test.err)
			}
		})
	}
}
