package sloglint

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/types/typeutil"
)

func typeName(info *types.Info, expr ast.Expr) string {
	if typ := info.TypeOf(expr); typ != nil {
		return typ.String()
	}
	return ""
}

func funcName(info *types.Info, call *ast.CallExpr) string {
	if fn := typeutil.StaticCallee(info, call); fn != nil {
		return fn.FullName()
	}
	return ""
}

// keyName returns the string value of the given key.
// The key must be a string literal or a constant, including constants declared
// in another file or package.
func keyName(info *types.Info, key ast.Expr) (string, bool) {
	if sel, ok := key.(*ast.SelectorExpr); ok {
		key = sel.Sel // The key is defined in another package, e.g. pkg.ConstKey.
	}

	if lit, ok := key.(*ast.BasicLit); ok {
		if lit.Kind != token.STRING {
			return "", false
		}

		name, err := strconv.Unquote(lit.Value)
		if err != nil {
			panic("unreachable") // String literals are always quoted.
		}

		return name, true
	}

	// The key is a constant, e.g. constKey. Unlike ast.Ident.Obj, go/types
	// resolves constants declared in another file or package.
	ident, ok := key.(*ast.Ident)
	if !ok {
		return "", false
	}

	obj, ok := info.ObjectOf(ident).(*types.Const)
	if !ok {
		return "", false
	}

	value := obj.Val()
	if value.Kind() != constant.String {
		return "", false
	}

	return constant.StringVal(value), true
}

func isGroup(info *types.Info, expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	name := funcName(info, call)
	return name == "log/slog.Group" || name == "log/slog.GroupAttrs"
}
