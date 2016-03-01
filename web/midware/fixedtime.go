package midware

import (
	"time"

	"github.com/ardanlabs/kit/web/app"
)

// FixedTime is useful when testing time dependent APIs as it sets the Now
// property of all Contexts to a fixed time.
func FixedTime(now time.Time) func(h app.Handler) app.Handler {
	return func(h app.Handler) app.Handler {
		return func(c *app.Context) error {
			log.Dev(c.SessionID, "FixedTime", "Setting context time to fixed time %v", now)
			c.Now = now

			return h(c)
		}
	}
}
