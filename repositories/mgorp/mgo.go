// Package mgo implements a mongo repository store defined by cqrskit. It allows usage of
// mongo has a read and write repository for aggregates.
//
//@mongo
package mgorp

import (
	"context"
	"time"

	"github.com/gokit/cqrskit"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// consts values of aggregate collection names.
const (
	AggregateCollection      = "aggregates"
	AggregateModelCollection = "aggregates_model"
	AggregateEventCollection = "aggregates_model_events"
)

// MongoDB defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type MongoDB interface {
	New(isread bool) (*mgo.Database, *mgo.Session, error)
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

// AggregateEvent embodies the data saved into the db, it batches
// events together to minimize issues with races and transactional
// consistency when attempting to insert multiple inverts in serial
// at once. It helps minimize consistency issues with transactions
// by batching event saves though at the cost of granularity.
type AggregateEvent struct {
	Version        int64           `bson:"version"`
	InstanceID     string          `bson:"instance_id"`
	AggregatedID   string          `bson:"aggregate_id"`
	LamportVersion string          `bson:"lamport_version"`
	Count          int64           `bson:"count"`
	Created        time.Time       `bson:"created"`
	Events         []cqrskit.Event `bson:"events"`
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

// Get attempts to retrieve aggregate WriteRepo instance for writing events for a giving
// instance of an aggregate model. If aggregate record does not exists, it will be created.
func (mw MgoWriteMaster) Get(aggregateID string, instanceID string) (cqrskit.WriteRepo, error) {
	zdb, _, err := mw.db.New(true)
	if err != nil {
		return nil, err
	}

	zcol := zdb.C(AggregateCollection)

	aggrQuery := bson.M{"aggregate_id": aggregateID}
	if total, err := zcol.Find(aggrQuery).Count(); err != nil || total == 0 {
		return mw.New(aggregateID, instanceID)
	}

	return &MgoWriteRepository{
		db:          mw.db,
		instanceID:  instanceID,
		aggregateID: aggregateID,
	}, nil
}

// New returns a new cqrskit.WriteRepo for a giving aggregateID and instanceID, if giving
// aggregate is not found then a record is created and same logic applies for the instance.
func (mw *MgoWriteMaster) New(aggregateID string, instanceID string) (cqrskit.WriteRepo, error) {
	zdb, zes, err := mw.db.New(false)
	if err != nil {
		return nil, err
	}

	defer zes.Close()

	// if we fail to get count or count is zero, then we have no aggregate record of such
	// so make one and appropriate set indexes.
	zcol := zdb.C(AggregateCollection)
	aggrQuery := bson.M{"aggregate_id": aggregateID}
	if total, err := zcol.Find(aggrQuery).Count(); err != nil || total == 0 {
		if err := zcol.EnsureIndex(mgo.Index{
			Key:    []string{"aggregate_id"},
			Unique: true,
			Name:   "aggregate_index",
		}); err != nil {
			return nil, err
		}

		var aggr Aggregate
		aggr.Id = bson.NewObjectId()
		aggr.AggregateID = aggregateID

		if err := zcol.Insert(aggr); err != nil {
			return nil, err
		}
	}

	icol := zdb.C(AggregateModelCollection)
	instQuery := bson.M{"aggregate_id": aggregateID, "instance_id": instanceID}
	if total, err := icol.Find(instQuery).Count(); err != nil || total == 0 {
		if err := icol.EnsureIndex(mgo.Index{
			Key:    []string{"instance_id"},
			Unique: true,
			Name:   "instance_index",
		}); err != nil {
			return nil, err
		}

		if err := icol.EnsureIndex(mgo.Index{
			Key:  []string{"aggregate_id"},
			Name: "aggregate_id_index",
		}); err != nil {
			return nil, err
		}

		var model AggregateModel
		model.Id = bson.NewObjectId()
		model.InstanceID = instanceID
		model.AggregatedID = aggregateID

		if err := icol.Insert(model); err != nil {
			return nil, err
		}

		// Add index to event collection for aggregate model.
		ecol := zdb.C(AggregateEventCollection)
		if err := ecol.EnsureIndex(mgo.Index{
			Key:    []string{"version"},
			Unique: true,
			Name:   "version",
		}); err != nil {
			return nil, err
		}

		if err := ecol.EnsureIndex(mgo.Index{
			Key:  []string{"instance_id", "aggregate_id"},
			Name: "instance_aggregate_index",
		}); err != nil {
			return nil, err
		}
	}

	return &MgoWriteRepository{
		db:          mw.db,
		instanceID:  instanceID,
		aggregateID: aggregateID,
	}, nil
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

	mc := zdb.C(AggregateEventCollection)

	lvQuery := bson.M{
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	info, err := mc.RemoveAll(lvQuery)
	if err != nil {
		return -1, err
	}

	return info.Removed, nil
}

// DeleteVersion removes event record associated with giving event version.
func (mwr *MgoWriteRepository) DeleteVersion(ctx context.Context, version int) error {
	zdb, _, err := mwr.db.New(true)
	if err != nil {
		return err
	}

	mc := zdb.C(AggregateEventCollection)

	lvQuery := bson.M{
		"version":      version,
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	return mc.Remove(lvQuery)
}

// Count returns total count of all records of events in db.
func (mwr *MgoWriteRepository) Count(ctx context.Context) (int, error) {
	zdb, _, err := mwr.db.New(true)
	if err != nil {
		return -1, err
	}

	mc := zdb.C(AggregateEventCollection)

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

// Save attempts to save slice of events as a single batch instance sacrificing a little
// granularity for transaction safety it ensures to consistently first retrieve last transaction
// version before attempting to insert.
func (mwr *MgoWriteRepository) SaveEvents(ctx context.Context, events []cqrskit.Event) error {
	if len(events) == 0 {
		return nil
	}

	zdb, zes, err := mwr.db.New(false)
	if err != nil {
		return err
	}

	defer zes.Close()

	mc := zdb.C(AggregateEventCollection)

	// Get last version sequence from last committed events.
	var last lastVersion

	lvQuery := bson.M{
		"aggregate_id": mwr.aggregateID,
		"instance_id":  mwr.instanceID,
	}

	if err := mc.Find(lvQuery).Sort("-version").One(&last); err != nil {
		if err != mgo.ErrNotFound {
			return err
		}

		last.Version = 0
	}

	initial := events[0]

	var newEvent AggregateEvent
	newEvent.Events = events
	newEvent.Created = initial.Created
	newEvent.Version = last.Version + 1
	newEvent.Count = int64(len(events))
	newEvent.InstanceID = mwr.instanceID
	newEvent.AggregatedID = mwr.aggregateID

	if err := mc.Insert(newEvent); err != nil {
		return err
	}

	return nil
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

// Get returns a reader which reads all events stored within a mongodb database based on
// specific criteria.
func (mgr MgoReadMaster) Get(aggregateID string, instanceID string) (cqrskit.ReadRepo, error) {
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

type totalEvents struct {
	Total int `bson:"total"`
}

// CountBatches returns total number of event batches saved, with total events across all batches
// available within db.
func (mrr *MgoReadRepository) CountBatches(ctx context.Context) (batch int, total int, err error) {
	zdb, _, zerr := mrr.db.New(true)
	if err != nil {
		err = zerr
		return
	}

	zcol := zdb.C(AggregateEventCollection)
	batch, err = zcol.Count()
	if err != nil {
		return
	}

	pipeline := zcol.Pipe([]bson.M{
		{
			"$match": bson.M{"$and": []bson.M{{"aggregate_id": mrr.aggregateID}, {"instance_id": mrr.instanceID}}},
		},
		{
			"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$count"}},
		},
	})

	var te totalEvents
	if err = pipeline.One(&te); err != nil {
		return
	}

	total = te.Total
	return
}

// CountBatchesFromVersion returns total number of event batches and total batches from version number.
func (mrr *MgoReadRepository) CountBatchesFromVersion(ctx context.Context, version int64, limit int) (batch int, total int, err error) {
	zdb, _, zerr := mrr.db.New(true)
	if err != nil {
		err = zerr
		return
	}

	zcol := zdb.C(AggregateEventCollection)
	batch, err = zcol.Count()
	if err != nil {
		return
	}

	pipelines := make([]bson.M, 0, 4)
	pipelines = append(pipelines, bson.M{
		"$match": bson.M{
			"$and": []bson.M{
				{"aggregate_id": mrr.aggregateID},
				{"instance_id": mrr.instanceID},
				{"version": bson.M{"$gte": version}},
			},
		},
	})

	pipelines = append(pipelines, bson.M{
		"$sort": bson.M{"version": 1},
	})

	if limit > 0 {
		pipelines = append(pipelines, bson.M{
			"$limit": limit,
		})
	}

	pipelines = append(pipelines, bson.M{
		"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$count"}},
	})

	pipeline := zcol.Pipe(pipelines)

	var te totalEvents
	if err = pipeline.One(&te); err != nil {
		return
	}

	total = te.Total
	return
}

// CountBatchesFromCount returns total number of event and batches after sorting and returns total from count starting
// from latest to provided count value.
func (mrr *MgoReadRepository) CountBatchesFromCount(ctx context.Context, count int) (batch int, total int, err error) {
	zdb, _, zerr := mrr.db.New(true)
	if err != nil {
		err = zerr
		return
	}

	zcol := zdb.C(AggregateEventCollection)
	batch, err = zcol.Count()
	if err != nil {
		return
	}

	pipelines := make([]bson.M, 0, 5)
	pipelines = append(pipelines, bson.M{
		"$match": bson.M{
			"$and": []bson.M{
				{"aggregate_id": mrr.aggregateID},
				{"instance_id": mrr.instanceID},
			},
		},
	})

	pipelines = append(pipelines, bson.M{
		"$sort": bson.M{"version": -1},
	})

	if count > 0 {
		pipelines = append(pipelines, bson.M{
			"$limit": count,
		})
	}

	pipelines = append(pipelines, bson.M{
		"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$count"}},
	})

	pipeline := zcol.Pipe(pipelines)

	var te totalEvents
	if err = pipeline.One(&te); err != nil {
		return
	}

	total = te.Total
	return
}

// CountBatchesFromTime returns total number of event and batches after sorting and returns total from time provided.
func (mrr *MgoReadRepository) CountBatchesFromTime(ctx context.Context, ts time.Time, limit int) (batch int, total int, err error) {
	zdb, _, zerr := mrr.db.New(true)
	if err != nil {
		err = zerr
		return
	}

	zcol := zdb.C(AggregateEventCollection)
	batch, err = zcol.Count()
	if err != nil {
		return
	}

	pipelines := make([]bson.M, 0, 4)
	pipelines = append(pipelines, bson.M{
		"$match": bson.M{
			"$and": []bson.M{
				{"aggregate_id": mrr.aggregateID},
				{"instance_id": mrr.instanceID},
				{"created": bson.M{"$gte": ts}},
			},
		},
	})

	pipelines = append(pipelines, bson.M{
		"$sort": bson.M{"version": 1},
	})

	if limit > 0 {
		pipelines = append(pipelines, bson.M{
			"$limit": limit,
		})
	}

	pipelines = append(pipelines, bson.M{
		"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$count"}},
	})

	pipeline := zcol.Pipe(pipelines[:len(pipelines)])

	var te totalEvents
	if err = pipeline.One(&te); err != nil {
		return
	}

	total = te.Total
	return
}

