package httpapi_test

import (
	"errors"

	"github.com/influx6/faux/context"

	"github.com/influx6/moz/examples/dap"
)

// MockCU defines a package level variable to access a instance of a mockCU.
var MockCU mockCU

// mockCU provides a struct which contains methods to create and update Ignitor;
type mockCU struct{}

// UpdateIgnitor updates a exiting dap.Ignitor record with the given update item.
func (mockCU) Update(ctx context.Context, record dap.Ignitor, updater dap.Ignitor) (dap.Ignitor, error) {

	// TODO(developer):
	// Override function contents with what should happen.

	return dap.Ignitor{}, errors.New("Not Implemented")
}

// CreateIgnitor creates a new dap.Ignitor record from the given create type.
func (mockCU) CreateIgnitor(ctx context.Context, elem dap.Ignitor) (dap.Ignitor, error) {

	// TODO(developer):
	// Override function contents with what should happen.

	return dap.Ignitor{}, errors.New("Not Implemented")
}
