

# tcp
`import "github.com/ardanlabs/kit/tcp"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package tcp provides the boilerpale code for working with TCP based data. The package
allows you to establish a TCP listener that can accept client connections on a specified IP address
and port. It also provides a function to send data back to the client. The processing
of received data and sending data happens on a configured routine pool, so concurrency
is handled.

There are three interfaces that need to be implemented to use the package. These
interfaces provide the API for processing data.

ConnHandler


	type ConnHandler interface {
	    Bind(context string, conn net.Conn) (io.Reader, io.Writer)
	}

The ConnHandler interface is implemented by the user to bind the client connection
to a reader and writer for processing.

ReqHandler


	type ReqHandler interface {
	    Read(context string, ipAddress string, reader io.Reader) ([]byte, int, error)
	    Process(context string, r *Request)
	}
	
	type Request struct {
	    TCP       *TCP
	    TCPAddr   *net.TCPAddr
	    Data      []byte
	    Length    int
	}

The ReqHandler interface is implemented by the user to implement the processing
of request messages from the client. Read is provided an ipaddress and the user-defined
reader and must return the data read off the wire and the length. Returning io.EOF or a non
temporary error will show down the listener.

RespHandler


	type RespHandler interface {
	    Write(context string, r *Response, writer io.Writer)
	}
	
	type Response struct {
	    TCPAddr   *net.TCPAddr
	    Data      []byte
	    Length    int
	}

The RespHandler interface is implemented by the user to implement the processing
of the response messages to the client. Write is provided the user-defined
writer and the data to write.

### Sample Application
After implementing the interfaces, the following code is all that is needed to
start processing messages.


	func main() {
	    log.Startf("TEST", "main", "Starting Test App")
	
	    cfg := tcp.Config{
	        NetType:      "tcp4",
	        Addr:         ":9000",
	        WorkRoutines: 2,
	        WorkStats:    time.Minute,
	        ConnHandler:  tcpConnHandler{},
	        ReqHandler:   udpReqHandler{},
	        RespHandler:  udpRespHandler{},
	    }
	
	    t, err := tcp.New("TEST", &cfg)
	    if err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    if err := t.Start("TEST"); err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    // Wait for a signal to shutdown.
	    sigChan := make(chan os.Signal, 1)
	    signal.Notify(sigChan, os.Interrupt)
	    <-sigChan
	
	    t.Stop("TEST")
	
	    log.Complete("TEST", "main")
	}




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Config](#Config)
  * [func (cfg *Config) Event(context interface{}, event string, format string, a ...interface{})](#Config.Event)
  * [func (cfg *Config) Validate() error](#Config.Validate)
* [type ConnHandler](#ConnHandler)
* [type OptEvent](#OptEvent)
* [type OptIntPool](#OptIntPool)
* [type OptRateLimit](#OptRateLimit)
* [type OptUserPool](#OptUserPool)
* [type ReqHandler](#ReqHandler)
* [type Request](#Request)
  * [func (r *Request) Work(context interface{}, id int)](#Request.Work)
* [type RespHandler](#RespHandler)
* [type Response](#Response)
  * [func (r *Response) Work(context interface{}, id int)](#Response.Work)
* [type TCP](#TCP)
  * [func New(context interface{}, name string, cfg Config) (*TCP, error)](#New)
  * [func (t *TCP) Addr() net.Addr](#TCP.Addr)
  * [func (t *TCP) Do(context interface{}, r *Response) error](#TCP.Do)
  * [func (t *TCP) DropConnections(context interface{}, drop bool)](#TCP.DropConnections)
  * [func (t *TCP) Start(context interface{}) error](#TCP.Start)
  * [func (t *TCP) StatsRecv() pool.Stat](#TCP.StatsRecv)
  * [func (t *TCP) StatsSend() pool.Stat](#TCP.StatsSend)
  * [func (t *TCP) Stop(context interface{}) error](#TCP.Stop)


#### <a name="pkg-files">Package files</a>
[client.go](/src/github.com/ardanlabs/kit/tcp/client.go) [doc.go](/src/github.com/ardanlabs/kit/tcp/doc.go) [handlers.go](/src/github.com/ardanlabs/kit/tcp/handlers.go) [tcp.go](/src/github.com/ardanlabs/kit/tcp/tcp.go) [tcp_config.go](/src/github.com/ardanlabs/kit/tcp/tcp_config.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrInvalidConfiguration     = errors.New("Invalid Configuration")
    ErrInvalidNetType           = errors.New("Invalid NetType Configuration")
    ErrInvalidConnHandler       = errors.New("Invalid Connection Handler Configuration")
    ErrInvalidReqHandler        = errors.New("Invalid Request Handler Configuration")
    ErrInvalidRespHandler       = errors.New("Invalid Response Handler Configuration")
    ErrInvalidPoolConfiguration = errors.New("Invalid Pool Configuration")
)
```
Set of error variables for start up.




