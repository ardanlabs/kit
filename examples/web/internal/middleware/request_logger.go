// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/ardanlabs/kit/examples/web/internal/sys/app"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// RequestLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(next web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) {
		v := ctx.Value(app.KeyValues).(*app.Values)

		start := time.Now()
		next(ctx, w, r, params)

		log.User(v.TraceID, "RL", "(%d) : %s %s -> %s (%s)",
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(start),
		)
	}
}
