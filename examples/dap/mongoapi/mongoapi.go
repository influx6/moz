// Package mongoapi provides a auto-generated package which contains a mongo CRUD API for the specific Ignitor struct in package dap.
//
//
package mongoapi

import (
	"encoding/json"

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
	col             string
	db              Mongod
	metrics         metrics.Metrics
	ensuredIndex    bool
	incompleteIndex bool
	indexes         []mgo.Index
}

// New returns a new instance of IgnitorDB.
func New(col string, m metrics.Metrics, mo Mongod, indexes ...mgo.Index) *IgnitorDB {
	return &IgnitorDB{
		db:      mo,
		col:     col,
		metrics: m,
		indexes: indexes,
	}
}

// ensureIndex attempts to ensure all provided indexes into the specific collection.
func (mdb *IgnitorDB) ensureIndex() error {
	m := stdout.Info("IgnitorDB.ensureIndex").Trace("IgnitorDB.ensureIndex")
	defer mdb.metrics.Emit(m.End())

	if mdb.ensuredIndex {
		return nil
	}

	if len(mdb.indexes) == 0 {
		return nil
	}

	// If we had an error before index was complete, then skip, we cant not
	// stop all ops because of failed index.
	if !mdb.ensuredIndex && mdb.incompleteIndex {
		return nil
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session for index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	collection := database.C(mdb.col)

	for _, index := range mdb.indexes {
		if err := collection.EnsureIndex(index); err != nil {
			mdb.metrics.Emit(stdout.Error("Failed to ensure session index").
				With("collection", mdb.col).
				With("index", index).
				With("error", err.Error()))

			mdb.incompleteIndex = true
			return err
		}

		mdb.metrics.Emit(stdout.Info("Succeeded in ensuring collection index").
			With("collection", mdb.col).
			With("index", index))
	}

	mdb.ensuredIndex = true

	mdb.metrics.Emit(stdout.Notice("Finished adding index").
		With("collection", mdb.col))

	return nil
}

// Count attempts to return the total number of record from the db.
func (mdb *IgnitorDB) Count(ctx context.Context) (int, error) {
	m := stdout.Info("IgnitorDB.Count").Trace("IgnitorDB.Count")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")

		mdb.metrics.Emit(stdout.Error("Failed to get record count").
			With("collection", mdb.col).
			With("error", err.Error()))
		return -1, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))

		return -1, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to get record count").
			With("collection", mdb.col).
			With("error", err.Error()))

		return -1, err
	}

	defer session.Close()

	query := bson.M{}

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to get record count").
			With("collection", mdb.col).
			With("error", err.Error()))

		return -1, err
	}

	total, err := database.C(mdb.col).Find(query).Count()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to get record count").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return -1, err
	}

	mdb.metrics.Emit(stdout.Notice("Deleted record").
		With("collection", mdb.col).
		With("query", query))

	return total, err
}

// Delete attempts to remove the record from the db using the provided publicID.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given dap.Ignitor struct.
func (mdb *IgnitorDB) Delete(ctx context.Context, publicID string) error {
	m := stdout.Info("IgnitorDB.Delete").With("publicID", publicID).Trace("IgnitorDB.Delete")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to delete record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to delete record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	query := bson.M{
		"publicID": publicID,
	}

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to delete record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := database.C(mdb.col).Remove(query); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to delete record").
			With("collection", mdb.col).
			With("query", query).
			With("publicID", publicID).With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Deleted record").
		With("collection", mdb.col).
		With("query", query).
		With("publicID", publicID))

	return nil
}

// Create attempts to add the record into the db using the provided instance of the
// dap.Ignitor.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Create(ctx context.Context, elem dap.Ignitor) error {
	m := stdout.Info("IgnitorDB.Create").With("publicID", elem.PublicID).Trace("IgnitorDB.Create")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to create record").With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session").
			With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to create record").With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(IgnitorBSON); ok {
		if err := database.C(mdb.col).Insert(fields.BSON()); err != nil {
			mdb.metrics.Emit(stdout.Error("Failed to create Ignitor record").
				With("collection", mdb.col).
				With("elem", elem).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("collection", mdb.col).
			With("elem", elem).
			With("error", err.Error()))

		return nil
	}

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		if err := database.C(mdb.col).Insert(bson.M(fields.Fields())); err != nil {
			mdb.metrics.Emit(stdout.Error("Failed to create Ignitor record").
				With("collection", mdb.col).
				With("elem", elem).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("collection", mdb.col).
			With("elem", elem).
			With("error", err.Error()))

		return nil
	}

	query := bson.M(map[string]interface{}{

		"hash": elem.Identity.Hash,

		"name": elem.Name,

		"public_id": elem.PublicID,

		"rack": elem.Rack,

		"rex": map[string]interface{}{

			"url": elem.Rex.URL,
		},
	})

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to create record").With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := database.C(mdb.col).Insert(query); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create Ignitor record").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Create record").
		With("collection", mdb.col).
		With("query", query))

	return nil
}

