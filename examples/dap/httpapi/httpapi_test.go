package httpapi_test

import (
	"fmt"

	"bytes"

	"testing"

	"encoding/json"

	"net/http"

	"net/http/httptest"

	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/tests"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics/sentries/stdout"

	"github.com/influx6/moz/examples/dap/httpapi"
)

var (
	events  = metrics.New(stdout.Stdout{})
	version = "v1"
)

// TestGetAllIgnitor validates the retrieval of a Ignitor
// record from httpapi.
func TestGetAllIgnitor(t *testing.T) {

	db := NewMockAPIOperator()

	api := httpapi.New(events, db)
	tree := httputil.NewMuxRouter(nil)

	// Register routes with router group.
	httpapi.Register(tree, api, version, "ignitors")

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
		tests.Failed("Should have successfully loaded JSON: ignitorCreateJSON : %+q.", err)
	}
	tests.Passed("Should have successfully loaded JSON: ignitorCreateJSON.")

	ctx := context.New()

	elem, err = db.Create(ctx, elem)
	if err != nil {
		tests.Failed("Should have successfully saved Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully saved Ignitor record.")

	req, err := http.NewRequest("GET", fmt.Sprintf("/%s/ignitors", version), nil)
	if err != nil {
		tests.Failed("Should have successfully created request Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully created request Ignitor record.")

	res := httptest.NewRecorder()

	tree.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		tests.Failed("Should have received status code %d from response but got %d.", http.StatusOK, res.Code)
	}
	tests.Passed("Should have received status code %d from response.", http.StatusOK)

	if res.Body.Len() == 0 {
		tests.Failed("Should have successfully received response body.")
	}
	tests.Passed("Should have successfully received response body.")
}

// TestGetIgnitor validates the retrieval of all Ignitor
// record from a httpapi.
func TestGetIgnitor(t *testing.T) {
	db := NewMockAPIOperator()

	api := httpapi.New(events, db)
	tree := httputil.NewMuxRouter(nil)

	// Register routes with router group.
	httpapi.Register(tree, api, version, "ignitors")

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
		tests.Failed("Should have successfully loaded JSON: ignitorCreateJSON : %+q.", err)
	}
	tests.Passed("Should have successfully loaded JSON: ignitorCreateJSON.")

	ctx := context.New()

	elem, err = db.Create(ctx, elem)
	if err != nil {
		tests.Failed("Should have successfully saved Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully saved Ignitor record.")

	req, err := http.NewRequest("GET", fmt.Sprintf("/%s/ignitors/%s", version, elem.PublicID), nil)
	if err != nil {
		tests.Failed("Should have successfully created request Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully created request Ignitor record.")

	res := httptest.NewRecorder()

	tree.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		tests.Failed("Should have received status code %d from response but got %d.", http.StatusOK, res.Code)
	}
	tests.Passed("Should have received status code %d from response.", http.StatusOK)

	if res.Body.Len() == 0 {
		tests.Failed("Should have successfully received response body.")
	}
	tests.Passed("Should have successfully received response body.")
}

// TestIgnitorCreate validates the creation of a Ignitor
// record with a httpapi.
func TestIgnitorCreate(t *testing.T) {

	db := NewMockAPIOperator()

	api := httpapi.New(events, db)
	tree := httputil.NewMuxRouter(nil)

	// Register routes with router group.
	httpapi.Register(tree, api, version, "ignitors")

	var body bytes.Buffer
	body.WriteString(ignitorCreateJSON)

	req, err := http.NewRequest("POST", fmt.Sprintf("/%s/ignitors", version), &body)
	if err != nil {
		tests.Failed("Should have successfully created request Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully created request Ignitor record.")

	res := httptest.NewRecorder()

	tree.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		tests.Failed("Should have received status code %d from response but got %d.", http.StatusCreated, res.Code)
	}
	tests.Passed("Should have received status code %d from response.", http.StatusCreated)

	if res.Body.Len() == 0 {
		tests.Failed("Should have successfully received response body.")
	}
	tests.Passed("Should have successfully received response body.")
}

// TestIgnitorUpdate validates the update of a Ignitor
// record with a httpapi.
func TestIgnitorUpdate(t *testing.T) {
	db := NewMockAPIOperator()

	api := httpapi.New(events, db)
	tree := httputil.NewMuxRouter(nil)

	// Register routes with router group.
	httpapi.Register(tree, api, version, "ignitors")

	var body bytes.Buffer
	body.WriteString(ignitorCreateJSON)

	req, err := http.NewRequest("POST", fmt.Sprintf("/%s/ignitors", version), &body)
	if err != nil {
		tests.Failed("Should have successfully created request Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully created request Ignitor record.")

	res := httptest.NewRecorder()

	tree.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		tests.Failed("Should have received status code %d from response but got %d.", http.StatusCreated, res.Code)
	}
	tests.Passed("Should have received status code %d from response.", http.StatusCreated)

	if res.Body.Len() == 0 {
		tests.Failed("Should have successfully received response body.")
	}
	tests.Passed("Should have successfully received response body.")

	elem, err := loadJSONFor(res.Body.String())
	if err != nil {
		tests.Failed("Should have successfully loaded JSON:  %+q.", err)
	}
	tests.Passed("Should have successfully loaded JSON.")

	elem.Name = "Joshua Lopez"

	var bu bytes.Buffer

	if err := json.NewEncoder(&bu).Encode(elem); err != nil {
		tests.Failed("Should have successfully encoded Ignitor:  %+q.", err)
	}
	tests.Passed("Should have successfully encoded Ignitor.")

	req, err = http.NewRequest("PUT", fmt.Sprintf("/%s/ignitors/%s", version, elem.PublicID), &bu)
	if err != nil {
		tests.Failed("Should have successfully created request Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully created request Ignitor record.")

	res = httptest.NewRecorder()

	tree.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		tests.Failed("Should have received status code %d from response but got %d.", http.StatusNoContent, res.Code)
	}
	tests.Passed("Should have received status code %d from response.", http.StatusNoContent)
}

// TestIgnitorDelete validates the removal of a Ignitor
// record from a httpapi.
func TestIgnitorDelete(t *testing.T) {
	db := NewMockAPIOperator()

	api := httpapi.New(events, db)
	tree := httputil.NewMuxRouter(nil)

	// Register routes with router group.
	httpapi.Register(tree, api, version, "ignitors")

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
		tests.Failed("Should have successfully loaded JSON: ignitorCreateJSON : %+q.", err)
	}
	tests.Passed("Should have successfully loaded JSON: ignitorCreateJSON.")

	ctx := context.New()

	elem, err = db.Create(ctx, elem)
	if err != nil {
		tests.Failed("Should have successfully saved Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully saved Ignitor record.")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s/ignitors/%s", version, elem.PublicID), nil)
	if err != nil {
		tests.Failed("Should have successfully created request Ignitor record : %+q.", err)
	}
	tests.Passed("Should have successfully created request Ignitor record.")

	res := httptest.NewRecorder()

	tree.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		tests.Failed("Should have received status code %d from response but got %d.", http.StatusNoContent, res.Code)
	}
	tests.Passed("Should have received status code %d from response.", http.StatusNoContent)
}
