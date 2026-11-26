package sloglint

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func preferTypedAttribute(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) != 2 {
		return
	}

	typ := pass.TypesInfo.TypeOf(call.Args[1])
	if typ == nil {
		return
	}

	var replacement string
	switch typ := types.Unalias(types.Default(typ)).(type) {
	case *types.Basic:
		switch typ.Kind() {
		case types.String:
			replacement = "String"
		case types.Int:
			replacement = "Int"
		case types.Int64:
			replacement = "Int64"
		case types.Uint64:
			replacement = "Uint64"
		case types.Float64:
			replacement = "Float64"
		case types.Bool:
			replacement = "Bool"
		}
	case *types.Named:
		obj := typ.Obj()
		if obj.Pkg() != nil && obj.Pkg().Path() == "time" {
			switch obj.Name() {
			case "Time", "Duration":
				replacement = obj.Name()
			}
		}
	}
	if replacement == "" {
		return
	}

	diagnostic := analysis.Diagnostic{
		Pos:     call.Fun.Pos(),
		End:     call.Fun.End(),
		Message: fmt.Sprintf("use slog.%s instead of slog.Any", replacement),
	}
	// A dot-imported replacement may be shadowed by another declaration.
	if sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr); ok {
		diagnostic.SuggestedFixes = []analysis.SuggestedFix{{
			Message: fmt.Sprintf("Replace slog.Any with slog.%s", replacement),
			TextEdits: []analysis.TextEdit{{
				Pos:     sel.Sel.Pos(),
				End:     sel.Sel.End(),
				NewText: []byte(replacement),
			}},
		}}
	}
	pass.Report(diagnostic)
}
