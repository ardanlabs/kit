// This program provides the kit project a basic cli tool to manage entities.
package main

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/cli/cmdauth"
	"github.com/ardanlabs/kit/cli/cmddb"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"github.com/spf13/cobra"
)

var kit = &cobra.Command{
	Use:   "kit",
	Short: "Kit provides the central cli housing for the kit tooling.",
}

func main() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.NONE
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	if err := cfg.Init("KIT"); err != nil {
		kit.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	err := mongo.Init()
	if err != nil {
		kit.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	kit.AddCommand(cmdauth.GetCommands(), cmddb.GetCommands())
	kit.Execute()
}
