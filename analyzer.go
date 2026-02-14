// Package sloglint implements the sloglint analyzer.
package sloglint

import (
	"go/ast"
	"go/types"
	"go/version"

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
		URL:      "https://go-simpler.org/sloglint",
		Flags:    flags(opts),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			if err := opts.validate(); err != nil {
				return nil, err
			}

			root := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Root()
			for cursor := range root.Preorder(new(ast.CallExpr), new(ast.CompositeLit)) {
				analyzeNode(pass, opts, cursor)
			}

			return nil, nil
		},
	}
}

var slogFuncs = map[string]struct {
	msgPos  int // The position of the "msg string" argument in the function signature, starting from 0.
	argsPos int // The position of the "args ...any" argument in the function signature, starting from 0.
}{
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
	"log/slog.Group":                  {msgPos: -1, argsPos: 1},
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

func analyzeNode(pass *analysis.Pass, opts *Options, cursor inspector.Cursor) {
	node := cursor.Node()
	if cl, ok := node.(*ast.CompositeLit); ok && typeName(pass.TypesInfo, cl) == "log/slog.Attr" {
		analyzeAttrKey(pass, opts, cl)
		return
	}

	call, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}

	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil {
		return
	}

	switch fn.FullName() {
	case "log/slog.Int",
		"log/slog.Int64",
		"log/slog.Uint64",
		"log/slog.Float64",
		"log/slog.String",
		"log/slog.Bool",
		"log/slog.Time",
		"log/slog.Duration",
		"log/slog.Any":
		analyzeKey(pass, opts, call.Args[0])
		return
	case "log/slog.Group":
		analyzeKey(pass, opts, call.Args[0])
		// Special case: don't return here, we also need to analyze the group's arguments.
	}

	info, ok := slogFuncs[fn.FullName()]
	if !ok {
		return
	}

	analyzeFunc(pass, opts, fn, call, cursor)

	if info.msgPos >= 0 && len(call.Args) > info.msgPos {
		analyzeMsg(pass, opts, call.Args[info.msgPos])
	}

	if info.argsPos >= 0 && len(call.Args) > info.argsPos {
		analyzeArgs(pass, opts, call.Args[info.argsPos:])
	}
}

func analyzeFunc(pass *analysis.Pass, opts *Options, fn *types.Func, call *ast.CallExpr, cursor inspector.Cursor) {
	if opts.NoGlobal != "" {
		noGlobal(pass, fn, call, opts.NoGlobal == noGlobalDefault)
	}
	if opts.ContextOnly != "" {
		contextOnly(pass, fn, call, cursor, opts.ContextOnly == contextOnlyScope)
	}
	v := pass.Module.GoVersion // Empty in test runs.
	if v == "" || version.Compare("go"+v, "go1.24") >= 0 {
		discardHandler(pass, call)
	}
}

func analyzeMsg(pass *analysis.Pass, opts *Options, msg ast.Expr) {
	if opts.StaticMsg {
		staticMsg(pass, msg)
	}
	if opts.MsgStyle != "" {
		msgStyle(pass, msg, opts.MsgStyle)
	}
}

func analyzeArgs(pass *analysis.Pass, opts *Options, args []ast.Expr) {
	var keys, attrs []ast.Expr

	for i := 0; i < len(args); i++ {
		typ := pass.TypesInfo.TypeOf(args[i])
		if typ == nil {
			continue
		}
		switch typ.String() {
		case "string":
			keys = append(keys, args[i])
			analyzeKey(pass, opts, args[i])
			i++ // Skip the value.
		case "log/slog.Attr":
			attrs = append(attrs, args[i])
		case "[]any", "[]log/slog.Attr":
			continue // The last argument may be an unpacked slice, skip it.
		}
	}

	if opts.NoMixedArgs {
		noMixedArgs(pass, keys, attrs)
	}
	if opts.KVOnly {
		kvOnly(pass, attrs)
	}
	if opts.AttrOnly {
		attrOnly(pass, keys)
	}
	if opts.ArgsOnSepLines {
		argsOnSepLines(pass, keys, attrs)
	}
}

func analyzeKey(pass *analysis.Pass, opts *Options, key ast.Expr) {
	if opts.NoRawKeys {
		noRawKeys(pass, key)
	}
	if opts.KeyNamingCase != "" {
		keyNamingCase(pass, key, opts.KeyNamingCase)
	}
	if len(opts.AllowedKeys) > 0 {
		allowedKeys(pass, key, opts.AllowedKeys)
	}
	if len(opts.ForbiddenKeys) > 0 {
		forbiddenKeys(pass, key, opts.ForbiddenKeys)
	}
}

func analyzeAttrKey(pass *analysis.Pass, opts *Options, attr *ast.CompositeLit) {
	switch len(attr.Elts) {
	case 1:
		if kv := attr.Elts[0].(*ast.KeyValueExpr); kv.Key.(*ast.Ident).Name == "Key" {
			analyzeKey(pass, opts, kv.Value) // slog.Attr{Key: ...}
		}
	case 2:
		if kv, ok := attr.Elts[0].(*ast.KeyValueExpr); ok && kv.Key.(*ast.Ident).Name == "Key" {
			analyzeKey(pass, opts, kv.Value) // slog.Attr{Key: ..., Value: ...}
		} else if kv, ok := attr.Elts[1].(*ast.KeyValueExpr); ok && kv.Key.(*ast.Ident).Name == "Key" {
			analyzeKey(pass, opts, kv.Value) // slog.Attr{Value: ..., Key: ...}
		} else {
			analyzeKey(pass, opts, attr.Elts[0]) // slog.Attr{..., ...}
		}
	}
}
