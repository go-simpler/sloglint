package sloglint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := map[string]struct {
		dir  string
		opts Options
	}{
		"no mixed arguments":          {dir: "no_mixed_args", opts: Options{NoMixedArgs: true}},
		"key-value pairs only":        {dir: "kv_only", opts: Options{KVOnly: true}},
		"attributes only":             {dir: "attr_only", opts: Options{AttrOnly: true}},
		"no global (all)":             {dir: "no_global_all", opts: Options{NoGlobal: noGlobalAll}},
		"no global (default)":         {dir: "no_global_default", opts: Options{NoGlobal: noGlobalDefault}},
		"context only (all)":          {dir: "context_only_all", opts: Options{ContextOnly: contextOnlyAll}},
		"context only (scope)":        {dir: "context_only_scope", opts: Options{ContextOnly: contextOnlyScope}},
		"static message":              {dir: "static_msg", opts: Options{StaticMsg: true}},
		"message style (lowercased)":  {dir: "msg_style_lowercased", opts: Options{MsgStyle: msgStyleLowercased}},
		"message style (capitalized)": {dir: "msg_style_capitalized", opts: Options{MsgStyle: msgStyleCapitalized}},
		"no raw keys":                 {dir: "no_raw_keys", opts: Options{NoRawKeys: true}},
		"key naming case":             {dir: "key_naming_case", opts: Options{KeyNamingCase: keyNamingCaseSnake}},
		"allowed keys":                {dir: "allowed_keys", opts: Options{AllowedKeys: []string{"foo"}}},
		"forbidden keys":              {dir: "forbidden_keys", opts: Options{ForbiddenKeys: []string{"bar"}}},
		"arguments on separate lines": {dir: "args_on_sep_lines", opts: Options{ArgsOnSepLines: true}},
		"slog.DiscardHandler":         {dir: "discard_handler", opts: Options{}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			analyzer := New(&test.opts)
			testdata := analysistest.TestData()
			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, test.dir)
		})
	}
}
