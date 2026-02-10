// Package sloglint implements the sloglint analyzer.
package sloglint

import (
	"go/ast"
	"go/token"
	"go/types"
	"iter"
	"strconv"

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

type slogFuncInfo struct {
	msgPos  int // The position of the "msg string" argument in the function signature, starting from 0.
	argsPos int // The position of the "args ...any" argument in the function signature, starting from 0.
}

var slogFuncs = map[string]slogFuncInfo{
	"log/slog.Log":                    {msgPos: 2, argsPos: 3},
	"log/slog.LogAttrs":               {msgPos: 2, argsPos: 3},
	"log/slog.Debug":                  {msgPos: 0, argsPos: 1},
	"log/slog.Info":                   {msgPos: 0, argsPos: 1},
	"log/slog.Warn":                   {msgPos: 0, argsPos: 1},
	"log/slog.Error":                  {msgPos: 0, argsPos: 1},
	"log/slog.DebugContext":           {msgPos: 1, argsPos: 2},
	"log/slog.InfoContext":            {msgPos: 1, argsPos: 2},
	"log/slog.WarnContext":            {msgPos: 1, argsPos: 2},
	"log/slog.ErrorContext":           {msgPos: 1, argsPos: 2},
	"log/slog.With":                   {msgPos: -1, argsPos: 0},
	"log/slog.NewTextHandler":         {msgPos: -1, argsPos: -1},
	"log/slog.NewJSONHandler":         {msgPos: -1, argsPos: -1},
	"(*log/slog.Logger).Log":          {msgPos: 2, argsPos: 3},
	"(*log/slog.Logger).LogAttrs":     {msgPos: 2, argsPos: 3},
	"(*log/slog.Logger).Debug":        {msgPos: 0, argsPos: 1},
	"(*log/slog.Logger).Info":         {msgPos: 0, argsPos: 1},
	"(*log/slog.Logger).Warn":         {msgPos: 0, argsPos: 1},
	"(*log/slog.Logger).Error":        {msgPos: 0, argsPos: 1},
	"(*log/slog.Logger).DebugContext": {msgPos: 1, argsPos: 2},
	"(*log/slog.Logger).InfoContext":  {msgPos: 1, argsPos: 2},
	"(*log/slog.Logger).WarnContext":  {msgPos: 1, argsPos: 2},
	"(*log/slog.Logger).ErrorContext": {msgPos: 1, argsPos: 2},
	"(*log/slog.Logger).With":         {msgPos: -1, argsPos: 0},
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

	info, ok := slogFuncs[fn.FullName()]
	if !ok {
		return
	}

	var msg ast.Expr
	if info.msgPos >= 0 && len(call.Args) > info.msgPos {
		msg = call.Args[info.msgPos]
	}

	var args []ast.Expr
	if info.argsPos >= 0 && len(call.Args) > info.argsPos {
		args = call.Args[info.argsPos:]
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

	for _, check := range checks {
		check(pass, opts, slogFuncCall{
			fn:     fn,
			expr:   call,
			cursor: cursor,
			msg:    msg,
			keys:   keys,
			attrs:  attrs,
		})
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
			case *ast.CallExpr: // slog.Int(..., ...)
				fn := typeutil.StaticCallee(info, attr)
				if fn == nil {
					continue
				}
				if _, ok := attrFuncs[fn.FullName()]; !ok {
					continue
				}
				if !yield(attr.Args[0]) {
					return
				}

			case *ast.CompositeLit: // slog.Attr{}
				switch len(attr.Elts) {
				case 1:
					if kv := attr.Elts[0].(*ast.KeyValueExpr); kv.Key.(*ast.Ident).Name == "Key" {
						if !yield(kv.Value) { // slog.Attr{Key: ...}
							return
						}
					}

				case 2:
					if kv, ok := attr.Elts[0].(*ast.KeyValueExpr); ok && kv.Key.(*ast.Ident).Name == "Key" {
						if !yield(kv.Value) { // slog.Attr{Key: ..., Value: ...}
							return
						}
					} else if kv, ok := attr.Elts[1].(*ast.KeyValueExpr); ok && kv.Key.(*ast.Ident).Name == "Key" {
						if !yield(kv.Value) { // slog.Attr{Value: ..., Key: ...}
							return
						}
					} else {
						if !yield(attr.Elts[0]) { // slog.Attr{..., ...}
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
			key = spec.Values[0] // TODO: Support len(spec.Values) > 1; e.g. const foo, bar = 1, 2.
		}
	}

	lit, ok := key.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	name, err := strconv.Unquote(lit.Value)
	if err != nil {
		panic("unreachable") // String literals are always quoted.
	}

	return name, true
}
