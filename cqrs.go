package cqrskit

import (
	"context"
	"time"
)

//*******************************************************************************
// Message Type
//*******************************************************************************

// Message embodies a generalized data used to contain to an incoming message
// or command.
type Message struct {
	Meta        interface{} `json:"meta"`
	Type        string      `json:"type"`
	Version     int         `json:"version"`
	Payload     interface{} `json:"payload"`
	AggregateID string      `json:"aggregate_id,omitempty"`
	InstanceID  string      `json:"instance_id,omitempty"`
}

//*******************************************************************************
// Saga Types
//*******************************************************************************

// SagaAction defines a function type which is called to by a Saga once
// its ready to deliver events to be commited and/or messages to be sent.
type SagaAction func([]Message, []EventCommit)

// Saga defines an interface which embodies a means for handling multistage
// processes. It acts as a abstraction for a process manager.
//
// With Sagas, every transition can lead to a request/command to be executed
// by another service or a new event to be stored and transmitted, hence every
// transition should have a giving or provided function if desired which provides
// a means to receive the commands/message and events which are to be translated.
// Even if a Saga, will only ever release messages and events after all it's state
// has being true.
type Saga interface {
	Transition(Message, SagaAction) error
}

// Transitions defines a type which transits its internal state due to a provided
// Message.
type Transitions interface {
	Transition(Message) error
}

// StaticSaga returns a type implementing the Transitions interface, where the provided
// Saga and associated SagaAction are used to perform the saga action.
func StaticSaga(saga Saga, action SagaAction) Transitions {
	return staticSaga{action: action, saga: saga}
}

// staticSage implements the Transitions interface, which wraps a Saga and Action
// which will be called every time a giving Message arrives.
type staticSaga struct {
	action SagaAction
	saga   Saga
}

// Transition implements the Transitions.Transition method which calls
// internal Saga with Message and provided Action.
func (ss staticSaga) Transition(m Message) error {
	return ss.saga.Transition(m, ss.action)
}

//*******************************************************************************
// Snapshot Type
//*******************************************************************************

// Snapshot embodies data representing a generated snapshot from a giving
// series of event-commits.
type Snapshot struct {
	Meta        interface{}            `json:"meta" bson:"meta" db:"meta"`
	SnapID      string                 `json:"snap_id" bson:"snap_id" db:"snap_id"`
	Revision    int                    `json:"revision" bson:"revision" db:"revision"`
	FromVersion int                    `json:"from_version" bson:"from_version" db:"from_version"`
	ToVersion   int                    `json:"to_version" bson:"to_version" db:"to_version"`
	InstanceID  string                 `json:"instance_id" bson:"instance_id" db:"instance_id"`
	AggregateID string                 `json:"aggregate_id" bson:"aggregate_id" db:"aggregate_id"`
	Payload     interface{}            `json:"payload" bson:"payload" db:"payload"`
	Header      map[string]interface{} `json:"header" bson:"header" db:"header"`
}

// SnapshotWriter defines a interface which expose a method to get a snapshot
// writer.
type SnapshotWriterRepository interface {
	Writer(aggregationID string, instanceID string) (SnapshotWriter, error)
}

// SnapshotReader defines a interface which expose a method to get a snapshot
// reader.
type SnapshotReaderRepoitory interface {
	Reader(aggregationID string, instanceID string) (SnapshotReader, error)
}

// SnapshotRepository defines a type which provides methods to gain access to
// snapshot readers and writers for a giving aggregate and instance.
type SnapshotRepository interface {
	SnapshotReaderRepoitory
	SnapshotWriterRepository
}

// SnapshotReader defines an interface that exposes a means to read snapshot details
// from a underline store.
type SnapshotReader interface {
	ReadAll(context.Context) ([]Snapshot, error)
	ReadID(context.Context, string) (Snapshot, error)
	ReadRevision(context.Context, int) (Snapshot, error)
	ReadVersion(ctx context.Context, fromVersion int, toVersion int) ([]Snapshot, error)
}

// SnapshotWriter defines an interface that exposes a means methods to write
// new snapshots into a underline store.
type SnapshotWriter interface {
	Write(context.Context, Snapshot) error
	Rewrite(ctx context.Context, revision int, snap Snapshot) error
}

//*******************************************************************************
// Event Type
//*******************************************************************************

// EventCommit embodies the collection of events that occurred from the execution
// of a specific command which produces a series of events which are then used
// used to regenerated said state.
type EventCommit struct {
	CommitID    string                 `json:"commit_id" bson:"commit_id" db:"commit_id"`
	InstanceID  string                 `json:"instance_id" bson:"instance_id" db:"instance_id"`
	AggregateID string                 `json:"aggregate_id" bson:"aggregate_id" db:"aggregate_id"`
	Version     int                    `json:"version" bson:"version" db:"version"`
	Command     string                 `json:"command" bson:"command" db:"command"`
	Created     time.Time              `json:"created" bson:"created" db:"created"`
	Events      []Event                `json:"events" bson:"events" db:"events"`
	Header      map[string]interface{} `json:"header" bson:"header" db:"header"`
}

// Event embodies the data stored correlating with a giving event
// that occurred on a giving aggregate model.
type Event struct {
	ID     string                 `json:"id" bson:"id" db:"created"`
	Type   string                 `json:"type" bson:"type" db:"type"`
	Meta   interface{}            `json:"meta" bson:"meta" db:"meta"`
	Data   interface{}            `json:"data" bson:"data" db:"data"`
	Header map[string]interface{} `json:"header" bson:"header" db:"header"`
}

