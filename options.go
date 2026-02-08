package sloglint

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Options contains options for the sloglint analyzer.
type Options struct {
	NoMixedArgs    bool     // Enforce not mixing key-value pairs and attributes (default true).
	KVOnly         bool     // Enforce using key-value pairs only (overrides NoMixedArgs, incompatible with AttrOnly).
	AttrOnly       bool     // Enforce using attributes only (overrides NoMixedArgs, incompatible with KVOnly).
	NoGlobal       string   // Enforce not using global loggers ("all" or "default").
	ContextOnly    string   // Enforce using methods that accept a context ("all" or "scope").
	StaticMsg      bool     // Enforce using static messages.
	MsgStyle       string   // Enforce message style ("lowercased" or "capitalized").
	NoRawKeys      bool     // Enforce using constants instead of raw keys.
	KeyNamingCase  string   // Enforce key naming convention ("snake", "kebab", "camel", or "pascal").
	AllowedKeys    []string // Enforce using only specific keys.
	ForbiddenKeys  []string // Enforce not using specific keys.
	ArgsOnSepLines bool     // Enforce putting arguments on separate lines.
}

// Possible values for [Options.NoGlobal].
const (
	noGlobalAll     = "all"
	noGlobalDefault = "default"
)

// Possible values for [Options.ContextOnly].
const (
	contextOnlyAll   = "all"
	contextOnlyScope = "scope"
)

// Possible values for [Options.MsgStyle].
const (
	msgStyleLowercased  = "lowercased"
	msgStyleCapitalized = "capitalized"
)

// Possible values for [Options.KeyNamingCase].
const (
	keyNamingCaseSnake  = "snake"
	keyNamingCaseKebab  = "kebab"
	keyNamingCaseCamel  = "camel"
	keyNamingCasePascal = "pascal"
)

var (
	errIncompatible = errors.New("incompatible")
	errInvalidValue = errors.New("invalid value")
)

func (opts *Options) validate() error {
	if opts.KVOnly && opts.AttrOnly {
		return fmt.Errorf("sloglint: Options.KVOnly and Options.AttrOnly are %w", errIncompatible)
	}

	switch opts.NoGlobal {
	case "", noGlobalAll, noGlobalDefault:
	default:
		return fmt.Errorf("sloglint: Options.NoGlobals has an %w %q", errInvalidValue, opts.NoGlobal)
	}

	switch opts.ContextOnly {
	case "", contextOnlyAll, contextOnlyScope:
	default:
		return fmt.Errorf("sloglint: Options.ContextOnlys has an %w %q", errInvalidValue, opts.ContextOnly)
	}

	switch opts.MsgStyle {
	case "", msgStyleLowercased, msgStyleCapitalized:
	default:
		return fmt.Errorf("sloglint: Options.MsgStyles has an %w %q", errInvalidValue, opts.MsgStyle)
	}

	switch opts.KeyNamingCase {
	case "", keyNamingCaseSnake, keyNamingCaseKebab, keyNamingCaseCamel, keyNamingCasePascal:
	default:
		return fmt.Errorf("sloglint: Options.KeyNamingCases has an %w %q", errInvalidValue, opts.KeyNamingCase)
	}

	return nil
}

func flags(opts *Options) flag.FlagSet {
	fset := flag.NewFlagSet("sloglint", flag.ContinueOnError)

	boolVar(fset, &opts.NoMixedArgs, "no-mixed-args", "enforce not mixing key-value pairs and attributes (default true)")
	boolVar(fset, &opts.KVOnly, "kv-only", "enforce using key-value pairs only (overrides -no-mixed-args, incompatible with -attr-only)")
	boolVar(fset, &opts.AttrOnly, "attr-only", "enforce using attributes only (overrides -no-mixed-args, incompatible with -kv-only)")
	stringVar(fset, &opts.NoGlobal, "no-global", "enforce not using global loggers (all|default)")
	stringVar(fset, &opts.ContextOnly, "context-only", "enforce using methods that accept a context (all|scope)")
	boolVar(fset, &opts.StaticMsg, "static-msg", "enforce using static messages")
	stringVar(fset, &opts.MsgStyle, "msg-style", "enforce message style (lowercased|capitalized)")
	boolVar(fset, &opts.NoRawKeys, "no-raw-keys", "enforce using constants instead of raw keys")
	stringVar(fset, &opts.KeyNamingCase, "key-naming-case", "enforce key naming convention (snake|kebab|camel|pascal)")
	sliceVar(fset, &opts.AllowedKeys, "allowed-keys", "enforce using specific keys only (comma-separated)")
	sliceVar(fset, &opts.ForbiddenKeys, "forbidden-keys", "enforce not using specific keys (comma-separated)")
	boolVar(fset, &opts.ArgsOnSepLines, "args-on-sep-lines", "enforce putting arguments on separate lines")

	return *fset
}

func boolVar(fset *flag.FlagSet, value *bool, name, usage string) {
	fset.BoolFunc(name, usage, func(s string) error {
		v, err := strconv.ParseBool(s)
		*value = v
		return err
	})
}

func stringVar(fset *flag.FlagSet, value *string, name, usage string) {
	fset.Func(name, usage, func(s string) error {
		*value = s
		return nil
	})
}

func sliceVar(fset *flag.FlagSet, value *[]string, name, usage string) {
	fset.Func(name, usage, func(s string) error {
		*value = append(*value, strings.Split(s, ",")...)
		return nil
	})
}
