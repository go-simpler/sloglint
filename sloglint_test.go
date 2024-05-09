package sloglint_test

import (
	"testing"

	"go-simpler.org/sloglint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	t.Run("no mixed arguments", func(t *testing.T) {
		analyzer := sloglint.New(nil)
		analysistest.Run(t, testdata, analyzer, "no_mixed_args")
	})

	t.Run("key-value pairs only", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{KVOnly: true})
		analysistest.Run(t, testdata, analyzer, "kv_only")
	})

	t.Run("attributes only", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{AttrOnly: true})
		analysistest.Run(t, testdata, analyzer, "attr_only")
	})

	t.Run("no global (all)", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{NoGlobal: "all"})
		analysistest.Run(t, testdata, analyzer, "no_global_all")
	})

	t.Run("no global (default)", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{NoGlobal: "default"})
		analysistest.Run(t, testdata, analyzer, "no_global_default")
	})

	t.Run("context only (all)", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{ContextOnly: "all"})
		analysistest.Run(t, testdata, analyzer, "context_only_all")
	})

	t.Run("context only (scope)", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{ContextOnly: "scope"})
		analysistest.Run(t, testdata, analyzer, "context_only_scope")
	})

	t.Run("static message", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{StaticMsg: true})
		analysistest.Run(t, testdata, analyzer, "static_msg")
	})

	t.Run("no raw keys", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{NoRawKeys: true})
		analysistest.Run(t, testdata, analyzer, "no_raw_keys")
	})

	t.Run("key naming case", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{KeyNamingCase: "snake"})
		analysistest.Run(t, testdata, analyzer, "key_naming_case")
	})

	t.Run("arguments on separate lines", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{ArgsOnSepLines: true})
		analysistest.Run(t, testdata, analyzer, "args_on_sep_lines")
	})

	t.Run("forbidden keys", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{ForbiddenKeys: []string{"foo_bar"}})
		analysistest.Run(t, testdata, analyzer, "forbidden_keys")
	})
}
