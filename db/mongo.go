package db

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ardanlabs/kit/db/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Multiple master sessions can be created for different databases. Each master
// session must be registered first. A new copy of the master session can be
// acquired. The session is provided through the DB type value.

// ErrInvalidDBProvided is returned in the event that an uninitialized db is
// used to perform actions against.
var ErrInvalidDBProvided = errors.New("Invalid DB provided")

//==============================================================================

// mgoDB maintains a master session for a given database.
type mgoDB struct {
	ses *mgo.Session
}

// masterMGO manages a set of different MongoDB master sessions.
var masterMGO = struct {
	sync.RWMutex
	ses map[string]mgoDB
}{
	ses: make(map[string]mgoDB),
}

// RegMasterSession adds a new master session to the set. If no url is provided,
// it will default to localhost:27017.
func RegMasterSession(context interface{}, name string, url string, timeout time.Duration) error {
	masterMGO.Lock()
	defer masterMGO.Unlock()

	if _, exists := masterMGO.ses[name]; exists {
		return errors.New("Master session already exists")
	}

	ses, err := mongo.New(url, timeout)
	if err != nil {
		return err
	}

	masterMGO.ses[name] = mgoDB{
		ses: ses,
	}

	return nil
}

//==============================================================================
// Factory function for acquiring a copy of a session.

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

	// As per the mgo documentation, if no database name is specified, then use
	// the default one, or the one that the connection was dialed with.
	mdb := ses.DB("")

	dbOut := DB{
		database: mdb,
		session:  ses,
	}

	return &dbOut, nil
}

//==============================================================================
// Methods for the DB struct type related to MongoDB.

// CloseMGO closes a DB value being used with MongoDB.
func (db *DB) CloseMGO(context interface{}) {
	db.session.Close()
}

// ExecuteMGO is used to execute MongoDB commands.
func (db *DB) ExecuteMGO(context interface{}, colName string, f func(*mgo.Collection) error) error {
	if db == nil || db.session == nil {
		return ErrInvalidDBProvided
	}

	return f(db.database.C(colName))
}

// ExecuteMGOTimeout is used to execute MongoDB commands with a timeout.
func (db *DB) ExecuteMGOTimeout(context interface{}, timeout time.Duration, colName string, f func(*mgo.Collection) error) error {
	if db == nil || db.session == nil {
		return ErrInvalidDBProvided
	}

	db.session.SetSocketTimeout(timeout)

	return f(db.database.C(colName))
}

// BatchedQueryMGO returns an iterator capable of iterating over
// all the results of a query in batches.
func (db *DB) BatchedQueryMGO(context interface{}, colName string, q bson.M) (*mgo.Iter, error) {
	if db == nil || db.session == nil {
		return nil, ErrInvalidDBProvided
	}

	c := db.database.C(colName)

	return c.Find(q).Iter(), nil
}

// BulkOperationMGO returns a bulk value that allows multiple orthogonal
// changes to be delivered to the server.
func (db *DB) BulkOperationMGO(context interface{}, colName string) (*mgo.Bulk, error) {
	if db == nil || db.session == nil {
		return nil, ErrInvalidDBProvided
	}

	c := db.database.C(colName)
	tx := c.Bulk()
	tx.Unordered()

	return tx, nil
}

// CollectionMGO is used to get a collection value.
func (db *DB) CollectionMGO(context interface{}, colName string) (*mgo.Collection, error) {
	if db == nil || db.session == nil {
		return nil, ErrInvalidDBProvided
	}

	return db.database.C(colName), nil
}

// CollectionMGOTimeout is used to get a collection value with a timeout.
func (db *DB) CollectionMGOTimeout(context interface{}, timeout time.Duration, colName string) (*mgo.Collection, error) {
	if db == nil || db.session == nil {
		return nil, ErrInvalidDBProvided
	}

	db.session.SetSocketTimeout(timeout)

	return db.database.C(colName), nil
}
