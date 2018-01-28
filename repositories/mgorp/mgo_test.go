package mgorp_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gokit/cqrskit"

	"github.com/influx6/faux/tests"

	"github.com/gokit/cqrskit/repositories/mgorp"
	"github.com/gokit/cqrskit/repositories/mgorp/mdb"
)

var (
	aggregateId = "43543b2323I"
	modelId     = "233JNIosd232"
	created     = time.Now()
	config      = mdb.Config{
		DB:       os.Getenv("MONGO_DB"),
		Host:     os.Getenv("MONGO_HOST"),
		User:     os.Getenv("MONGO_USER"),
		AuthDB:   os.Getenv("MONGO_AUTHDB"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}
)

func TestMongoRepository(t *testing.T) {
	hostdb := mdb.NewMongoDB(config)
	writeRepo := mgorp.NewWriteMaster(hostdb)
	readRepo := mgorp.NewReadMaster(hostdb)

	testWriteMaster_New(t, hostdb, writeRepo)
	testWriteRepository_SaveEvents(t, hostdb, writeRepo)
	testReadRepository_ReadAll(t, hostdb, readRepo)
	testReadRepository_CountBatches(t, hostdb, readRepo)
	testReadRepository_CountBatchesFromTime(t, hostdb, readRepo)
	testReadRepository_CountBatchesFromTimeWithLimit(t, hostdb, readRepo)
	testReadRepository_CountBatchesFromCount(t, hostdb, readRepo)
	testReadRepository_CountBatchesFromCountWithLimit(t, hostdb, readRepo)
	testReadRepository_CountBatchesFromVersion(t, hostdb, readRepo)
	testReadRepository_CountBatchesFromVersionWithLimit(t, hostdb, readRepo)
	testReadRepository_ReadFromVersion(t, hostdb, readRepo)
	testReadRepository_ReadFromVersionWithLimit(t, hostdb, readRepo)
	testReadRepository_ReadFromCount(t, hostdb, readRepo)
	testReadRepository_ReadFromCountWithLimit(t, hostdb, readRepo)
	testReadRepository_ReadVersion(t, hostdb, readRepo)
	dropCollection(t, hostdb)
}

func dropCollection(t *testing.T, db mdb.MongoDB) {
	zdb, zses, err := db.New(false)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten db session")
	}
	tests.Passed("Should have successfully gotten db session")

	defer zses.Close()

	if err := zdb.C(mgorp.AggregateCollection).DropCollection(); err != nil {
		tests.FailedWithError(err, "Should have successfully dropped 'aggregate' collection")
	}
	tests.Passed("Should have successfully dropped 'aggregate' collection")

	if err := zdb.C(mgorp.AggregateModelCollection).DropCollection(); err != nil {
		tests.FailedWithError(err, "Should have successfully dropped 'aggregate_model' collection")
	}
	tests.Passed("Should have successfully dropped 'aggregate_model' collection")

	if err := zdb.C(mgorp.AggregateEventCollection).DropCollection(); err != nil {
		tests.FailedWithError(err, "Should have successfully dropped 'aggregate_events_model' collection")
	}
	tests.Passed("Should have successfully dropped 'aggregate_events_model' collection")
}

func testWriteMaster_New(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoWriteMaster) {
	repo, err := hostRepo.New(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully created new aggregate repository")
	}
	tests.Passed("Should have successfully created new aggregate repository")

	count, err := repo.Count(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved event count")
	}
	tests.Passed("Should have successfully retrieved event count")

	if count != 0 {
		tests.Failed("Should have total event record of 0 in db")
	}
	tests.Passed("Should have total event record of 0 in db")
}

func testWriteRepository_SaveEvents(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoWriteMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully created new aggregate repository")
	}
	tests.Passed("Should have successfully created new aggregate repository")

	events := []struct {
		Event []cqrskit.Event
		Done  func(error)
	}{
		{
			Event: []cqrskit.Event{
				{
					AggregateID: aggregateId,
					InstanceID:  modelId,
					Created:     created,
					EventType:   "UserCreated",
					EventData:   map[string]interface{}{"name": "bob", "email": "bob@bob.com"},
				},
				{
					AggregateID: aggregateId,
					InstanceID:  modelId,
					Created:     created,
					EventType:   "UserEmailUpdated",
					EventData:   map[string]interface{}{"new_email": "bob_quatz@bob.com"},
				},
			},
			Done: func(e error) {
				if e != nil {
					tests.FailedWithError(e, "Should have successfully saved user events")
				}
			},
		},
		{
			Event: []cqrskit.Event{
				{
					AggregateID: aggregateId,
					InstanceID:  modelId,
					Created:     created,
					EventType:   "UsernameUpdated",
					EventData:   map[string]interface{}{"vid": 1},
				},
				{
					AggregateID: aggregateId,
					InstanceID:  modelId,
					Created:     created,
					EventType:   "UserAccountUpgraded",
					EventData:   map[string]interface{}{"plan": "gold"},
				},
			},
			Done: func(e error) {
				if e != nil {
					tests.Failed("Should have successfully to saved user events")
				}
			},
		},
	}

	for _, event := range events {
		if err := repo.SaveEvents(context.Background(), event.Event); err != nil && event.Done != nil {
			event.Done(err)
		}
	}

	count, err := repo.Count(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved event count")
	}
	tests.Passed("Should have successfully retrieved event count")

	if count != 2 {
		tests.Failed("Should have total event record of 2 in db")
	}
	tests.Passed("Should have total event record of 2 in db")

	tests.Passed("Should have successfully saved all events")
}

