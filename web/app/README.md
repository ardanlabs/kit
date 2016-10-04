

# app
`import "github.com/ardanlabs/kit/web/app"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package app provides a thin layer of support for writing web services. It
integrates with the ardanlabs kit repo to provide support for logging,
configuration, database, routing and application context. The base things
you need to write a web service is provided.

Package app provides application support for context and MongoDB access.
Current Status Codes:


	200 OK           : StatusOK                  : Call is success and returning data.
	204 No Content   : StatusNoContent           : Call is success and returns no data.
	400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
	401 Unauthorized : StatusUnauthorized        : Authentication failure.
	404 Not Found    : StatusNotFound            : Invalid URL or identifier.
	500 Internal     : StatusInternalServerError : Application specific beyond scope of user.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Init(p cfg.Provider)](#Init)
* [func Run(host string, routes http.Handler, readTimeout, writeTimeout time.Duration) error](#Run)
* [type App](#App)
  * [func New(mw ...Middleware) *App](#New)
  * [func (a *App) CORS()](#App.CORS)
  * [func (a *App) Group(mw ...Middleware) *Group](#App.Group)
  * [func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware)](#App.Handle)
  * [func (a *App) Use(mw ...Middleware)](#App.Use)
* [type Context](#Context)
  * [func (c *Context) Error(err error)](#Context.Error)
  * [func (c *Context) Proxy(targetURL string, rewrite func(req *http.Request)) error](#Context.Proxy)
  * [func (c *Context) Respond(data interface{}, code int)](#Context.Respond)
  * [func (c *Context) RespondError(error string, code int)](#Context.RespondError)
  * [func (c *Context) RespondInvalid(fields []Invalid)](#Context.RespondInvalid)
* [type Group](#Group)
  * [func (g *Group) Handle(verb, path string, handler Handler, mw ...Middleware)](#Group.Handle)
  * [func (g *Group) Use(mw ...Middleware)](#Group.Use)
* [type Handler](#Handler)
* [type Invalid](#Invalid)
* [type Middleware](#Middleware)
* [type ProxyResponseWriter](#ProxyResponseWriter)
  * [func (prw *ProxyResponseWriter) Header() http.Header](#ProxyResponseWriter.Header)
  * [func (prw *ProxyResponseWriter) Write(data []byte) (int, error)](#ProxyResponseWriter.Write)
  * [func (prw *ProxyResponseWriter) WriteHeader(status int)](#ProxyResponseWriter.WriteHeader)


#### <a name="pkg-files">Package files</a>
[app.go](/src/github.com/ardanlabs/kit/web/app/app.go) [context.go](/src/github.com/ardanlabs/kit/web/app/context.go) [proxy.go](/src/github.com/ardanlabs/kit/web/app/proxy.go) 


## <a name="pkg-constants">Constants</a>
``` go
const TraceIDHeader = "X-Trace-ID"
```
TraceIDHeader is the header added to outgoing requests which adds the
traceID to it.


## <a name="pkg-variables">Variables</a>
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


## <a name="Init">func</a> [Init](/src/target/app.go?s=6348:6373#L185)
``` go
func Init(p cfg.Provider)
```
Init is called to initialize the application.



## <a name="Run">func</a> [Run](/src/target/app.go?s=7231:7320#L219)
``` go
func Run(host string, routes http.Handler, readTimeout, writeTimeout time.Duration) error
```
Run is called to start the web service.




## <a name="App">type</a> [App](/src/target/app.go?s=2347:2434#L66)
``` go
type App struct {
    *httptreemux.TreeMux
    Ctx map[string]interface{}
    // contains filtered or unexported fields
}
```
App is the entrypoint into our application and what configures our context
object for each of our http handlers. Feel free to add any configuration
data/logic on this App struct







### <a name="New">func</a> [New](/src/target/app.go?s=2610:2641#L76)
``` go
func New(mw ...Middleware) *App
```
New create an App value that handle a set of routes for the application.
You can provide any number of middleware and they'll be used to wrap every
request handler.





### <a name="App.CORS">func</a> (\*App) [CORS](/src/target/app.go?s=4897:4917#L142)
``` go
func (a *App) CORS()
```
CORS providing support for Cross-Origin Resource Sharing.
<a href="https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/">https://metajack.im/2010/01/19/crossdomain-ajax-for-xmpp-http-binding-made-easy/</a>




### <a name="App.Group">func</a> (\*App) [Group](/src/target/app.go?s=2836:2880#L86)
``` go
func (a *App) Group(mw ...Middleware) *Group
```
Group creates a new App Group based on the current App and provided
middleware.




### <a name="App.Handle">func</a> (\*App) [Handle](/src/target/app.go?s=3363:3437#L102)
``` go
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle is our mechanism for mounting Handlers for a given HTTP verb and path
pair, this makes for really easy, convenient routing.




