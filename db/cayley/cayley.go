// Package cayley provides support for the cayley Graph DB with a Mongo backend.
package cayley

import (
	"net/url"
	"strings"

	"github.com/cayleygraph/cayley"

	// Blank import the mongo library for cayley.
	_ "github.com/cayleygraph/cayley/graph/mongo"
)

//==============================================================================

// New creates a new cayley handle.
func New(mongoURL string) (*cayley.Handle, error) {

	cfg, err := url.Parse(mongoURL)
	if err != nil {
		return nil, err
	}

	// Form the Cayley connection options.
	opts := make(map[string]interface{})

	// Load the database name from the path, but the path will contain the
	// leading slash as well.
	opts["database_name"] = strings.TrimPrefix(cfg.Path, "/")

	if cfg.User != nil {
		if password, ok := cfg.User.Password(); ok {
			opts["password"] = password
		}

		opts["username"] = cfg.User.Username()
	}

	// Create the cayley handle that maintains a connection to the
	// Cayley graph database in Mongo.
	store, err := cayley.NewGraph("mongo", cfg.Host, opts)
	if err != nil {
		return store, err
	}

	return store, nil
}
