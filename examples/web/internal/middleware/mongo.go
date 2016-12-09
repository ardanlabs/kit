package middleware

import (
	"context"
	"net/http"

	"github.com/ardanlabs/kit/examples/web/internal/sys/app"
	"github.com/ardanlabs/kit/examples/web/internal/sys/db"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// Mongo initializes the master session and wires in the connection middleware.
func Mongo() web.Middleware {

	// session contains the master session for accessing MongoDB.
	session := db.Init()

	// Return this middleware to be chained together.
	return func(next web.Handler) web.Handler {

		// Wrap this handler around the next one provided.
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) {
			v := ctx.Value(app.KeyValues).(*app.Values)

			// Get a MongoDB session connection.
			log.User(v.TraceID, "Mongo", "*****> Capture Mongo Session")
			v.DB = session.Copy()

			// Defer releasing the db session connection.
			defer func() {
				log.User(v.TraceID, "Mongo", "*****> Release Mongo Session")
				v.DB.Close()
			}()

			next(ctx, w, r, params)
		}
	}
}
