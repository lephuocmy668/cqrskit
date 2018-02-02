CQRSKit
--------
[![Go Report Card](https://goreportcard.com/badge/github.com/gokit/cqrskit)](https://goreportcard.com/report/github.com/gokit/cqrskit)
[![Travis Build Status](https://travis-ci.org/gokit/cqrskit.svg?branch=master)](https://travis-ci.org/gokit/cqrskit#)
[![CircleCI](https://circleci.com/gh/gokit/cqrskit.svg?style=svg)](https://circleci.com/gh/gokit/cqrskit)

CQRSKit provides a base libary that implements a simple but extensive Event Sourcing + CQRS library. It's purpose is to provide multi-store event stores with a simple but robust API for handling CQRS style APIs. It is heavily inspired by John Oliver's [EventStore](https://github.com/NEventStore/NEventStore).

## Install

```
go get github.com/gokit/cqrskit/...
```

## Storage Supported

Below are the available database able to be used as the event store technologies:

- MongoDB
- BadgerDB (Planned)
- PostgreSQL (Planned)


## Publisher Supported

Below are the implemented queueing technology available for use 

- NATS 
- Amazon SQS
- NSQ (Planned)
- Redis (Planned)


## CLI Tooling

CQRSKit comes bundled with a command line code generation tool, that provides a means of avoiding the usage of reflect by generating
an appropriate method on structs annotated with `@escqrs` to handle different events types, based on methods that have the `HandlePrefix`, except in cases there such methods have a `@escqrs-skip` annotation in comments.

```go
//@escqrs
type User struct {
	Version  int
	Email    string
	Username string
}

type UserEmailUpdated struct {
	New string
}

func (u *User) HandleUserEmailUpdated(ev UserEmailUpdated) error {
	return nil
}

func (u *User) HandleUserNameUpdated(ev events.UserNameUpdated) error {
	return nil
}

//@escqrs-method-skip
func (u *User) HandleUserRackUpdated(ev UserEmailUpdated) error {
	return nil
}
```

Where the above produces the following after running `cqrskit generate` in terminal:

```go
var (
	// UserAggregateID represents the unique aggregate id for all events
	// related to the User type. It is the typeName hashed using a md5 sum.
	UserAggregateID = "f7091ac77d9b52a3ec5609891cd9f54f"
)

//*******************************************************************************
// User Event Applier
//*******************************************************************************

// Apply embodies the internal logic necessary to apply specific events to a User by
// calling appropriate methods.
func (u *User) Apply(evs cqrskit.EventCommit) error {
	for _, event := range evs.Events {
		switch ev := event.Data.(type) {
		case UserEmailUpdated:
			return u.HandleUserEmailUpdated(ev)
		case events.UserNameUpdated:
			return u.HandleUserNameUpdated(ev)

		}
	}
	return nil
}
```
