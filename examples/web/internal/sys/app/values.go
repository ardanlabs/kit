package app

import (
	"github.com/ardanlabs/kit/web"
	"gopkg.in/mgo.v2"
)

// Key represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

//==============================================================================

// Values extends the kit's value type stored inside
// the context.
type Values struct {
	*web.Values

	DB   *mgo.Session
	Auth string
}
