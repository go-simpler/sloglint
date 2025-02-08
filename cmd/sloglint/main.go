package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"go-simpler.org/sloglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var version = "dev" // injected at build time.

func main() {
	// override the builtin -V flag.
	flag.Var(versionFlag{}, "V", "print version and exit")
	singlechecker.Main(sloglint.New(nil))
}

type versionFlag struct{}

func (versionFlag) String() string   { return "" }
func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) Set(string) error {
	fmt.Printf("sloglint version %s %s/%s (built with %s)\n", version, runtime.GOOS, runtime.GOARCH, runtime.Version())
	os.Exit(0)
	return nil
}
