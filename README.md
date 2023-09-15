# sloglint

[![checks](https://github.com/tmzane/sloglint/actions/workflows/checks.yml/badge.svg)](https://github.com/tmzane/sloglint/actions/workflows/checks.yml)
[![pkg.go.dev](https://pkg.go.dev/badge/go.tmz.dev/sloglint.svg)](https://pkg.go.dev/go.tmz.dev/sloglint)
[![goreportcard](https://goreportcard.com/badge/go.tmz.dev/sloglint)](https://goreportcard.com/report/go.tmz.dev/sloglint)
[![codecov](https://codecov.io/gh/tmzane/sloglint/branch/main/graph/badge.svg)](https://codecov.io/gh/tmzane/sloglint)

Ensure consistent code-style when using `log/slog`.

## 🚀 Features

* Forbid mixing key-value pairs and attributes in a single function call (default)
* Enforce using either key-value pairs or attributes for the entire project (optional)
* Enforce using constants (or custom `slog.Attr` constructors) instead of raw keys (optional)
* [WIP] Enforce putting arguments on separate lines (optional)
