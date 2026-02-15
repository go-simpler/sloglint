package sloglint

import (
	"errors"
	"testing"
)

func TestOptions_validate(t *testing.T) {
	tests := map[string]struct {
		opts Options
		err  error
	}{
		"KVOnly+AttrOnly":       {Options{KVOnly: true, AttrOnly: true}, errIncompatible},
		"invalid NoGlobal":      {Options{NoGlobal: "-"}, errInvalidValue},
		"invalid ContextOnly":   {Options{ContextOnly: "-"}, errInvalidValue},
		"invalid MsgStyle":      {Options{MsgStyle: "-"}, errInvalidValue},
		"invalid KeyNamingCase": {Options{KeyNamingCase: "-"}, errInvalidValue},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if err := test.opts.validate(); !errors.Is(err, test.err) {
				t.Errorf("got: %v; want: %v", err, test.err)
			}
		})
	}
}
