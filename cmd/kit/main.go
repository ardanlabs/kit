// This program provides the kit project a basic cli tool to manage entities.
package main

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/cmd/kit/cmdauth"
	"github.com/ardanlabs/kit/cmd/kit/cmddb"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"github.com/spf13/cobra"
)

// Config environmental variables.
const (
	cfgLoggingLevel  = "LOGGING_LEVEL"
	cfgMongoHost     = "MONGO_HOST"
	cfgMongoAuthDB   = "MONGO_AUTHDB"
	cfgMongoDB       = "MONGO_DB"
	cfgMongoUser     = "MONGO_USER"
	cfgMongoPassword = "MONGO_PASS"
)

var kit = &cobra.Command{
	Use:   "kit",
	Short: "Kit provides the central command housing for the kit tooling.",
}

//==============================================================================

func main() {
	if err := cfg.Init("KIT"); err != nil {
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
	log.Init(os.Stderr, logLevel)

	cfg := mongo.Config{
		Host:     cfg.MustString(cfgMongoHost),
		AuthDB:   cfg.MustString(cfgMongoAuthDB),
		DB:       cfg.MustString(cfgMongoDB),
		User:     cfg.MustString(cfgMongoUser),
		Password: cfg.MustString(cfgMongoPassword),
	}

	if err := mongo.Init(cfg); err != nil {
		kit.Println("Unable to initialize MongoDB")
		os.Exit(1)
	}

	kit.AddCommand(cmdauth.GetCommands(), cmddb.GetCommands())
	kit.Execute()
}
