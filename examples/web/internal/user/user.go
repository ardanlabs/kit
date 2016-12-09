// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package user

import (
	"context"

	mgo "gopkg.in/mgo.v2"

	"github.com/ardanlabs/kit/examples/web/internal/sys/db"
	"github.com/ardanlabs/kit/log"
)

const usersCollection = "users"

//==============================================================================

// List retrieves a list of existing users from the database.
func List(ctx context.Context, traceID string, dbSes *mgo.Session) ([]User, error) {
	log.User(traceID, "List", "Started")

	u := []User{}
	f := func(collection *mgo.Collection) error {
		log.User(traceID, "List", "MGO :\n\ndb.users.find()\n\n")
		return collection.Find(nil).All(&u)
	}

	if err := db.Execute(dbSes, usersCollection, f); err != nil {
		log.Error(traceID, "List", err, "Executing DB")
		return nil, err
	}

	log.User(traceID, "List", "Completed")
	return u, nil
}
