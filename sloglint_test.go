package sloglint_test

import (
	"testing"

	"go.tmz.dev/sloglint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := sloglint.New()
	analysistest.Run(t, testdata, analyzer, "tests")
}
