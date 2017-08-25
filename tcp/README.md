

# tcp
`import "github.com/ardanlabs/kit/tcp"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package tcp provides the boilerpale code for working with TCP based data. The package
allows you to establish a TCP listener that can accept client connections on a specified IP address
and port. It also provides a function to send data back to the client.

There are three interfaces that need to be implemented to use the package. These
interfaces provide the API for processing data.

ConnHandler


	type ConnHandler interface {
	    Bind(conn net.Conn) (io.Reader, io.Writer)
	}

The ConnHandler interface is implemented by the user to bind the client connection
to a reader and writer for processing.

ReqHandler


	type ReqHandler interface {
	    Read(ipAddress string, reader io.Reader) ([]byte, int, error)
	    Process(r *Request)
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
	    Write(r *Response, writer io.Writer) error
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
	    log.Println("Starting Test App")
	
	    cfg := tcp.Config{
	        NetType:      "tcp4",
	        Addr:         ":9000",
	        WorkRoutines: 2,
	        WorkStats:    time.Minute,
	        ConnHandler:  tcpConnHandler{},
	        ReqHandler:   udpReqHandler{},
	        RespHandler:  udpRespHandler{},
	    }
	
	    t, err := tcp.New(&cfg)
	    if err != nil {
	        log.Println(err)
	         return
	    }
	
	    if err := t.Start(); err != nil {
	        log.Println(err)
	         return
	    }
	
	    // Wait for a signal to shutdown.
	    sigChan := make(chan os.Signal, 1)
	    signal.Notify(sigChan, os.Interrupt)
	    <-sigChan
	
	    t.Stop()
	    log.Println("down")
	}




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [type CltError](#CltError)
  * [func (ce CltError) Error() string](#CltError.Error)
* [type Config](#Config)
  * [func (cfg *Config) Event(evt, typ int, ipAddress string, format string, a ...interface{})](#Config.Event)
  * [func (cfg *Config) Validate() error](#Config.Validate)
* [type ConnHandler](#ConnHandler)
* [type OptEvent](#OptEvent)
* [type OptRateLimit](#OptRateLimit)
* [type ReqHandler](#ReqHandler)
* [type Request](#Request)
* [type RespHandler](#RespHandler)
* [type Response](#Response)
* [type Stat](#Stat)
* [type TCP](#TCP)
  * [func New(name string, cfg Config) (*TCP, error)](#New)
  * [func (t *TCP) Addr() net.Addr](#TCP.Addr)
  * [func (t *TCP) ClientStats() []Stat](#TCP.ClientStats)
  * [func (t *TCP) Clients() int](#TCP.Clients)
  * [func (t *TCP) Connections() int](#TCP.Connections)
  * [func (t *TCP) Drop(tcpAddr *net.TCPAddr) error](#TCP.Drop)
  * [func (t *TCP) DropConnections(drop bool)](#TCP.DropConnections)
  * [func (t *TCP) Groom(d time.Duration)](#TCP.Groom)
  * [func (t *TCP) Send(ctx context.Context, r *Response) error](#TCP.Send)
  * [func (t *TCP) SendAll(ctx context.Context, r *Response) error](#TCP.SendAll)
  * [func (t *TCP) Start() error](#TCP.Start)
  * [func (t *TCP) Stop() error](#TCP.Stop)


#### <a name="pkg-files">Package files</a>
[client.go](/src/github.com/ardanlabs/kit/tcp/client.go) [doc.go](/src/github.com/ardanlabs/kit/tcp/doc.go) [handlers.go](/src/github.com/ardanlabs/kit/tcp/handlers.go) [tcp.go](/src/github.com/ardanlabs/kit/tcp/tcp.go) [tcp_config.go](/src/github.com/ardanlabs/kit/tcp/tcp_config.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    EvtAccept = iota + 1
    EvtJoin
    EvtRead
    EvtRemove
    EvtDrop
    EvtGroom
)
```
Set of event types.