func testReadRepository_CountBatches(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatches(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved records count")
	}
	tests.Passed("Should have successfully retrieved records count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	if totalEvents != 4 {
		tests.Failed("Should have a total of 4 events with 2 per record")
	}
	tests.Passed("Should have a total of 4 events with 2 per record")
}

func testReadRepository_ReadAll(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatches(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	events, err := repo.ReadAll(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != totalEvents {
		tests.Info("Expected Count: %d", totalEvents)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_CountBatchesFromTime(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromTime(context.Background(), created, -1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved record count")
	}
	tests.Passed("Should have successfully retrieved record count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	if totalEvents != 4 {
		tests.Failed("Should have a total of 4 events with 2 per record")
	}
	tests.Passed("Should have a total of 4 events with 2 per record")
}

func testReadRepository_CountBatchesFromTimeWithLimit(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromTime(context.Background(), created, 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved record count")
	}
	tests.Passed("Should have successfully retrieved record count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	if totalEvents != 2 {
		tests.Failed("Should have a total of 2 events with 2 per record")
	}
	tests.Passed("Should have a total of 2 events with 2 per record")
}

func testReadRepository_CountBatchesFromCount(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromCount(context.Background(), 2)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved record count")
	}
	tests.Passed("Should have successfully retrieved record count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	if totalEvents != 4 {
		tests.Failed("Should have a total of 4 events with 2 per record")
	}
	tests.Passed("Should have a total of 4 events with 2 per record")
}

func testReadRepository_CountBatchesFromCountWithLimit(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromCount(context.Background(), 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved record count")
	}
	tests.Passed("Should have successfully retrieved record count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	if totalEvents != 2 {
		tests.Failed("Should have a total of 2 events with 2 per record")
	}
	tests.Passed("Should have a total of 2 events with 2 per record")
}

func testReadRepository_CountBatchesFromVersion(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromVersion(context.Background(), 1, -1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved record count")
	}
	tests.Passed("Should have successfully retrieved record count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	if totalEvents != 4 {
		tests.Failed("Should have a total of 4 events with 2 per record")
	}
	tests.Passed("Should have a total of 4 events with 2 per record")
}

func testReadRepository_CountBatchesFromVersionWithLimit(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromVersion(context.Background(), 1, 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved record count")
	}
	tests.Passed("Should have successfully retrieved record count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	if totalEvents != 2 {
		tests.Failed("Should have a total of 2 events with 2 per record")
	}
	tests.Passed("Should have a total of 2 events with 2 per record")
}

func testReadRepository_ReadFromVersion(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromVersion(context.Background(), 1, -1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	events, err := repo.ReadFromVersion(context.Background(), 1, -1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != totalEvents {
		tests.Info("Expected Count: %d", totalEvents)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadFromVersionWithLimit(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromVersion(context.Background(), 1, 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	events, err := repo.ReadFromVersion(context.Background(), 1, 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != totalEvents {
		tests.Info("Expected Count: %d", totalEvents)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadFromCount(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromCount(context.Background(), -1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	events, err := repo.ReadFromLastCount(context.Background(), -1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != totalEvents {
		tests.Info("Expected Count: %d", totalEvents)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadFromCountWithLimit(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, totalEvents, err := mgoRepo.CountBatchesFromCount(context.Background(), 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")

	events, err := repo.ReadFromLastCount(context.Background(), 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != totalEvents {
		tests.Info("Expected Count: %d", totalEvents)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadVersion(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Get(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	events, err := repo.ReadVersion(context.Background(), 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != 2 {
		tests.Info("Expected Count: %d", 2)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")

	events2, err := repo.ReadVersion(context.Background(), 2)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events2); recCount != 2 {
		tests.Info("Expected Count: %d", 2)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}
