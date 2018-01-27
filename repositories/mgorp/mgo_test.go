package mgorp_test

import (
	"context"
	"os"
	"testing"

	"github.com/gokit/cqrskit"

	"github.com/influx6/faux/tests"

	"github.com/gokit/cqrskit/repositories/mgorp"
	"github.com/gokit/cqrskit/repositories/mgorp/mdb"
)

var (
	config = mdb.Config{
		DB:       os.Getenv("MONGO_DB"),
		Host:     os.Getenv("MONGO_HOST"),
		User:     os.Getenv("MONGO_USER"),
		AuthDB:   os.Getenv("MONGO_AUTHDB"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}
)

func TestWriteMaster_New(t *testing.T) {
	hostdb := mdb.NewMongoDB(config)
	hostRepo := mgorp.NewWriteMaster(hostdb)

	aggregateId := "43543b2323I"
	modelId := "233JNIosd232"
	_, err := hostRepo.New(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully created new aggregate repository")
	}
	tests.Passed("Should have successfully created new aggregate repository")

}

func TestWriteRepository_SaveEvents(t *testing.T) {
	hostdb := mdb.NewMongoDB(config)
	hostRepo := mgorp.NewWriteMaster(hostdb)

	aggregateId := "43543b2323I"
	modelId := "233JNIosd232"
	repo, err := hostRepo.New(aggregateId, modelId)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully created new aggregate repository")
	}
	tests.Passed("Should have successfully created new aggregate repository")

	repo.SaveEvents(context.Background(), []cqrskit.Event{
		{},
	})
}

func TestReadRepository_CountBatches(t *testing.T) {

}

func TestReadRepository_ReadAll(t *testing.T) {

}
