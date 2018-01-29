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
	testReadRepository_CountCommits(t, hostdb, readRepo)
	testReadRepository_ReadFromVersion(t, hostdb, readRepo)
	testReadRepository_ReadFromVersionWithLimit(t, hostdb, readRepo)
	testReadRepository_ReadSinceCount(t, hostdb, readRepo)
	testReadRepository_ReadSinceCountWithLimit(t, hostdb, readRepo)
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

	if err := zdb.C(mgorp.AggregateEventCommitCollection).DropCollection(); err != nil {
		tests.FailedWithError(err, "Should have successfully dropped 'aggregate_events_model' collection")
	}
	tests.Passed("Should have successfully dropped 'aggregate_events_model' collection")
}

func testWriteMaster_New(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoWriteMaster) {
	repo, err := hostRepo.Writer(aggregateId, modelId)
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
	repo, err := hostRepo.Writer(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully created new aggregate repository")
	}
	tests.Passed("Should have successfully created new aggregate repository")

	events := []struct {
		Event cqrskit.EventCommitRequest
		Done  func(error)
	}{
		{
			Event: cqrskit.EventCommitRequest{
				Command: "CreateUser",
				ID:      "433436577674674574567575675",
				Created: created,
				Events: []cqrskit.Event{
					{
						Type: "UserCreated",
						Data: map[string]interface{}{"name": "bob", "email": "bob@bob.com"},
					},
					{
						Type: "UserPlanChange",
						Data: map[string]interface{}{"email": "bob@bob.com", "plan": "gold"},
					},
				},
			},
			Done: func(e error) {
				if e != nil {
					tests.FailedWithError(e, "Should have successfully saved user events")
				}
			},
		},
		{
			Event: cqrskit.EventCommitRequest{
				Command: "UpdateUserEmail",
				ID:      "436895577674674574567575675",
				Created: created,
				Events: []cqrskit.Event{
					{
						Type: "UserEmailUpdated",
						Data: map[string]interface{}{"email": "bob@bob.com"},
					},
				},
			},
			Done: func(e error) {
				if e != nil {
					tests.FailedWithError(e, "Should have successfully saved user events")
				}
			},
		},
	}

	for _, event := range events {
		if _, err := repo.Write(context.Background(), event.Event); err != nil && event.Done != nil {
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

func testReadRepository_CountCommits(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Reader(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, err := mgoRepo.CountCommits(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved records count")
	}
	tests.Passed("Should have successfully retrieved records count")

	if totalRecords != 2 {
		tests.Failed("Should have a total of 2 records")
	}
	tests.Passed("Should have a total of 2 records")
}

func testReadRepository_ReadAll(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Reader(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, err := mgoRepo.CountCommits(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	events, err := repo.ReadAll(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != totalRecords {
		tests.Info("Expected Count: %d", totalRecords)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadFromVersion(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Reader(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	mgoRepo, ok := repo.(*mgorp.MgoReadRepository)
	if !ok {
		tests.Failed("Should have a underline *mgorp.MgoReadMaster")
	}
	tests.Passed("Should have a underline *mgorp.MgoReadMaster")

	totalRecords, err := mgoRepo.CountCommits(context.Background())
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	events, err := repo.ReadSinceVersion(context.Background(), 1, -1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != totalRecords {
		tests.Info("Expected Count: %d", totalRecords)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadFromVersionWithLimit(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Reader(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	events, err := repo.ReadSinceVersion(context.Background(), 1, 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != 1 {
		tests.Info("Expected Count: %d", 1)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadSinceCount(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Reader(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	events, err := repo.ReadSinceCount(context.Background(), -1)
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
}

func testReadRepository_ReadSinceCountWithLimit(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Reader(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	events, err := repo.ReadSinceCount(context.Background(), 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events); recCount != 1 {
		tests.Info("Expected Count: %d", 1)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}

func testReadRepository_ReadVersion(t *testing.T, db mdb.MongoDB, hostRepo mgorp.MgoReadMaster) {
	repo, err := hostRepo.Reader(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully gotten aggregate read repository")
	}
	tests.Passed("Should have successfully gotten aggregate read repository")

	events, err := repo.ReadVersion(context.Background(), 1)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved all records")
	}
	tests.Passed("Should have successfully retrieved all records")

	if recCount := len(events.Events); recCount != 2 {
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

	if recCount := len(events2.Events); recCount != 1 {
		tests.Info("Expected Count: %d", 1)
		tests.Info("Received Count: %d", recCount)
		tests.Failed("Should have retrieved expected records in count")
	}
	tests.Passed("Should have retrieved expected records in count")
}
