package midware

import (
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// DB simulates a midware for connections.
func DB(h web.Handler) web.Handler {

	// Pretend this get a DB connection
	f := func(c *web.Ctx) error {
		log.Dev(c.SessionID, "DB", "******> Capture DB Connection")
		c.Values["DB"] = "CONN"
		defer func() {
			log.Dev(c.SessionID, "DB", "******> Release DB Connection")
		}()

		return h(c)
	}

	return f
}
