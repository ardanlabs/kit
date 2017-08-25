

# web
`import "github.com/ardanlabs/kit/web"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package web provides a thin layer of support for writing web services. It
integrates with the ardanlabs kit repo to provide support for routing and
application ctx. The base things you need to write a web service is
provided.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Error(cxt context.Context, w http.ResponseWriter, err error)](#Error)
* [func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, code int) error](#Respond)
* [func RespondError(ctx context.Context, w http.ResponseWriter, err error, code int) error](#RespondError)
* [func Run(host string, routes http.Handler, readTimeout, writeTimeout time.Duration) error](#Run)
* [func Unmarshal(r io.Reader, v interface{}) error](#Unmarshal)
* [type App](#App)
  * [func New(mw ...Middleware) *App](#New)
  * [func (a *App) Group(mw ...Middleware) *Group](#App.Group)
  * [func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware)](#App.Handle)
  * [func (a *App) Use(mw ...Middleware)](#App.Use)
* [type Group](#Group)
  * [func (g *Group) Handle(verb, path string, handler Handler, mw ...Middleware)](#Group.Handle)
  * [func (g *Group) Use(mw ...Middleware)](#Group.Use)
* [type Handler](#Handler)
* [type Invalid](#Invalid)
* [type InvalidError](#InvalidError)
  * [func (err InvalidError) Error() string](#InvalidError.Error)
* [type JSONError](#JSONError)
* [type Middleware](#Middleware)
* [type ProxyResponseWriter](#ProxyResponseWriter)
  * [func (prw *ProxyResponseWriter) Header() http.Header](#ProxyResponseWriter.Header)
  * [func (prw *ProxyResponseWriter) Write(data []byte) (int, error)](#ProxyResponseWriter.Write)
  * [func (prw *ProxyResponseWriter) WriteHeader(status int)](#ProxyResponseWriter.WriteHeader)
* [type Values](#Values)


#### <a name="pkg-files">Package files</a>
[proxy.go](/src/github.com/ardanlabs/kit/web/proxy.go) [reponse.go](/src/github.com/ardanlabs/kit/web/reponse.go) [web.go](/src/github.com/ardanlabs/kit/web/web.go) 


## <a name="pkg-constants">Constants</a>
``` go
const KeyValues ctxKey = 1
```
KeyValues is how request values or stored/retrieved.

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


## <a name="Error">func</a> [Error](/src/target/reponse.go?s=2013:2078#L55)
``` go
func Error(cxt context.Context, w http.ResponseWriter, err error)
```
Error handles all error responses for the API.



## <a name="Respond">func</a> [Respond](/src/target/reponse.go?s=2963:3053#L95)
``` go
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, code int) error
```
Respond sends JSON to the client.
If code is StatusNoContent, v is expected to be nil.



## <a name="RespondError">func</a> [RespondError](/src/target/reponse.go?s=2715:2803#L89)
``` go
func RespondError(ctx context.Context, w http.ResponseWriter, err error, code int) error
```
RespondError sends JSON describing the error



## <a name="Run">func</a> [Run](/src/target/web.go?s=4848:4937#L154)
``` go
func Run(host string, routes http.Handler, readTimeout, writeTimeout time.Duration) error
```
Run is called to start the web service.



## <a name="Unmarshal">func</a> [Unmarshal](/src/target/web.go?s=845:893#L24)
``` go
func Unmarshal(r io.Reader, v interface{}) error
```
Unmarshal decodes the input to the struct type and checks the
fields to verify the value is in a proper state.




## <a name="App">type</a> [App](/src/target/web.go?s=2020:2110#L64)
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







### <a name="New">func</a> [New](/src/target/web.go?s=2286:2317#L74)
``` go
func New(mw ...Middleware) *App
```
New create an App value that handle a set of routes for the application.
You can provide any number of middleware and they'll be used to wrap every
request handler.





### <a name="App.Group">func</a> (\*App) [Group](/src/target/web.go?s=2512:2556#L84)
``` go
func (a *App) Group(mw ...Middleware) *Group
```
Group creates a new App Group based on the current App and provided
middleware.




### <a name="App.Handle">func</a> (\*App) [Handle](/src/target/web.go?s=3039:3113#L100)
``` go
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle is our mechanism for mounting Handlers for a given HTTP verb and path
pair, this makes for really easy, convenient routing.




