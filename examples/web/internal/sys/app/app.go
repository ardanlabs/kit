package app

import (
	"fmt"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
)

const (

	// Namespace is the key that is the prefix for configuration in the
	// environment.
	Namespace = "KIT"

	// CFGLoggingLevel is the logging level to use.
	CFGLoggingLevel = "LOGGING_LEVEL"
)

//==============================================================================

// Init sets up the configuration and logging systems.
func Init(traceID string, p cfg.Provider) {

	// Init the configuration system.
	if err := cfg.Init(p); err != nil {
		fmt.Println("error initalizing configuration system", err)
		os.Exit(1)
	}

	// Init the log system.
	logLevel := func() int {
		ll, err := cfg.Int(CFGLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	// Log all the configuration options
	log.User(traceID, "main", "\n\nConfig Settings:\n%s\n", cfg.Log())
}
