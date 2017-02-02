// Sample program to show how to use the tcp package to build servers that
// can accept tcp connections and send messages.
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/tcp"
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
	if err := cfg.Init(cfg.EnvProvider{Namespace: configKey}); err != nil {
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
	log.User("startup", "init", "\n\nConfig Settings: %s\n%s\n", configKey, cfg.Log())
}

// Event writes tcp events.
func Event(event string, format string, a ...interface{}) {
	log.User("*EVENT*", event, format, a...)
}

func main() {
	const traceID = "startup"

	// Create the configuration.
	cfg := tcp.Config{
		NetType: "tcp4",
		Addr:    ":6000",

		ConnHandler: tcpConnHandler{},
		ReqHandler:  tcpReqHandler{},
		RespHandler: tcpRespHandler{},

		OptEvent: tcp.OptEvent{
			Event: Event,
		},
	}

	// Create a new TCP value.
	t, err := tcp.New("Sample", cfg)
	if err != nil {
		log.Error(traceID, "main", err, "Creating tcp")
		return
	}

	// Start accepting client data.
	if err := t.Start(); err != nil {
		log.Error(traceID, "main", err, "Starting tcp")
		return
	}

	// Defer the stop on shutdown.
	defer t.Stop()

	log.User(traceID, "main", "Waiting for data on: %s", t.Addr())

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Use telnet to test the server.
	// telnet localhost 6000

	log.User(traceID, "main", "Shutting down")
}
