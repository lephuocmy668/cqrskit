// Package mgo implements a mongo repository store defined by cqrskit. It allows usage of
// mongo has a read and write repository for aggregates.
//
//@mongo
package mgorp

import (
	"context"
	"errors"
	"time"

	"github.com/gokit/cqrskit"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// errors ....
var (
	ErrInvalidAggregateID     = errors.New("invalid aggregate id")
	ErrInvalidInstanceID      = errors.New("invalid instance id")
	ErrNoCommitsYet           = errors.New("no commits has being made")
	ErrConcurrentWrites       = errors.New("concurrent write occured; version used")
	ErrDuplicateCommitRequest = errors.New("request commit id handld, duplicate request")
)

// consts values of aggregate collection names.
const (
	AggregateCollection             = "aggregates"
	AggregateModelCollection        = "aggregates_model"
	AggregateDispatchCollection     = "aggregates_model_event_dispatch"
	AggregateEventCommitCollection  = "aggregates_model_event_commits"
	AggregateCommitHeaderCollection = "aggregates_model_event_commit_header"
)

// MongoDB defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type MongoDB interface {
	New(isread bool) (*mgo.Database, *mgo.Session, error)
}

// CommitDispatchHeader embodies data stored by commit about undispatched commits
// which have being persisted into underline event store.
type CommitDispatchHeader struct {
	cqrskit.PendingDispatch
	ID bson.ObjectId `json:"id" bson:"id" db:"id"`
}

// CommitHeader embodies data stored within event store to handle transaction
// idempotency by providing a last-write-win guarantee for a giving EventCommit.
// How:
//  When a new Event commit occurs on the event store, first the CommitHeader is written
//  with the allotted incremented version from the last commit header, which then allows
//  said EventCommit to attempt to save it's itself into store, if commit version has not
//  being used by another EventCommit concurrently storing itself into db, then after saving
//  CommitHeader is updated with successfully EventCommit commit id.
//  Else if EventCommit failed due to version taking by another concurrent write with succesffuly
//  save, then commit must be denied. This ensures we maintain transaction writes even on underline
//  stores not supporting such.
type CommitHeader struct {
	cqrskit.CommitHeader
	ID bson.ObjectId `json:"_id" bson:"_id" db:"_id"`
}

// Aggregate embodies the data stored in a db to represent
// a unique model class or struct .eg a User or a Article.
// It represents the unique identifier for that type alone.
type Aggregate struct {
	Id          bson.ObjectId `bson:"_id"`
	AggregateID string        `bson:"aggregate_id"`
}

// AggregateModel embodies the data stored of a object of a
// Aggregate class instance, it is used to reference a record
// in db of the associated aggregate type identified by it's
// AggregateID.
type AggregateModel struct {
	Id           bson.ObjectId `bson:"_id"`
	InstanceID   string        `bson:"instance_id"`
	AggregatedID string        `bson:"aggregate_id"`
}

// MgoWriteMaser implements the cqrskit.WriteRepository interface exposing
// methods to have a direct writer for a giving aggregate and and instance.
type MgoWriteMaster struct {
	db MongoDB
}

// NewWriteMaster returns a new instance of MgoWriteMaster.
func NewWriteMaster(db MongoDB) MgoWriteMaster {
	return MgoWriteMaster{db: db}
}

// Writer attempts to retrieve aggregate WriteRepo instance for writing events for a giving
// instance of an aggregate model. If aggregate record does not exists, it will be created.
func (mw MgoWriteMaster) Writer(aggregateID string, instanceID string) (cqrskit.WriteRepo, error) {
	zdb, _, err := mw.db.New(true)
	if err != nil {
		return nil, err
	}

	zcol := zdb.C(AggregateCollection)
	icol := zdb.C(AggregateModelCollection)

	aggrQuery := bson.M{"aggregate_id": aggregateID}
	total, err := zcol.Find(aggrQuery).Count()
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	if total == 0 {
		return mw.new(aggregateID, instanceID)
	}

	instQuery := bson.M{"aggregate_id": aggregateID, "instance_id": instanceID}
	itotal, err := icol.Find(instQuery).Count()
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	if itotal == 0 {
		return mw.new(aggregateID, instanceID)
	}

	return &MgoWriteRepository{
		db:          mw.db,
		instanceID:  instanceID,
		aggregateID: aggregateID,
	}, nil
}

