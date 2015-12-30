// Package db provides a thin layer of abstraction for any database system
// being use. This will allow service layer API to remain consistent.
package db

import (
	"errors"

	"github.com/ardanlabs/kit/db/mongo"

	"gopkg.in/mgo.v2"
)

// DB abstracts different database systems we can use.
type DB struct {
	MGOConn *mgo.Session
}

// NewMGO return a new DB value for use with MongoDB.
func NewMGO() *DB {
	return &DB{
		MGOConn: mongo.GetSession(),
	}
}

// CloseMGO closes a DB value being used with MongoDB.
func (db *DB) CloseMGO() {
	db.MGOConn.Close()
}

// ExecuteMGO is used to execute MongoDB commands.
func (db *DB) ExecuteMGO(context interface{}, collection string, f func(*mgo.Collection) error) error {
	if db == nil || db.MGOConn == nil {
		return errors.New("Invalid DB provided")
	}

	if err := mongo.ExecuteDB(context, db.MGOConn, collection, f); err != nil {
		return err
	}

	return nil
}
