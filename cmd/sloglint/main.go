package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"go.tmz.dev/sloglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var version = "dev" // injected at build time.

func main() {
	// override the builtin -V flag.
	flag.Var(versionFlag{}, "V", "print version and exit")
	singlechecker.Main(sloglint.New())
}

type versionFlag struct{}

func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) Get() interface{} { return nil }
func (versionFlag) String() string   { return "" }
func (versionFlag) Set(string) error {
	fmt.Printf("sloglint version %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
	return nil
}
