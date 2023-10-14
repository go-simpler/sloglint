# sloglint

[![checks](https://github.com/go-simpler/sloglint/actions/workflows/checks.yml/badge.svg)](https://github.com/go-simpler/sloglint/actions/workflows/checks.yml)
[![pkg.go.dev](https://pkg.go.dev/badge/go-simpler.org/sloglint.svg)](https://pkg.go.dev/go-simpler.org/sloglint)
[![goreportcard](https://goreportcard.com/badge/go-simpler.org/sloglint)](https://goreportcard.com/report/go-simpler.org/sloglint)
[![codecov](https://codecov.io/gh/go-simpler/sloglint/branch/main/graph/badge.svg)](https://codecov.io/gh/go-simpler/sloglint)

Ensure consistent code-style when using `log/slog`.

## 🚀 Features

* Forbid mixing key-value pairs and attributes in a single function call (default)
* Enforce using either key-value pairs or attributes for the entire project (optional)
* Enforce using constants (or custom `slog.Attr` constructors) instead of raw keys (optional)
* Enforce putting arguments on separate lines (optional)
