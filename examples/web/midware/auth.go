package midware

import (
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
)

// Auth simulates a midware for authentication.
func Auth(h app.Handler) app.Handler {

	// Pretend this get a DB connection
	f := func(c *app.Context) error {
		log.Dev(c.SessionID, "Auth", "******> Authorized")
		c.Ctx["UserID"] = "123"

		return h(c)
	}

	return f
}
