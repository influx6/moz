// Package mongoapi provides a auto-generated package which contains a mongo base pkg for db operations.
//
//
package mongo

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics"
)

// Mongod defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type Mongod interface {
	New() (*mgo.Database, *mgo.Session, error)
}

// DB defines a structure which provide DB CRUD operations
// using mongo as the underline db.
type DB struct {
	db      Mongod
	metrics metrics.Metrics
}

// New returns a new instance of DB.
func New(m metrics.Metrics, mo Mongod) *DB {
	return &DB{
		db:      mo,
		metrics: m,
	}
}

// WithIndex applies the provided index slice to the provided collection configuration.
func (mdb *DB) WithIndex(ctx context.Context, col string, indexes ...mgo.Index) error {
	m := metrics.NewTrace("DB.WithIndex")
	defer mdb.metrics.Emit(metrics.Info("DB.WithIndex"), metrics.WithTrace(m.End()))

	if len(indexes) == 0 {
		return nil
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session for indexes"), metrics.With("collection", col), metrics.With("error", err.Error()))
		return err
	}

	defer session.Close()

	collection := database.C(col)

	for _, index := range indexes {

		if err := collection.EnsureIndex(index); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to ensure session index"), metrics.With("collection", col), metrics.With("index", index), metrics.With("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(metrics.Info("Succeeded in ensuring collection index"), metrics.With("collection", col), metrics.With("index", index))
	}

	mdb.metrics.Emit(metrics.Info("Finished adding index").
		With("collection", col))

	return nil
}

// Exec provides a function which allows the execution of a custom function against the collection.
func (mdb *DB) Exec(ctx context.Context, col string, fx func(col *mgo.Collection) error) error {
	m := metrics.NewTrace("DB.Exec")
	defer mdb.metrics.Emit(metrics.Info("DB.Exec"), metrics.WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to execute operation"), metrics.With("collection", col), metrics.With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session"), metrics.With("collection", col), metrics.With("error", err.Error()))
		return err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to finish, context has expired"), metrics.With("collection", col), metrics.With("error", err.Error()))
		return err
	}

	if err := fx(database.C(col)); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to execute operation"), metrics.With("collection", col), metrics.With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Operation executed"), metrics.With("collection", col))

	return nil
}
