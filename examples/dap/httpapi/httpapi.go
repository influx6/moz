// Package httpapi provides a auto-generated package which contains a http restful CRUD API for the specific Ignitor struct in package dap.
//
//
package httpapi

import (
	"errors"
	"fmt"
	"strconv"

	"net/http"

	"encoding/json"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics"

	httputil "github.com/influx6/faux/httputil"

	"github.com/influx6/faux/metrics/sentries/stdout"

	"github.com/influx6/moz/examples/dap"
)

// Register registers the giving route into the provided httptreemux function with the
// provided router and prefixed path.
func Register(router httputil.Router, api *HTTPApi, version string, resource string) {
	router.Handle("GET", fmt.Sprintf("/%s/%s", version, resource), api.GetAll)
	router.Handle("GET", fmt.Sprintf("/%s/%s/:public_id", version, resource), api.Get)

	router.Handle("POST", fmt.Sprintf("/%s/%s", version, resource), api.Create)

	router.Handle("PUT", fmt.Sprintf("/%s/%s/:public_id", version, resource), api.Update)
	router.Handle("DELETE", fmt.Sprintf("/%s/%s/:public_id", version, resource), api.Delete)
}

//================================================================================================

// APIOperator defines an interface which allows the HTTPApi to divert the final operation of
// the given CRUD request for the Unconvertible Type type. This is provided by the user.
type APIOperator interface {
	Delete(context.Context, string) error
	Get(context.Context, string) (dap.Ignitor, error)
	Update(context.Context, string, dap.Ignitor) error
	Create(context.Context, dap.Ignitor) (dap.Ignitor, error)
	GetAll(context.Context, string, string, int, int) ([]dap.Ignitor, int, error)
}

// IgnitorRecords defines a type to represent the response given to a request for
// all records of the type dap.Ignitor.
type IgnitorRecords struct {
	Page            int           `json:"page"`
	ResponsePerPage int           `json:"responsePerPage"`
	TotalRecords    int           `json:"total_records"`
	Records         []dap.Ignitor `json:"records"`
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
func (api *HTTPApi) Create(ctx *httputil.Context) error {
	m := stdout.Info("HTTPApi.Create").Trace()
	defer api.metrics.Emit(m.End())

	ctx.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("Create request received").WithFields(metrics.Fields{
		"url": ctx.Request().URL.String(),
	}))

	var incoming dap.Ignitor

	if err := json.NewDecoder(ctx.Body()).Decode(&incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Fields{
			"error": err,
			"url":   ctx.Request().URL.String(),
		}))

		return err
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"data": incoming,
		"url":  ctx.Request().URL.String(),
	}))

	response, err := api.operator.Create(ctx, incoming)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to create record").WithFields(metrics.Fields{
			"error": err,
			"url":   ctx.Request().URL.String(),
		}))

		return err
	}

	api.metrics.Emit(stdout.Info("Response Delivered").WithFields(metrics.Fields{
		"url":    ctx.Request().URL.String(),
		"status": http.StatusCreated,
	}))

	if err := ctx.JSON(http.StatusCreated, response); err != nil {
		api.metrics.Emit(stdout.Error("Failed to deliver response").WithFields(metrics.Fields{
			"error": err,
			"url":   ctx.Request().URL.String(),
		}))
		return err
	}

	return nil
}

// Update receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: PUT
// BODY: JSON
//
func (api *HTTPApi) Update(ctx *httputil.Context) error {
	m := stdout.Info("HTTPApi.Update").Trace()
	defer api.metrics.Emit(m.End())

	ctx.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("Update request received").WithFields(metrics.Fields{
		"url": ctx.Request().URL.String(),
	}))

	publicID, ok := ctx.GetString("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Fields{
			"url": ctx.Request().URL.String(),
		}))

		return errors.New("piblicId parameter not found")
	}

	var incoming dap.Ignitor

	if err := json.NewDecoder(ctx.Body()).Decode(&incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to decode request body").WithFields(metrics.Fields{
			"error":     err.Error(),
			"public_id": publicID,
			"url":       ctx.Request().URL.String(),
		}))

		return err
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"data":      incoming,
		"url":       ctx.Request().URL.String(),
		"public_id": publicID,
	}))

	if err := api.operator.Update(ctx, publicID, incoming); err != nil {
		api.metrics.Emit(stdout.Error("Failed to parse params and url.Values").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       ctx.Request().URL.String(),
		}))

		return err
	}

	api.metrics.Emit(stdout.Info("Response Delivered").WithFields(metrics.Fields{
		"url":       ctx.Request().URL.String(),
		"public_id": publicID,
		"status":    http.StatusNoContent,
	}))

	return ctx.NoContent(http.StatusNoContent)
}