### <a name="App.Use">func</a> (\*App) [Use](/src/target/app.go?s=3157:3192#L96)
``` go
func (a *App) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Application middleware
chain. Any route running off of this App will use all the middleware provided
this way always regardless of the ordering of the Handle/Use functions.




## <a name="Context">type</a> [Context](/src/target/context.go?s=1226:1428#L28)
``` go
type Context struct {
    http.ResponseWriter
    Request   *http.Request
    Now       time.Time
    Params    map[string]string
    SessionID string
    Status    int
    Ctx       map[string]interface{}
    App       *App
}
```
Context contains data associated with a single request.










### <a name="Context.Error">func</a> (\*Context) [Error](/src/target/context.go?s=1480:1514#L40)
``` go
func (c *Context) Error(err error)
```
Error handles all error responses for the API.




### <a name="Context.Proxy">func</a> (\*Context) [Proxy](/src/target/context.go?s=3799:3879#L122)
``` go
func (c *Context) Proxy(targetURL string, rewrite func(req *http.Request)) error
```
Proxy will setup a direct proxy inbetween this service and the destination
service.




### <a name="Context.Respond">func</a> (\*Context) [Respond](/src/target/context.go?s=1998:2051#L57)
``` go
func (c *Context) Respond(data interface{}, code int)
```
Respond sends JSON to the client.
If code is StatusNoContent, v is expected to be nil.




### <a name="Context.RespondError">func</a> (\*Context) [RespondError](/src/target/context.go?s=3607:3661#L116)
``` go
func (c *Context) RespondError(error string, code int)
```
RespondError sends JSON describing the error




### <a name="Context.RespondInvalid">func</a> (\*Context) [RespondInvalid](/src/target/context.go?s=3390:3440#L107)
``` go
func (c *Context) RespondInvalid(fields []Invalid)
```
RespondInvalid sends JSON describing field validation errors.




## <a name="Group">type</a> [Group](/src/target/app.go?s=5600:5649#L161)
``` go
type Group struct {
    // contains filtered or unexported fields
}
```
Group allows a segment of middleware to be shared amongst handlers.










### <a name="Group.Handle">func</a> (\*Group) [Handle](/src/target/app.go?s=5865:5941#L172)
``` go
func (g *Group) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle proxies the Handle function of the underlying App.




### <a name="Group.Use">func</a> (\*Group) [Use](/src/target/app.go?s=5733:5770#L167)
``` go
func (g *Group) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Application middleware chain.




## <a name="Handler">type</a> [Handler](/src/target/app.go?s=1881:1914#L55)
``` go
type Handler func(*Context) error
```
A Handler is a type that handles an http request within our own little mini
framework. The fun part is that our context is fully controlled and
configured by us so we can extend the functionality of the Context whenever
we want.










## <a name="Invalid">type</a> [Invalid](/src/target/context.go?s=830:912#L14)
``` go
type Invalid struct {
    Fld string `json:"field_name"`
    Err string `json:"error"`
}
```
Invalid describes a validation error belonging to a specific field.










## <a name="Middleware">type</a> [Middleware](/src/target/app.go?s=2039:2076#L59)
``` go
type Middleware func(Handler) Handler
```
A Middleware is a type that wraps a handler to remove boilerplate or other
concerns not direct to any given Handler.










## <a name="ProxyResponseWriter">type</a> [ProxyResponseWriter](/src/target/proxy.go?s=228:296#L1)
``` go
type ProxyResponseWriter struct {
    Status int
    http.ResponseWriter
}
```
ProxyResponseWriter records the status code written by a call to the
WriteHeader function on a http.ResponseWriter interface. This type also
implements the http.ResponseWriter interface.










### <a name="ProxyResponseWriter.Header">func</a> (\*ProxyResponseWriter) [Header](/src/target/proxy.go?s=387:439#L5)
``` go
func (prw *ProxyResponseWriter) Header() http.Header
```
Header implements the http.ResponseWriter interface and simply relays the
request.




### <a name="ProxyResponseWriter.Write">func</a> (\*ProxyResponseWriter) [Write](/src/target/proxy.go?s=569:632#L11)
``` go
func (prw *ProxyResponseWriter) Write(data []byte) (int, error)
```
Write implements the http.ResponseWriter interface and simply relays the
request.




### <a name="ProxyResponseWriter.WriteHeader">func</a> (\*ProxyResponseWriter) [WriteHeader](/src/target/proxy.go?s=807:862#L17)
``` go
func (prw *ProxyResponseWriter) WriteHeader(status int)
```
WriteHeader implements the http.ResponseWriter interface and simply relays
the request and records the status code written.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
