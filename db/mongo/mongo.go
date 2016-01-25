// Package mongo provides support for using MongoDB.
package mongo

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

// Config provides configuration values.
type Config struct {
	Host     string
	AuthDB   string
	DB       string
	User     string
	Password string
}

//==============================================================================

// New creates a new master session.
func New(cfg Config) (*mgo.Session, error) {

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := mgo.DialInfo{
		Addrs:    []string{cfg.Host},
		Timeout:  60 * time.Second,
		Database: cfg.AuthDB,
		Username: cfg.User,
		Password: cfg.Password,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := mgo.DialWithInfo(&mongoDBDialInfo)
	if err != nil {
		return nil, err
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	ses.SetMode(mgo.Monotonic, true)

	return ses, nil
}

//==============================================================================

// Holds global state for mongo access.
var m struct {
	dbName string
	ses    *mgo.Session
	mu     sync.RWMutex
}

// Init sets up the MongoDB environment. This expects that the
// cfg package has been initialized first.
func Init(cfg Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ses != nil {
		return nil
	}

	ses, err := New(cfg)
	if err != nil {
		return err
	}

	// Save the settings for this single connection.
	m.ses = ses
	m.dbName = cfg.DB

	return nil
}

// Query provides a string version of the value
func Query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}

// GetSession returns a copy of the master session for use.
func GetSession() *mgo.Session {
	var ses *mgo.Session
	m.mu.RLock()
	{
		ses = m.ses.Copy()
	}
	m.mu.RUnlock()

	return ses
}

// GetDatabase returns a mgo database value based on configuration.
func GetDatabase(ses *mgo.Session) *mgo.Database {
	var db *mgo.Database
	m.mu.RLock()
	{
		db = ses.DB(m.dbName)
	}
	m.mu.RUnlock()

	return db
}

// GetDatabaseName returns the name of the database being used.
func GetDatabaseName() string {
	var name string
	m.mu.RLock()
	{
		name = m.dbName
	}
	m.mu.RUnlock()

	return name
}

// GetCollection returns a mgo collection value based on configuration.
func GetCollection(ses *mgo.Session, colName string) *mgo.Collection {
	var col *mgo.Collection
	m.mu.RLock()
	{
		col = ses.DB(m.dbName).C(colName)
	}
	m.mu.RUnlock()

	return col
}

// ExecuteDB the MongoDB literal function.
func ExecuteDB(context interface{}, ses *mgo.Session, collectionName string, f func(*mgo.Collection) error) error {

	// Validate we have a valid session.
	if ses == nil {
		return errors.New("Invalid session provided")
	}

	var col *mgo.Collection
	m.mu.RLock()
	{
		col = ses.DB(m.dbName).C(collectionName)
	}
	m.mu.RUnlock()

	// Execute the MongoDB call.
	return f(col)
}

// CollectionExists returns true if the collection name exists in the specified database.
func CollectionExists(context interface{}, ses *mgo.Session, useCollection string) bool {

	// Validate we have a valid session.
	if ses == nil {
		return false
	}

	var db *mgo.Database
	m.mu.RLock()
	{
		db = ses.DB(m.dbName)
	}
	m.mu.RUnlock()

	// Capture the list of collection names.
	cols, err := db.CollectionNames()
	if err != nil {
		return false
	}

	// Find it in the list.
	for _, col := range cols {
		if col == useCollection {
			return true
		}
	}

	return false
}
