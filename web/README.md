
# web
    import "github.com/ardanlabs/kit/web"

Package app provides application support for context and MongoDB access.
Current Status Codes:


	200 OK           : StatusOK                  : Call is success and returning data.
	204 No Content   : StatusNoContent           : Call is success and returns no data.
	400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
	401 Unauthorized : StatusUnauthorized        : Authentication failure.
	404 Not Found    : StatusNotFound            : Invalid URL or identifier.
	500 Internal     : StatusInternalServerError : Weblication specific beyond scope of user.

Package app provides a thin layer of support for writing web services. It
integrates with the ardanlabs kit repo to provide support for logging,
configuration, database, routing and application context. The base things
you need to write a web service is provided.




## Constants
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

## func Run
``` go
func Run(host string, routes http.Handler, readTimeout, writeTimeout time.Duration) error
```
Run is called to start the web service.



## type Context
``` go
type Context struct {
    http.ResponseWriter
    Request   *http.Request
    Now       time.Time
    Params    map[string]string
    SessionID string
    Status    int
    Ctx       map[string]interface{}
    Web       *Web
}
```
Context contains data associated with a single request.











### func (\*Context) Error
``` go
func (c *Context) Error(err error)
```
Error handles all error responses for the API.



### func (\*Context) Proxy
``` go
func (c *Context) Proxy(targetURL string, rewrite func(req *http.Request)) error
```
Proxy will setup a direct proxy inbetween this service and the destination
service.



### func (\*Context) Respond
``` go
func (c *Context) Respond(data interface{}, code int) error
```
Respond sends JSON to the client.
If code is StatusNoContent, v is expected to be nil.



### func (\*Context) RespondError
``` go
func (c *Context) RespondError(error string, code int)
```
RespondError sends JSON describing the error



### func (\*Context) RespondInvalid
``` go
func (c *Context) RespondInvalid(fields []Invalid)
```
RespondInvalid sends JSON describing field validation errors.



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
Handle proxies the Handle function of the underlying Web.



### func (\*Group) Use
``` go
func (g *Group) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Weblication middleware chain.



## type Handler
``` go
type Handler func(*Context) error
```
A Handler is a type that handles an http request within our own little mini
framework. The fun part is that our context is fully controlled and
configured by us so we can extend the functionality of the Context whenever
we want.











## type Invalid
``` go
type Invalid struct {
    Fld string `json:"field_name"`
    Err string `json:"error"`
}
```
Invalid describes a validation error belonging to a specific field.











## type Middleware
``` go
type Middleware func(Handler) Handler
```
A Middleware is a type that wraps a handler to remove boilerplate or other
concerns not direct to any given Handler.











## type ProxyResponseWriter
``` go
type ProxyResponseWriter struct {
    Status int
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
the request and records the status code written.



## type Web
``` go
type Web struct {
    *httptreemux.TreeMux
    Ctx map[string]interface{}
    // contains filtered or unexported fields
}
```
Web is the entrypoint into our application and what configures our context
object for each of our http handlers. Feel free to add any configuration
data/logic on this Web struct









### func New
``` go
func New(mw ...Middleware) *Web
```
New create an Web value that handle a set of routes for the application.
You can provide any number of middleware and they'll be used to wrap every
request handler.




### func (\*Web) CORS
``` go
func (a *Web) CORS()
```
CORS providing support for Cross-Origin Resource Sharing.
<a href="https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/">https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/</a>



### func (\*Web) Group
``` go
func (a *Web) Group(mw ...Middleware) *Group
```
Group creates a new Web Group based on the current Web and provided
middleware.



### func (\*Web) Handle
``` go
func (a *Web) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle is our mechanism for mounting Handlers for a given HTTP verb and path
pair, this makes for really easy, convenient routing.



### func (\*Web) Use
``` go
func (a *Web) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Weblication middleware
chain. Any route running off of this Web will use all the middleware provided
this way always regardless of the ordering of the Handle/Use functions.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)