``` go
const (
    TypError = iota + 1
    TypInfo
    TypTrigger
)
```
Set of event sub types.


## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrInvalidConfiguration = errors.New("invalid configuration")
    ErrInvalidNetType       = errors.New("invalid net type configuration")
    ErrInvalidConnHandler   = errors.New("invalid connection handler configuration")
    ErrInvalidReqHandler    = errors.New("invalid request handler configuration")
    ErrInvalidRespHandler   = errors.New("invalid response handler configuration")
)
```
Set of error variables for start up.




## <a name="CltError">type</a> [CltError](/src/target/tcp.go?s=795:816#L32)
``` go
type CltError []error
```
CltError provides support for multi client operations that might error.










### <a name="CltError.Error">func</a> (CltError) [Error](/src/target/tcp.go?s=871:904#L35)
``` go
func (ce CltError) Error() string
```
Error implments the error interface for CltError.




## <a name="Config">type</a> [Config](/src/target/tcp_config.go?s=470:1102#L7)
``` go
type Config struct {
    NetType string // "tcp", tcp4" or "tcp6"
    Addr    string // "host:port" or "[ipv6-host%zone]:port"

    ConnHandler ConnHandler // Support for binding new connections to a reader and writer.
    ReqHandler  ReqHandler  // Support for handling the specific request workflow.
    RespHandler RespHandler // Support for handling the specific response workflow.

    OptRateLimit
    OptEvent
}
```
Config provides a data structure of required configuration parameters.










### <a name="Config.Event">func</a> (\*Config) [Event](/src/target/tcp_config.go?s=1626:1715#L49)
``` go
func (cfg *Config) Event(evt, typ int, ipAddress string, format string, a ...interface{})
```
Event fires events back to the user for important events.




### <a name="Config.Validate">func</a> (\*Config) [Validate](/src/target/tcp_config.go?s=1160:1195#L24)
``` go
func (cfg *Config) Validate() error
```
Validate checks the configuration to required items.




## <a name="ConnHandler">type</a> [ConnHandler](/src/target/handlers.go?s=485:609#L20)
``` go
type ConnHandler interface {

    // Bind is called to set the reader and writer.
    Bind(conn net.Conn) (io.Reader, io.Writer)
}
```
ConnHandler is implemented by the user to bind the connection
to a reader and writer for processing.










## <a name="OptEvent">type</a> [OptEvent](/src/target/tcp_config.go?s=293:394#L2)
``` go
type OptEvent struct {
    Event func(evt, typ int, ipAddress string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.










## <a name="OptRateLimit">type</a> [OptRateLimit](/src/target/tcp_config.go?s=128:236#L1)
``` go
type OptRateLimit struct {
    RateLimit func() time.Duration // Connection rate limit per single connection.
}
```
OptRateLimit declares fields for the user to provide configuration
for connection rate limit.










## <a name="ReqHandler">type</a> [ReqHandler](/src/target/handlers.go?s=720:1356#L28)
``` go
type ReqHandler interface {

    // Read is provided an ipaddress and the user-defined reader and must return
    // the data read off the wire and the length. Returning io.EOF or a non
    // temporary error will show down the listener.
    Read(ipAddress string, reader io.Reader) ([]byte, int, error)

    // Process is used to handle the processing of the request.
    Process(r *Request)
}
```
ReqHandler is implemented by the user to implement the processing
of request messages from the client.










## <a name="Request">type</a> [Request](/src/target/handlers.go?s=107:253#L1)
``` go
type Request struct {
    TCP     *TCP
    TCPAddr *net.TCPAddr
    IsIPv6  bool
    ReadAt  time.Time
    Context context.Context
    Data    []byte
    Length  int
}
```
Request is the message received by the client.










## <a name="RespHandler">type</a> [RespHandler](/src/target/handlers.go?s=1471:1619#L46)
``` go
type RespHandler interface {

    // Write is provided the response to write and the user-defined writer.
    Write(r *Response, writer io.Writer) error
}
```
RespHandler is implemented by the user to implement the processing
of the response messages to the client.










## <a name="Response">type</a> [Response](/src/target/handlers.go?s=301:376#L12)
``` go
type Response struct {
    TCPAddr *net.TCPAddr
    Data    []byte
    Length  int
}
```
Response is message to send to the client.










## <a name="Stat">type</a> [Stat](/src/target/tcp.go?s=8518:8623#L376)
``` go
type Stat struct {
    IP       string
    Reads    int
    Writes   int
    TimeConn time.Time
    LastAct  time.Time
}
```
Stat represents a client statistic.










## <a name="TCP">type</a> [TCP](/src/target/tcp.go?s=1084:1384#L45)
``` go
type TCP struct {
    Config
    Name string
    // contains filtered or unexported fields
}
```
TCP contains a set of networked client connections.







### <a name="New">func</a> [New](/src/target/tcp.go?s=1435:1482#L68)
``` go
func New(name string, cfg Config) (*TCP, error)
```
New creates a new manager to service clients.





### <a name="TCP.Addr">func</a> (\*TCP) [Addr](/src/target/tcp.go?s=8106:8135#L352)
``` go
func (t *TCP) Addr() net.Addr
```
Addr returns the listener's network address. This may be different than the values
provided in the configuration, for example if configuration port value is 0.




### <a name="TCP.ClientStats">func</a> (\*TCP) [ClientStats](/src/target/tcp.go?s=8679:8713#L385)
``` go
func (t *TCP) ClientStats() []Stat
```
ClientStats return details for all active clients.




### <a name="TCP.Clients">func</a> (\*TCP) [Clients](/src/target/tcp.go?s=9132:9159#L410)
``` go
func (t *TCP) Clients() int
```
Clients returns the number of active clients connected.




### <a name="TCP.Connections">func</a> (\*TCP) [Connections](/src/target/tcp.go?s=8350:8381#L363)
``` go
func (t *TCP) Connections() int
```
Connections returns the number of client connections.




### <a name="TCP.Drop">func</a> (\*TCP) [Drop](/src/target/tcp.go?s=6106:6152#L269)
``` go
func (t *TCP) Drop(tcpAddr *net.TCPAddr) error
```
Drop will close the socket connection.




### <a name="TCP.DropConnections">func</a> (\*TCP) [DropConnections](/src/target/tcp.go?s=7797:7837#L341)
``` go
func (t *TCP) DropConnections(drop bool)
```
DropConnections sets a flag to tell the accept routine to immediately
drop connections that come in.




### <a name="TCP.Groom">func</a> (\*TCP) [Groom](/src/target/tcp.go?s=9343:9379#L422)
``` go
func (t *TCP) Groom(d time.Duration)
```
Groom drops connections that are not active for the specified duration.




### <a name="TCP.Send">func</a> (\*TCP) [Send](/src/target/tcp.go?s=6659:6717#L291)
``` go
func (t *TCP) Send(ctx context.Context, r *Response) error
```
Send will deliver the response back to the client.




### <a name="TCP.SendAll">func</a> (\*TCP) [SendAll](/src/target/tcp.go?s=7237:7298#L314)
``` go
func (t *TCP) SendAll(ctx context.Context, r *Response) error
```
SendAll will deliver the response back to all connected clients.




### <a name="TCP.Start">func</a> (\*TCP) [Start](/src/target/tcp.go?s=2178:2205#L102)
``` go
func (t *TCP) Start() error
```
Start creates the accept routine and begins to accept connections.




### <a name="TCP.Stop">func</a> (\*TCP) [Stop](/src/target/tcp.go?s=5082:5108#L221)
``` go
func (t *TCP) Stop() error
```
Stop shuts down the manager and closes all connections.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
