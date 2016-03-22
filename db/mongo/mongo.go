// Package mongo provides support for using MongoDB.
package mongo

import (
	"encoding/json"
	"strings"
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

	// Can be provided a comma delimited set of hosts.
	hosts := strings.Split(cfg.Host, ",")

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := mgo.DialInfo{
		Addrs:    hosts,
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

// Query provides a string version of the value
func Query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}
