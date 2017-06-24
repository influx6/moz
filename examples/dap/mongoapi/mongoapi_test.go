package mongoapi_test

import (
	"testing"
	"time"

	"encoding/json"

	mgo "gopkg.in/mgo.v2"

	"github.com/influx6/faux/tests"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/db/mongo"

	"github.com/influx6/faux/metrics/sentries/stdout"

	"github.com/influx6/moz/examples/dap/mongoapi"
)

var (
	events = metrics.New(stdout.Stdout{})

	config = mongo.Config{
		Mode:     mgo.Monotonic,
		User:     "<USERNAME>",
		DB:       "DATABASE_NAME",
		Password: "<PASSWORD>",
		Host:     "MONGO_DB_HOST_ADDR",
		AuthDB:   "AUTH_DATABASE_NAME",
	}

	testCol = "ignitor_test_collection"
)

// TestGetIgnitor validates the retrieval of a Ignitor
// record from a mongodb.
func TestGetIgnitor(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.New().WithDeadline(10*time.Second, false)

	elem, err := loadJSONFor(IgnitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")

	record, err := api.Get(ctx, elem.PublicID)
	if err != nil {
		tests.Failed("Successfully retrieved stored record for Ignitor from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved stored record for Ignitor from db.")

	if record.Name != elem.Name {
		tests.Failed("Successfully matched retrieved record with inserted record for Ignitor.")
	}
	tests.Passed("Successfully matched retrieved record with inserted record for Ignitor.")
}

// TestGetAllIgnitor validates the retrieval of all Ignitor
// record from a mongodb.
func TestGetAllIgnitor(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.New().WithDeadline(10*time.Second, false)

	elem, err := loadJSONFor(IgnitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")

	records, err := api.GetAll(ctx)
	if err != nil {
		tests.Failed("Successfully retrieved all records for Ignitor from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved all records for Ignitor from db.")

	if len(records) == 0 {
		tests.Failed("Successfully retrieved atleast 1 record for Ignitor from db.")
	}
	tests.Passed("Successfully retrieved atleast 1 record for Ignitor from db.")
}

// TestIgnitorCreate validates the creation of a Ignitor
// record with a mongodb.
func TestIgnitorCreate(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.New().WithDeadline(10*time.Second, false)

	elem, err := loadJSONFor(IgnitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")
}

// TestIgnitorUpdate validates the update of a Ignitor
// record with a mongodb.
func TestIgnitorUpdate(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.New().WithDeadline(10*time.Second, false)

	elem, err := loadJSONFor(IgnitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")

	elem.Name = "Bob Marley"

	if err := api.Update(ctx, elem); err != nil {
		tests.Failed("Successfully updated record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully updated record for Ignitor into db.")

	record, err := api.Get(ctx, elem.PublicID)
	if err != nil {
		tests.Failed("Successfully retrieved stored record for Ignitor from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved stored record for Ignitor from db.")

	if elem.Name != record.Name {
		tests.Failed("Successfully matched updated record with inserted record for Ignitor.")
	}
	tests.Passed("Successfully matched updated record with inserted record for Ignitor.")
}

// TestIgnitorDelete validates the removal of a Ignitor
// record from a mongodb.
func TestIgnitorDelete(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.New().WithDeadline(10*time.Second, false)

	elem, err := loadJSONFor(IgnitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")

	if err := api.Delete(ctx, elem.PublicID); err != nil {
		tests.Failed("Successfully removed record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully removed record for Ignitor into db.")
}

// loadJSONFor returns a new instance of a Ignitor from the provide json content.
func loadJSONForIgnitor(content string) (Ignitor, error) {
	var elem Ignitor

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return nil, err
	}

	return elem, nil
}
