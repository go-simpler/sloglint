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
		"invalid NoGlobalLogger":           {Options{NoGlobalLogger: "-"}, errInvalidValue},
		"invalid ContextOnly":              {Options{ContextOnly: "-"}, errInvalidValue},
		"invalid MessageStyle":             {Options{MessageStyle: "-"}, errInvalidValue},
		"invalid KeyNamingCase":            {Options{KeyNamingCase: "-"}, errInvalidValue},
		"KeyValuePairsOnly+AttributesOnly": {Options{KeyValuePairsOnly: true, AttributesOnly: true}, errIncompatible},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if err := test.opts.validate(); !errors.Is(err, test.err) {
				t.Errorf("got: %v; want: %v", err, test.err)
			}
		})
	}
}

func TestPreferTypedAttributesFlag(t *testing.T) {
	opts := &Options{}
	fs := flags(opts)
	if err := fs.Parse([]string{"-prefer-typed-attrs"}); err != nil {
		t.Fatal(err)
	}
	if !opts.PreferTypedAttributes {
		t.Error("the -prefer-typed-attrs flag did not enable PreferTypedAttributes")
	}
}
