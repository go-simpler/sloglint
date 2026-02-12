package sloglint

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"go/version"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/ettle/strcase"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

type slogFuncCall struct {
	fn     *types.Func
	expr   *ast.CallExpr
	cursor inspector.Cursor
	msg    ast.Expr // Optional, nil if does not exist.
	keys   []ast.Expr
	attrs  []ast.Expr
}

var checks = []func(*analysis.Pass, *Options, *slogFuncCall){
	noMixedArgs,
	kvOnly,
	attrOnly,
	noGlobal,
	contextOnly,
	staticMsg,
	msgStyle,
	noRawKeys,
	keyNamingCase,
	allowedKeys,
	forbiddenKeys,
	argsOnSepLines,
	discardHandler,
}

func noMixedArgs(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if opts.NoMixedArgs && len(call.keys) > 0 && len(call.attrs) > 0 {
		pass.Reportf(call.expr.Pos(), "key-value pairs and attributes should not be mixed")
	}
}

func kvOnly(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if opts.KVOnly && len(call.attrs) > 0 {
		pass.Reportf(call.expr.Pos(), "attributes should not be used")
	}
}

func attrOnly(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if opts.AttrOnly && len(call.keys) > 0 {
		pass.Reportf(call.expr.Pos(), "key-value pairs should not be used")
	}
}

func noGlobal(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	switch call.fn.Name() {
	case "Log", "LogAttrs",
		"Debug", "Info", "Warn", "Error",
		"DebugContext", "InfoContext", "WarnContext", "ErrorContext",
		"With":
	default:
		return
	}

	switch opts.NoGlobal {
	case "":
		return
	case noGlobalAll:
		if strings.HasPrefix(call.fn.FullName(), "log/slog.") {
			pass.Reportf(call.expr.Pos(), "global logger should not be used")
			return
		}
	case noGlobalDefault:
		if strings.HasPrefix(call.fn.FullName(), "log/slog.") {
			pass.Reportf(call.expr.Pos(), "default logger should not be used")
		}
		return
	}

	sel, ok := call.expr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	obj := pass.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return
	}

	if obj.Parent() == obj.Pkg().Scope() {
		pass.Reportf(call.expr.Pos(), "global logger should not be used")
	}
}

func contextOnly(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	switch call.fn.Name() {
	case "Debug", "Info", "Warn", "Error":
	default:
		return
	}

	switch opts.ContextOnly {
	case "":
		return
	case contextOnlyAll:
		pass.Reportf(call.expr.Pos(), "%sContext should be used instead", call.fn.Name())
		return
	}

	for cursor := range call.cursor.Enclosing(new(ast.FuncDecl), new(ast.FuncLit)) {
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

		typ := pass.TypesInfo.TypeOf(params[0].Names[0])
		if typ == nil {
			continue
		}

		if typ.String() == "context.Context" || typ.String() == "*net/http.Request" {
			pass.Reportf(call.expr.Pos(), "%sContext should be used instead", call.fn.Name())
			return
		}
	}
}

func staticMsg(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if !opts.StaticMsg || call.msg == nil {
		return
	}

	var isStatic func(msg ast.Expr) bool
	isStatic = func(msg ast.Expr) bool {
		switch msg := msg.(type) {
		case *ast.BasicLit: // e.g. slog.Info("msg")
			return msg.Kind == token.STRING
		case *ast.Ident: // e.g. slog.Info(constMsg)
			_, isConst := pass.TypesInfo.ObjectOf(msg).(*types.Const)
			return isConst
		case *ast.BinaryExpr: // e.g. slog.Info("x" + "y")
			if msg.Op != token.ADD {
				panic("unreachable") // Only "+" can be applied to strings.
			}
			return isStatic(msg.X) && isStatic(msg.Y)
		default:
			return false
		}
	}

	if !isStatic(call.msg) {
		pass.Reportf(call.msg.Pos(), "message should be a string literal or a constant")
	}
}

func msgStyle(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if opts.MsgStyle == "" || call.msg == nil {
		return
	}

	lit, ok := call.msg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		panic("unreachable") // String literals are always quoted.
	}

	runes := []rune(strings.TrimSpace(msg))
	if len(runes) < 2 {
		return
	}

	first, second := runes[0], runes[1]

	switch opts.MsgStyle {
	case msgStyleLowercased:
		if unicode.IsLower(first) {
			return
		}
		if unicode.IsPunct(second) {
			return // e.g. "U.S."
		}
		if unicode.IsUpper(second) {
			return // e.g. "HTTP"
		}
	case msgStyleCapitalized:
		if unicode.IsUpper(first) {
			return
		}
		if unicode.IsUpper(second) {
			return // e.g. "iPhone"
		}
	}

	pass.Reportf(call.msg.Pos(), "message should be %s", opts.MsgStyle)
}