// GetAllPerPage retrieves all records from the db and returns a slice of dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) GetAllPerPage(ctx context.Context, order string, orderBy string, page int, responsePerPage int) ([]dap.Ignitor, error) {
	m := stdout.Info("IgnitorDB.GetAll").Trace("IgnitorDB.GetAll")
	defer mdb.metrics.Emit(m.End())

	switch strings.ToLower(order) {
	case "dsc", "desc":
		orderBy = "-" + orderBy
	}

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	if page <= 0 && responsePerPage <= 0 {
		records, err := mdb.GetAll(table, order, orderBy)
		return records, len(records), err
	}

	// Get total number of records.
	totalRecords, err := mdb.Count(ctx)
	if err != nil {
		return nil, -1, err
	}

	var totalWanted, indexToStart int

	if page <= 1 && responsePerPage > 0 {
		totalWanted = responsePerPage
		indexToStart = 0
	} else {
		totalWanted = responsePerPage * page
		indexToStart = totalWanted / 2

		if page > 1 {
			indexToStart++
		}
	}

	mdb.metrics.Emit(stdout.Info("DB:Query:GetAllPerPage").WithFields(metrics.Fields{
		"starting_index":       indexToStart,
		"total_records_wanted": totalWanted,
		"order":                order,
		"orderBy":              orderBy,
		"page":                 page,
		"responsePerPage":      responsePerPage,
	}))

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	query := bson.M{}

	var items []dap.Ignitor

	if err := database.C(mdb.col).Find(query).Skip(indexToStart).Limit(totalWanted).Sort(orderBy).All(&items); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to retrieve all records of Ignitor type from db").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return nil, err
	}

	return items, nil

}

// GetAll retrieves all records from the db and returns a slice of dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) GetAll(ctx context.Context, order, orderBy string) ([]dap.Ignitor, error) {
	m := stdout.Info("IgnitorDB.GetAll").Trace("IgnitorDB.GetAll")
	defer mdb.metrics.Emit(m.End())

	switch strings.ToLower(order) {
	case "dsc", "desc":
		orderBy = "-" + orderBy
	}

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	query := bson.M{}

	var items []dap.Ignitor

	if err := database.C(mdb.col).Find(query).Sort(orderBy).All(&items); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to retrieve all records of Ignitor type from db").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return nil, err
	}

	return items, nil

}

// Get retrieves a record from the db using the publicID and returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Get(ctx context.Context, publicID string) (dap.Ignitor, error) {
	m := stdout.Info("IgnitorDB.Get").With("publicID", publicID).Trace("IgnitorDB.Get")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return dap.Ignitor{}, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return dap.Ignitor{}, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session").
			With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return dap.Ignitor{}, err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return dap.Ignitor{}, err
	}

	query := bson.M{"public_id": publicID}

	var item dap.Ignitor

	if err := database.C(mdb.col).Find(query).One(&item); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to retrieve all records of Ignitor type from db").
			With("query", query).
			With("collection", mdb.col).
			With("error", err.Error()))

		return dap.Ignitor{}, err
	}

	return item, nil

}

// Update uses a record from the db using the publicID and returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Update(ctx context.Context, publicID string, elem dap.Ignitor) error {
	m := stdout.Info("IgnitorDB.Update").With("publicID", publicID).Trace("IgnitorDB.Update")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to finish, context has expired").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session").
			With("publicID", publicID).
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to finish, context has expired").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(IgnitorBSON); ok {
		query := fields.BSON()

		if err := database.C(mdb.col).Insert(query); err != nil {
			mdb.metrics.Emit(stdout.Error("Failed to update Ignitor record").
				With("collection", mdb.col).
				With("public_id", publicID).
				With("query", query).
				With("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(stdout.Notice("Update record").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("query", query).
			With("error", err.Error()))

		return nil
	}

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		query := bson.M(fields.Fields())

		if err := database.C(mdb.col).Insert(query); err != nil {
			mdb.metrics.Emit(stdout.Error("Failed to update Ignitor record").
				With("query", query).
				With("public_id", publicID).
				With("collection", mdb.col).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("collection", mdb.col).
			With("query", query).
			With("public_id", publicID).
			With("error", err.Error()))

		return nil
	}

	query := bson.M{"publicID": publicID}
	queryData := bson.M(map[string]interface{}{

		"hash": elem.Identity.Hash,

		"name": elem.Name,

		"public_id": elem.PublicID,

		"rack": elem.Rack,

		"rex": map[string]interface{}{

			"url": elem.Rex.URL,
		},
	})

	if err := database.C(mdb.col).Update(query, queryData); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to update Ignitor record").
			With("collection", mdb.col).
			With("query", query).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Update record").
		With("collection", mdb.col).
		With("public_id", publicID).
		With("query", query))

	return nil
}

// Exec provides a function which allows the execution of a custom function against the collection.
func (mdb *IgnitorDB) Exec(ctx context.Context, fx func(col *mgo.Collection) error) error {
	m := stdout.Info("IgnitorDB.Exec").Trace("IgnitorDB.Exec")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to execute operation").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to finish, context has expired").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := fx(database.C(mdb.col)); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to execute operation").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Operation executed").
		With("collection", mdb.col))

	return nil
}
