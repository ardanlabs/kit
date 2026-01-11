package routes

import (
	"net/http"

	"github.com/kurisu1024/kit/examples/web/cmd/app/handlers"
	"github.com/kurisu1024/kit/examples/web/internal/middleware"
	"github.com/kurisu1024/kit/web"
)

// API returns a handler for a set of routes.
func API() http.Handler {

	// Look at /kit/web/midware for middleware options and
	// patterns for writing middleware.
	a := web.New(middleware.Values, middleware.RequestLogger, middleware.Mongo())

	// Set a handler that only needs DB Connection.
	a.Handle("GET", "/v1/test/noauth", handlers.UserList)

	// Create a group for handlers that need auth as well.
	ag := a.Group(middleware.Auth)
	ag.Handle("GET", "/v1/test/names", handlers.UserList)

	return a
}
