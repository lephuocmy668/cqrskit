package cqrskit

import (
	"context"
	"time"
)

//*******************************************************************************
// Event Type and Functions
//*******************************************************************************

// PendingDispatch embodies data stored by commit about undispatched commits
// which have being persisted into underline event store.
type PendingDispatch struct {
	DispatchID  string `json:"dispatch_id" bson:"dispatch_id" db:"dispatch_id"`
	CommitID    string `json:"commit_id" bson:"commit_id" db:"commit_id"`
	InstanceID  string `json:"instance_id" bson:"instance_id" db:"instance_id"`
	AggregateID string `json:"aggregate_id" bson:"aggregate_id" db:"aggregate_id"`
}

// CommitHeader embodies data stored within db about the commit of a event commit request.
type CommitHeader struct {
	Version     int       `json:"version" bson:"version" db:"version"`
	CommitID    string    `json:"commit_id" bson:"commit_id" db:"commit_id"`
	InstanceID  string    `json:"instance_id" bson:"instance_id" db:"instance_id"`
	AggregateID string    `json:"aggregate_id" bson:"aggregate_id" db:"aggregate_id"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp" db:"timestamp"`
}

// EventCommit embodies the collection of events that occurred from the execution
// of a specific command which produces a series of events which are then used
// used to regenerated said state.
type EventCommit struct {
	CommitID    string                 `json:"commit_id" bson:"commit_id" db:"commit_id"`
	InstanceID  string                 `json:"instance_id" bson:"instance_id" db:"instance_id"`
	AggregateID string                 `json:"aggregate_id" bson:"aggregate_id" db:"aggregate_id"`
	Version     int                    `json:"version" bson:"version" db:"version"`
	Command     string                 `json:"command" bson:"command" db:"command"`
	Header      map[string]interface{} `json:"header" bson:"header" db:"header"`
	Events      []Event                `json:"events" bson:"events" db:"events"`
	Created     time.Time              `json:"created" bson:"created" db:"created"`
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

//*******************************************************************************
// CQRS Repository Interface and Implementation
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
// Write Repository Interface and Implementation
//*******************************************************************************

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
// DispatchRepo Repository Interface and Implementation
//*******************************************************************************

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

//*******************************************************************************
// Read Repository Interface and Implementation
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
