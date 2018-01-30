package nats

import (
	"sync"

	"github.com/gokit/cqrskit"

	nats "github.com/nats-io/go-nats"
	gnats "github.com/nats-io/go-nats-streaming"
)

var (
	ackPool = sync.Pool{
		New: func() interface{} {
			return new(pendingRequest)
		},
	}
)

type pendingRequest struct {
	next   *pendingRequest
	commit cqrskit.EventCommit
	Ack    cqrskit.AckHandler
}

//*******************************************************************************
// NATS Publisher
//*******************************************************************************

// NATSPublisher implements the cqrskit.Publisher using nats as the
// underline transport.
type NATSPublisher struct {
	addr    string
	ops     []nats.Option
	encoder cqrskit.Encoder
	cl      sync.Mutex
	conn    *nats.Conn
}

// NATSPublisherFrom returns a new instance of NATSPublisher using the provided
// nats.Conn instance.
func NATSPublisherFrom(conn *nats.Conn, encoder cqrskit.Encoder) *NATSPublisher {
	return &NATSPublisher{
		encoder: encoder,
		addr:    conn.ConnectedUrl(),
	}
}

// NewNATSPublisher returns a new instance of NATSPublisher.
func NewNATSPublisher(addr string, encoder cqrskit.Encoder, ops ...nats.Option) *NATSPublisher {
	return &NATSPublisher{
		addr:    addr,
		ops:     ops,
		encoder: encoder,
	}
}

// Publish implements the cqrskit.Publisher and sends the event into the nats pubsub
// and calls acknowledged function if successful in both the encoding of the commit and
// it's delivery into nats event stream, else returning an error of why.
func (np *NATSPublisher) Publish(ns string, commit cqrskit.EventCommit, fn cqrskit.AckHandler) error {
	conn, err := np.getConnection()
	if err != nil {
		return err
	}

	commitBytes, err := np.encoder.Encode(commit)
	if err != nil {
		return err
	}

	if err := conn.Publish(ns, commitBytes); err != nil {
		return err
	}

	fn(cqrskit.PubAck{
		Namespace:   ns,
		Version:     commit.Version,
		CommitID:    commit.CommitID,
		InstanceID:  commit.InstanceID,
		AggregateID: commit.AggregateID,
	})

	return nil
}

// Close ends the underline nats connection.
func (np *NATSPublisher) Close() error {
	np.cl.Lock()
	defer np.cl.Unlock()

	if np.conn == nil {
		return nil
	}

	np.conn.Flush()
	np.conn.Close()
	np.conn = nil
	return nil
}

// getConnection returns the current nats client for delivery messages.
func (np *NATSPublisher) getConnection() (*nats.Conn, error) {
	np.cl.Lock()
	defer np.cl.Unlock()

	if np.conn != nil {
		return np.conn, nil
	}

	var err error
	np.conn, err = nats.Connect(np.addr, np.ops...)
	return np.conn, err
}

//*******************************************************************************
// NATS Streaming Publisher
//*******************************************************************************

type NATStreamingPublisher struct {
	addr       string
	clientID   string
	clusterID  string
	cl         sync.Mutex
	conn       gnats.Conn
	nativeConn *nats.Conn
	encoder    cqrskit.Encoder
	ops        []nats.Option
}

// NATStreamingPublisherFrom returns a new instance of NATStreamingPulisher using the provided
// nats.Conn address.
func NATStreamingPublisherFrom(clusterID string, clientID string, encoder cqrskit.Encoder, conn *nats.Conn) *NATStreamingPublisher {
	return &NATStreamingPublisher{
		nativeConn: conn,
		encoder:    encoder,
		clientID:   clientID,
		clusterID:  clusterID,
		addr:       conn.ConnectedUrl(),
	}
}

// NewNATStreamingPublisher returns a new instance of NATStreamingPulisher.
func NewNATStreamingPublisher(addr string, clusterID string, clientID string, encoder cqrskit.Encoder, ops ...nats.Option) *NATStreamingPublisher {
	return &NATStreamingPublisher{
		ops:       ops,
		addr:      addr,
		encoder:   encoder,
		clientID:  clientID,
		clusterID: clusterID,
	}
}

// Publish implements the cqrskit.Publisher and sends the event into the nats pubsub
// and calls acknowledged function if successful in both the encoding of the commit and
// it's delivery into nats event stream, else returning an error of why.
func (np *NATStreamingPublisher) Publish(ns string, commit cqrskit.EventCommit, fn cqrskit.AckHandler) error {
	conn, err := np.getConnection()
	if err != nil {
		return err
	}

	commitBytes, err := np.encoder.Encode(commit)
	if err != nil {
		return err
	}

	if err := conn.Publish(ns, commitBytes); err != nil {
		return err
	}

	fn(cqrskit.PubAck{
		Namespace:   ns,
		Version:     commit.Version,
		CommitID:    commit.CommitID,
		InstanceID:  commit.InstanceID,
		AggregateID: commit.AggregateID,
	})

	return nil
}

// Close ends the underline nats connection.
func (np *NATStreamingPublisher) Close() error {
	np.cl.Lock()
	defer np.cl.Unlock()

	if np.conn == nil {
		return nil
	}

	np.nativeConn.Flush()
	np.conn.Close()
	np.nativeConn.Close()
	np.conn = nil
	return nil
}

// getConnection returns the current nats client for delivery messages.
func (np *NATStreamingPublisher) getConnection() (gnats.Conn, error) {
	if np.conn != nil {
		return np.conn, nil
	}

	nativeConn, err := np.getNativeConnection()
	if err != nil {
		return nil, err
	}

	np.conn, err = gnats.Connect(np.clusterID, np.clientID, gnats.NatsConn(nativeConn))
	return np.conn, err
}

// getNativeConnection returns underline nats client nats.Conn.
func (np *NATStreamingPublisher) getNativeConnection() (*nats.Conn, error) {
	if np.nativeConn != nil {
		return np.nativeConn, nil
	}

	var err error
	np.nativeConn, err = nats.Connect(np.addr, np.ops...)
	return np.nativeConn, err
}
