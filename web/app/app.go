// Package app provides a thin layer of support for writing web services. It
// integrates with the ardanlabs kit repo to provide support for logging,
// configuration, database, routing and application context. The base things
// you need to write a web service is provided.
package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/braintree/manners"
	"github.com/dimfeld/httptreemux"
	"github.com/pborman/uuid"
)

// Web config environmental variables.
const (
	cfgLoggingLevel = "LOGGING_LEVEL"
	cfgHost         = "HOST"

	// TraceIDHeader is the header added to outgoing requests which adds the
	// traceID to it.
	TraceIDHeader = "X-Trace-ID"
)

var (
	// ErrNotAuthorized occurs when the call is not authorized.
	ErrNotAuthorized = errors.New("Not authorized")

	// ErrDBNotConfigured occurs when the DB is not initialized.
	ErrDBNotConfigured = errors.New("DB not initialized")

	// ErrNotFound is abstracting the mgo not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in it's proper form")

	// ErrValidation occurs when there are validation errors.
	ErrValidation = errors.New("Validation errors occurred")
)

type (
	// A Handler is a type that handles an http request within our own little mini
	// framework. The fun part is that our context is fully controlled and
	// configured by us so we can extend the functionality of the Context whenever
	// we want.
	Handler func(*Context) error

	// A Middleware is a type that wraps a handler to remove boilerplate or other
	// concerns not direct to any given Handler.
	Middleware func(Handler) Handler
)

// app maintains some framework state.
var app = struct {
	userHeaders map[string]string
}{
	userHeaders: make(map[string]string),
}

//==============================================================================

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*httptreemux.TreeMux
	Ctx map[string]interface{}

	mw []Middleware
}

// New create an App value that handle a set of routes for the application.
// You can provide any number of middleware and they'll be used to wrap every
// request handler.
func New(mw ...Middleware) *App {
	return &App{
		TreeMux: httptreemux.New(),
		Ctx:     make(map[string]interface{}),
		mw:      mw,
	}
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {

	// Wrap up the application-wide first, this will call the first function
	// of each middleware which will return a function of type Handler. Each
	// Handler will then be wrapped up with the other handlers from the chain.
	for i := len(a.mw) - 1; i >= 0; i-- {
		handler = a.mw[i](handler)
	}

	// Then wrap with our route specific ones.
	for i := len(mw) - 1; i >= 0; i-- {
		handler = mw[i](handler)
	}

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		c := Context{
			ResponseWriter: w,
			Request:        r,
			Now:            time.Now(),
			Params:         p,
			SessionID:      uuid.New(),
			Ctx:            make(map[string]interface{}),
			App:            a,
		}

		log.User(c.SessionID, "Request", "Started : Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)

		// Set the request id on the outgoing requests before any other header to
		// ensure that the trace id is ALWAYS added to the request regardless of
		// any error occuring or not.
		c.Header().Set(TraceIDHeader, c.SessionID)

		// Call the wrapped handler and handle any possible error.
		if err := handler(&c); err != nil {
			c.Error(err)
		}

		log.User(c.SessionID, "Request", "Completed : Status[%d] Duration[%s]", c.Status, time.Since(c.Now))
	}

	// Add this handler for the specified verb and route.
	a.TreeMux.Handle(verb, path, h)
}

// CORS providing support for Cross-Origin Resource Sharing.
// https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/
func (a *App) CORS() {
	h := func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
	}

	a.TreeMux.OptionsHandler = h

	app.userHeaders["Access-Control-Allow-Origin"] = "*"
}

//==============================================================================

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

	// Load user defined custom headers. HEADERS should be key:value,key:value
	if hs, err := cfg.String("HEADERS"); err == nil {
		hdrs := strings.Split(hs, ",")
		for _, hdr := range hdrs {
			if kv := strings.Split(hdr, ":"); len(kv) == 2 {
				log.User("startup", "Init", "User Headers : %s:%s", kv[0], kv[1])
				app.userHeaders[kv[0]] = kv[1]
			}
		}
	}
}

// Run is called to start the web service.
func Run(host string, routes http.Handler, readTimeout, writeTimeout time.Duration) error {

	// Check for a configured host value.
	useHost, err := cfg.String(cfgHost)
	if err != nil {
		useHost = host
	}

	log.Dev("startup", "Run", "Start : Using Host[%s]", useHost)

	// Create a new server and set timeout values.
	server := manners.NewWithServer(&http.Server{
		Addr:           useHost,
		Handler:        routes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	})

	go func() {

		// Listen for an interrupt signal from the OS.
		osSignals := make(chan os.Signal)
		signal.Notify(osSignals, os.Interrupt)

		sig := <-osSignals
		log.User("shutdown", "Run", "Captured %v. Shutting Down...", sig)

		// Shut down the API server.
		server.Close()
	}()

	return server.ListenAndServe()
}
