// Package mongoapi provides a auto-generated package which contains a mongo CRUD API for the specific Ignitor struct in package dap.
//
//
package mongoapi

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/metrics/sentries/stdout"

	"github.com/influx6/moz/examples/dap"
)

// IgnitorFields defines an interface which exposes method to return a map of all
// attributes associated with the defined structure as decided by the structure.
type IgnitorFields interface {
	Fields() map[string]interface{}
}

// IgnitorBSON defines an interface which exposes method to return a bson.M type
// which contains all related fields for the giving  object.
type IgnitorBSON interface {
	BSON() bson.M
}

// IgnitorBSONConsumer defines an interface which accepts a map of data which will be consumed
// into the giving implementing structure as decided by the structure.
type IgnitorBSONConsumer interface {
	BSONConsume(bson.M) error
}

// IgnitorConsumer defines an interface which accepts a map of data which will be consumed
// into the giving implementing structure as decided by the structure.
type IgnitorConsumer interface {
	Consume(map[string]interface{}) error
}

// Mongod defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type Mongod interface {
	New() (*mgo.Database, *mgo.Session, error)
}

// IgnitorDB defines a structure which provide DB CRUD operations
// using mongo as the underline db.
type IgnitorDB struct {
	col     string
	db      Mongod
	metrics metrics.metrics
}

// New returns a new instance of IgnitorDB.
func New(col string, m metrics.metrics, mo Mongod) *IgnitorDB {
	return &IgnitorDB{
		db:      mo,
		col:     col,
		metrics: m,
	}
}

// Delete attempts to remove the record from the db using the provided publicID.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given dap.Ignitor struct.
func (mdb *IgnitorDB) Delete(ctx context.Context, publicID string) error {
	m := stdout.Info("IgnitorDB.Delete").With("public_id", publicID).Trace("IgnitorDB.Delete")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	database, session, err := mdb.DB.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	defer session.Close()

	query := bson.M{
		"public_id": publicID,
	}

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	if err := database.C(mdb.col).Remove(query); err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").
			With("query", query).
			With("public_id", public_id).With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Deleted record").
		With("query", query).
		With("public_id", public_id).With("error", err.Error()))

	return nil
}

// Create attempts to add the record into the db using the provided instance of the
// dap.Ignitor.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Create(ctx context.Context, elem dap.Ignitor) error {
	m := stdout.Info("IgnitorDB.Create").With("public_id", publicID).Trace("IgnitorDB.Create")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to create record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	database, session, err := mdb.DB.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to create record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(IgnitorBSON); ok {
		if err := database.C(mdb.col).Insert(fields.BSON()); err != nil {
			mdb.metrics.Emit(stdout.Errorf("Failed to create Ignitor record").
				With("query", query).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("query", query).
			With("error", err.Error()))

		return nil
	}

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		if err := database.C(mdb.col).Insert(bson.M(fields.Fields())); err != nil {
			mdb.metrics.Emit(stdout.Errorf("Failed to create Ignitor record").
				With("query", query).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("query", query).
			With("error", err.Error()))

		return nil
	}

	query := bson.M(map[string]interface{}{

		"hash": elem.Identity.Hash,

		"name": elem.Name,

		"public_id": elem.PublicID,

		"rex": map[string]interface{}{

			"url": elem.Rex.URL,
		},
	})

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to create record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	if err := database.C(mdb.col).Insert(query); err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to create Ignitor record").
			With("query", query).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Create record").
		With("query", query).
		With("error", err.Error()))

	return nil
}

// GetAll retrieves all records from the db and returns a slice of dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) GetAll(ctx context.Context) ([]dap.Ignitor, error) {
	m := stdout.Info("IgnitorDB.GetAll").Trace("IgnitorDB.GetAll")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to retrieve record").With("error", err.Error()))
		return err
	}

	database, session, err := mdb.DB.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").With("public_id", public_id).With("error", err.Error()))
		return nil, err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to retrieve record").With("error", err.Error()))
		return err
	}

	query := bson.M{
		"public_id": publicID,
	}

	var items []dap.Ignitor

	if err := database.C(mdb.col).Find(query).All(&items); err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to retrieve all records of Ignitor type from db").
			With("query", query).
			With("error", err.Error()))

		return nil, err
	}

	return items, nil

}

// Get retrieves a record from the db using the public_id and returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Get(ctx context.Context, publicID string) (dap.Ignitor, error) {
	m := stdout.Info("IgnitorDB.Get").With("public_id", publicID).Trace("IgnitorDB.Get")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to retrieve record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	database, session, err := mdb.DB.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to retrieve record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	var item dap.Ignitor

	if err := database.C(mdb.col).Find(query).All(&items); err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to retrieve all records of Ignitor type from db").
			With("query", query).
			With("error", err.Error()))

		return nil, err
	}

	return item, nil

}

// Update uses a record from the db using the public_id and returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Update(ctx context.Context, publicID string, elem dap.Ignitor) error {
	m := stdout.Info("IgnitorDB.Update").With("public_id", publicID).Trace("IgnitorDB.Update")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to update record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	database, session, err := mdb.DB.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to delete record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Errorf("Failed to update record").With("public_id", public_id).With("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(IgnitorBSON); ok {
		query := fields.BSON()

		if err := database.C(mdb.col).Insert(query); err != nil {
			mdb.metrics.Emit(stdout.Errorf("Failed to create Ignitor record").
				With("query", query).
				With("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("query", query).
			With("error", err.Error()))

		return nil
	}

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		query := bson.M(fields.Fields())

		if err := database.C(mdb.col).Insert(query); err != nil {
			mdb.metrics.Emit(stdout.Errorf("Failed to create Ignitor record").
				With("query", query).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("query", query).
			With("error", err.Error()))

		return nil
	}

	query := bson.M(map[string]interface{}{

		"hash": elem.Identity.Hash,

		"name": elem.Name,

		"public_id": elem.PublicID,

		"rex": map[string]interface{}{

			"url": elem.Rex.URL,
		},
	})

	if err := database.C(mdb.col).Insert(query); err != nil {
		mdb.metrics.Emit(stdout.Errorf("Failed to create Ignitor record").
			With("query", query).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Create record").
		With("query", query).
		With("error", err.Error()))

	return nil
}
