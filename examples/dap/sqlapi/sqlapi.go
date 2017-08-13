// Package sqlapi provides a auto-generated package which contains a sql CRUD API for the specific Ignitor struct in package dap.
//
//
package sqlapi

import (
	"fmt"

	"github.com/influx6/faux/db"

	"github.com/influx6/faux/db/sql"
	"github.com/influx6/faux/db/sql/tables"

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
	m := stdout.Info("IgnitorDB.Delete").With("publicID", publicID).Trace("IgnitorDB.Delete")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to delete record").
			With("publicID", publicID).
			With("table", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := mdb.dx.Delete(mdb.table, "public_id", publicID); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to delete record").
			With("table", mdb.col).
			With("publicID", publicID).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Deleted record").
		With("table", mdb.col).
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
		mdb.metrics.Emit(stdout.Error("Failed to create record").
			With("publicID", elem.PublicID).
			With("table", mdb.col).
			With("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		if err := mdb.dx.Save(mdb.table, mapFields(fields.Fields())); err != nil {
			mdb.metrics.Emit(stdout.Error("Failed to create Ignitor record").
				With("table", mdb.col).
				With("elem", elem).
				With("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(stdout.Notice("Create record").
			With("table", mdb.col).
			With("elem", elem))

		return nil
	}

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")

		mdb.metrics.Emit(stdout.Error("Failed to create record").
			With("publicID", elem.PublicID).
			With("table", mdb.col).
			With("error", err.Error()))
		return err
	}

	content := mapFields(map[string]interface{}{

		"hash": elem.Identity.Hash,

		"name": elem.Name,

		"public_id": elem.PublicID,

		"rack": elem.Rack,

		"rex": map[string]interface{}{

			"url": elem.Rex.URL,
		},
	})

	if err := mdb.dx.Save(mdb.table, content); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create Ignitor record").
			With("table", mdb.col).
			With("query", content).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Create record").
		With("table", mdb.col).
		With("query", content))

	return nil
}

// GetAll retrieves all records from the db and returns a slice of dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) GetAll(ctx context.Context, order string, orderby string, page int, responsePerPage int) ([]dap.Ignitor, int, error) {
	m := stdout.Info("IgnitorDB.GetAll").Trace("IgnitorDB.GetAll")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").
			With("table", mdb.col).
			With("error", err.Error()))

		return nil, -1, err
	}

	var ritems []dap.Ignitor

	total, err := mdb.dx.GetAllPerPageBy(mdb.table, order, orderby, page, responsePerPage, func(rows *sqlx.Rows) error {
		for rows.Next() {
			var ritem dap.Ignitor

			if err := rows.StructScan(&ritem); err != nil {
				mdb.metrics.Emit(stdout.Error(err).WithFields(metrics.Fields{
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
		mdb.metrics.Emit(stdout.Error("Failed to consume record data for Ignitor from db").
			With("table", mdb.col).
			With("error", err.Error()))

		return nil, total, err
	}

	return ritems, total, nil

}

// Get retrieves a record from the db using the publicID and returns the dap.Ignitor type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Ignitor struct.
func (mdb *IgnitorDB) Get(ctx context.Context, publicID string) (dap.Ignitor, error) {
	m := stdout.Info("IgnitorDB.Get").With("publicID", publicID).Trace("IgnitorDB.Get")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to retrieve record").
			With("publicID", publicID).
			With("table", mdb.col).
			With("error", err.Error()))
		return dap.Ignitor{}, err
	}

	var elem dap.Ignitor

	if err := mdb.dx.GetBy(mdb.table, func(row *sqlx.Row) error {
		if err := row.StructScan(&elem); err != nil {
			mdb.metrics.Emit(stdout.Error(err).WithFields(metrics.Fields{
				"err":   err,
				"table": mdb.table.Table(),
			}))

			return err
		}

		return nil
	}, "public_id", publicID); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to consume record data for Ignitor from db").
			With("table", mdb.col).
			With("error", err.Error()))

		return dap.Ignitor{}, err
	}

	return elem, nil

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
			With("table", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	var data mapFields

	if fields, ok := interface{}(elem).(IgnitorFields); ok {
		data = mapFields(fields.Fields())
	} else {
		data = mapFields(map[string]interface{}{

			"hash": elem.Identity.Hash,

			"name": elem.Name,

			"public_id": elem.PublicID,

			"rack": elem.Rack,

			"rex": map[string]interface{}{

				"url": elem.Rex.URL,
			},
		})
	}

	if err := mdb.dx.Update(mdb.table, data, "public_id", publicID); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to update Ignitor record").
			With("query", data).
			With("table", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Update record").
		With("table", mdb.col).
		With("public_id", publicID).
		With("query", data))

	return nil
}

// Exec provides a function which allows the execution of a custom function against the table.
func (mdb *IgnitorDB) Exec(ctx context.Context, fx func(*sql.SQL, sql.DB) error) error {
	m := stdout.Info("IgnitorDB.Exec").Trace("IgnitorDB.Exec")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to execute operation").
			With("table", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := fx(mdb.dx, mdb.sx); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to execute operation").
			With("table", mdb.col).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Operation executed").
		With("table", mdb.col))

	return nil
}
