// Package session provides a level of security for web tokens by giving them an
// expiration date and a lookup point for the user accessing the API. The session
// is what is used to look up the user performing the web call. Though the user's
// PublicID is not private, it is not used directly.
package session

import (
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the user collection.
const Collection = "auth_sessions"

//==============================================================================

// Create adds a new session for the specified user to the database.
func Create(context interface{}, db *db.DB, publicID string, expires time.Duration) (*Session, error) {
	log.Dev(context, "Create", "Started : PublicID[%s]", publicID)

	s := Session{
		SessionID:   uuid.New(),
		PublicID:    publicID,
		DateExpires: time.Now().Add(expires),
		DateCreated: time.Now(),
	}

	f := func(c *mgo.Collection) error {
		log.Dev(context, "Create", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(&s))
		return c.Insert(&s)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Create", "Completed : SessionID[%s]", s.SessionID)
	return &s, nil
}

//==============================================================================

// GetBySessionID retrieves a session from the session store.
func GetBySessionID(context interface{}, db *db.DB, sessionID string) (*Session, error) {
	log.Dev(context, "GetBySessionID", "Started : SessionID[%s]", sessionID)

	var s Session
	f := func(c *mgo.Collection) error {
		q := bson.M{"session_id": sessionID}
		log.Dev(context, "GetBySessionID", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&s)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "GetBySessionID", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetBySessionID", "Completed")
	return &s, nil
}

// GetByLatest retrieves the latest session for the specified user.
func GetByLatest(context interface{}, db *db.DB, publicID string) (*Session, error) {
	log.Dev(context, "GetByLatest", "Started : PublicID[%s]", publicID)

	var s Session
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		log.Dev(context, "GetByLatest", "MGO : db.%s.find(%s).sort({\"date_created\": -1}).limit(1)", c.Name, mongo.Query(q))
		return c.Find(q).Sort("-date_created").One(&s)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "GetByLatest", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetByLatest", "Completed")
	return &s, nil
}
