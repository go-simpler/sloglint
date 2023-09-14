// Package sloglint implements the sloglint analyzer.
package sloglint

import (
	"errors"
	"flag"
	"go/ast"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

type Options struct {
	KVOnly    bool
	AttrOnly  bool
	NoRawKeys bool
}

// New creates a new sloglint analyzer.
func New(opts *Options) *analysis.Analyzer {
	if opts == nil {
		opts = new(Options)
	}
	return &analysis.Analyzer{
		Name:     "sloglint",
		Doc:      "ensure consistent code style when using log/slog",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    flags(opts),
		Run: func(pass *analysis.Pass) (any, error) {
			if opts.KVOnly && opts.AttrOnly {
				return nil, errors.New("sloglint: incompatible options provided")
			}
			run(pass, opts)
			return nil, nil
		},
	}
}

func flags(opts *Options) flag.FlagSet {
	fs := flag.NewFlagSet("sloglint", flag.ContinueOnError)

	boolVar := func(value *bool, name, usage string) {
		fs.BoolFunc(name, usage, func(s string) error {
			v, err := strconv.ParseBool(s)
			*value = v
			return err
		})
	}

	boolVar(&opts.KVOnly, "kv-only", "enforce using key-value pairs only (incompatible with -attr-only)")
	boolVar(&opts.AttrOnly, "attr-only", "enforce using attributes only (incompatible with -kv-only)")
	boolVar(&opts.NoRawKeys, "no-raw-keys", "forbid using raw keys")

	return *fs
}

// mapping: function name -> arguments position.
var funcs = map[string]int{
	"log/slog.Log":                    3,
	"log/slog.Debug":                  1,
	"log/slog.Info":                   1,
	"log/slog.Warn":                   1,
	"log/slog.Error":                  1,
	"log/slog.DebugContext":           2,
	"log/slog.InfoContext":            2,
	"log/slog.WarnContext":            2,
	"log/slog.ErrorContext":           2,
	"(*log/slog.Logger).Log":          3,
	"(*log/slog.Logger).Debug":        1,
	"(*log/slog.Logger).Info":         1,
	"(*log/slog.Logger).Warn":         1,
	"(*log/slog.Logger).Error":        1,
	"(*log/slog.Logger).DebugContext": 2,
	"(*log/slog.Logger).InfoContext":  2,
	"(*log/slog.Logger).WarnContext":  2,
	"(*log/slog.Logger).ErrorContext": 2,
}

func run(pass *analysis.Pass, opts *Options) {
	visit := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	types := []ast.Node{(*ast.CallExpr)(nil)}

	visit.Preorder(types, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		fn := typeutil.StaticCallee(pass.TypesInfo, call)
		if fn == nil {
			return
		}

		argsPos, ok := funcs[fn.FullName()]
		if !ok {
			return
		}

		// NOTE: we assume that the arguments have already been validated by govet.
		args := call.Args[argsPos:]
		if len(args) == 0 {
			return
		}

		keys := make([]ast.Expr, 0)
		attrs := make([]ast.Expr, 0)

		for i := 0; i < len(args); i++ {
			switch pass.TypesInfo.TypeOf(args[i]).String() {
			case "string":
				keys = append(keys, args[i])
				i++ // skip the value.
			case "log/slog.Attr":
				attrs = append(attrs, args[i])
			}
		}

		switch {
		case opts.KVOnly && len(attrs) > 0:
			pass.Reportf(call.Pos(), "attributes should not be used")
		case opts.AttrOnly && len(attrs) < len(args):
			pass.Reportf(call.Pos(), "key-value pairs should not be used")
		case 0 < len(attrs) && len(attrs) < len(args):
			pass.Reportf(call.Pos(), "key-value pairs and attributes should not be mixed")
		}

		const rawKeysReport = "raw keys should not be used"

		if !opts.NoRawKeys {
			return
		}

		for _, key := range keys {
			if !isConst(key) {
				pass.Reportf(call.Pos(), rawKeysReport)
				return
			}
		}

		for _, attr := range attrs {
			switch attr := attr.(type) {
			case *ast.CallExpr: // e.g. slog.Int()
				builtins := map[string]struct{}{
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
				fn := typeutil.StaticCallee(pass.TypesInfo, attr)
				if _, ok := builtins[fn.FullName()]; ok && !isConst(attr.Args[0]) {
					pass.Reportf(call.Pos(), rawKeysReport)
					return
				}
			}
		}
	})
}

func isConst(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Obj != nil && ident.Obj.Kind == ast.Con
}