// New returns a new cqrskit.WriteRepo for a giving aggregateID and instanceID, if giving
// aggregate is not found then a record is created and same logic applies for the instance.
func (mw *MgoWriteMaster) new(aggregateID string, instanceID string) (cqrskit.WriteRepo, error) {
	zdb, zes, err := mw.db.New(false)
	if err != nil {
		return nil, err
	}

	defer zes.Close()

	zcol := zdb.C(AggregateCollection)

	// if we fail to get count or count is zero, then we have no aggregate record of such
	// so make one and appropriate set indexes.
	aggrQuery := bson.M{"aggregate_id": aggregateID}
	atotal, err := zcol.Find(aggrQuery).Count()
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	if atotal == 0 {
		if err := mw.createAggregate(aggregateID, zdb); err != nil {
			return nil, err
		}
	}

	icol := zdb.C(AggregateModelCollection)
	instQuery := bson.M{"aggregate_id": aggregateID, "instance_id": instanceID}
	itotal, err := icol.Find(instQuery).Count()
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	if itotal == 0 {
		if err := mw.createAggregateModel(aggregateID, instanceID, zdb); err != nil {
			return nil, err
		}
	}

	return &MgoWriteRepository{
		db:          mw.db,
		instanceID:  instanceID,
		aggregateID: aggregateID,
	}, nil
}

func (mw *MgoWriteMaster) createAggregateModel(aggregateID string, instanceID string, zdb *mgo.Database) error {
	icol := zdb.C(AggregateModelCollection)
	if err := icol.EnsureIndex(mgo.Index{
		Key:    []string{"instance_id"},
		Unique: true,
		Name:   "instance_index",
	}); err != nil {
		return err
	}

	if err := icol.EnsureIndex(mgo.Index{
		Key:  []string{"aggregate_id"},
		Name: "aggregate_id_index",
	}); err != nil {
		return err
	}

	var model AggregateModel
	model.Id = bson.NewObjectId()
	model.InstanceID = instanceID
	model.AggregatedID = aggregateID

	if err := icol.Insert(model); err != nil {
		return err
	}

	// Add index to event collection for aggregate model.
	cmCol := zdb.C(AggregateEventCommitCollection)
	if err := cmCol.EnsureIndex(mgo.Index{
		Key:    []string{"commit_id"},
		Unique: true,
		Name:   "commit_id",
	}); err != nil {
		return err
	}

	if err := cmCol.EnsureIndex(mgo.Index{
		Key:    []string{"version"},
		Unique: true,
		Name:   "version",
	}); err != nil {
		return err
	}

	if err := cmCol.EnsureIndex(mgo.Index{
		Key:  []string{"created"},
		Name: "created",
	}); err != nil {
		return err
	}

	if err := cmCol.EnsureIndex(mgo.Index{
		Key:  []string{"instance_id", "aggregate_id"},
		Name: "instance_aggregate_index",
	}); err != nil {
		return err
	}

	cmhCol := zdb.C(AggregateCommitHeaderCollection)
	if err := cmhCol.EnsureIndex(mgo.Index{
		Key:    []string{"commit_id"},
		Unique: true,
		Name:   "commit_id",
	}); err != nil {
		return err
	}

	if err := cmhCol.EnsureIndex(mgo.Index{
		Key:    []string{"version"},
		Unique: true,
		Name:   "version",
	}); err != nil {
		return err
	}

	if err := cmhCol.EnsureIndex(mgo.Index{
		Key:  []string{"timestamp"},
		Name: "timestamp",
	}); err != nil {
		return err
	}

	if err := cmhCol.EnsureIndex(mgo.Index{
		Key:  []string{"instance_id", "aggregate_id"},
		Name: "instance_aggregate_index",
	}); err != nil {
		return err
	}

	dhCol := zdb.C(AggregateDispatchCollection)
	if err := dhCol.EnsureIndex(mgo.Index{
		Key:    []string{"commit_id"},
		Unique: true,
		Name:   "commit_id",
	}); err != nil {
		return err
	}

	if err := dhCol.EnsureIndex(mgo.Index{
		Key:  []string{"instance_id", "aggregate_id"},
		Name: "instance_aggregate_index",
	}); err != nil {
		return err
	}

	return nil
}

