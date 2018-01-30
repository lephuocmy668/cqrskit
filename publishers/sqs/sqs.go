package sqs

import (
	"errors"
	"strings"
	"sync"

	"github.com/gokit/cqrskit"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// errors ...
var (
	ErrTargetNameAssigned  = errors.New("target name already assigned")
	ErrURLHasNoRegion      = errors.New("sqs url has no region, possible invalid")
	ErrNoRegionWithTargget = errors.New("target name has no associated sqs region registered")
)

//*******************************************************************************
// Utils
//*******************************************************************************

// RegionFromURL parses an sqs url and returns the aws region
func RegionFromURL(url string) string {
	pieces := strings.Split(url, ".")
	if len(pieces) > 2 {
		return pieces[1]
	}

	return ""
}

//*******************************************************************************
// Amazon SQS Publisher
//*******************************************************************************

// NewServiceFunc defines a function type which is used to create a SQSAPI service
// instance for delivery events.
type NewServiceFunc func(region string) (sqsiface.SQSAPI, error)

// DefaultNewServiceFunc defines a default version of the NewServiceFunc type for
// using in creating a sqs publisher.
func DefaultNewServiceFunc(region string) (sqsiface.SQSAPI, error) {
	awsConfig := aws.NewConfig().WithRegion(region)
	s, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	return sqs.New(s), nil
}

// sqsRegion identifies a giving region and associated sqs service.
type sqsRegion struct {
	URL     string
	Service sqsiface.SQSAPI
}

// SQSPublisher implements the cqrskit.Publisher for using amazon SQS queue has a
// means of publishing event commits.
type SQSPublisher struct {
	encoder    cqrskit.Encoder
	NewService NewServiceFunc
	rl         sync.Mutex
	regions    map[string]sqsRegion
}

// New returns a new instance of a SQSPublisher using the DefaultNewServiceFunc has
// the new service generator.
func New(encoder cqrskit.Encoder) *SQSPublisher {
	return NewSQSPublisher(DefaultNewServiceFunc, encoder)
}

// NewSQSPublisher returns a new instance of the SQSPublisher using provided
// NewServiceFunc function.
func NewSQSPublisher(newService NewServiceFunc, encoder cqrskit.Encoder) *SQSPublisher {
	return &SQSPublisher{
		encoder:    encoder,
		NewService: newService,
		regions:    make(map[string]sqsRegion),
	}
}

// AddSQSRegion adds a region URL into the queue dictionary, this allows you to refer
// to said region url through it's targetName when publisher. This allows you to
// easily associated different event name to same or different sqs queue region urls.
// eg. url like http://sqs.us-east-2.amazonaws.com/123456789012/MyQueue
func (np *SQSPublisher) AddSQSRegion(targetName string, queueURL string) error {
	np.rl.Lock()
	defer np.rl.Unlock()

	if _, ok := np.regions[targetName]; ok {
		return ErrTargetNameAssigned
	}

	region := RegionFromURL(queueURL)
	if region == "" {
		return ErrURLHasNoRegion
	}

	svc, err := np.NewService(region)
	if err != nil {
		return err
	}

	np.regions[targetName] = sqsRegion{
		Service: svc,
		URL:     queueURL,
	}

	return nil
}

// getSQSRegion returns associated sqsRegion value pointed to by targetName.
func (np *SQSPublisher) getSQSRegion(targetName string) (sqsRegion, error) {
	np.rl.Lock()
	defer np.rl.Unlock()

	region, ok := np.regions[targetName]
	if !ok {
		return sqsRegion{}, ErrNoRegionWithTargget
	}

	return region, nil
}

// Publish implements the cqrskit.Publisher interface and sends a giving
func (np *SQSPublisher) Publish(targetName string, commit cqrskit.EventCommit, fn cqrskit.AckHandler) error {
	region, err := np.getSQSRegion(targetName)
	if err != nil {
		return err
	}

	commitBytes, err := np.encoder.Encode(commit)
	if err != nil {
		return err
	}

	var message sqs.SendMessageInput
	message.QueueUrl = aws.String(region.URL)
	message.MessageBody = aws.String(string(commitBytes))

	output, err := region.Service.SendMessage(&message)
	if err != nil {
		return err
	}

	fn(cqrskit.PubAck{
		Response:    output,
		Namespace:   region.URL,
		Version:     commit.Version,
		CommitID:    commit.CommitID,
		InstanceID:  commit.InstanceID,
		AggregateID: commit.AggregateID,
	})

	return nil
}
