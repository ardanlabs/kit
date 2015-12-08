package session

import (
	"time"

	"github.com/ardanlabs/kit/log"

	"gopkg.in/mgo.v2/bson"
)

// Session denotes a user's session within the system.
type Session struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	SessionID   string        `bson:"session_id" json:"session_id"`
	PublicID    string        `bson:"public_id" json:"public_id"`
	DateExpires time.Time     `bson:"date_expires" json:"date_expires"`
	DateCreated time.Time     `bson:"date_created" json:"date_created"`
}

// IsExpired returns true if the session is expired.
func (s *Session) IsExpired(context interface{}) bool {
	if s.DateExpires.Before(time.Now()) {
		log.Dev(context, "IsExpired", "Expired : Date[%v] Duration[%v]", s.DateExpires, time.Since(s.DateExpires))
		return true
	}
	return false
}
