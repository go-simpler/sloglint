# sloglint

[![checks](https://github.com/go-simpler/sloglint/actions/workflows/checks.yaml/badge.svg)](https://github.com/go-simpler/sloglint/actions/workflows/checks.yaml)
[![docs](https://pkg.go.dev/badge/go-simpler.org/sloglint.svg)](https://pkg.go.dev/go-simpler.org/sloglint)
[![codecov](https://codecov.io/gh/go-simpler/sloglint/branch/main/graph/badge.svg)](https://codecov.io/gh/go-simpler/sloglint)

A Go linter that ensures consistent code style when using `log/slog`.

## Install

`sloglint` is integrated into [golangci-lint](https://golangci-lint.run), and this is the recommended way to use it.

```yaml
# .golangci.yaml
linters:
  enable:
    - sloglint
```

Alternatively, you can download a prebuilt binary from the [Releases](https://github.com/go-simpler/sloglint/releases) page to use `sloglint` standalone.

## Supported checks

For `log/slog` functions:
- [No global logger](#no-global-logger)
- [Context only](#context-only)
- [Discard handler](#discard-handler)

For log messages:
- [Static message](#static-message)
- [Message style](#message-style)

For log arguments:
- [No mixed arguments](#no-mixed-arguments)
- [Key-value pairs only](#key-value-pairs-only)
- [Attributes only](#attributes-only)
- [Arguments on separate lines](#arguments-on-separate-lines)

For log keys:
- [No literal keys](#no-literal-keys)
- [Key naming case](#key-naming-case)
- [Allowed keys](#allowed-keys)
- [Forbidden keys](#forbidden-keys)

The checks for log messages, arguments, and keys can also be used to analyze [custom functions](#custom-functions-analysis).

### No global logger

Report the use of global loggers.
Alternatively, only report the use of the `slog.Default()` logger.

```go
// sloglint: global logger should not be used
slog.Info("a user has logged in", "user_id", 42)
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      no-global: "all" # Or "default".
```

### Context only

Report the use of functions without a `context.Context`.
Alternatively, only report their use if a context exists within the scope of the outermost function.

```go
// sloglint: InfoContext should be used instead
slog.Info("a user has logged in")
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      context: "all" # Or "scope".
```

This check partially supports autofix.

### Discard handler

Suggest using `slog.DiscardHandler` when possible.

```go
// sloglint: use slog.DiscardHandler instead
slog.NewJSONHandler(io.Discard, nil)
```

This check is enabled by default and supports autofix.

### Static message

Report dynamic log messages, such as those that are built with `fmt.Sprintf`.

```go
// sloglint: message should be a string literal or a constant
slog.Info(fmt.Sprintf("a user with id %d has logged in", 42))
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      static-msg: true
```

### Message style

Report log messages that do not match a particular style.
The supported styles are `lowercased` (the first letter is lowercase) and `capitalized` (the first letter is uppercase).

```go
// sloglint: message should be lowercased
slog.Info("A user has logged in")
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      msg-style: "lowercased" # Or "capitalized".
```

### No mixed arguments

Report the use of both key-value pairs and attributes within a single function call.

```go
// sloglint: key-value pairs and attributes should not be mixed
slog.Info("a user has logged in", "user_id", 42, slog.String("ip_address", "192.0.2.0"))
```

This check is enabled by default.

### Key-value pairs only

Report any use of attributes as function call arguments.

```go
// sloglint: attributes should not be used
slog.Info("a user has logged in", slog.Int("user_id", 42))
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      kv-only: true
```

### Attributes only

Report any use of key-value pairs as function call arguments.

```go
// sloglint: key-value pairs should not be used
slog.Info("a user has logged in", "user_id", 42)
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      attr-only: true
```

### Arguments on separate lines

Report two or more arguments on the same line.
A key-value pair is considered a single argument.

```go
// sloglint: arguments should be put on separate lines
slog.Info("a user has logged in", "user_id", 42, "ip_address", "192.0.2.0")
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      args-on-sep-lines: true
```

### No literal keys

Report the use of literal strings as log keys.

```go
// sloglint: literal keys should not be used
slog.Info("a user has logged in", "user_id", 42)
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      no-raw-keys: true
```

### Key naming case

Report log keys that do not match a particular naming case.
The supported cases are `snake_case`, `kebab-case`, `camelCase`, and `PascalCase`.

```go
// sloglint: keys should be written in snake_case
slog.Info("a user has logged in", "user-id", 42)
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      key-naming-case: "snake" # Or "kebab", "camel", "pascal".
```

This check supports autofix.

### Allowed keys

Report the use of log keys that are not explicitly allowed.

```go
// sloglint: the "id" key is not allowed and should not be used
slog.Info("a user has logged in", "id", 42)
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      allowed-keys:
        - user_id
        - ip_address
```

### Forbidden keys

Report the use of forbidden log keys.
When using the standard `slog.JSONHandler` or `slog.TextHandler`,
you may want to forbid the `time`, `level`, `msg`, and `source` keys,
as these will be written by the handler.

```go
// sloglint: the "time" key is forbidden and should not be used
slog.Info("a user has logged in", "time", time.Now())
```

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      forbidden-keys:
        - time
        - level
        - msg
        - source
```

## Custom functions analysis

Analyze custom functions in addition to the standard `log/slog` functions.

The following properties must be specified:
1. The full name of the function, including the package, e.g. `log/slog.Info`.
If the function is a method, the receiver type must be wrapped in parentheses, e.g. `(*log/slog.Logger).Info`.
2. The position of the `msg string` argument in the function signature, starting from 0.
If there is no message in the function, a negative value must be passed.
3. The position of the `args ...any` argument in the function signature, starting from 0.
If there are no arguments in the function, a negative value must be passed.

Here's an example for the [exp/slog](https://pkg.go.dev/golang.org/x/exp/slog) package, the predecessor of `log/slog`.

```yaml
# .golangci.yaml
linters:
  settings:
    sloglint:
      custom-funcs:
        - name: "(*golang.org/x/exp/slog.Logger).InfoContext"
          msg-pos: 1
          args-pos: 2
```
