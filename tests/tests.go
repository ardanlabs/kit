// Package tests provides the generic support all tests require.
package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
)

// Context provides a base context for tests.
var Context = "Test"

// TestSession is the name used to register the MongoDB session.
var TestSession = "test"

// Success and failure markers.
var (
	Success = "\u2713"
	Failed  = "\u2717"
)

// logdash is the central buffer where all logs are stored.
var logdash bytes.Buffer
var loglock sync.RWMutex

//==============================================================================

// ResetLog resets the contents of logdash.
func ResetLog() {
	loglock.Lock()
	defer loglock.Unlock()
	logdash.Reset()
}

// DisplayLog writes the logdash data to standand out, if testing in verbose mode
// was turned on.
func DisplayLog() {
	if !testing.Verbose() {
		return
	}

	loglock.RLock()
	defer loglock.RUnlock()
	logdash.WriteTo(os.Stdout)
}

// Init initializes the log package.
func Init(cfgKey string) {
	cfg.Init(cfg.EnvProvider{Namespace: cfgKey})

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(&logdash, logLevel)
}

// InitMongo initializes the mongodb connections for testing.
func InitMongo(cfg mongo.Config) {
	if err := db.RegMasterSession("Test", TestSession, cfg); err != nil {
		log.Error("Test", "Init", err, "Completed")
		logdash.WriteTo(os.Stdout)
		os.Exit(1)
	}
}

// NewRequest used to setup a request for mocking API calls with httptreemux.
func NewRequest(method, path string, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, path, body)
	u, _ := url.Parse(path)
	r.URL = u
	r.RequestURI = path

	return r
}