### <a name="App.Use">func</a> (\*App) [Use](/src/target/web.go?s=2833:2868#L94)
``` go
func (a *App) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Application middleware
chain. Any route running off of this App will use all the middleware provided
this way always regardless of the ordering of the Handle/Use functions.




## <a name="Group">type</a> [Group](/src/target/web.go?s=4188:4237#L132)
``` go
type Group struct {
    // contains filtered or unexported fields
}
```
Group allows a segment of middleware to be shared amongst handlers.










### <a name="Group.Handle">func</a> (\*Group) [Handle](/src/target/web.go?s=4453:4529#L143)
``` go
func (g *Group) Handle(verb, path string, handler Handler, mw ...Middleware)
```
Handle proxies the Handle function of the underlying App.




### <a name="Group.Use">func</a> (\*Group) [Use](/src/target/web.go?s=4321:4358#L138)
``` go
func (g *Group) Use(mw ...Middleware)
```
Use adds the set of provided middleware onto the Application middleware chain.




## <a name="Handler">type</a> [Handler](/src/target/web.go?s=1559:1669#L55)
``` go
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error
```
A Handler is a type that handles an http request within our own little mini
framework.










## <a name="Invalid">type</a> [Invalid](/src/target/reponse.go?s=1403:1485#L31)
``` go
type Invalid struct {
    Fld string `json:"field_name"`
    Err string `json:"error"`
}
```
Invalid describes a validation error belonging to a specific field.










## <a name="InvalidError">type</a> [InvalidError](/src/target/reponse.go?s=1546:1573#L37)
``` go
type InvalidError []Invalid
```
InvalidError is a custom error type for invalid fields.










### <a name="InvalidError.Error">func</a> (InvalidError) [Error](/src/target/reponse.go?s=1633:1671#L40)
``` go
func (err InvalidError) Error() string
```
Error implements the error interface for InvalidError.




## <a name="JSONError">type</a> [JSONError](/src/target/reponse.go?s=1853:1961#L49)
``` go
type JSONError struct {
    Error  string       `json:"error"`
    Fields InvalidError `json:"fields,omitempty"`
}
```
JSONError is the response for errors that occur within the API.










## <a name="Middleware">type</a> [Middleware](/src/target/web.go?s=1794:1831#L59)
``` go
type Middleware func(Handler) Handler
```
A Middleware is a type that wraps a handler to remove boilerplate or other
concerns not direct to any given Handler.










## <a name="ProxyResponseWriter">type</a> [ProxyResponseWriter](/src/target/proxy.go?s=228:334#L1)
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










### <a name="ProxyResponseWriter.Header">func</a> (\*ProxyResponseWriter) [Header](/src/target/proxy.go?s=425:477#L6)
``` go
func (prw *ProxyResponseWriter) Header() http.Header
```
Header implements the http.ResponseWriter interface and simply relays the
request.




### <a name="ProxyResponseWriter.Write">func</a> (\*ProxyResponseWriter) [Write](/src/target/proxy.go?s=607:670#L12)
``` go
func (prw *ProxyResponseWriter) Write(data []byte) (int, error)
```
Write implements the http.ResponseWriter interface and simply relays the
request.




### <a name="ProxyResponseWriter.WriteHeader">func</a> (\*ProxyResponseWriter) [WriteHeader](/src/target/proxy.go?s=891:946#L19)
``` go
func (prw *ProxyResponseWriter) WriteHeader(status int)
```
WriteHeader implements the http.ResponseWriter interface and simply relays
the request after cleaning up the request headers. It theb records the status
code written.




## <a name="Values">type</a> [Values](/src/target/web.go?s=1385:1464#L47)
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
