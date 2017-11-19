package sqldb_test

import (
	"os"
	"time"

	"testing"

	"github.com/influx6/faux/tests"

	"github.com/influx6/faux/db/sql"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics/custom"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"

	_ "github.com/mattn/go-sqlite3"

	"github.com/influx6/moz/examples/dap/sqldb"
)

var (
	events = metrics.New(custom.StackDisplay(os.Stdout))

	config = sql.Config{
		DBName:       os.Getenv("dap_SQL_DB"),
		User:         os.Getenv("dap_SQL_USER"),
		DBIP:         os.Getenv("dap_SQL_ADDR"),
		DBPort:       os.Getenv("dap_SQL_PORT"),
		DBDriver:     os.Getenv("dap_SQL_Driver"),
		UserPassword: os.Getenv("dap_SQL_PASSWORD"),
	}

	testCol = "ignitor_test_collection"
)

// TestGetIgnitor validates the retrieval of a Ignitor
// record from a sqldb.
func TestGetIgnitor(t *testing.T) {
	api := sqldb.New(testCol, events, sql.NewDB(config, events))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")

	_, err = api.Get(ctx, elem.PublicID)
	if err != nil {
		tests.Failed("Successfully retrieved stored record for Ignitor from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved stored record for Ignitor from db.")
}

// TestGetAllIgnitor validates the retrieval of all Ignitor
// record from a sqldb.
func TestGetAllIgnitor(t *testing.T) {
	api := sqldb.New(testCol, events, sql.NewDB(config, events))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")

	records, _, err := api.GetAll(ctx, "asc", "public_id", -1, -1)
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
// record with a sqldb.
func TestIgnitorCreate(t *testing.T) {
	api := sqldb.New(testCol, events, sql.NewDB(config, events))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(ignitorCreateJSON)
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
// record with a sqldb.
func TestIgnitorUpdate(t *testing.T) {
	api := sqldb.New(testCol, events, sql.NewDB(config, events))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Ignitor record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Ignitor record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Ignitor into db.")

	//TODO: Update something.

	if err := api.Update(ctx, elem.PublicID, elem); err != nil {
		tests.Failed("Successfully updated record for Ignitor into db: %+q.", err)
	}
	tests.Passed("Successfully updated record for Ignitor into db.")

	_, err = api.Get(ctx, elem.PublicID)
	if err != nil {
		tests.Failed("Successfully retrieved stored record for Ignitor from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved stored record for Ignitor from db.")
}

// TestIgnitorDelete validates the removal of a Ignitor
// record from a sqldb.
func TestIgnitorDelete(t *testing.T) {
	api := sqldb.New(testCol, events, sql.NewDB(config, events))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(ignitorCreateJSON)
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