## <a name="Config">type</a> [Config](/src/target/tcp_config.go?s=1221:2324#L27)
``` go
type Config struct {
    NetType string // "tcp", tcp4" or "tcp6"
    Addr    string // "host:port" or "[ipv6-host%zone]:port"

    ConnHandler ConnHandler // Support for binding new connections to a reader and writer.
    ReqHandler  ReqHandler  // Support for handling the specific request workflow.
    RespHandler RespHandler // Support for handling the specific response workflow.

    OptUserPool
    OptIntPool

    OptRateLimit
    OptEvent
}
```
Config provides a data structure of required configuration parameters.










### <a name="Config.Event">func</a> (\*Config) [Event](/src/target/tcp_config.go?s=2988:3080#L84)
``` go
func (cfg *Config) Event(context interface{}, event string, format string, a ...interface{})
```
Event fires events back to the user for important events.




### <a name="Config.Validate">func</a> (\*Config) [Validate](/src/target/tcp_config.go?s=2382:2417#L55)
``` go
func (cfg *Config) Validate() error
```
Validate checks the configuration to required items.




## <a name="ConnHandler">type</a> [ConnHandler](/src/target/handlers.go?s=153:297#L1)
``` go
type ConnHandler interface {
    // Bind is called to set the reader and writer.
    Bind(context interface{}, conn net.Conn) (io.Reader, io.Writer)
}
```
ConnHandler is implemented by the user to bind the connection
to a reader and writer for processing.










