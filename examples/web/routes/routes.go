package routes

import (
	"net/http"
	"os"

	"github.com/ardanlabs/kit/examples/web/handlers"

	"github.com/ardanlabs/kit/web/app"

	// If you want to include middleware such as basic authentication.
	// "github.com/ardanlabs/kit/web/midware"
)

func init() {
	os.Setenv("ENV_PREFIX_LOGGING_LEVEL", "1")

	set := app.Settings{
		ConfigKey: "ENV_PREFIX",
		UseMongo:  false,
	}

	app.Init(&set)
}

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {
	a := app.New()

	// If you want to include middleware such as basic authentication.
	// a := app.New(midware.Auth)

	// Initialize the routes for the API.

	// http://localhost:4000/1.0/test/names
	a.Handle("GET", "/1.0/test/names", handlers.Test.List)

	return a
}
