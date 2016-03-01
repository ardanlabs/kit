package midware

import (
	"github.com/ardanlabs/kit/auth"
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
)

// Auth config environmental variables.
const cfgAuth = "AUTH"

// Auth handles token authentication.
func Auth(h app.Handler) app.Handler {

	// Check if authentication is turned off.
	on, err := cfg.Bool(cfgAuth)
	if err == nil && !on {
		return func(c *app.Context) error {
			log.Dev(c.SessionID, "Auth", "******> Authentication Off")
			return h(c)
		}
	}

	// Turn authentication on.
	return func(c *app.Context) error {
		token := c.Request.Header.Get("Authorization")
		log.Dev(c.SessionID, "Auth", "Started : Token[%s]", token)

		if len(token) < 5 || token[0:5] != "Basic" {
			log.Error(c.SessionID, "Auth", app.ErrNotAuthorized, "Validating token")
			return app.ErrNotAuthorized
		}

		var err error
		if c.Ctx["User"], err = auth.ValidateWebToken(c.SessionID, c.Ctx["DB"].(*db.DB), token[6:]); err != nil {
			log.Error(c.SessionID, "Auth", err, "Validating token")
			return app.ErrNotAuthorized
		}

		log.Dev(c.SessionID, "Auth", "Completed")
		return h(c)
	}
}
