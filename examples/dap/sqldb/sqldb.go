// Package sqldb provides a auto-generated package which contains a sql CRUD API for the specific Ignitor struct in package dap.
//
//
package sqldb

import (
	"errors"
	"fmt"

	"github.com/influx6/faux/db"
	"github.com/jmoiron/sqlx"

	"github.com/influx6/faux/db/sql"
	"github.com/influx6/faux/db/sql/tables"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/moz/examples/dap"
)

// IgnitorFields defines an interface which exposes method to return a map of all
// attributes associated with the defined structure as decided by the structure.
type IgnitorFields interface {
	Fields() map[string]interface{}
}

// IgnitorConsumer defines an interface which accepts a map of data which will be consumed
// into the giving implementing structure as decided by the structure.
type IgnitorConsumer interface {
	Consume(map[string]interface{}) error
}

// mapFields defines a type for a map that exposes a Fields() method.
type mapFields map[string]interface{}

// Fields returns the map itself and provides a method to match the sql.TableField interface.
func (m mapFields) Fields() map[string]interface{} {
	return m
}

// IgnitorDB defines a structure which provide DB CRUD operations
// using sql as the underline db.
type IgnitorDB struct {
	col     string
	sx      sql.DB
	dx      *sql.SQL
	metrics metrics.Metrics
	table   db.TableIdentity
}

// New returns a new instance of IgnitorDB.
func New(table string, m metrics.Metrics, sx sql.DB, tm ...tables.TableMigration) *IgnitorDB {
	dx := sql.New(m, sx, tm...)

	return &IgnitorDB{
		sx:      sx,
		dx:      dx,
		col:     table,
		metrics: m,
		table:   db.TableName{Name: table},
	}
}

// Delete attempts to remove the record from the db using the provided publicID.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given dap.Ignitor struct.
func (mdb *IgnitorDB) Delete(ctx context.Context, publicID string) error {
	m := metrics.NewTrace("IgnitorDB.Delete")
	defer mdb.metrics.Emit(metrics.Info("IgnitorDB.Delete"), metrics.With("publicID", publicID), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to delete record"), metrics.WithFields("publicID", publicID), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))
		return err
	}

	if err := mdb.dx.Delete(mdb.table, "public_id", publicID); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to delete record"), metrics.WithFields("table", mdb.col), metrics.WithFields("publicID", publicID), metrics.WithFields("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Deleted record"), metrics.WithFields("table", mdb.col), metrics.WithFields("publicID", publicID))

	return nil
}

// Create attempts to add the record into the db using the provided instance of the
// dap.Ignitor.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Create(ctx context.Context, elem dap.Ignitor) error {
	m := metrics.NewTrace("IgnitorDB.Create")
	defer mdb.metrics.Emit(metrics.Info("IgnitorDB.Create"), metrics.With("publicID", elem.PublicID), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to create record"), metrics.WithFields("publicID", elem.PublicID), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		if err := mdb.dx.Save(mdb.table, mapFields(fields.Fields())); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to create Ignitor record"), metrics.WithFields("table", mdb.col), metrics.WithFields("elem", elem), metrics.WithFields("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(metrics.Info("Create record"), metrics.WithFields("table", mdb.col), metrics.WithFields("elem", elem))

		return nil
	}

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")

		mdb.metrics.Emit(metrics.Errorf("Failed to create record"), metrics.WithFields("publicID", elem.PublicID), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))
		return err
	}

	content := mapFields(map[string]interface{}{

		"hash": elem.var1.Hash,

		"name": elem.Name,

		"public_id": elem.PublicID,

		"rack": elem.Rack,

		"rex": map[string]interface{}{

			"url": elem.Rex.URL,
		},
	})

	if err := mdb.dx.Save(mdb.table, content); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create Ignitor record"), metrics.WithFields("table", mdb.col), metrics.WithFields("query", content), metrics.WithFields("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Create record"), metrics.WithFields("table", mdb.col), metrics.WithFields("query", content))
	return nil
}

