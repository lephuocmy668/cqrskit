CQRSKit
--------
[![Go Report Card](https://goreportcard.com/badge/github.com/gokit/cqrskit)](https://goreportcard.com/report/github.com/gokit/cqrskit)
[![Travis Build Status](https://travis-ci.org/gokit/cqrskit.svg?branch=master)](https://travis-ci.org/gokit/cqrskit#)
[![CircleCI](https://circleci.com/gh/gokit/cqrskit.svg?style=svg)](https://circleci.com/gh/gokit/cqrskit)

CQRSKit provides a package that combines code generation and a simple ES+CQRS-style API, to quickly scaffold around a giving struct. It removes the need of using reflection to figure out which events and which methods to call by generating an applier which uses Go AST to typed-safe alternative that knows which events are to be handled by which struct method.

## Install

```
go get github.com/gokit/cqrskit
```

## Examples

See [Examples](./examples) for demonstrations of using cqrskit which a package. 

## CLI

```bash
> cqrskit generate
```

```bash
> cqrskit
Usage: cqrskit [flags] [command] 

⡿ COMMANDS:
	⠙ generate        Generates a es+cqrs scaffolding API for struct types.

⡿ HELP:
	Run [command] help

⡿ OTHERS:
	Run 'cqrskit flags' to print all flags of all commands.

⡿ WARNING:
	Uses internal flag package so flags must precede command name. 
	e.g 'cqrskit -cmd.flag=4 run'

```


## How It works

1. Annotate all chosen interface with `@escqrs`
2. Run `cqrskit generate` within package of annotated interfaces.

### Interface Based Annotation

You annotate any giving interface with `@escqrs` which marks giving interface has a target for code generation.


Sample below:

```go
// @escqrs
type Users struct {
	Name string `json:"name"`
	Email string `json:"email"`
}
```


