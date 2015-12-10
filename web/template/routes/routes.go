package routes

import (
	"net/http"

	"github.com/ardanlabs/kit/web/app"
	"github.com/ardanlabs/kit/web/template/handlers"

	// If you want to include middleware such as basic authentication.
	// "github.com/ardanlabs/kit/web/midware"
)

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