func (mw *MgoWriteMaster) createAggregate(aggregateID string, zdb *mgo.Database) error {
	zcol := zdb.C(AggregateCollection)
	if err := zcol.EnsureIndex(mgo.Index{
		Key:    []string{"aggregate_id"},
		Unique: true,
		Name:   "aggregate_index",
	}); err != nil {
		return err
	}

	var aggr Aggregate
	aggr.Id = bson.NewObjectId()
	aggr.AggregateID = aggregateID

	if err := zcol.Insert(aggr); err != nil {
		return err
	}

	return nil
}

// MgoWriteRepository implements the cqrskit.WriteRepo
// using mongodb has the underline store.
type MgoWriteRepository struct {
	db          MongoDB
	aggregateID string
	instanceID  string
}

type lastVersion struct {
	Version        int64  `bson:"version"`
	LamportVersion string `bson:"lamport_version"`
}

// DeleteAll removes all record associated with giving event and returns total
// records of all event records removed.
func (mwr *MgoWriteRepository) DeleteAll(ctx context.Context) (int, error) {
	zdb, zes, err := mwr.db.New(false)
	if err != nil {
		return -1, err
	}

	defer zes.Close()

	lvQuery := bson.M{
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	mc := zdb.C(AggregateEventCommitCollection)
	info, err := mc.RemoveAll(lvQuery)
	if err != nil {
		return -1, err
	}

	if _, err = zdb.C(AggregateCommitHeaderCollection).RemoveAll(lvQuery); err != nil {
		return info.Removed, err
	}

	if _, err = zdb.C(AggregateDispatchCollection).RemoveAll(lvQuery); err != nil {
		return info.Removed, err
	}

	return info.Removed, nil
}

// Count returns total count of all commited events for giving aggregate and instance.
func (mwr *MgoWriteRepository) Count(ctx context.Context) (int, error) {
	zdb, _, err := mwr.db.New(true)
	if err != nil {
		return -1, err
	}

	mc := zdb.C(AggregateEventCommitCollection)

	lvQuery := bson.M{
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	total, err := mc.Find(lvQuery).Count()
	if err != nil {
		return -1, err
	}

	return total, nil
}

// LastCommitVersion returns the last event version successfully committed into the
// CommitHeaderCollection and returns it's version.
func (mwr *MgoWriteRepository) LastCommitVersion(ctx context.Context) (cqrskit.CommitHeader, error) {
	var header cqrskit.CommitHeader

	zdb, zes, err := mwr.db.New(true)
	if err != nil {
		return header, err
	}

	defer zes.Close()

	commitCollection := zdb.C(AggregateEventCommitCollection)

	lvQuery := bson.M{
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	if err := commitCollection.Find(lvQuery).Sort("-version").One(&header); err != nil {
		if err != mgo.ErrNotFound {
			return header, err
		}

		return header, ErrNoCommitsYet
	}

	return header, nil
}

// Write receives a EventCommitRequest and attempts to had the giving commit into store if
// the ff follows true:
// 1. Request CommitID has not being seen or handled before.
// 2. Request does not attempt to conflict with version that has already being taking.
// In each case an appropriate error is returned to indicate status of request.
func (mwr *MgoWriteRepository) Write(ctx context.Context, req cqrskit.EventCommitRequest) (cqrskit.CommitHeader, error) {
	zdb, zes, err := mwr.db.New(false)
	if err != nil {
		return cqrskit.CommitHeader{}, err
	}

	defer zes.Close()

	dispatchCollection := zdb.C(AggregateDispatchCollection)
	commitCollection := zdb.C(AggregateEventCommitCollection)
	commitHeaderCollection := zdb.C(AggregateCommitHeaderCollection)

	probeQuery := bson.M{
		"commit_id":    req.ID,
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	var header CommitHeader

	totalFound, err := commitCollection.Find(probeQuery).Count()
	if err != nil && err != mgo.ErrNotFound {
		return header.CommitHeader, err
	}

	if err == nil && totalFound == 1 {
		return header.CommitHeader, ErrDuplicateCommitRequest
	}

	// Attempt to get current leased header.
	leaseQuery := bson.M{
		"commit_id":    "",
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	var dispatchHeader CommitDispatchHeader

	if err := commitHeaderCollection.Find(leaseQuery).One(&header); err != nil {
		if err != mgo.ErrNotFound {
			return header.CommitHeader, err
		}

		// Get last version number and lease out next version number for self.
		lastHeader, err := mwr.LastCommitVersion(ctx)
		if err != nil && err != ErrNoCommitsYet {
			return header.CommitHeader, err
		}

		header.ID = bson.NewObjectId()
		header.CommitID = ""
		header.Version = lastHeader.Version + 1
		header.InstanceID = mwr.instanceID
		header.AggregateID = mwr.aggregateID

		dispatchHeader.ID = bson.NewObjectId()
		dispatchHeader.CommitID = ""
		dispatchHeader.InstanceID = mwr.instanceID
		dispatchHeader.AggregateID = mwr.aggregateID
		dispatchHeader.DispatchID = dispatchHeader.ID.Hex()

		// Register new lease into commit header, requesting weak lock on version.
		if err := commitHeaderCollection.Insert(bson.M{
			"_id":          header.ID,
			"commit_id":    header.CommitID,
			"version":      header.Version,
			"instance_id":  header.InstanceID,
			"aggregate_id": header.AggregateID,
		}); err != nil {
			return header.CommitHeader, err
		}

		if err := dispatchCollection.Insert(bson.M{
			"_id":          dispatchHeader.ID,
			"commit_id":    dispatchHeader.CommitID,
			"dispatch_id":  dispatchHeader.DispatchID,
			"instance_id":  dispatchHeader.InstanceID,
			"aggregate_id": dispatchHeader.AggregateID,
		}); err != nil {
			return header.CommitHeader, err
		}
	}

	total, err := dispatchCollection.Find(leaseQuery).Count()
	if err != nil && err != mgo.ErrNotFound {
		return header.CommitHeader, err
	}

	if total == 0 {
		dispatchHeader.ID = bson.NewObjectId()
		dispatchHeader.InstanceID = mwr.instanceID
		dispatchHeader.AggregateID = mwr.aggregateID
		dispatchHeader.DispatchID = dispatchHeader.ID.String()

		if err := dispatchCollection.Insert(bson.M{
			"_id":          dispatchHeader.ID,
			"commit_id":    dispatchHeader.CommitID,
			"dispatch_id":  dispatchHeader.DispatchID,
			"instance_id":  dispatchHeader.InstanceID,
			"aggregate_id": dispatchHeader.AggregateID,
		}); err != nil {
			return header.CommitHeader, err
		}
	}

	var eventCommit cqrskit.EventCommit
	eventCommit.CommitID = req.ID
	eventCommit.Events = req.Events
	eventCommit.Header = req.Header
	eventCommit.Command = req.Command
	eventCommit.Created = req.Created
	eventCommit.Version = header.Version
	eventCommit.InstanceID = mwr.instanceID
	eventCommit.AggregateID = mwr.aggregateID

	if err := commitCollection.Insert(eventCommit); err != nil {
		if lastErr, ok := err.(*mgo.LastError); ok {
			return header.CommitHeader, lastErr
		}

		if mgo.IsDup(err) {
			return header.CommitHeader, ErrConcurrentWrites
		}

		return header.CommitHeader, err
	}

	commited := time.Now()
	if err := commitHeaderCollection.UpdateId(header.ID, bson.M{
		"$set": bson.M{
			"timestamp": commited,
			"commit_id": eventCommit.CommitID,
		},
	}); err != nil {
		if lastErr, ok := err.(*mgo.LastError); ok {
			return header.CommitHeader, lastErr
		}

		return header.CommitHeader, err
	}

	if err := dispatchCollection.UpdateId(dispatchHeader.ID, bson.M{
		"$set": bson.M{
			"commit_id": eventCommit.CommitID,
		},
	}); err != nil {
		if lastErr, ok := err.(*mgo.LastError); ok {
			return header.CommitHeader, lastErr
		}

		return header.CommitHeader, err
	}

	header.Timestamp = commited
	header.CommitID = eventCommit.CommitID

	return header.CommitHeader, nil
}

// MgoReadMaser implements the cqrskit.ReadRepository interface exposing
// methods to have a direct reader for a giving aggregate and related events instance.
type MgoReadMaster struct {
	db MongoDB
}

// NewReadMaster returns a new instance of MgoReadMaster.
func NewReadMaster(db MongoDB) MgoReadMaster {
	return MgoReadMaster{db: db}
}

// Reader returns a new store reader which provides methods for reading commited events from
// underline mongo store.
func (mgr MgoReadMaster) Reader(aggregateID string, instanceID string) (cqrskit.ReadRepo, error) {
	return &MgoReadRepository{
		db:          mgr.db,
		instanceID:  instanceID,
		aggregateID: aggregateID,
	}, nil
}

// MgoReadRepository implements the cqrskit.ReadRepo
// using mongodb has the underline store.
type MgoReadRepository struct {
	db          MongoDB
	aggregateID string
	instanceID  string
}

// CountBatches returns total number of event batches saved, with total events across all batches
// available within db.
func (mrr *MgoReadRepository) CountCommits(ctx context.Context) (int, error) {
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return -1, zerr
	}

	fnQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
	}

	zcol := zdb.C(AggregateEventCommitCollection)
	return zcol.Find(fnQuery).Count()
}

// ReadAll returns all events for giving aggregate and events for aggregate model.
func (mrr *MgoReadRepository) ReadAll(ctx context.Context) ([]cqrskit.EventCommit, error) {
	var events []cqrskit.EventCommit
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	zcol := zdb.C(AggregateEventCommitCollection)
	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
	}

	if err := zcol.Find(rmQuery).Sort("version").All(&events); err != nil {
		return nil, err
	}

	return events, nil
}

// ReadVersion returns all events for giving aggregate and instance model for requested version
// if found.
func (mrr *MgoReadRepository) ReadVersion(ctx context.Context, version int64) (cqrskit.EventCommit, error) {
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return cqrskit.EventCommit{}, zerr
	}

	zcol := zdb.C(AggregateEventCommitCollection)
	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
		"version":      version,
	}

	var commit cqrskit.EventCommit
	if err := zcol.Find(rmQuery).One(&commit); err != nil {
		return cqrskit.EventCommit{}, err
	}

	return commit, nil
}

