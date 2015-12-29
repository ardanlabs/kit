// Sample program to show how to use the udp package to build servers that
// can accept udp connections and send messages.
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/udp"
)

// Configuation settings.
const (
	configKey       = "KIT"
	cfgLoggingLevel = "LOGGING_LEVEL"
)

func init() {

	// This is being added to showcase configuration.
	os.Setenv("KIT_LOGGING_LEVEL", "1")

	// Init the configuration system.
	if err := cfg.Init(configKey); err != nil {
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
	log.Init(os.Stderr, logLevel)

	// Log all the configuration options
	log.User("startup", "init", "\n\nConfig Settings: %s\n%s\n", configKey, cfg.Log())
}

//==============================================================================

func main() {
	const context = "startup"

	// Create the configuration.
	cfg := udp.Config{
		NetType: "udp4",
		Addr:    ":6000",

		ConnHandler: udpConnHandler{},
		ReqHandler:  udpReqHandler{},
		RespHandler: udpRespHandler{},

		OptIntPool: udp.OptIntPool{
			RecvMinPoolSize: func() int { return 2 },
			RecvMaxPoolSize: func() int { return 100 },
			SendMinPoolSize: func() int { return 2 },
			SendMaxPoolSize: func() int { return 100 },
		},
	}

	// Create a new UDP value.
	u, err := udp.New(context, "Sample", cfg)
	if err != nil {
		log.Error(context, "main", err, "Creating udp")
		return
	}

	// Start accepting client data.
	if err := u.Start(context); err != nil {
		log.Error(context, "main", err, "Starting udp")
		return
	}

	// Defer the stop on shutdown.
	defer u.Stop(context)

	log.User(context, "main", "Waiting for data on: %s", u.Addr())

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Use netcat to test the server.
	// nc -4u localhost 6000 < test.hex

	log.User(context, "main", "Shutting down")
}