// ReadAll returns all events for giving aggregate and events for aggregate model.
// ReadAll uses the mongodb mgo.Iterator and will iterate through all records till
// it has full covered the total records or meets an error, if an error occured, then
// an error and the collected record slice is returned.
func (mrr *MgoReadRepository) ReadAll(ctx context.Context) ([]cqrskit.Event, error) {
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	_, totalEvents, err := mrr.CountBatches(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]cqrskit.Event, totalEvents)

	next := struct {
		Events []cqrskit.Event `bson:"events"`
	}{Events: events}

	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
	}

	zcol := zdb.C(AggregateEventCollection)
	itr := zcol.Find(rmQuery).Sort("version").Iter()

	for itr.Next(&next) {
		if err := itr.Err(); err != nil {
			return events, err
		}

		last := len(next.Events)
		next.Events = events[last:]
	}

	if err = itr.Close(); err != nil {
		return events, err
	}

	return events, nil
}

// ReadFromLastCount returns all events from the last saved to the total count which acts as a limit.
// If count is -1, then all events are returned in a descending order where the lastest is first.
func (mrr *MgoReadRepository) ReadFromLastCount(ctx context.Context, count int) ([]cqrskit.Event, error) {
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	zcol := zdb.C(AggregateEventCollection)

	var aggevents []AggregateEvent

	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
	}

	if count > 0 {
		if err := zcol.Find(rmQuery).Sort("-version").Limit(count).All(&aggevents); err != nil {
			return nil, err
		}
	} else {
		if err := zcol.Find(rmQuery).Sort("-version").All(&aggevents); err != nil {
			return nil, err
		}
	}

	var totalEvents int64
	for _, aggr := range aggevents {
		totalEvents += aggr.Count
	}

	events := make([]cqrskit.Event, int(totalEvents))

	var n int
	for _, aggr := range aggevents {
		n = copy(events[:n], aggr.Events)
	}

	return events, nil
}

