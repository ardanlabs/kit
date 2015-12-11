package app

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"github.com/dimfeld/httptreemux"
	"github.com/pborman/uuid"
)

var (
	// ErrNotAuthorized occurs when the call is not authorized.
	ErrNotAuthorized = errors.New("Not authorized")

	// ErrNotFound is abstracting the mgo not found error.
	ErrNotFound = errors.New("Entity Not found")

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
var app struct {
	useMongo bool
}

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*httptreemux.TreeMux
	mw []Middleware
}

// New create an App value that handle a set of routes for the application.
// You can provide any number of middleware and they'll be used to wrap every
// request handler.
func New(mw ...Middleware) *App {
	return &App{
		TreeMux: httptreemux.New(),
		mw:      mw,
	}
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {
	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		start := time.Now()

		var dbConn *db.DB
		if app.useMongo {
			dbConn = db.NewMGO()
		}

		c := Context{
			DB:             dbConn,
			ResponseWriter: w,
			Request:        r,
			Params:         p,
			SessionID:      uuid.New(),
		}

		if app.useMongo {
			defer c.DB.CloseMGO()
		}

		log.User(c.SessionID, "Request", "Started : Method[%s] URL[%s] RADDR[%s]", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)

		// Wrap the handler in all associated middleware.
		wrap := func(h Handler) Handler {
			// Wrap up the application-wide first...
			for i := len(a.mw) - 1; i >= 0; i-- {
				h = a.mw[i](h)
			}

			// Then wrap with our route specific ones.
			for i := len(mw) - 1; i >= 0; i-- {
				h = mw[i](h)
			}

			return h
		}

		// Call the wrapped handler and handle any possible error.
		if err := wrap(handler)(&c); err != nil {
			c.Error(err)
		}

		log.User(c.SessionID, "Request", "Completed : Status[%d] Duration[%s]", c.Status, time.Since(start))
	}

	// Add this handler for the specified verb and route.
	a.TreeMux.Handle(verb, path, h)
}

// Settings represents things required to initialize the app.
type Settings struct {
	ConfigKey string // The based environment variable key for all variables.
	UseMongo  bool   // If MongoDB should be initialized and used.
}

// Init is called to initialize the application.
func Init(set *Settings) {
	app.useMongo = set.UseMongo

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	if err := cfg.Init(set.ConfigKey); err != nil {
		log.Error("startup", "Init", err, "Initializing config")
		os.Exit(1)
	}

	if set.UseMongo {
		err := mongo.Init()
		if err != nil {
			log.Error("startup", "Init", err, "Initializing MongoDB")
			os.Exit(1)
		}
	}
}

// Run is called to start the web service.
func Run(cfgHost string, defaultHost string, routes http.Handler) {
	log.Dev("startup", "Run", "Start : cfgHost[%s] defaultHost[%s]", cfgHost, defaultHost)

	// Check for a configured host value.
	host, err := cfg.String(cfgHost)
	if err != nil {
		host = defaultHost
	}

	// Create this goroutine to run the web server.
	go func() {
		log.Dev("listener", "Run", "Listening on: %s", host)
		http.ListenAndServe(host, routes)
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Dev("shutdown", "Run", "Complete")
}
