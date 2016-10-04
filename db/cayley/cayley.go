// Package cayley provides support for the cayley Graph DB with a Mongo backend.
package cayley

import (
	"github.com/cayleygraph/cayley"

	// Blank import the mongo library for cayley.
	_ "github.com/cayleygraph/cayley/graph/mongo"
)

// Config provides configuration values.
type Config struct {
	Host     string
	DB       string
	User     string
	Password string
}

//==============================================================================

// New creates a new cayley handle.
func New(cfg Config) (*cayley.Handle, error) {

	// Form the Cayley connection options.
	opts := map[string]interface{}{
		"database_name": cfg.DB,
		"username":      cfg.User,
		"password":      cfg.Password,
	}

	// Create the cayley handle that maintains a connection to the
	// Cayley graph database in Mongo.
	store, err := cayley.NewGraph("mongo", cfg.Host, opts)
	if err != nil {
		return store, err
	}

	return store, nil
}