## <a name="OptEvent">type</a> [OptEvent](/src/target/tcp_config.go?s=1041:1145#L22)
``` go
type OptEvent struct {
    Event func(context interface{}, event string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.










## <a name="OptIntPool">type</a> [OptIntPool](/src/target/tcp_config.go?s=434:774#L8)
``` go
type OptIntPool struct {
    RecvMinPoolSize func() int // Min number of routines the recv pool must have.
    RecvMaxPoolSize func() int // Max number of routines the recv pool can have.
    SendMinPoolSize func() int // Min number of routines the send pool must have.
    SendMaxPoolSize func() int // Max number of routines the send pool can have.
}
```
OptIntPool declares fields for the user to provide configuration
for an internally configured pool.










## <a name="OptRateLimit">type</a> [OptRateLimit](/src/target/tcp_config.go?s=876:984#L17)
``` go
type OptRateLimit struct {
    RateLimit func() time.Duration // Connection rate limit per single connection.
}
```
OptRateLimit declares fields for the user to provide configuration
for connection rate limit.










## <a name="OptUserPool">type</a> [OptUserPool](/src/target/tcp_config.go?s=162:326#L1)
``` go
type OptUserPool struct {
    RecvPool *pool.Pool // User provided work pool for the receive work.
    SendPool *pool.Pool // User provided work pool for the send work.
}
```
OptUserPool declares fields for the user to pass their own
work pools for configuration.










## <a name="ReqHandler">type</a> [ReqHandler](/src/target/handlers.go?s=490:1231#L10)
``` go
type ReqHandler interface {

    // Read is provided an ipaddress and the user-defined reader and must return
    // the data read off the wire and the length. Returning io.EOF or a non
    // temporary error will show down the listener.
    Read(context interface{}, ipAddress string, reader io.Reader) ([]byte, int, error)

    // Process is used to handle the processing of the request. This method
    // is called on a routine from a pool of routines.
    Process(context interface{}, r *Request)
}
```
ReqHandler is implemented by the user to implement the processing
of request messages from the client.










## <a name="Request">type</a> [Request](/src/target/handlers.go?s=1283:1404#L27)
``` go
type Request struct {
    TCP     *TCP
    TCPAddr *net.TCPAddr
    IsIPv6  bool
    ReadAt  time.Time
    Data    []byte
    Length  int
}
```
Request is the message received by the client.










### <a name="Request.Work">func</a> (\*Request) [Work](/src/target/handlers.go?s=1531:1582#L38)
``` go
func (r *Request) Work(context interface{}, id int)
```
Work implements the worker interface for processing received messages.
This is called from a routine in the work pool.




## <a name="RespHandler">type</a> [RespHandler](/src/target/handlers.go?s=1821:1983#L46)
``` go
type RespHandler interface {
    // Write is provided the response to write and the user-defined writer.
    Write(context interface{}, r *Response, writer io.Writer)
}
```
RespHandler is implemented by the user to implement the processing
of the response messages to the client.










## <a name="Response">type</a> [Response](/src/target/handlers.go?s=2031:2190#L52)
``` go
type Response struct {
    TCPAddr  *net.TCPAddr
    Data     []byte
    Length   int
    Complete func(r *Response)
    // contains filtered or unexported fields
}
```
Response is message to send to the client.










### <a name="Response.Work">func</a> (\*Response) [Work](/src/target/handlers.go?s=2319:2371#L65)
``` go
func (r *Response) Work(context interface{}, id int)
```
Work implements the worker interface for sending messages to the client.
This is called from a routine in the work pool.




## <a name="TCP">type</a> [TCP](/src/target/tcp.go?s=778:1139#L18)
``` go
type TCP struct {
    Config
    Name string
    // contains filtered or unexported fields
}
```
TCP contains a set of networked client connections.







### <a name="New">func</a> [New](/src/target/tcp.go?s=1190:1258#L45)
``` go
func New(context interface{}, name string, cfg Config) (*TCP, error)
```
New creates a new manager to service clients.





### <a name="TCP.Addr">func</a> (\*TCP) [Addr](/src/target/tcp.go?s=8329:8358#L342)
``` go
func (t *TCP) Addr() net.Addr
```
Addr returns the listener's network address. This may be different than the values
provided in the configuration, for example if configuration port value is 0.




### <a name="TCP.Do">func</a> (\*TCP) [Do](/src/target/tcp.go?s=7053:7109#L294)
``` go
func (t *TCP) Do(context interface{}, r *Response) error
```
Do will post the request to be sent by the client worker pool.




### <a name="TCP.DropConnections">func</a> (\*TCP) [DropConnections](/src/target/tcp.go?s=7739:7800#L321)
``` go
func (t *TCP) DropConnections(context interface{}, drop bool)
```
DropConnections sets a flag to tell the accept routine to immediately
drop connections that come in.




### <a name="TCP.Start">func</a> (\*TCP) [Start](/src/target/tcp.go?s=2928:2974#L121)
``` go
func (t *TCP) Start(context interface{}) error
```
Start creates the accept routine and begins to accept connections.




### <a name="TCP.StatsRecv">func</a> (\*TCP) [StatsRecv](/src/target/tcp.go?s=7969:8004#L331)
``` go
func (t *TCP) StatsRecv() pool.Stat
```
StatsRecv returns the current snapshot of the recv pool stats.




### <a name="TCP.StatsSend">func</a> (\*TCP) [StatsSend](/src/target/tcp.go?s=8099:8134#L336)
``` go
func (t *TCP) StatsSend() pool.Stat
```
StatsSend returns the current snapshot of the send pool stats.




### <a name="TCP.Stop">func</a> (\*TCP) [Stop](/src/target/tcp.go?s=5876:5921#L241)
``` go
func (t *TCP) Stop(context interface{}) error
```
Stop shuts down the manager and closes all connections.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