// Delete receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: DELETE
//
func (api *HTTPApi) Delete(ctx *httputil.Context) error {
	m := stdout.Info("HTTPApi.Delete").Trace()
	defer api.metrics.Emit(m.End())

	api.metrics.Emit(stdout.Info("Delete request received").WithFields(metrics.Fields{
		"url": ctx.Request().URL.String(),
	}))

	publicID, ok := ctx.GetString("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Fields{
			"url": ctx.Request().URL.String(),
		}))

		return fmt.Errorf("No public_id provided in params")
	}

	api.metrics.Emit(stdout.Info("JSON received").WithFields(metrics.Fields{
		"url":       ctx.Request().URL.String(),
		"public_id": publicID,
	}))

	if err := api.operator.Delete(ctx, publicID); err != nil {
		api.metrics.Emit(stdout.Error("Failed to delete dap.Ignitor record").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       ctx.Request().URL.String(),
		}))

		return err
	}

	api.metrics.Emit(stdout.Info("Response Delivered").WithFields(metrics.Fields{
		"url":       ctx.Request().URL.String(),
		"public_id": publicID,
		"status":    http.StatusNoContent,
	}))

	return ctx.NoContent(http.StatusNoContent)
}

// Get receives an http request to create a new Unconvertible Type.
//
// Route: /{Route}/:public_id
// Method: GET
// RESPONSE-BODY: JSON
func (api *HTTPApi) Get(ctx *httputil.Context) error {
	m := stdout.Info("HTTPApi.Get").Trace()
	defer api.metrics.Emit(m.End())

	ctx.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("Get request received").WithFields(metrics.Fields{
		"url": ctx.Request().URL.String(),
	}))

	publicID, ok := ctx.GetString("public_id")
	if !ok {
		api.metrics.Emit(stdout.Error("No public_id provided in params").WithFields(metrics.Fields{
			"url": ctx.Request().URL.String(),
		}))

		return errors.New("public_id parameter not found")
	}

	requested, err := api.operator.Get(ctx, publicID)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to get dap.Ignitor record").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       ctx.Request().URL.String(),
		}))

		return err
	}

	if err := ctx.JSON(http.StatusOK, requested); err != nil {
		api.metrics.Emit(stdout.Error("Failed to get serialized dap.Ignitor record to response writer").WithFields(metrics.Fields{
			"error":     err,
			"public_id": publicID,
			"url":       ctx.Request().URL.String(),
		}))

		return err
	}

	api.metrics.Emit(stdout.Info("Response Delivered").WithFields(metrics.Fields{
		"url":       ctx.Request().URL.String(),
		"public_id": publicID,
		"status":    http.StatusOK,
	}))

	return nil
}

// GetAll receives an http request to return all Unconvertible Type records.
//
// Route: /{Route}/
// Method: GET
// RESPONSE-BODY: JSON
func (api *HTTPApi) GetAll(ctx *httputil.Context) error {
	m := stdout.Info("HTTPApi.GetAll").Trace()
	defer api.metrics.Emit(m.End())

	ctx.Header().Set("Content-Type", "application/json")

	api.metrics.Emit(stdout.Info("GetAll request received").WithFields(metrics.Fields{
		"url": ctx.Request().URL.String(),
	}))

	var order, orderBy string

	if od, ok := ctx.Get("order"); ok {
		if ordr, ok := od.(string); ok {
			order = ordr
		} else {
			order = "asc"
		}
	}

	if od, ok := ctx.Get("orderBy"); ok {
		if ordr, ok := od.(string); ok {
			orderBy = ordr
		} else {
			orderBy = "public_id"
		}
	}

	var err error
	var pageNo, responsePerPage int

	if rpp, ok := ctx.GetString("responsePerPage"); ok {
		responsePerPage, err = strconv.Atoi(rpp)
		if err != nil {
			api.metrics.Emit(stdout.Error("Failed to retrieve responserPerPage number details").WithFields(metrics.Fields{
				"error": err,
				"url":   ctx.Request().URL.String(),
			}))
		}
	} else {
		responsePerPage = -1
	}

	if pg, ok := ctx.GetString("page"); ok {
		pageNo, err = strconv.Atoi(pg)
		if err != nil {
			api.metrics.Emit(stdout.Error("Failed to retrieve page number details").WithFields(metrics.Fields{
				"error": err,
				"url":   ctx.Request().URL.String(),
			}))
		}
	} else {
		pageNo = -1
	}

	requested, total, err := api.operator.GetAll(ctx, order, orderBy, pageNo, responsePerPage)
	if err != nil {
		api.metrics.Emit(stdout.Error("Failed to get all dap.Ignitor record").WithFields(metrics.Fields{
			"error": err,
			"url":   ctx.Request().URL.String(),
		}))

		return err
	}

	if err := ctx.JSON(http.StatusOK, IgnitorRecords{
		Page:            pageNo,
		Records:         requested,
		TotalRecords:    total,
		ResponsePerPage: responsePerPage,
	}); err != nil {
		api.metrics.Emit(stdout.Error("Failed to get serialized dap.Ignitor record to response writer").WithFields(metrics.Fields{
			"error": err,
			"url":   ctx.Request().URL.String(),
		}))

		return err
	}

	api.metrics.Emit(stdout.Info("Response Delivered").WithFields(metrics.Fields{
		"url":    ctx.Request().URL.String(),
		"status": http.StatusOK,
	}))

	return nil
}

//================================================================================================
