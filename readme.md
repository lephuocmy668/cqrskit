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