//*******************************************************************************
// Encoding and Decoding Types
//*******************************************************************************

// Encoder defines a type which embodies the conversion or serialization of a event
// commit into a byte slice.
type Encoder interface {
	Encode(EventCommit) ([]byte, error)
}

// Decoder defines a type which embodies the deserialization of a byte slice into
// a EventCommit into a byte slice.
type Decoder interface {
	Decode([]byte) (EventCommit, error)
}

//*******************************************************************************
// CQRS Repository Interface
//*******************************************************************************

// CQRSEventStore defines a central interface which is exposed by an implementing
// event store which provides access to readers, writers and dispatchers for
// storing and retrieving records for all event stored within.
type CQRSEventStore interface {
	ReadRepository
	WriteRepository
	DispatchRepository
}

//*******************************************************************************
// Write Repository Interface
//*******************************************************************************

// CommitHeader embodies data stored within db about the commit of a event commit request.
type CommitHeader struct {
	Version     int       `json:"version" bson:"version" db:"version"`
	CommitID    string    `json:"commit_id" bson:"commit_id" db:"commit_id"`
	InstanceID  string    `json:"instance_id" bson:"instance_id" db:"instance_id"`
	AggregateID string    `json:"aggregate_id" bson:"aggregate_id" db:"aggregate_id"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp" db:"timestamp"`
}

// EventCommitRequest embodies the data sent into the store to have a set of
// events committed for a giving aggregated and instance.
// All ID values must be UUID and must commit from client of request
// and not any other valid to ensure and maximize idempotent insertions and
// de-duplication of event commit requests.
type EventCommitRequest struct {
	ID      string
	Command string
	Events  []Event
	Created time.Time
	Header  map[string]interface{}
}

// WriteRepository defines the interface which provides
// a single method to retrieve a WriteRepository which
// stores all events for a particular  identified by it's instanceID.
type WriteRepository interface {
	Writer(aggregationID string, instanceID string) (WriteRepo, error)
}

// WriteRepo embodies a repository which houses the store
// of events for giving type .
type WriteRepo interface {
	Count(context.Context) (int, error)
	LastCommitVersion(context.Context) (CommitHeader, error)
	Write(context.Context, EventCommitRequest) (CommitHeader, error)
}

//*******************************************************************************
// Read Repository Interface
//*******************************************************************************

// ReadRepository defines the interface which provides
// a single method to retrieve a ReadRepos to read
// events that occur for a giving type through
// it's instanceID which identifies that records events.
type ReadRepository interface {
	Reader(aggregationID string, instanceID string) (ReadRepo, error)
}

// ReadRepo embodies a repository which reads the store
// of events for giving type , returning an Applier
// to apply said events to target.
type ReadRepo interface {
	CountCommits(context.Context) (int, error)
	ReadAll(context.Context) ([]EventCommit, error)
	ReadVersion(ctx context.Context, version int64) (EventCommit, error)
	ReadSinceCount(ctx context.Context, count int) ([]EventCommit, error)
	ReadSinceTime(ctx context.Context, last time.Time, limit int) ([]EventCommit, error)
	ReadSinceVersion(ctx context.Context, version int64, limit int) ([]EventCommit, error)
}

//*******************************************************************************
// Publisher Repository Interface
//*******************************************************************************

// PubAck defines the data used for responding to a publish request
type PubAck struct {
	Version     int         `json:"version"`
	Namespace   string      `json:"namespace"`
	CommitID    string      `json:"commit_id"`
	InstanceID  string      `json:"instance_id"`
	AggregateID string      `json:"aggregate_id"`
	Response    interface{} `json:"response"`
}

// AckHandler defines a function type used to received a PubAck acknowledge
// response for the publishing of an EventCommit.
type AckHandler func(ack PubAck)

// Publisher defines an interface which defines the implementation to be done
// for the publishing of a committed EventCommit using a desired namespace or tag.
// It's expects the Publish method returns an error if the giving EventCommit failed
// to be pushed into the underline queue else will call the handler once said request
// is added successfully into the queue.
type Publisher interface {
	Publish(string, EventCommit, AckHandler) error
}

//*******************************************************************************
// DispatchRepo Repository Interface
//*******************************************************************************

// PendingDispatch embodies data stored by commit about undispatched commits
// which have being persisted into underline event store.
type PendingDispatch struct {
	DispatchID  string `json:"dispatch_id" bson:"dispatch_id" db:"dispatch_id"`
	CommitID    string `json:"commit_id" bson:"commit_id" db:"commit_id"`
	InstanceID  string `json:"instance_id" bson:"instance_id" db:"instance_id"`
	AggregateID string `json:"aggregate_id" bson:"aggregate_id" db:"aggregate_id"`
}

// DispatchRepository exposes a interface which defines a mechanism for
// implementers to present meta-details related to the dispatch state of
// giving events within a event store.
type DispatchRepository interface {
	Dispatcher(aggregationID string, instanceID string) (DispatchRepo, error)
}

// DispatchRepo defines the interface representing the dispatch tables
// for a giving aggregate and instance type, it provides access to all
// dispatched and non-dispatched EventCommits, allowing the marking of
// non-dispatched as dispatched.
type DispatchRepo interface {
	Dispatch(ctx context.Context, id string) error
	Undispatched(context.Context) ([]PendingDispatch, error)
}
