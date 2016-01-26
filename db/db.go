// Package db abstracts different database systems we can use. I want to be
// able to access the raw database support so an interface does not work. Each
// database is too different.
package db

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
)

// DB provides access to a session that is already tied to a particular
// database for use.
type DB struct {
	database *mgo.Database
	session  *mgo.Session
}

//==============================================================================
// MongoDB support

// NewMGO returns a new DB value for use with MongoDB based on a registered
// master session.
func NewMGO(context interface{}, name string) (*DB, error) {
	var db mgoDB
	var exists bool
	masterMGO.Lock()
	{
		db, exists = masterMGO.ses[name]
	}
	masterMGO.Unlock()

	if !exists {
		return nil, fmt.Errorf("Master sesssion %q does not exist", name)
	}

	ses := db.ses.Copy()
	mdb := ses.DB(db.dbName)

	return &DB{mdb, ses}, nil
}

// CloseMGO closes a DB value being used with MongoDB.
func (db *DB) CloseMGO(context interface{}) {
	db.session.Close()
}

// ExecuteMGO is used to execute MongoDB commands.
func (db *DB) ExecuteMGO(context interface{}, colName string, f func(*mgo.Collection) error) error {
	if db == nil || db.session == nil {
		return errors.New("Invalid DB provided")
	}

	return f(db.database.C(colName))
}

// CollectionMGO is used to get a collection value..
func (db *DB) CollectionMGO(context interface{}, colName string) (*mgo.Collection, error) {
	if db == nil || db.session == nil {
		return nil, errors.New("Invalid DB provided")
	}

	return db.database.C(colName), nil
}
