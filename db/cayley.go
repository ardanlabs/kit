package db

import (
	"errors"

	kitcayley "github.com/ardanlabs/kit/db/cayley"
	"github.com/cayleygraph/cayley"
)

// ErrGraphHandle is returned when a graph handle is not initialized.
const ErrGraphHandle = errors.New("Graph handle not initialized.")

//==============================================================================
// Methods for the DB struct type related to Cayley.

// OpenCayley opens a connection to Cayley and adds that support to the
// database value.
func (db *DB) OpenCayley(context interface{}, cfg kitcayley.Config) error {
	store, err := kitcayley.New(cfg)
	if err != nil {
		return err
	}
	db.graphHandle = store
	return nil
}

// GraphHandle returns the Cayley graph handle for graph interactions.
func (db *DB) GraphHandle(context interface{}) (*cayley.Handle, error) {
	if db.graphHandle != nil {
		return db.graphHandle, nil
	}
	return nil, ErrGraphHandle
}

// CloseCayley closes a graph handle value.
func (db *DB) CloseCayley(context interface{}) {
	db.graphHandle.Close()
}
