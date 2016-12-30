package middleware

import (
	"context"
	"net/http"

	"github.com/ardanlabs/kit/examples/web/internal/sys/app"
	"github.com/ardanlabs/kit/web"
)

// Values adds the system values to the context.
func Values(next web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)

		values := app.Values{
			Values: v,
		}
		ctx = context.WithValue(ctx, app.KeyValues, &values)

		return next(ctx, w, r, params)
	}
}
