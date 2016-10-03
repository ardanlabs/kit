package routes

import (
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/examples/web/handlers"
	"github.com/ardanlabs/kit/examples/web/midware"
	"github.com/ardanlabs/kit/web/app"
)

// Configuation settings.
const configKey = "KIT"

func init() {

	// This is being added to showcase configuration.
	os.Setenv("KIT_LOGGING_LEVEL", "1")

	app.Init(cfg.EnvProvider{Namespace: configKey})
}

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {

	// Look at /kit/web/midware for middleware options and
	// patterns for writing middleware.
	a := app.New(midware.DB)

	// Set a handler that only needs DB Connection.
	a.Handle("GET", "/v1/test/noauth", handlers.Test.List)

	// Create a group for handlers that need auth as well.
	ag := a.Group(midware.Auth)
	ag.Handle("GET", "/v1/test/names", handlers.Test.List)

	return a
}
