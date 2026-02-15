package sloglint

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Func describes a function to analyze, e.g. [slog.Info].
type Func struct {
	// The full name of the function, including the package, e.g. "log/slog.Info".
	// If the function is a method, the receiver type must be wrapped in parentheses, e.g. "(*log/slog.Logger).Info".
	Name string
	// The position of the "msg string" argument in the function signature, starting from 0.
	// If there is no message in the function, a negative value must be passed.
	MsgPos int
	// The position of the "args ...any" argument in the function signature, starting from 0.
	// If there are no arguments in the function, a negative value must be passed.
	ArgsPos int
	// Whether this is a function from the standard [log/slog] package.
	standard bool
}

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
	CustomFuncs    []Func   // Custom functions to analyze in addition to the standard [log/slog] functions.
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
	fs := flag.NewFlagSet("sloglint", flag.ContinueOnError)

	fs.BoolVar(&opts.NoMixedArgs, "no-mixed-args", opts.NoMixedArgs, "enforce not mixing key-value pairs and attributes")
	fs.BoolVar(&opts.KVOnly, "kv-only", opts.KVOnly, "enforce using key-value pairs only (overrides -no-mixed-args, incompatible with -attr-only)")
	fs.BoolVar(&opts.AttrOnly, "attr-only", opts.AttrOnly, "enforce using attributes only (overrides -no-mixed-args, incompatible with -kv-only)")
	fs.StringVar(&opts.NoGlobal, "no-global", opts.NoGlobal, "enforce not using global loggers (all|default)")
	fs.StringVar(&opts.ContextOnly, "context-only", opts.ContextOnly, "enforce using methods that accept a context (all|scope)")
	fs.BoolVar(&opts.StaticMsg, "static-msg", opts.StaticMsg, "enforce using static messages")
	fs.StringVar(&opts.MsgStyle, "msg-style", opts.MsgStyle, "enforce message style (lowercased|capitalized)")
	fs.BoolVar(&opts.NoRawKeys, "no-raw-keys", opts.NoRawKeys, "enforce using constants instead of raw keys")
	fs.StringVar(&opts.KeyNamingCase, "key-naming-case", opts.KeyNamingCase, "enforce key naming convention (snake|kebab|camel|pascal)")
	fs.BoolVar(&opts.ArgsOnSepLines, "args-on-sep-lines", opts.ArgsOnSepLines, "enforce putting arguments on separate lines")

	fs.Func("allowed-keys", "enforce using specific keys only (comma-separated)", func(s string) error {
		opts.AllowedKeys = append(opts.AllowedKeys, strings.Split(s, ",")...)
		return nil
	})

	fs.Func("forbidden-keys", "enforce not using specific keys (comma-separated)", func(s string) error {
		opts.ForbiddenKeys = append(opts.ForbiddenKeys, strings.Split(s, ",")...)
		return nil
	})

	fs.Func("fn", "analyze a custom function (name:msg-pos:args-pos)", func(s string) error {
		name, rest, _ := strings.Cut(s, ":")
		fn := Func{Name: name}
		_, err := fmt.Sscanf(rest, "%d:%d", &fn.MsgPos, &fn.ArgsPos)
		opts.CustomFuncs = append(opts.CustomFuncs, fn)
		return err
	})

	return *fs
}
