package middleware

import (
	"context"
	"net/http"

	"github.com/ardanlabs/kit/examples/web/internal/sys/app"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// Auth simulates a midware for authentication.
func Auth(next web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) {
		v := ctx.Value(app.KeyValues).(*app.Values)

		log.Dev(v.TraceID, "Auth", "******> Authorized")
		v.Auth = "1234"

		next(ctx, w, r, params)
	}
}
