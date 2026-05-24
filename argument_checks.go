package sloglint

import (
	"go/ast"
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

func noMixedArguments(pass *analysis.Pass, keys, attrs []ast.Expr) {
	if len(keys) == 0 {
		return
	}
	for _, attr := range attrs {
		if isGroup(pass.TypesInfo, attr) {
			continue // Special case: slog.Group/GroupAttrs should always be allowed.
		}
		pass.ReportRangef(attr, "key-value pairs and attributes should not be mixed")
		return
	}
}

func keyValuePairsOnly(pass *analysis.Pass, call *ast.CallExpr, attrs []ast.Expr) {
	fnName := typeutil.StaticCallee(pass.TypesInfo, call).FullName()

	if replacement, ok := map[string]string{
		"log/slog.GroupAttrs":         "slog.Group",
		"log/slog.LogAttrs":           "slog.Log",
		"(*log/slog.Logger).LogAttrs": "slog.Logger.Log",
	}[fnName]; ok {
		pass.ReportRangef(call, "use %s with key-value pairs instead", replacement)
		return
	}

	for _, attr := range attrs {
		if isGroup(pass.TypesInfo, attr) {
			continue // Special case: slog.Group should always be allowed.
		}
		pass.ReportRangef(attr, "attributes should not be used")
		return
	}
}

func attributesOnly(pass *analysis.Pass, call *ast.CallExpr, keys []ast.Expr) {
	fnName := typeutil.StaticCallee(pass.TypesInfo, call).FullName()

	if replacement, ok := map[string]string{
		"log/slog.Group":         "slog.GroupAttrs",
		"log/slog.Log":           "slog.LogAttrs",
		"(*log/slog.Logger).Log": "slog.Logger.LogAttrs",
	}[fnName]; ok {
		pass.ReportRangef(call, "use %s with attributes instead", replacement)
		return
	}

	for _, key := range keys {
		pass.ReportRangef(key, "key-value pairs should not be used")
		return
	}
}

func argumentsOnSeparateLines(pass *analysis.Pass, keys, attrs []ast.Expr) {
	args := slices.Concat(keys, attrs)
	if len(args) <= 1 {
		return // Special case: slog.Info("msg", "key", "value") is fine.
	}

	prevLine := pass.Fset.Position(args[0].Pos()).Line
	for _, arg := range args[1:] {
		currLine := pass.Fset.Position(arg.Pos()).Line
		if currLine == prevLine {
			pass.Reportf(arg.Pos(), "arguments should be put on separate lines")
			return
		}
		prevLine = currLine
	}
}
