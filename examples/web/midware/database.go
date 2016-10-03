package midware

import (
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
)

// DB simulates a midware for connections.
func DB(h app.Handler) app.Handler {

	// Pretend this get a DB connection
	f := func(c *app.Context) error {
		log.Dev(c.SessionID, "DB", "******> Capture DB Connection")
		c.Ctx["DB"] = "CONN"
		defer func() {
			log.Dev(c.SessionID, "DB", "******> Release DB Connection")
		}()

		return h(c)
	}

	return f
}
