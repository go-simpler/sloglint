// Package sloglint implements the sloglint analyzer.
package sloglint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

// New creates a new sloglint analyzer.
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "sloglint",
		Doc:      "ensure consistent code style when using log/slog",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
}

// mapping: function name -> arguments position.
var funcs = map[string]int{
	"log/slog.Log":                    3,
	"log/slog.Debug":                  1,
	"log/slog.Info":                   1,
	"log/slog.Warn":                   1,
	"log/slog.Error":                  1,
	"log/slog.DebugContext":           2,
	"log/slog.InfoContext":            2,
	"log/slog.WarnContext":            2,
	"log/slog.ErrorContext":           2,
	"(*log/slog.Logger).Log":          3,
	"(*log/slog.Logger).Debug":        1,
	"(*log/slog.Logger).Info":         1,
	"(*log/slog.Logger).Warn":         1,
	"(*log/slog.Logger).Error":        1,
	"(*log/slog.Logger).DebugContext": 2,
	"(*log/slog.Logger).InfoContext":  2,
	"(*log/slog.Logger).WarnContext":  2,
	"(*log/slog.Logger).ErrorContext": 2,
}

func run(pass *analysis.Pass) (any, error) {
	visit := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	types := []ast.Node{(*ast.CallExpr)(nil)}

	visit.Preorder(types, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		callee := typeutil.StaticCallee(pass.TypesInfo, call)
		if callee == nil {
			return
		}

		argsPos, ok := funcs[callee.FullName()]
		if !ok {
			return
		}

		args := call.Args[argsPos:]
		if len(args) == 0 {
			return
		}

		var attrsCount int
		for _, arg := range args {
			if pass.TypesInfo.TypeOf(arg).String() == "log/slog.Attr" {
				attrsCount++
			}
		}

		if attrsCount != 0 && attrsCount != len(args) {
			pass.Reportf(call.Pos(), "key-value pairs and attributes should not be mixed")
		}
	})

	return nil, nil
}
