// Package sloglint implements the sloglint analyzer.
package sloglint

import (
	"go/ast"
	"go/version"
	"slices"

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

var standardFuncs = []Func{
	{Name: "log/slog.Log", MsgPos: 2, ArgsPos: 3, standard: true},
	{Name: "log/slog.LogAttrs", MsgPos: 2, ArgsPos: 3, standard: true},
	{Name: "log/slog.Debug", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "log/slog.Info", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "log/slog.Warn", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "log/slog.Error", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "log/slog.DebugContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "log/slog.InfoContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "log/slog.WarnContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "log/slog.ErrorContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "log/slog.With", MsgPos: -1, ArgsPos: 0, standard: true},
	{Name: "log/slog.Group", MsgPos: -1, ArgsPos: 1, standard: true},
	{Name: "log/slog.NewTextHandler", MsgPos: -1, ArgsPos: -1, standard: true},
	{Name: "log/slog.NewJSONHandler", MsgPos: -1, ArgsPos: -1, standard: true},
	{Name: "(*log/slog.Logger).Log", MsgPos: 2, ArgsPos: 3, standard: true},
	{Name: "(*log/slog.Logger).LogAttrs", MsgPos: 2, ArgsPos: 3, standard: true},
	{Name: "(*log/slog.Logger).Debug", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "(*log/slog.Logger).Info", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "(*log/slog.Logger).Warn", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "(*log/slog.Logger).Error", MsgPos: 0, ArgsPos: 1, standard: true},
	{Name: "(*log/slog.Logger).DebugContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "(*log/slog.Logger).InfoContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "(*log/slog.Logger).WarnContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "(*log/slog.Logger).ErrorContext", MsgPos: 1, ArgsPos: 2, standard: true},
	{Name: "(*log/slog.Logger).With", MsgPos: -1, ArgsPos: 0, standard: true},
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

	funcs := slices.Concat(standardFuncs, opts.CustomFuncs)
	idx := slices.IndexFunc(funcs, func(f Func) bool {
		return f.Name == fn.FullName()
	})
	if idx == -1 {
		return
	}

	if funcs[idx].standard {
		analyzeFunc(pass, opts, call, cursor)
	}
	if pos := funcs[idx].MsgPos; pos >= 0 && len(call.Args) > pos {
		analyzeMsg(pass, opts, call.Args[pos])
	}
	if pos := funcs[idx].ArgsPos; pos >= 0 && len(call.Args) > pos {
		analyzeArgs(pass, opts, call.Args[pos:])
	}
}

func analyzeFunc(pass *analysis.Pass, opts *Options, call *ast.CallExpr, cursor inspector.Cursor) {
	if opts.NoGlobal != "" {
		noGlobal(pass, call, opts.NoGlobal == noGlobalDefault)
	}
	if opts.ContextOnly != "" {
		contextOnly(pass, call, cursor, opts.ContextOnly == contextOnlyScope)
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