// ReadSinceCount returns all available events stored within mongdb, returning desired count from last saved
// upward. If count is below zero that is -1, then all records are returned in reverse, that is in descending
// order of version number.
func (mrr *MgoReadRepository) ReadSinceCount(ctx context.Context, count int) ([]cqrskit.EventCommit, error) {
	var events []cqrskit.EventCommit
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	zcol := zdb.C(AggregateEventCommitCollection)
	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
	}

	if count > 0 {
		if err := zcol.Find(rmQuery).Limit(count).Sort("version").All(&events); err != nil {
			return nil, err
		}
	} else {
		if err := zcol.Find(rmQuery).Sort("version").All(&events); err != nil {
			return nil, err
		}
	}

	return events, nil
}

// ReadSinceVersion returns all events that have occured around giving version and upwards.
func (mrr *MgoReadRepository) ReadSinceVersion(ctx context.Context, version int64, limit int) ([]cqrskit.EventCommit, error) {
	var events []cqrskit.EventCommit
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	zcol := zdb.C(AggregateEventCommitCollection)
	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
		"version":      bson.M{"$gte": version},
	}

	if limit > 0 {
		if err := zcol.Find(rmQuery).Limit(limit).Sort("version").All(&events); err != nil {
			return nil, err
		}
	} else {
		if err := zcol.Find(rmQuery).Sort("version").All(&events); err != nil {
			return nil, err
		}
	}

	return events, nil
}

// ReadSinceTime returns all events for giving aggregate and events for aggregate model for the time of creation
// of event upwards till required limit. Limit of -1 returns all events from said time.
func (mrr *MgoReadRepository) ReadSinceTime(ctx context.Context, ts time.Time, limit int) ([]cqrskit.EventCommit, error) {
	var events []cqrskit.EventCommit
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	zcol := zdb.C(AggregateEventCommitCollection)
	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
		"created":      ts,
	}

	if limit > 0 {
		if err := zcol.Find(rmQuery).Limit(limit).All(&events); err != nil {
			return nil, err
		}
	} else {
		if err := zcol.Find(rmQuery).All(&events); err != nil {
			return nil, err
		}
	}

	return events, nil
}
