// Package sloglint implements the sloglint analyzer.
package sloglint

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"go/version"
	"iter"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/ettle/strcase"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

// New creates a new sloglint analyzer.
func New(opts *Options) *analysis.Analyzer {
	if opts == nil {
		opts = &Options{NoMixedArgs: true}
	}

	return &analysis.Analyzer{
		Name:     "sloglint",
		Doc:      "Ensures consistent code style when using log/slog.",
		URL:      "https://github.com/go-simpler/sloglint",
		Flags:    flags(opts),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			if err := opts.validate(); err != nil {
				return nil, err
			}

			root := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Root()
			for cursor := range root.Preorder(new(ast.CallExpr)) {
				analyze(pass, opts, cursor)
			}

			return nil, nil
		},
	}
}

var slogFuncs = map[string]struct {
	argsPos   int  // The position of key-value arguments in the function signature, starting from 0.
	hasMsgArg bool // Whether the function has the "msg" argument that should be analyzed.
	hasCtxAlt bool // Whether an alternative function that accepts a context as the first argument exists.
}{
	"log/slog.With":                   {argsPos: 0},
	"log/slog.Log":                    {argsPos: 3, hasMsgArg: true},
	"log/slog.LogAttrs":               {argsPos: 3, hasMsgArg: true},
	"log/slog.Debug":                  {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"log/slog.Info":                   {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"log/slog.Warn":                   {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"log/slog.Error":                  {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"log/slog.DebugContext":           {argsPos: 2, hasMsgArg: true},
	"log/slog.InfoContext":            {argsPos: 2, hasMsgArg: true},
	"log/slog.WarnContext":            {argsPos: 2, hasMsgArg: true},
	"log/slog.ErrorContext":           {argsPos: 2, hasMsgArg: true},
	"(*log/slog.Logger).With":         {argsPos: 0},
	"(*log/slog.Logger).Log":          {argsPos: 3, hasMsgArg: true},
	"(*log/slog.Logger).LogAttrs":     {argsPos: 3, hasMsgArg: true},
	"(*log/slog.Logger).Debug":        {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"(*log/slog.Logger).Info":         {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"(*log/slog.Logger).Warn":         {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"(*log/slog.Logger).Error":        {argsPos: 1, hasMsgArg: true, hasCtxAlt: true},
	"(*log/slog.Logger).DebugContext": {argsPos: 2, hasMsgArg: true},
	"(*log/slog.Logger).InfoContext":  {argsPos: 2, hasMsgArg: true},
	"(*log/slog.Logger).WarnContext":  {argsPos: 2, hasMsgArg: true},
	"(*log/slog.Logger).ErrorContext": {argsPos: 2, hasMsgArg: true},
}

var attrFuncs = map[string]struct{}{
	"log/slog.String":   {},
	"log/slog.Int64":    {},
	"log/slog.Int":      {},
	"log/slog.Uint64":   {},
	"log/slog.Float64":  {},
	"log/slog.Bool":     {},
	"log/slog.Time":     {},
	"log/slog.Duration": {},
	"log/slog.Group":    {},
	"log/slog.Any":      {},
}

func analyze(pass *analysis.Pass, opts *Options, cursor inspector.Cursor) {
	call := cursor.Node().(*ast.CallExpr)

	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil {
		return
	}

	name := fn.FullName()

	if name == "log/slog.NewTextHandler" || name == "log/slog.NewJSONHandler" {
		checkDiscardHandler(pass, call)
	}

	info, ok := slogFuncs[name]
	if !ok {
		return
	}

	switch opts.NoGlobal {
	case noGlobalAll:
		if strings.HasPrefix(name, "log/slog.") || isGlobalLoggerUsed(pass.TypesInfo, call.Fun) {
			pass.Reportf(call.Pos(), "global logger should not be used")
		}
	case noGlobalDefault:
		if strings.HasPrefix(name, "log/slog.") {
			pass.Reportf(call.Pos(), "default logger should not be used")
		}
	}

	if info.hasCtxAlt {
		switch opts.ContextOnly {
		case contextOnlyAll:
			typ := pass.TypesInfo.TypeOf(call.Args[0])
			if typ != nil && typ.String() != "context.Context" {
				pass.Reportf(call.Pos(), "%sContext should be used instead", fn.Name())
			}
		case contextOnlyScope:
			if isContextInScope(pass.TypesInfo, cursor) {
				typ := pass.TypesInfo.TypeOf(call.Args[0])
				if typ != nil && typ.String() != "context.Context" {
					pass.Reportf(call.Pos(), "%sContext should be used instead", fn.Name())
				}
			}
		}
	}

	if info.hasMsgArg {
		msgPos := info.argsPos - 1

		if opts.StaticMsg && msgPos >= 0 && !isStaticMsg(pass.TypesInfo, call.Args[msgPos]) {
			pass.Reportf(call.Args[msgPos].Pos(), "message should be a string literal or a constant")
		}

		if opts.MsgStyle != "" && msgPos >= 0 {
			if lit, ok := call.Args[msgPos].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				value, err := strconv.Unquote(lit.Value)
				if err != nil {
					panic("unreachable") // String literals are always quoted.
				}
				if ok := isValidMsgStyle(value, opts.MsgStyle); !ok {
					pass.Reportf(call.Args[msgPos].Pos(), "message should be %s", opts.MsgStyle)
				}
			}
		}

	}

	if len(call.Args) < info.argsPos {
		return
	}

	args := call.Args[info.argsPos:]
	if len(args) == 0 {
		return
	}

	var keys []ast.Expr
	var attrs []ast.Expr

	for i := 0; i < len(args); i++ {
		typ := pass.TypesInfo.TypeOf(args[i])
		if typ == nil {
			continue
		}

		switch typ.String() {
		case "string":
			keys = append(keys, args[i])
			i++ // Skip the value.
		case "log/slog.Attr":
			if call, ok := args[i].(*ast.CallExpr); ok {
				if fn := typeutil.StaticCallee(pass.TypesInfo, call); fn != nil && fn.FullName() == "log/slog.Group" {
					continue // Skip slog.Group() calls.
				}
			}
			attrs = append(attrs, args[i])
		case "[]any", "[]log/slog.Attr":
			continue // The last argument may be an unpacked slice, skip it.
		}
	}

	switch {
	case opts.KVOnly && len(attrs) > 0:
		pass.Reportf(call.Pos(), "attributes should not be used")
	case opts.AttrOnly && len(keys) > 0:
		pass.Reportf(call.Pos(), "key-value pairs should not be used")
	case opts.NoMixedArgs && len(attrs) > 0 && len(keys) > 0:
		pass.Reportf(call.Pos(), "key-value pairs and attributes should not be mixed")
	}

	if opts.NoRawKeys {
		for key := range allKeys(pass.TypesInfo, keys, attrs) {
			if sel, ok := key.(*ast.SelectorExpr); ok {
				key = sel.Sel // The key is defined in another package, e.g. pkg.ConstKey.
			}

			isConst := false

			if ident, ok := key.(*ast.Ident); ok {
				if obj := pass.TypesInfo.ObjectOf(ident); obj != nil {
					if _, ok := obj.(*types.Const); ok {
						isConst = true
					}
				}
			}

			if !isConst {
				pass.Reportf(key.Pos(), "raw keys should not be used")
			}
		}
	}

	checkKeysNaming(opts, pass, keys, attrs)

	if len(opts.AllowedKeys) > 0 {
		for key := range allKeys(pass.TypesInfo, keys, attrs) {
			if name, ok := getKeyName(key); ok && !slices.Contains(opts.AllowedKeys, name) {
				pass.Reportf(key.Pos(), "%q key is not in the allowed keys list and should not be used", name)
			}
		}
	}

	if len(opts.ForbiddenKeys) > 0 {
		for key := range allKeys(pass.TypesInfo, keys, attrs) {
			if name, ok := getKeyName(key); ok && slices.Contains(opts.ForbiddenKeys, name) {
				pass.Reportf(key.Pos(), "%q key is forbidden and should not be used", name)
			}
		}
	}

	if opts.ArgsOnSepLines && areArgsOnSameLine(pass.Fset, call, keys, attrs) {
		pass.Reportf(call.Pos(), "arguments should be put on separate lines")
	}
}

func checkKeysNaming(opts *Options, pass *analysis.Pass, keys, attrs []ast.Expr) {
	checkKeyNamingCase := func(caseFn func(string) string, caseName string) {
		for key := range allKeys(pass.TypesInfo, keys, attrs) {
			name, ok := getKeyName(key)
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

	switch opts.KeyNamingCase {
	case keyNamingCaseSnake:
		checkKeyNamingCase(strcase.ToSnake, "snake_case")
	case keyNamingCaseKebab:
		checkKeyNamingCase(strcase.ToKebab, "kebab-case")
	case keyNamingCaseCamel:
		checkKeyNamingCase(strcase.ToCamel, "camelCase")
	case keyNamingCasePascal:
		checkKeyNamingCase(strcase.ToPascal, "PascalCase")
	}
}

func checkDiscardHandler(pass *analysis.Pass, call *ast.CallExpr) {
	switch v := pass.Module.GoVersion; {
	case v == "": // Empty in test runs.
	case version.Compare("go"+v, "go1.24") >= 0:
	default:
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

func isGlobalLoggerUsed(info *types.Info, call ast.Expr) bool {
	sel, ok := call.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := info.ObjectOf(ident)
	return obj.Parent() == obj.Pkg().Scope()
}

func isContextInScope(info *types.Info, cursor inspector.Cursor) bool {
	for c := range cursor.Enclosing(new(ast.FuncDecl), new(ast.FuncLit)) {
		var params []*ast.Field

		switch fn := c.Node().(type) {
		case *ast.FuncDecl:
			params = fn.Type.Params.List
		case *ast.FuncLit:
			params = fn.Type.Params.List
		}

		if len(params) == 0 || len(params[0].Names) == 0 {
			continue
		}

		typ := info.TypeOf(params[0].Names[0])
		if typ != nil && (typ.String() == "context.Context" || typ.String() == "*net/http.Request") {
			return true
		}
	}

	return false
}

func isStaticMsg(info *types.Info, msg ast.Expr) bool {
	switch msg := msg.(type) {
	case *ast.BasicLit: // slog.Info("msg")
		return msg.Kind == token.STRING
	case *ast.Ident: // const msg = "msg"; slog.Info(msg)
		_, isConst := info.ObjectOf(msg).(*types.Const)
		return isConst
	case *ast.BinaryExpr: // slog.Info("x" + "y")
		if msg.Op != token.ADD {
			panic("unreachable") // Only "+" can be applied to strings.
		}
		return isStaticMsg(info, msg.X) && isStaticMsg(info, msg.Y)
	default:
		return false
	}
}

func isValidMsgStyle(msg, style string) bool {
	runes := []rune(strings.TrimSpace(msg))
	if len(runes) < 2 {
		return true
	}

	first, second := runes[0], runes[1]

	switch style {
	case msgStyleLowercased:
		if unicode.IsLower(first) {
			return true
		}
		if unicode.IsPunct(second) {
			return true // e.g. "U.S.A.".
		}
		return unicode.IsUpper(second) // e.g. "HTTP".
	case msgStyleCapitalized:
		if unicode.IsUpper(first) {
			return true
		}
		return unicode.IsUpper(second) // e.g. "iPhone".
	default:
		panic("unreachable")
	}
}

func allKeys(info *types.Info, keys, attrs []ast.Expr) iter.Seq[ast.Expr] {
	return func(yield func(key ast.Expr) bool) {
		for _, key := range keys {
			if !yield(key) {
				return
			}
		}

		for _, attr := range attrs {
			switch attr := attr.(type) {
			case *ast.CallExpr: // slog.Int()
				callee := typeutil.StaticCallee(info, attr)
				if callee == nil {
					continue
				}
				if _, ok := attrFuncs[callee.FullName()]; !ok {
					continue
				}
				if !yield(attr.Args[0]) {
					return
				}

			case *ast.CompositeLit: // slog.Attr{}
				switch len(attr.Elts) {
				case 1: // slog.Attr{Key: ...} | slog.Attr{Value: ...}
					if kv := attr.Elts[0].(*ast.KeyValueExpr); kv.Key.(*ast.Ident).Name == "Key" {
						if !yield(kv.Value) {
							return
						}
					}

				case 2: // slog.Attr{Key: ..., Value: ...} | slog.Attr{Value: ..., Key: ...} | slog.Attr{..., ...}
					if kv, ok := attr.Elts[0].(*ast.KeyValueExpr); ok && kv.Key.(*ast.Ident).Name == "Key" {
						if !yield(kv.Value) {
							return
						}
					} else if kv, ok := attr.Elts[1].(*ast.KeyValueExpr); ok && kv.Key.(*ast.Ident).Name == "Key" {
						if !yield(kv.Value) {
							return
						}
					} else {
						if !yield(attr.Elts[0]) {
							return
						}
					}
				}
			}
		}
	}
}

func getKeyName(key ast.Expr) (string, bool) {
	if ident, ok := key.(*ast.Ident); ok {
		if ident.Obj == nil || ident.Obj.Decl == nil || ident.Obj.Kind != ast.Con {
			return "", false
		}
		if spec, ok := ident.Obj.Decl.(*ast.ValueSpec); ok && len(spec.Values) > 0 {
			// TODO: Support len(spec.Values) > 1; e.g. const foo, bar = 1, 2.
			key = spec.Values[0]
		}
	}
	if lit, ok := key.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		value, err := strconv.Unquote(lit.Value)
		if err != nil {
			panic("unreachable") // String literals are always quoted.
		}
		return value, true
	}
	return "", false
}

func areArgsOnSameLine(fset *token.FileSet, call ast.Expr, keys, attrs []ast.Expr) bool {
	if len(keys)+len(attrs) <= 1 {
		return false // Special case: slog.Info("msg", "key", "value") is fine.
	}

	args := slices.Concat([]ast.Expr{call}, keys, attrs)

	lines := make(map[int]struct{}, len(args))
	for _, arg := range args {
		line := fset.Position(arg.Pos()).Line
		if _, ok := lines[line]; ok {
			return true
		}
		lines[line] = struct{}{}
	}

	return false
}
