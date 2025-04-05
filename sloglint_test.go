package sloglint

import (
	"errors"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := map[string]struct {
		opts Options
		dir  string
	}{
		"no mixed arguments":          {Options{NoMixedArgs: true}, "no_mixed_args"},
		"key-value pairs only":        {Options{KVOnly: true}, "kv_only"},
		"attributes only":             {Options{AttrOnly: true}, "attr_only"},
		"no global (all)":             {Options{NoGlobal: "all"}, "no_global_all"},
		"no global (default)":         {Options{NoGlobal: "default"}, "no_global_default"},
		"context only (all)":          {Options{ContextOnly: "all"}, "context_only_all"},
		"context only (scope)":        {Options{ContextOnly: "scope"}, "context_only_scope"},
		"static message":              {Options{StaticMsg: true}, "static_msg"},
		"no raw keys":                 {Options{NoRawKeys: true}, "no_raw_keys"},
		"key naming case":             {Options{KeyNamingCase: "snake"}, "key_naming_case"},
		"arguments on separate lines": {Options{ArgsOnSepLines: true}, "args_on_sep_lines"},
		"forbidden keys":              {Options{ForbiddenKeys: []string{"foo_bar"}}, "forbidden_keys"},
		"message style (lowercased)":  {Options{MsgStyle: "lowercased"}, "msg_style_lowercased"},
		"message style (capitalized)": {Options{MsgStyle: "capitalized"}, "msg_style_capitalized"},
		"slog.DiscardHandler":         {Options{go124: true}, "discard_handler"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			analyzer := New(&tt.opts)
			testdata := analysistest.TestData()
			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, tt.dir)
		})
	}
}

func TestOptions(t *testing.T) {
	tests := map[string]struct {
		opts Options
		err  error
	}{
		"KVOnly+AttrOnly: incompatible": {Options{KVOnly: true, AttrOnly: true}, errIncompatible},
		"NoGlobal: invalid value":       {Options{NoGlobal: "-"}, errInvalidValue},
		"ContextOnly: invalid value":    {Options{ContextOnly: "-"}, errInvalidValue},
		"MsgStyle: invalid value":       {Options{MsgStyle: "-"}, errInvalidValue},
		"KeyNamingCase: invalid value":  {Options{KeyNamingCase: "-"}, errInvalidValue},
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
