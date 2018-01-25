CQRSKit
--------
[![Go Report Card](https://goreportcard.com/badge/github.com/gokit/cqrskit)](https://goreportcard.com/report/github.com/gokit/cqrskit)
[![Travis Build Status](https://travis-ci.org/gokit/cqrskit.svg?branch=master)](https://travis-ci.org/gokit/cqrskit#)
[![CircleCI](https://circleci.com/gh/gokit/cqrskit.svg?style=svg)](https://circleci.com/gh/gokit/cqrskit)

CQRSKit implements a code generator which automatically generates a ES+CQRS-style API package, which provides code scaffolding for quick event sourcing based APIs.

## Install

```
go get github.com/gokit/cqrskit
```

## Examples

See [Examples](./examples) for demonstrations of packages generated using cqrskit which creates a full RPC-style API with client code for a interface API definition/declaration.

## CLI

```bash
> cqrskit generate
```

```bash
> cqrskit
Usage: cqrskit [flags] [command] 

⡿ COMMANDS:
	⠙ generate        Generates a rpc like API for interface types.

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


