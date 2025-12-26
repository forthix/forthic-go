# Forthic Go Runtime

A Go implementation of the Forthic stack-based concatenative programming language.

## Overview

Forthic is a stack-based, concatenative language designed for composable transformations. This is the official Go runtime implementation, providing full compatibility with other Forthic runtimes (TypeScript, Python, Rust, Ruby, Swift, Java, Erlang, Zig, .NET).

## Features

- ✅ Complete Forthic language implementation
- ✅ All 8 standard library modules
- ✅ gRPC support for multi-runtime execution
- ✅ CLI with REPL, script execution, and eval modes
- ✅ Native Go idioms (goroutines for parallelism)
- ✅ Comprehensive test suite

## Installation

```bash
go get github.com/forthix/forthic-go
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "github.com/forthix/forthic-go/forthic"
)

func main() {
    interp := forthic.NewStandardInterpreter()

    err := interp.Run("[1 2 3] \"2 *\" MAP")
    if err != nil {
        panic(err)
    }

    result := interp.StackPop()
    fmt.Println(result) // [2 4 6]
}
```

### CLI

```bash
# REPL mode
forthic-go repl

# Execute a script
forthic-go run script.forthic

# Eval mode (one-liner)
forthic-go eval "[1 2 3] LENGTH"
```

## Development

```bash
# Run tests
go test ./...

# Build
go build ./...

# Run specific test
go test -v -run TestInterpreter ./forthic
```

## Project Structure

```
forthic-go/
├── forthic/              # Core runtime
│   ├── interpreter.go    # Main interpreter
│   ├── tokenizer.go      # Lexical analysis
│   ├── module.go         # Module system
│   ├── literals.go       # Literal parsers
│   └── modules/
│       └── standard/     # Standard library (8 modules)
├── grpc/                 # gRPC support
├── cmd/forthic/          # CLI tool
└── tests/                # Test suites
```

## Standard Library Modules

- **core**: Stack operations, variables, control flow
- **array**: Data transformation (MAP, SELECT, SORT, etc.)
- **record**: Dictionary operations
- **string**: Text processing
- **math**: Arithmetic operations
- **boolean**: Logical operations
- **datetime**: Date/time manipulation
- **json**: JSON serialization

## Multi-Runtime Execution

This runtime supports calling words from other Forthic runtimes via gRPC:

```go
// Call a Python word from Go
result, err := interp.ExecuteRemoteWord("python-runtime", "MY-WORD", args)
```

## License

Apache 2.0

## Links

- [Forthic Language Specification](https://github.com/forthix/forthic)
- [TypeScript Runtime](https://github.com/forthix/forthic-ts) (reference implementation)
- [Documentation](https://forthix.github.io/forthic-go)
