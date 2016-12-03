
# web
    import "github.com/ardanlabs/kit/web"

Package web provides web application support for ctx and MongoDB access.
Current Status Codes:


	200 OK           : StatusOK                  : Call is success and returning data.
	204 No Content   : StatusNoContent           : Call is success and returns no data.
	400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
	401 Unauthorized : StatusUnauthorized        : Authentication failure.
	404 Not Found    : StatusNotFound            : Invalid URL or identifier.
	500 Internal     : StatusInternalServerError : Weblication specific beyond scope of user.

Package web provides a thin layer of support for writing web services. It
integrates with the ardanlabs kit repo to provide support for routing and
application ctx. The base things you need to write a web service is
provided.




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



## type Ctx
``` go
type Ctx struct {
    http.ResponseWriter
    Request   *http.Request
    Now       time.Time
    Params    map[string]string
    SessionID string
    Status    int
    Values    map[string]interface{}
    Web       *Web
}
```
Ctx contains data associated with a single request.











### func (\*Ctx) Error
``` go
func (c *Ctx) Error(err error)
```
Error handles all error responses for the API.



### func (\*Ctx) Proxy
``` go
func (c *Ctx) Proxy(targetURL string, rewrite func(req *http.Request)) error
```
Proxy will setup a direct proxy inbetween this service and the destination
service.



### func (\*Ctx) Respond
``` go
func (c *Ctx) Respond(data interface{}, code int) error
```
Respond sends JSON to the client.
If code is StatusNoContent, v is expected to be nil.



### func (\*Ctx) RespondError
``` go
func (c *Ctx) RespondError(error string, code int)
```
RespondError sends JSON describing the error



### func (\*Ctx) RespondInvalid
``` go
func (c *Ctx) RespondInvalid(fields []Invalid)
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
type Handler func(*Ctx) error
```
A Handler is a type that handles an http request within our own little mini
framework. The fun part is that our Ctx is fully controlled and
configured by us so we can extend the functionality of the ctx whenever
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



## type Web
``` go
type Web struct {
    *httptreemux.TreeMux
    Ctx map[string]interface{}
    // contains filtered or unexported fields
}
```
Web is the entrypoint into our application and what configures our ctx
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
func (w *Web) CORS() Middleware
```
CORS providing support for Cross-Origin Resource Sharing.
<a href="https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/">https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/</a>



### func (\*Web) Group
``` go
func (w *Web) Group(mw ...Middleware) *Group
```
Group creates a new Web Group based on the current Web and provided
middleware.



### func (\*Web) Handle
``` go
func (w *Web) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle is our mechanism for mounting Handlers for a given HTTP verb and path
pair, this makes for really easy, convenient routing.



### func (\*Web) Use
``` go
func (w *Web) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Weblication middleware
chain. Any route running off of this Web will use all the middleware provided
this way always regardless of the ordering of the Handle/Use functions.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)