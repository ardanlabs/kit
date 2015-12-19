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
	// Setting this here but should be done outside the app.
	os.Setenv("ENV_PREFIX_LOGGING_LEVEL", "1")

	// If you want to add a custom header to every request.
	// Add the following env variable:
	// ENV_PRRFIX_HEADERS=key:value,key:value

	app.Init("ENV_PREFIX")
}

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {
	a := app.New()

	// If you want to include middleware such as basic authentication.
	// a := app.New(midware.Auth)

	// If you have configured auth and want to turn it off.
	// ENV_PRRFIX_AUTH=false

	// Initialize the routes for the API.

	// http://localhost:4000/1.0/test/names
	a.Handle("GET", "/1.0/test/names", handlers.Test.List)

	return a
}
