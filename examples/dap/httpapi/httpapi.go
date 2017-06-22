// Package http provides a auto-generated package which contains a http restful CRUD API for the specific Ignitor struct in package dap.
//
//
package httpapi

import (
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
	grp.Get(fmt.Sprintf("%s/%s", version, resource), WrapTreemux(api.GetAll))
	grp.Post(fmt.Sprintf("%s/%s", version, resource), WrapTreemux(api.Create))

	grp.Get(fmt.Sprintf("%s/%s/:public_id", version, resource), WrapTreemux(api.Get))
	grp.Put(fmt.Sprintf("%s/%s/:public_id", version, resource), WrapTreemux(api.Update))
	grp.Delete(fmt.Sprintf("%s/%s/:public_id", version, resource), WrapTreemux(api.Delete))
}

//================================================================================================

// HTTPContextHandler defines a function which is used to service a request with a
// context
type HTTPContextHandler func(ctx context.Context, w http.ResponseWriter, r http.Request)

// WrapTreemux defines the function to meet the httptreemux.Handler interface to appropriately
// parse all request to the appropriate handler.
func WrapTreemux(fn HTTPContextHandler) httptreemux.Handler {
	return func(w http.ResponseWriter, r http.Request, params map[string]interface{}) {
		ctx := context.From(r.Context())

		for name, value := range params {
			ctx.Set(name, value)
		}

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
func New(m metrics.Metric, operator APIOperator) *HTTPApi {
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
func (api *HTTPApi) Create(ctx context.Context, w http.ResponseWriter, r http.Request) {
	api.metrics.Emit(stdout.Info("Create request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	if err := httputil.Params(ctx, r, 0); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to parse params"), http.StatusBadRequest)
		return
	}

	var incoming dap.Ignitor

	if err := json.NewDecoder(w).Decode(&incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to decode json body"), http.StatusInternalServerError)
		return
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"data": incoming,
		"url":  r.URL.String(),
	}))

	response, err := api.operator.Create(ctx, incoming)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to create dap.Ignitor object"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.metrics.Emit(stdout.Error("Failed to write response body").WithFields(metrics.Field{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to create dap.Ignitor object"), http.StatusInternalServerError)
		return
	}
}

// Update receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: PUT
// BODY: JSON
//
func (api *HTTPApi) Update(ctx context.Context, w http.ResponseWriter, r http.Request) {
	api.metrics.Emit(stdout.Info("Update request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	if err := httputil.Params(ctx, r, 0); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to parse params"), http.StatusBadRequest)
		return
	}

	publicID, ok := ctx.Get("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Field{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Error("No public_id provided in params"), http.StatusBadRequest)
		return
	}

	var incoming dap.Ignitor

	if err := json.NewDecoder(w).Decode(&incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to decode json body"), http.StatusInternalServerError)
		return
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"data":      incoming,
		"url":       r.URL.String(),
		"public_id": publicID,
	}))

	if err := api.operator.Update(ctx, publicID, incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to create dap.Ignitor object"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: DELETE
//
func (api *HTTPApi) Delete(ctx context.Context, w http.ResponseWriter, r http.Request) {
	api.metrics.Emit(stdout.Info("Delete request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	if err := httputil.Params(ctx, r, 0); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to parse params"), http.StatusBadRequest)
		return
	}

	publicID, ok := ctx.Get("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Field{
			"url": r.URL.String(),
		}))

		http.Error(w, fmt.Error("No public_id provided in params"), http.StatusBadRequest)
		return
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"url":       r.URL.String(),
		"public_id": publicID,
	}))

	if err := api.metrics.Delete(publicID); err != nil {
		api.metrics.Emit(stdout.Error("Failed to delete dap.Ignitor record").WithFields(metrics.Field{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to parse params"), http.StatusBadRequest)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}

// Get receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: GET
// RESPONSE-BODY: JSON
func (api *HTTPApi) Get(ctx context.Context, w http.ResponseWriter, r http.Request) {
	api.metrics.Emit(stdout.Info("Get request received").WithFields(metrics.Fields{
		"url": r.URL.String(),
	}))

	if err := httputil.Params(ctx, r, 0); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Field{
			"error": err,
			"url":   r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to parse params"), http.StatusBadRequest)
		return
	}

	publicID, ok := ctx.Get("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Field{
			"url": r.URL.String(),
		})).Error(w, fmt.Error("No public_id provided in params"), http.StatusBadRequest)
		return
	}

	requested, err := api.operator.Get(publicID)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to get dap.Ignitor record").WithFields(metrics.Field{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to parse params"), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(requested); err != nil {
		api.metrics.Emit(stdout.Error("Failed to get serialized dap.Ignitor record to response writer").WithFields(metrics.Field{
			"error":     err,
			"public_id": publicID,
			"url":       r.URL.String(),
		}))

		http.Error(w, fmt.Error("Failed to parse params"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//================================================================================================
