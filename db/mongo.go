package db

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ardanlabs/kit/db/mongo"

	"gopkg.in/mgo.v2"
)

// Multiple master sessions can be created for different databases. Each master
// session must be registered first. A new copy of the master session can be
// acquired. The session is provided through the DB type value.

//==============================================================================

// mgoDB maintains a master session for a given database.
type mgoDB struct {
	dbName string
	ses    *mgo.Session
}

// masterMGO manages a set of different MongoDB master sessions.
var masterMGO = struct {
	sync.RWMutex
	ses map[string]mgoDB
}{
	ses: make(map[string]mgoDB),
}

// RegMasterSession adds a new master session to the set.
func RegMasterSession(context interface{}, name string, cfg mongo.Config) error {
	masterMGO.Lock()
	defer masterMGO.Unlock()

	if _, exists := masterMGO.ses[name]; exists {
		return errors.New("Master session already exists")
	}

	ses, err := mongo.New(cfg)
	if err != nil {
		return err
	}

	masterMGO.ses[name] = mgoDB{cfg.DB, ses}

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
	mdb := ses.DB(db.dbName)

	return &DB{mdb, ses}, nil
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
