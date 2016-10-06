package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/examples/web/handlers"
	"github.com/ardanlabs/kit/examples/web/midware"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

// Configuation settings.
const (
	configKey       = "KIT"
	cfgLoggingLevel = "LOGGING_LEVEL"
)

func init() {

	// This is being added to showcase configuration.
	os.Setenv("KIT_LOGGING_LEVEL", "1")

	Init(cfg.EnvProvider{Namespace: configKey})
}

// Init is called to initialize the application.
func Init(p cfg.Provider) {

	// Init the configuration system.
	if err := cfg.Init(p); err != nil {
		fmt.Println("Error initalizing configuration system", err)
		os.Exit(1)
	}

	// Init the log system.
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	// Log all the configuration options
	log.User("startup", "Init", "\n\nConfig Settings:\n%s\n", cfg.Log())
}

//==============================================================================

// API returns a handler for a set of routes.
func API() http.Handler {

	// Look at /kit/web/midware for middleware options and
	// patterns for writing middleware.
	a := web.New(midware.DB)

	// Set a handler that only needs DB Connection.
	a.Handle("GET", "/v1/test/noauth", handlers.Test.List)

	// Create a group for handlers that need auth as well.
	ag := a.Group(midware.Auth)
	ag.Handle("GET", "/v1/test/names", handlers.Test.List)

	return a
}
