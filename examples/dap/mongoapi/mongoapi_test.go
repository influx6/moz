package mongoapi_test

import (
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/db/mongo"

	"github.com/influx6/faux/metrics/sentries/stdout"
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
)

// TestGetIgnitor validates the retrieval of a Ignitor
// record from a mongodb.
func TestGetIgnitor(t *testing.T) {
	db := mongo.New(config)

}

// TestGetAllIgnitor validates the retrieval of all Ignitor
// record from a mongodb.
func TestGetAllIgnitor(t *testing.T) {

}

// TestIgnitorCreate validates the creation of a Ignitor
// record with a mongodb.
func TestIgnitorCreate(t *testing.T) {

	var ignitorVar Ignitor

	var identityVar Identity
	identityVar.Hash = ""

	ignitorVar.Identity = identityVar

	ignitorVar.PublicID = ""

	ignitorVar.Name = ""

	var rexVar Repo
	rexVar.URL = ""

	ignitorVar.Rex = rexVar

	ignitorVar.Rack = 0

	ignitorVar.version = ""

}

// TestIgnitorUpdate validates the update of a Ignitor
// record with a mongodb.
func TestIgnitorUpdate(t *testing.T) {

	var ignitorVar Ignitor

	var identityVar Identity
	identityVar.Hash = ""

	ignitorVar.Identity = identityVar

	ignitorVar.PublicID = ""

	ignitorVar.Name = ""

	var rexVar Repo
	rexVar.URL = ""

	ignitorVar.Rex = rexVar

	ignitorVar.Rack = 0

	ignitorVar.version = ""

}

// TestIgnitorDelete validates the removal of a Ignitor
// record from a mongodb.
func TestIgnitorDelete(t *testing.T) {

}
