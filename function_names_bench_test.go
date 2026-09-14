package sloglint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

func BenchmarkAnalyzeNodeFunctionName(b *testing.B) {
	const source = `package p
func plain() {}
type receiver struct{}
func (receiver) method() {}
func (*receiver) pointerMethod() {}
type generic[T any] struct{}
func (generic[T]) method() {}
func calls(r receiver, g generic[int]) {
	plain()
	r.method()
	r.pointerMethod()
	g.method()
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p.go", source, 0)
	require.NoError(b, err)

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Instances:  make(map[*ast.Ident]types.Instance),
	}
	var config types.Config
	files := []*ast.File{file}
	pkg, err := config.Check("example.com/project/p", fset, files, info)
	require.NoError(b, err)

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     files,
		Pkg:       pkg,
		TypesInfo: info,
	}
	var calls []inspector.Cursor
	for cursor := range inspector.New(files).Root().Preorder(new(ast.CallExpr)) {
		calls = append(calls, cursor)
	}
	require.Len(b, calls, 4)

	for _, count := range []int{0, 25, 100} {
		b.Run(fmt.Sprintf("custom=%d", count), func(b *testing.B) {
			var opts Options
			for i := range count {
				opts.CustomFuncs = append(opts.CustomFuncs, Func{FullName: fmt.Sprintf("custom.Log%d", i)})
			}

			for i, name := range []string{"function", "value-method", "pointer-method", "generic-method"} {
				b.Run(name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						analyzeNode(pass, &opts, calls[i])
					}
				})
			}
		})
	}
}