func noRawKeys(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if !opts.NoRawKeys {
		return
	}

	isConst := func(key ast.Expr) bool {
		ident, ok := key.(*ast.Ident)
		if !ok {
			return false
		}
		_, ok = pass.TypesInfo.ObjectOf(ident).(*types.Const)
		return ok
	}

	for key := range allKeys(pass.TypesInfo, call.keys, call.attrs) {
		if sel, ok := key.(*ast.SelectorExpr); ok {
			key = sel.Sel // The key is defined in another package, e.g. pkg.ConstKey.
		}
		if !isConst(key) {
			pass.Reportf(key.Pos(), "raw keys should not be used")
		}
	}
}

func keyNamingCase(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	var caseFn func(string) string
	var caseName string

	switch opts.KeyNamingCase {
	case "":
		return
	case keyNamingCaseSnake:
		caseFn, caseName = strcase.ToSnake, "snake_case"
	case keyNamingCaseKebab:
		caseFn, caseName = strcase.ToKebab, "kebab-case"
	case keyNamingCaseCamel:
		caseFn, caseName = strcase.ToCamel, "camelCase"
	case keyNamingCasePascal:
		caseFn, caseName = strcase.ToPascal, "PascalCase"
	}

	for key := range allKeys(pass.TypesInfo, call.keys, call.attrs) {
		name, ok := keyName(key)
		if !ok || name == caseFn(name) {
			return
		}
		pass.Report(analysis.Diagnostic{
			Pos:     key.Pos(),
			Message: fmt.Sprintf("keys should be written in %s", caseName),
			SuggestedFixes: []analysis.SuggestedFix{{
				TextEdits: []analysis.TextEdit{{
					Pos:     key.Pos(),
					End:     key.End(),
					NewText: []byte(strconv.Quote(caseFn(name))),
				}},
			}},
		})
	}
}

func allowedKeys(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if len(opts.AllowedKeys) == 0 {
		return
	}

	for key := range allKeys(pass.TypesInfo, call.keys, call.attrs) {
		if name, ok := keyName(key); ok && !slices.Contains(opts.AllowedKeys, name) {
			pass.Reportf(key.Pos(), "%q key is not allowed and should not be used", name)
		}
	}
}

func forbiddenKeys(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if len(opts.ForbiddenKeys) == 0 {
		return
	}

	for key := range allKeys(pass.TypesInfo, call.keys, call.attrs) {
		if name, ok := keyName(key); ok && slices.Contains(opts.ForbiddenKeys, name) {
			pass.Reportf(key.Pos(), "%q key is forbidden and should not be used", name)
		}
	}
}

func argsOnSepLines(pass *analysis.Pass, opts *Options, call *slogFuncCall) {
	if !opts.ArgsOnSepLines {
		return
	}

	if len(call.keys)+len(call.attrs) <= 1 {
		return // Special case: slog.Info("msg", "key", "value") is fine.
	}

	prevLine := pass.Fset.Position(call.expr.Pos()).Line

	for _, arg := range slices.Concat(call.keys, call.attrs) {
		currLine := pass.Fset.Position(arg.Pos()).Line
		if currLine == prevLine {
			pass.Reportf(call.expr.Pos(), "arguments should be put on separate lines")
			return
		}
		prevLine = currLine
	}
}

func discardHandler(pass *analysis.Pass, _ *Options, call *slogFuncCall) {
	switch call.fn.Name() {
	case "NewTextHandler", "NewJSONHandler":
	default:
		return
	}

	switch v := pass.Module.GoVersion; {
	case v == "": // Empty in test runs.
	case version.Compare("go"+v, "go1.24") >= 0:
	default:
		return
	}

	sel, ok := call.expr.Args[0].(*ast.SelectorExpr)
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
		Pos:     call.expr.Pos(),
		Message: "use slog.DiscardHandler instead",
		SuggestedFixes: []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{{
				Pos:     call.expr.Pos(),
				End:     call.expr.End(),
				NewText: []byte("slog.DiscardHandler"),
			}},
		}},
	})
}
