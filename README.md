<div align="center">
  <h1 style="font-size: 2rem;margin: 0;"> pretty-go-log </h1>
  <p>
    <a href="https://github.com/canefe/pretty-go-log/releases"><img src="https://img.shields.io/github/v/release/canefe/pretty-go-log?include_prereleases&style=flat-square" alt="Latest Version" /></a>&nbsp;
    <a href="https://github.com/canefe/pretty-go-log/actions"><img src="https://img.shields.io/github/actions/workflow/status/canefe/pretty-go-log/ci.yml?style=flat-square" alt="Build Status" /></a>&nbsp;
    <a href="https://goreportcard.com/report/github.com/canefe/pretty-go-log"><img src="https://goreportcard.com/badge/github.com/canefe/pretty-go-log?style=flat-square" alt="Go Report Card" /></a>
  </p>
</div>

A small Go package that gives logrus a cleaner, structured, and visually aligned output style with bracketed tags, optional colors, and multi-output support.

Use it as a drop-in logger with sane defaults, or fully customize formatting, output, and caller behavior using functional options.

![pretty-go-log](https://i.imgur.com/zK3tD3M.png)

## What is Pretty Go Log?

Pretty Go Log wraps logrus with a formatter designed for readability:

- **Aligned Tag Output**. Bracketed tags align to a common column for fast scanning
- **Multiple Tag Styles**. Default, centered, or right-aligned tag padding
- **Multi Output**. Console + rotating file output with separate formatters
- **Configurable Colors**. Colorized levels and tags, dim gray padding
- **Caller Awareness**. Toggle caller output for warnings and errors
- **Environment Config**. Use `LOG_LEVEL`, `LOG_OUTPUT`, `LOG_FORMAT` defaults

## Installation

```bash
go get github.com/canefe/pretty-go-log
```

## Quick Start

```go
package main

import (
    "github.com/canefe/pretty-go-log/logrus/pretty"
)

func main() {
    log := pretty.New()

    log.Info("[Server] Started on :8080")
    log.Warn("[Cache] Miss for key user:123")
}
```

## Usage

### Functional Options

```go
package main

import (
    "github.com/canefe/pretty-go-log/logrus/pretty"
    "github.com/sirupsen/logrus"
)

func main() {
    log := pretty.New(
        pretty.WithLevel(logrus.DebugLevel),
        pretty.WithOutput(pretty.OutputMulti),
        pretty.WithFormat(pretty.FormatPlain),
        pretty.WithNamespace("App"),
        pretty.WithFile("logs/service.log"),
    )

    log.Debug("[Init] Starting")
}
```

### Custom Formatter

```go
package main

import (
    "github.com/canefe/pretty-go-log/logrus/pretty"
    "github.com/sirupsen/logrus"
)

func main() {
    formatter := pretty.NewCustomFormatter(
        pretty.WithTagStyle(pretty.StyleRight, "."),
        pretty.WithBracketPadding(15),
        pretty.WithColorBrackets(true),
    )

    log := pretty.New(
        pretty.WithLevel(logrus.InfoLevel),
        pretty.WithOutput(pretty.OutputMulti),
        pretty.WithCustomFormat(*formatter),
        pretty.WithFile("logs/service.log"),
    )

    log.Info("[API] Request completed")
}
```

### Environment Configuration

```go
package main

import (
    "os"

    "github.com/canefe/pretty-go-log/logrus/pretty"
)

func main() {
    os.Setenv("LOG_LEVEL", "debug")
    os.Setenv("LOG_OUTPUT", "console")
    os.Setenv("LOG_FORMAT", "plain")

    log := pretty.New()
    log.Debug("[Env] Logger configured via env vars")
}
```

## Options and Types

### Output Types

- `pretty.OutputConsole`
- `pretty.OutputFile`
- `pretty.OutputMulti`

### Format Types

- `pretty.FormatRaw`
- `pretty.FormatPlain`
- `pretty.FormatJSON`

### Common Options

- `pretty.WithLevel(level logrus.Level)`
- `pretty.WithOutput(output pretty.OutputType)`
- `pretty.WithFormat(format pretty.FormatType)`
- `pretty.WithNamespace(name string)`
- `pretty.WithFile(path string)`
- `pretty.WithoutCaller()`
- `pretty.WithCustomFormat(formatter pretty.CustomFormatter)`

## Examples

- `examples/logrus/basic`
- `examples/logrus/custom-config`
- `examples/logrus/env-config`
- `examples/logrus/file-output`
- `examples/logrus/json-format`
- `examples/logrus/multi-output`
- `examples/showcase`