// ReadVersion returns all events for giving aggregate and aggregate model for requested version
// if found.
func (mrr *MgoReadRepository) ReadVersion(ctx context.Context, version int64) ([]cqrskit.Event, error) {
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	zcol := zdb.C(AggregateEventCollection)
	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
		"version":      version,
	}

	next := struct {
		Events []cqrskit.Event `bson:"events"`
	}{}

	if err := zcol.Find(rmQuery).One(&next); err != nil {
		return nil, err
	}

	return next.Events, nil
}

// ReadFromVersion returns all events for giving aggregate and events for aggregate model from the version
// till the provided limit count in batch versions.
// Note Limit works by individual record batch and not by total events in batch, it helps avoid partial document
// update.
func (mrr *MgoReadRepository) ReadFromVersion(ctx context.Context, version int64, limit int) ([]cqrskit.Event, error) {
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	_, totalEvents, err := mrr.CountBatchesFromVersion(ctx, version, limit)
	if err != nil {
		return nil, err
	}

	events := make([]cqrskit.Event, totalEvents)

	next := struct {
		Events []cqrskit.Event `bson:"events"`
	}{Events: events}

	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
		"version":      bson.M{"$gte": version},
	}

	zcol := zdb.C(AggregateEventCollection)
	var itr *mgo.Iter

	if limit > 0 {
		itr = zcol.Find(rmQuery).Limit(limit).Sort("version").Iter()
	} else {
		itr = zcol.Find(rmQuery).Sort("version").Iter()
	}

	for itr.Next(&next) {
		if err := itr.Err(); err != nil {
			return events, err
		}

		last := len(next.Events)
		next.Events = events[last:]
	}

	if err = itr.Close(); err != nil {
		return events, err
	}

	return events, nil
}

// ReadFromTime returns all events for giving aggregate and events for aggregate model for the time of creation
// of event batch until the provided limit.
func (mrr *MgoReadRepository) ReadFromTime(ctx context.Context, ts time.Time, limit int) ([]cqrskit.Event, error) {
	zdb, _, zerr := mrr.db.New(true)
	if zerr != nil {
		return nil, zerr
	}

	_, totalEvents, err := mrr.CountBatchesFromTime(ctx, ts, limit)
	if err != nil {
		return nil, err
	}

	events := make([]cqrskit.Event, totalEvents)

	next := struct {
		Events []cqrskit.Event `bson:"events"`
	}{Events: events}

	rmQuery := bson.M{
		"aggregate_id": mrr.aggregateID,
		"instance_id":  mrr.instanceID,
		"created":      bson.M{"$gte": ts.Format(time.RFC3339)},
	}

	zcol := zdb.C(AggregateEventCollection)
	itr := zcol.Find(rmQuery).Limit(limit).Sort("version").Iter()

	for itr.Next(&next) {
		if err := itr.Err(); err != nil {
			return events, err
		}

		last := len(next.Events)
		next.Events = events[last:]
	}

	if err = itr.Close(); err != nil {
		return events, err
	}

	return events, nil
}
