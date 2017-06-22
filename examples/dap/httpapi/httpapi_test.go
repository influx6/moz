package httpapi_test

import (
	"testing"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/metrics/sentries/stdout"
)

var (
	events = metrics.New(stdout.Stdout{})
)

// TestGetIgnitor validates the retrieval of a Ignitor
// record from httpapi.
func TestGetIgnitor(t *testing.T) {

}

// TestGetAllIgnitor validates the retrieval of all Ignitor
// record from a httpapi.
func TestGetAllIgnitor(t *testing.T) {

}

// TestIgnitorCreate validates the creation of a Ignitor
// record with a httpapi.
func TestIgnitorCreate(t *testing.T) {

}

// TestIgnitorUpdate validates the update of a Ignitor
// record with a httpapi.
func TestIgnitorUpdate(t *testing.T) {

}

// TestIgnitorDelete validates the removal of a Ignitor
// record from a httpapi.
func TestIgnitorDelete(t *testing.T) {

}
