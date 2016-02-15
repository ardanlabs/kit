package midware

import (
	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
)

// cfgMongoDB config environmental variables.
const cfgMongoDB = "MONGO_DB"

// Mongo handles session management.
func Mongo(h app.Handler) app.Handler {

	// Check if mongodb is configured.
	dbName, err := cfg.String(cfgMongoDB)
	if err != nil {
		return func(c *app.Context) error {
			log.Dev(c.SessionID, "Mongo", "******> Mongo Not Configured")
			return h(c)
		}
	}

	// Wrap the handlers inside a session copy/close.
	return func(c *app.Context) error {
		mgoDB, err := db.NewMGO("Mongo", dbName)
		if err != nil {
			log.Error(c.SessionID, "Mongo", err, "Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)
			return app.ErrDBNotConfigured
		}

		log.Dev(c.SessionID, "Mongo", "******> Capture Mongo Session")
		c.Ctx["DB"] = mgoDB
		defer func() {
			log.Dev(c.SessionID, "Mongo", "******> Release Mongo Session")
			mgoDB.CloseMGO("Mongo")
		}()

		return h(c)
	}
}
