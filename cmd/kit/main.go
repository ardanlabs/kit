// This program provides the kit project a basic cli tool to manage entities.
package main

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/cmd/kit/cmdauth"
	"github.com/ardanlabs/kit/cmd/kit/cmddb"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"

	"github.com/spf13/cobra"
)

// Config environmental variables.
const (
	cfgNamespace    = "KIT"
	cfgLoggingLevel = "LOGGING_LEVEL"
	cfgMongoURI     = "MONGO_URI"
)

var kit = &cobra.Command{
	Use:   "kit",
	Short: "Kit provides the central command housing for the kit tooling.",
}

//==============================================================================

func main() {
	if err := cfg.Init(cfg.EnvProvider{Namespace: cfgNamespace}); err != nil {
		kit.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.NONE
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	cfg := cfg.MustURL(cfgMongoURI)

	// Here we use the path of the mongo uri as the master session name, as we
	// just need to specify a unique identifier for this session as it has no real
	// relation to the actual database name.
	if err := db.RegMasterSession("startup", cfg.Path, cfg.String(), 0); err != nil {
		kit.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	// Here we will load the session out of the master session using the
	// unique identifier as above (the path).
	db, err := db.NewMGO("", cfg.Path)
	if err != nil {
		kit.Println("Unable to get MongoDB session")
		os.Exit(1)
	}
	defer db.CloseMGO("")

	kit.AddCommand(
		cmdauth.GetCommands(db),
		cmddb.GetCommands(db),
	)
	kit.Execute()
}
