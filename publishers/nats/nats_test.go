package nats_test

import (
	"testing"

	"github.com/gokit/cqrskit"
	"github.com/influx6/faux/tests"

	pubnats "github.com/gokit/cqrskit/publishers/nats"
)

const (
	defaultURL = "nats://0.0.0.0:4222"
)

func TestNATSPublisher(t *testing.T) {
	publisher := pubnats.NewNATSPublisher(defaultURL, cqrskit.JSONEncoder{})
	defer publisher.Close()

	if err := publisher.Publish("users.events", cqrskit.EventCommit{}, func(ack cqrskit.PubAck) {}); err != nil {
		tests.FailedWithError(err, "Should have successfully published event commit")
	}
	tests.Passed("Should have successfully published event commit")
}

func TestNATStreamingPublisher(t *testing.T) {
	publisher := pubnats.NewNATStreamingPublisher(defaultURL, "test-daddy", "test-child", cqrskit.JSONEncoder{})
	defer publisher.Close()

	if err := publisher.Publish("users.events", cqrskit.EventCommit{}, func(ack cqrskit.PubAck) {}); err != nil {
		tests.FailedWithError(err, "Should have successfully published event commit")
	}
	tests.Passed("Should have successfully published event commit")
}
