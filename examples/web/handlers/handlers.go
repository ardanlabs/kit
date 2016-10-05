// Package handlers contains the handler logic for processing requests.
package handlers

import (
	"net/http"

	"github.com/ardanlabs/kit/web"
)

// testHandle maintains the set of handlers for the test api.
type testHandle struct{}

// Test fronts the access to the test service functionality.
var Test testHandle

//==============================================================================

// List returns all the existing test names in the system.
// 200 Success, 404 Not Found, 500 Internal
func (testHandle) List(c *web.Context) error {
	names := []string{"Apple", "Orange", "Banana"}

	c.Respond(names, http.StatusOK)
	return nil
}