// GetAll retrieves all records from the db and returns a slice of dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) GetAll(ctx context.Context, order string, orderby string, page int, responsePerPage int) ([]dap.Ignitor, int, error) {
	m := metrics.NewTrace("IgnitorDB.GetAll")
	defer mdb.metrics.Emit(metrics.Info("IgnitorDB.GetAll"), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record"), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))

		return nil, -1, err
	}

	var ritems []dap.Ignitor

	total, err := mdb.dx.GetAllPerPageBy(mdb.table, order, orderby, page, responsePerPage, func(rows *sqlx.Rows) error {
		for rows.Next() {
			var ritem dap.Ignitor

			if err := rows.StructScan(&ritem); err != nil {
				mdb.metrics.Emit(metrics.Errorf(err), metrics.WithFields(metrics.Field{
					"err":   err,
					"table": mdb.table.Table(),
				}))

				return err
			}

			ritems = append(ritems, ritem)
		}

		return nil
	})

	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to consume record data for Ignitor from db"), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))

		return nil, total, err
	}

	return ritems, total, nil

}

// GetByField retrieves a record from the db using the field key and value,
// returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) GetByField(ctx context.Context, key string, value string) (dap.Ignitor, error) {
	m := metrics.NewTrace("IgnitorDB.Get")
	defer mdb.metrics.Emit(metrics.Info("IgnitorDB.Get"), metrics.With(key, value), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record"), metrics.WithFields(key, value), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))

		return dap.Ignitor{}, err
	}

	item, err := mdb.dx.Get(mdb.table, key, value)
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of Ignitor type from db"), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))

		return dap.Ignitor{}, err
	}

	var elem dap.Ignitor

	if _, ok := elem.(IgnitorConsumer); !ok {
		return elem, errors.New("Only IgnitorConsumer allowed")
	}

	if err := elem.Consume(item); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to consume record data for Ignitor from db"), metrics.WithFields("table", mdb.col), metrics.WithFields("data", item), metrics.WithFields("error", err.Error()))

		return dap.Ignitor{}, err
	}

	return elem, nil
}

// Get retrieves a record from the db using the publicID and returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Get(ctx context.Context, publicID string) (dap.Ignitor, error) {
	m := metrics.NewTrace("IgnitorDB.Get")
	defer mdb.metrics.Emit(metrics.Info("IgnitorDB.Get"), metrics.With("publicID", publicID), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record"), metrics.WithFields("publicID", publicID), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))
		return dap.Ignitor{}, err
	}

	var elem dap.Ignitor

	if err := mdb.dx.GetBy(mdb.table, func(row *sqlx.Row) error {
		if err := row.StructScan(&elem); err != nil {
			mdb.metrics.Emit(metrics.Errorf(err), metrics.WithFields(metrics.Field{
				"err":   err,
				"table": mdb.table.Table(),
			}))

			return err
		}

		return nil
	}, "public_id", publicID); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to consume record data for Ignitor from db"), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))

		return dap.Ignitor{}, err
	}

	return elem, nil

}

// Update uses a record from the db using the publicID and returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Update(ctx context.Context, publicID string, elem dap.Ignitor) error {
	m := metrics.NewTrace("IgnitorDB.Update")
	defer mdb.metrics.Emit(metrics.Info("IgnitorDB.Update"), metrics.With("publicID", publicID), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to finish, context has expired"), metrics.WithFields("table", mdb.col), metrics.WithFields("public_id", publicID), metrics.WithFields("error", err.Error()))
		return err
	}

	var data mapFields

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		data = mapFields(fields.Fields())
	} else {
		data = mapFields(map[string]interface{}{

			"hash": elem.var1.Hash,

			"name": elem.Name,

			"public_id": elem.PublicID,

			"rack": elem.Rack,

			"rex": map[string]interface{}{

				"url": elem.Rex.URL,
			},
		})
	}

	if err := mdb.dx.Update(mdb.table, data, "public_id", publicID); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to update Ignitor record"), metrics.WithFields("query", data), metrics.WithFields("table", mdb.col), metrics.WithFields("public_id", publicID), metrics.WithFields("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Update record"), metrics.WithFields("table", mdb.col), metrics.WithFields("public_id", publicID), metrics.WithFields("query", data))

	return nil
}

// Exec provides a function which allows the execution of a custom function against the table.
func (mdb *IgnitorDB) Exec(ctx context.Context, fx func(*sql.SQL, sql.DB) error) error {
	m := metrics.NewTrace("IgnitorDB.Exec")
	defer mdb.metrics.Emit(metrics.Info("IgnitorDB.Exec"), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to execute operation"), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))
		return err
	}

	if err := fx(mdb.dx, mdb.sx); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to execute operation"), metrics.WithFields("table", mdb.col), metrics.WithFields("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Operation executed"), metrics.WithFields("table", mdb.col))

	return nil
}
