package sloglint

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

func noGlobal(pass *analysis.Pass, fn *types.Func, call *ast.CallExpr, defaultOnly bool) {
	switch fn.Name() {
	case "Log", "LogAttrs",
		"Debug", "Info", "Warn", "Error",
		"DebugContext", "InfoContext", "WarnContext", "ErrorContext",
		"With":
	default:
		return
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	if ident.Name == "slog" {
		pass.ReportRangef(sel.X, "default logger should not be used")
		return
	}

	if defaultOnly {
		return
	}

	if obj := pass.TypesInfo.ObjectOf(ident); obj != nil && obj.Parent() == obj.Pkg().Scope() {
		pass.ReportRangef(sel.X, "global logger should not be used")
	}
}

func contextOnly(pass *analysis.Pass, fn *types.Func, call *ast.CallExpr, cursor inspector.Cursor, scopeOnly bool) {
	switch fn.Name() {
	case "Debug", "Info", "Warn", "Error":
	default:
		return
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     sel.Sel.Pos(),
		End:     sel.Sel.End(),
		Message: fmt.Sprintf("%sContext should be used instead", fn.Name()),
		SuggestedFixes: []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{{
				Pos:     sel.Sel.Pos(),
				End:     sel.Sel.End(),
				NewText: fmt.Appendf(nil, "%sContext", fn.Name()),
			}},
		}},
	}

	if !scopeOnly {
		pass.Report(diag)
		return
	}

	for cursor := range cursor.Enclosing(new(ast.FuncDecl), new(ast.FuncLit)) {
		var params []*ast.Field

		switch fn := cursor.Node().(type) {
		case *ast.FuncDecl:
			params = fn.Type.Params.List
		case *ast.FuncLit:
			params = fn.Type.Params.List
		}

		if len(params) == 0 || len(params[0].Names) == 0 {
			continue
		}

		name := typeName(pass.TypesInfo, params[0].Names[0])
		if name == "context.Context" || name == "*net/http.Request" {
			pass.Report(diag)
			return
		}
	}
}

func discardHandler(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	sel, ok := call.Args[0].(*ast.SelectorExpr)
	if !ok {
		return
	}

	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil {
		return
	}

	if obj.Pkg().Name() != "io" || obj.Name() != "Discard" {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:     call.Pos(),
		End:     call.Pos(),
		Message: "use slog.DiscardHandler instead",
		SuggestedFixes: []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{{
				Pos:     call.Pos(),
				End:     call.End(),
				NewText: []byte("slog.DiscardHandler"),
			}},
		}},
	})
}
