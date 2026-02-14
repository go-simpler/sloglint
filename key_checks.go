package sloglint

import (
	"fmt"
	"go/ast"
	"go/types"
	"slices"
	"strconv"

	"github.com/ettle/strcase"
	"golang.org/x/tools/go/analysis"
)

func noRawKeys(pass *analysis.Pass, key ast.Expr) {
	if sel, ok := key.(*ast.SelectorExpr); ok {
		key = sel.Sel // The key is defined in another package, e.g. pkg.ConstKey.
	}
	if ident, ok := key.(*ast.Ident); ok {
		if _, ok := pass.TypesInfo.ObjectOf(ident).(*types.Const); ok {
			return
		}
	}
	pass.ReportRangef(key, "raw keys should not be used")
}

func keyNamingCase(pass *analysis.Pass, key ast.Expr, caseName string) {
	name, ok := keyName(key)
	if !ok {
		return
	}

	var newName string
	switch caseName {
	case keyNamingCaseSnake:
		newName = strcase.ToSnake(name)
	case keyNamingCaseKebab:
		newName = strcase.ToKebab(name)
	case keyNamingCaseCamel:
		newName = strcase.ToCamel(name)
	case keyNamingCasePascal:
		newName = strcase.ToPascal(name)
	}

	if name == newName {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:     key.Pos(),
		End:     key.End(),
		Message: fmt.Sprintf("keys should be written in %s case", caseName),
		SuggestedFixes: []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{{
				Pos:     key.Pos(),
				End:     key.End(),
				NewText: strconv.AppendQuote(nil, newName),
			}},
		}},
	})
}

func allowedKeys(pass *analysis.Pass, key ast.Expr, allowed []string) {
	if name, ok := keyName(key); ok && !slices.Contains(allowed, name) {
		pass.ReportRangef(key, "%q key is not allowed and should not be used", name)
	}
}

func forbiddenKeys(pass *analysis.Pass, key ast.Expr, forbidden []string) {
	if name, ok := keyName(key); ok && slices.Contains(forbidden, name) {
		pass.ReportRangef(key, "%q key is forbidden and should not be used", name)
	}
}
