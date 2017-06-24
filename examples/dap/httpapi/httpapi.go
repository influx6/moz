// Package httpapi provides a auto-generated package which contains a http restful CRUD API for the specific Ignitor struct in package dap.
//
//
package httpapi

import (
	"fmt"

	"net/http"

	"encoding/json"

	"github.com/dimfeld/httptreemux"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics"

	httputil "github.com/influx6/faux/httputil"

	"github.com/influx6/faux/metrics/sentries/stdout"

	"github.com/influx6/moz/examples/dap"
)

// RegisterRouteGroup registers the giving route into the provided httptreemux function with the
// provided router and prefixed path.
func RegisterRouteGroup(grp *httptreemux.Group, api *HTTPApi, version string, resource string) {
	grp.GET(fmt.Sprintf("/%s/%s", version, resource), WrapTreemux(api.GetAll))
	grp.GET(fmt.Sprintf("/%s/%s/:public_id", version, resource), WrapTreemux(api.Get))

	grp.POST(fmt.Sprintf("/%s/%s", version, resource), WrapTreemux(api.Create))

	grp.PUT(fmt.Sprintf("/%s/%s/:public_id", version, resource), WrapTreemux(api.Update))
	grp.DELETE(fmt.Sprintf("/%s/%s/:public_id", version, resource), WrapTreemux(api.Delete))
}

// RegisterRoute registers the giving route into the provided httptreemux function with the
// provided router and prefixed path.
func RegisterRoute(router *httptreemux.TreeMux, api *HTTPApi, version string, resource string) {
	router.GET(fmt.Sprintf("/%s/%s", version, resource), WrapTreemux(api.GetAll))
	router.GET(fmt.Sprintf("/%s/%s/:public_id", version, resource), WrapTreemux(api.Get))

	router.POST(fmt.Sprintf("/%s/%s", version, resource), WrapTreemux(api.Create))

	router.PUT(fmt.Sprintf("/%s/%s/:public_id", version, resource), WrapTreemux(api.Update))
	router.DELETE(fmt.Sprintf("/%s/%s/:public_id", version, resource), WrapTreemux(api.Delete))
}

//================================================================================================

// HTTPContextHandler defines a function which is used to service a request with a
// context
type HTTPContextHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// WrapTreemux defines the function to meet the httptreemux.Handler interface to appropriately
// parse all request to the appropriate handler.
func WrapTreemux(fn HTTPContextHandler) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		ctx := context.From(r.Context())

		for name, value := range params {
			ctx.Set(name, value)
		}

		fn(ctx, w, r)
	}
}

// HTTPParams provides a function which returns a http.Handler which calls httputil.Params
// to collect parameters from a request into a context.
func HTTPParams(fn HTTPContextHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.From(r.Context())

		httputil.Params(ctx, r, 0)

		fn(ctx, w, r)
	}
}

//================================================================================================

// APIOperator defines an interface which allows the HTTPApi to divert the final operation of
// the given CRUD request for the Unconvertible Type type. This is provided by the user.
type APIOperator interface {
	Delete(context.Context, string) error
	GetAll(context.Context) ([]dap.Ignitor, error)
	Get(context.Context, string) (dap.Ignitor, error)
	Update(context.Context, string, dap.Ignitor) error
	Create(context.Context, dap.Ignitor) (dap.Ignitor, error)
}

//================================================================================================

// HTTPApi defines a struct which holds the http api handlers for providing CRUD
// operations for the provided Unconvertible Type type.
type HTTPApi struct {
	operator APIOperator
	metrics  metrics.Metrics
}

// New returns a new HTTPApi instance using the provided operator and
// metric.
func New(m metrics.Metrics, operator APIOperator) *HTTPApi {
	return &HTTPApi{
		operator: operator,
		metrics:  m,
	}
}

// Create receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: POST
// BODY: JSON
//
func (api *HTTPApi) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("Create request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	var incoming dap.Ignitor

	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Fields{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to decode json body"), http.StatusInternalServerError)
		return
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"data": incoming,
		"url":  r.URL.String(),
	}))

	response, err := api.operator.Create(ctx, incoming)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to create record").WithFields(metrics.Fields{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to create dap.Ignitor object"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.metrics.Emit(stdout.Error("Failed to write response body").WithFields(metrics.Fields{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to write response of dap.Ignitor object"), http.StatusInternalServerError)
		return
	}
}

// Update receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: PUT
// BODY: JSON
//
func (api *HTTPApi) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("Update request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	publicIDVal, ok := ctx.Get("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Fields{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("No public_id provided in params"), http.StatusBadRequest)
		return
	}

	publicID, ok := publicIDVal.(string)
	if !ok {
		api.metrics.Emit(stdout.Error("public_id param is not a string").WithFields(metrics.Fields{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("public_id param is not a string"), http.StatusBadRequest)
		return
	}

	var incoming dap.Ignitor

	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to decode json body"), http.StatusInternalServerError)
		return
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"data":      incoming,
		"url":       r.URL.String(),
		"public_id": publicID,
	}))

	if err := api.operator.Update(ctx, publicID, incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to update record of dap.Ignitor object"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: DELETE
//
func (api *HTTPApi) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	api.metrics.Emit(stdout.Info("Delete request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	publicIDVal, ok := ctx.Get("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Fields{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("No public_id provided in params"), http.StatusBadRequest)
		return
	}

	publicID, ok := publicIDVal.(string)
	if !ok {
		api.metrics.Emit(stdout.Error("public_id param is not a string").WithFields(metrics.Fields{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("public_id param is not a string"), http.StatusBadRequest)
		return
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"url":       r.URL.String(),
		"public_id": publicID,
	}))

	if err := api.operator.Delete(ctx, publicID); err != nil {
		api.metrics.Emit(stdout.Error("Failed to delete dap.Ignitor record").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to delete record"), http.StatusBadRequest)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}

// Get receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: GET
// RESPONSE-BODY: JSON
func (api *HTTPApi) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("Get request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	publicIDVal, ok := ctx.Get("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Fields{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("No public_id provided in params"), http.StatusBadRequest)
		return
	}

	publicID, ok := publicIDVal.(string)
	if !ok {
		api.metrics.Emit(stdout.Error("public_id param is not a string").WithFields(metrics.Fields{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("public_id param is not a string"), http.StatusBadRequest)
		return
	}

	requested, err := api.operator.Get(ctx, publicID)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to get dap.Ignitor record").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to retrieve record"), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(requested); err != nil {
		api.metrics.Emit(stdout.Error("Failed to get serialized dap.Ignitor record to response writer").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to write response params"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetAll receives an http request to return all Unconvertible Type records.
//
// Route: /{Route}/
// Method: GET
// RESPONSE-BODY: JSON
func (api *HTTPApi) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("GetAll request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	requested, err := api.operator.GetAll(ctx)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to get all dap.Ignitor record").WithFields(metrics.Fields{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to retrieve all records"), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(requested); err != nil {
		api.metrics.Emit(stdout.Error("Failed to get serialized dap.Ignitor record to response writer").WithFields(metrics.Fields{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Sprintf("Failed to write response"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//================================================================================================
