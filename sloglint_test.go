package sloglint_test

import (
	"testing"

	"go-simpler.org/sloglint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	t.Run("mixed arguments", func(t *testing.T) {
		analyzer := sloglint.New(nil)
		analysistest.Run(t, testdata, analyzer, "mixed_args")
	})

	t.Run("key-value pairs only", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{KVOnly: true})
		analysistest.Run(t, testdata, analyzer, "kv_only")
	})

	t.Run("attributes only", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{AttrOnly: true})
		analysistest.Run(t, testdata, analyzer, "attr_only")
	})

	t.Run("no raw keys", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{NoRawKeys: true})
		analysistest.Run(t, testdata, analyzer, "no_raw_keys")
	})

	t.Run("arguments on separate lines", func(t *testing.T) {
		analyzer := sloglint.New(&sloglint.Options{ArgsOnSepLines: true})
		analysistest.Run(t, testdata, analyzer, "args_on_sep_lines")
	})
}
