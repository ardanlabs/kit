// Package mongo provides support for using MongoDB.
package mongo

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2"
)

//==============================================================================

// New creates a new master session. If no url is provided, it will defaul to
// localhost:27017. If a zero value timeout is specified, a timeout of 60sec
// will be used instead.
func New(url string, timeout time.Duration) (*mgo.Session, error) {

	// Set the default timeout for the session.
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := mgo.DialWithTimeout(url, timeout)
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
