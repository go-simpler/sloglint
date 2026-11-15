package sloglint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	custom := []Func{
		{FullName: "no_mixed_args.customLog", MessagePos: 0, ArgumentsPos: 1},
	}

	tests := map[string]struct {
		dir  string
		opts Options
	}{
		"no global logger (all)":      {dir: "no_global_all", opts: Options{NoGlobalLogger: noGlobalLoggerAll}},
		"no global logger (default)":  {dir: "no_global_default", opts: Options{NoGlobalLogger: noGlobalLoggerDefault}},
		"context only (all)":          {dir: "context_only_all", opts: Options{ContextOnly: contextOnlyAll}},
		"context only (scope)":        {dir: "context_only_scope", opts: Options{ContextOnly: contextOnlyScope}},
		"discard handler":             {dir: "discard_handler", opts: Options{}},
		"static message":              {dir: "static_msg", opts: Options{StaticMessage: true}},
		"message style (lowercased)":  {dir: "msg_style_lowercased", opts: Options{MessageStyle: messageStyleLowercased}},
		"message style (capitalized)": {dir: "msg_style_capitalized", opts: Options{MessageStyle: messageStyleCapitalized}},
		"no mixed arguments":          {dir: "no_mixed_args", opts: Options{NoMixedArguments: true, CustomFuncs: custom}},
		"key-value pairs only":        {dir: "kv_only", opts: Options{KeyValuePairsOnly: true}},
		"attributes only":             {dir: "attr_only", opts: Options{AttributesOnly: true}},
		"arguments on separate lines": {dir: "args_on_sep_lines", opts: Options{ArgumentsOnSeparateLines: true}},
		"constant keys":               {dir: "no_raw_keys", opts: Options{ConstantKeys: true}},
		"allowed keys":                {dir: "allowed_keys", opts: Options{AllowedKeys: []string{"foo"}}},
		"forbidden keys":              {dir: "forbidden_keys", opts: Options{ForbiddenKeys: []string{"bar"}}},
		"key naming case":             {dir: "key_naming_case", opts: Options{KeyNamingCase: keyNamingCaseSnake}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			analyzer := New(&test.opts)
			testdata := analysistest.TestData()
			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, test.dir)
		})
	}
}
