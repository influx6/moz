package httpapi_test

import (
	"fmt"

	"golang.com/x/sync/syncmap"

	"github.com/influx6/faux/context"

	"github.com/influx6/moz/examples/dap"
)

// MockAPIOperator defines an structure which implements the APIOperator for providing
// a mock usage in tests for use with the Unconvertible Type http API.
type MockAPIOperator struct {
	store *syncmap.Map
}

// NewMockAPIOperator returns a new instance of a MockAPIOperator.
func NewMockAPIOperator() *MockAPIOperator {
	return &MockAPIOperator{
		store: new(syncmap.Map),
	}
}

// Delete provides the operation to remove a giving record identified by ID.
func (m *MockAPIOperator) Delete(ctx context.Context, publicID string) error {
	if _, has := m.store.Load(publicID); !has {
		return fmt.Errorf("Record does not exists with id %q", publicID)
	}

	m.store.Delete(publicID)
	return nil
}

// GetAll returns a slice of all available record of type dap.Ignitor.
func (m *MockAPIOperator) GetAll(ctx context.Context) ([]dap.Ignitor, error) {
	var records []dap.Ignitor

	m.store.Range(func(k, v interface{}) bool {
		if elem, ok := v.(dap.Ignitor); ok {
			records = append(records, elem)
		}

		return true
	})

	return records
}

// Get retrieves a record based on the provided publicID.
func (m *MockAPIOperator) Get(ctx context.Context, publicID string) (dap.Ignitor, error) {
	elem, found := m.store.Load(publicID)
	if !found {
		return nil, fmt.Errorf("Record does not exists with id %q", publicID)
	}

	return elem, nil
}

// Update updates a giving record with the given new value.
func (m *MockAPIOperator) Update(ctx context.Context, publicID string, elem dap.Ignitor) error {

	m.store.Store(publicID, elem)
	return nil

}

// Create adds a new record into the giving record store.
func (m *MockAPIOperator) Create(ctx context.Context, elem dap.Ignitor) (dap.Ignitor, error) {

	m.store.Store(elem.PublicID, elem)
	return elem, nil

}
