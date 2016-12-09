// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// This program provides a sample web service that implements a
// RESTFul CRUD API against a MongoDB database.
package main

import (
	"os"
	"os/signal"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/examples/web/cmd/app/routes"
	"github.com/ardanlabs/kit/examples/web/internal/sys/app"
	"github.com/ardanlabs/kit/log"
	"github.com/braintree/manners"
)

// init is called before main. We are using init to customize logging output.
func init() {

	// This is being added to showcase configuration.
	os.Setenv("KIT_LOGGING_LEVEL", "1")
}

//==============================================================================

// main is the entry point for the application.
func main() {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	app.Init("startup", cfg.EnvProvider{Namespace: app.Namespace})

	// Check the environment for a configured port value.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create this goroutine to run the web server.
	go func() {
		log.User("startup", "main", "Started : Listening on: http://localhost:"+port)
		manners.ListenAndServe(":"+port, routes.API())
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.User("shutdown", "main", "Shutting down...")
	manners.Close()

	log.User("shutdown", "main", "Down")
}
