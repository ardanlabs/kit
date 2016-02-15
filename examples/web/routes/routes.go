package routes

import (
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/examples/web/handlers"

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

	// Look at /kit/web/midware for middleware options.
	a := app.New()

	// Initialize the routes for the API.
	a.Handle("GET", "/1.0/test/names", handlers.Test.List)

	return a
}
