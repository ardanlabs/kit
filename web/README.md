
# web
    import "github.com/ardanlabs/kit/web"

Package web provides a thin layer of support for writing web services. It
integrates with the ardanlabs kit repo to provide support for routing and
application ctx. The base things you need to write a web service is
provided.




## Constants
``` go
const KeyValues ctxKey = 1
```
KeyValues is how request values or stored/retrieved.

``` go
const TraceIDHeader = "X-Trace-ID"
```
TraceIDHeader is the header added to outgoing requests which adds the
traceID to it.


## Variables
``` go
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
```

## func Error
``` go
func Error(cxt context.Context, w http.ResponseWriter, err error)
```
Error handles all error responses for the API.


## func Respond
``` go
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, code int) error
```
Respond sends JSON to the client.
If code is StatusNoContent, v is expected to be nil.


## func RespondError
``` go
func RespondError(ctx context.Context, w http.ResponseWriter, err error, code int) error
```
RespondError sends JSON describing the error


## func Run
``` go
func Run(host string, routes http.Handler, readTimeout, writeTimeout time.Duration) error
```
Run is called to start the web service.


## func Unmarshal
``` go
func Unmarshal(r io.Reader, v interface{}) error
```
Unmarshal decodes the input to the struct type and checks the
fields to verify the value is in a proper state.



## type App
``` go
type App struct {
    *httptreemux.TreeMux
    Values map[string]interface{}
    // contains filtered or unexported fields
}
```
App is the entrypoint into our application and what configures our context
object for each of our http handlers. Feel free to add any configuration
data/logic on this App struct









### func New
``` go
func New(mw ...Middleware) *App
```
New create an App value that handle a set of routes for the application.
You can provide any number of middleware and they'll be used to wrap every
request handler.




### func (\*App) Group
``` go
func (a *App) Group(mw ...Middleware) *Group
```
Group creates a new App Group based on the current App and provided
middleware.



### func (\*App) Handle
``` go
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle is our mechanism for mounting Handlers for a given HTTP verb and path
pair, this makes for really easy, convenient routing.



### func (\*App) Use
``` go
func (a *App) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Application middleware
chain. Any route running off of this App will use all the middleware provided
this way always regardless of the ordering of the Handle/Use functions.



## type Group
``` go
type Group struct {
    // contains filtered or unexported fields
}
```
Group allows a segment of middleware to be shared amongst handlers.











### func (\*Group) Handle
``` go
func (g *Group) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle proxies the Handle function of the underlying App.



### func (\*Group) Use
``` go
func (g *Group) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Application middleware chain.



## type Handler
``` go
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error
```
A Handler is a type that handles an http request within our own little mini
framework.











## type Invalid
``` go
type Invalid struct {
    Fld string `json:"field_name"`
    Err string `json:"error"`
}
```
Invalid describes a validation error belonging to a specific field.











## type InvalidError
``` go
type InvalidError []Invalid
```
InvalidError is a custom error type for invalid fields.











### func (InvalidError) Error
``` go
func (err InvalidError) Error() string
```
Error implements the error interface for InvalidError.



## type JSONError
``` go
type JSONError struct {
    Error  string       `json:"error"`
    Fields InvalidError `json:"fields,omitempty"`
}
```
JSONError is the response for errors that occur within the API.











## type Middleware
``` go
type Middleware func(Handler) Handler
```
A Middleware is a type that wraps a handler to remove boilerplate or other
concerns not direct to any given Handler.











## type ProxyResponseWriter
``` go
type ProxyResponseWriter struct {
    Status          int
    UpstreamHeaders http.Header
    http.ResponseWriter
}
```
ProxyResponseWriter records the status code written by a call to the
WriteHeader function on a http.ResponseWriter interface. This type also
implements the http.ResponseWriter interface.











### func (\*ProxyResponseWriter) Header
``` go
func (prw *ProxyResponseWriter) Header() http.Header
```
Header implements the http.ResponseWriter interface and simply relays the
request.



### func (\*ProxyResponseWriter) Write
``` go
func (prw *ProxyResponseWriter) Write(data []byte) (int, error)
```
Write implements the http.ResponseWriter interface and simply relays the
request.



### func (\*ProxyResponseWriter) WriteHeader
``` go
func (prw *ProxyResponseWriter) WriteHeader(status int)
```
WriteHeader implements the http.ResponseWriter interface and simply relays
the request after cleaning up the request headers. It theb records the status
code written.



## type Values
``` go
type Values struct {
    TraceID    string
    Now        time.Time
    StatusCode int
}
```
Values represent state for each request.

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)