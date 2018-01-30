package sqs_test

import (
	"testing"

	"github.com/influx6/faux/tests"

	"github.com/gokit/cqrskit"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	pubsqs "github.com/gokit/cqrskit/publishers/sqs"
)

func TestSQSPublisher(t *testing.T) {
	mockSVC := mockSQSClient{
		actions: make(chan struct{}, 0),
	}

	publisher := pubsqs.NewSQSPublisher(func(region string) (sqsiface.SQSAPI, error) {
		return mockSVC, nil
	}, cqrskit.JSONEncoder{})

	if err := publisher.AddSQSRegion(
		"bomb.events",
		"http://sqs.us-east-2.amazonaws.com/123456789012/MyQueue",
	); err != nil {
		tests.FailedWithError(err, "Should have successfully added new queue region")
	}
	tests.Passed("Should have successfully added new quue region")

	if err := publisher.AddSQSRegion(
		"user.events",
		"http://sqs.us-east-2.amazonaws.com/123456789012/MyQueue",
	); err != nil {
		tests.FailedWithError(err, "Should have successfully added new queue region")
	}
	tests.Passed("Should have successfully added new queue region")

	if err := publisher.Publish("bomb.events", cqrskit.EventCommit{}, func(ack cqrskit.PubAck) {}); err != nil {
		tests.FailedWithError(err, "Should have successfully published event commit")
	}
	tests.Passed("Should have successfully published event commit")

	<-mockSVC.actions

	if err := publisher.Publish("user.events", cqrskit.EventCommit{}, func(ack cqrskit.PubAck) {}); err != nil {
		tests.FailedWithError(err, "Should have successfully published event commit")
	}
	tests.Passed("Should have successfully published event commit")

	<-mockSVC.actions
}

type mockSQSClient struct {
	sqsiface.SQSAPI
	actions chan struct{}
}

func (m mockSQSClient) SendMessage(s *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	go func() { m.actions <- struct{}{} }()
	return &sqs.SendMessageOutput{}, nil
}

func (m mockSQSClient) ReceiveMessage(s *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	out := &sqs.ReceiveMessageOutput{}
	return out, nil
}
