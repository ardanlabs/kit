// Package db abstracts different database systems we can use.
package db

import (
	"gopkg.in/mgo.v2"
)

// DB is a collection of support for different DB technologies. Currently
// only MongoDB has been implemented. We want to be able to access the raw
// database support for the given DB so an interface does not work. Each
// database is too different.
type DB struct {

	// MongoDB Support
	database *mgo.Database
	session  *mgo.Session
}
