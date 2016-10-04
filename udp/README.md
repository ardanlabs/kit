

# udp
`import "github.com/ardanlabs/kit/udp"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package udp provides the boilerpale code for working with UDP based data. The package
allows you to establish a UDP listener that can accept data on a specified IP address
and port. It also provides a function to send data back to the client. The processing
of received data and sending data happens on a configured routine pool, so concurrency
is handled.

There are three interfaces that need to be implemented to use the package. These
interfaces provide the API for processing data.

ConnHandler


	type ConnHandler interface {
	    Bind(context string, listener *net.UDPConn) (io.Reader, io.Writer)
	}

The ConnHandler interface is implemented by the user to bind the listener
to a reader and writer for processing.

ReqHandler


	type ReqHandler interface {
	    Read(context string, reader io.Reader) (*net.UDPAddr, []byte, int, error)
	    Process(context string, r *Request)
	}
	
	type Request struct {
	    UDP     *UDP
	    UDPAddr *net.UDPAddr
	    Data    []byte
	    Length  int
	}

The ReqHandler interface is implemented by the user to implement the processing
of request messages from the client. Read is provided the user-defined reader
and must return the data read off the wire and the length. Returning io.EOF or
a non temporary error will show down the listener.

RespHandler


	type RespHandler interface {
	    Write(context string, r *Response, writer io.Writer)
	}
	
	type Response struct {
	    UDPAddr *net.UDPAddr
	    Data    []byte
	    Length  int
	}

The RespHandler interface is implemented by the user to implement the processing
of the response messages to the client. Write is provided the user-defined
writer and the data to write.

### Sample Application
After implementing the interfaces, the following code is all that is needed to
start processing messages.


	func main() {
	    log.Startf("TEST", "main", "Starting Test App")
	
	    cfg := udp.Config{
	        NetType:      "udp4",
	        Addr:         ":9000",
	        WorkRoutines: 2,
	        WorkStats:    time.Minute,
	        ConnHandler:  udpConnHandler{},
	        ReqHandler:   udpReqHandler{},
	        RespHandler:  udpRespHandler{},
	    }
	
	    u, err := udp.New("TEST", &cfg)
	    if err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    if err := u.Start("TEST"); err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    // Wait for a signal to shutdown.
	    sigChan := make(chan os.Signal, 1)
	    signal.Notify(sigChan, os.Interrupt)
	    <-sigChan
	
	    u.Stop("TEST")
	
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
* [type OptUserPool](#OptUserPool)
* [type ReqHandler](#ReqHandler)
* [type Request](#Request)
  * [func (r *Request) Work(context interface{}, id int)](#Request.Work)
* [type RespHandler](#RespHandler)
* [type Response](#Response)
  * [func (r *Response) Work(context interface{}, id int)](#Response.Work)
* [type UDP](#UDP)
  * [func New(context interface{}, name string, cfg Config) (*UDP, error)](#New)
  * [func (d *UDP) Addr() net.Addr](#UDP.Addr)
  * [func (d *UDP) Do(context interface{}, r *Response) error](#UDP.Do)
  * [func (d *UDP) Start(context interface{}) error](#UDP.Start)
  * [func (d *UDP) StatsRecv() pool.Stat](#UDP.StatsRecv)
  * [func (d *UDP) StatsSend() pool.Stat](#UDP.StatsSend)
  * [func (d *UDP) Stop(context interface{}) error](#UDP.Stop)


#### <a name="pkg-files">Package files</a>
[doc.go](/src/github.com/ardanlabs/kit/udp/doc.go) [handlers.go](/src/github.com/ardanlabs/kit/udp/handlers.go) [udp.go](/src/github.com/ardanlabs/kit/udp/udp.go) [udp_config.go](/src/github.com/ardanlabs/kit/udp/udp_config.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrInvalidConfiguration = errors.New("Invalid Configuration")
    ErrInvalidNetType       = errors.New("Invalid NetType Configuration")
    ErrInvalidConnHandler   = errors.New("Invalid Connection Handler Configuration")
    ErrInvalidReqHandler    = errors.New("Invalid Request Handler Configuration")
    ErrInvalidRespHandler   = errors.New("Invalid Response Handler Configuration")
)
```
Set of error variables for start up.




## <a name="Config">type</a> [Config](/src/target/udp_config.go?s=997:2086#L17)
``` go
type Config struct {
    NetType string // "udp", udp4" or "udp6"
    Addr    string // "host:port" or "[ipv6-host%zone]:port"

    ConnHandler ConnHandler // Support for binding new connections to a reader and writer.
    ReqHandler  ReqHandler  // Support for handling the specific request workflow.
    RespHandler RespHandler // Support for handling the specific response workflow.

    OptUserPool
    OptIntPool

    OptEvent
}
```
Config provides a data structure of required configuration parameters.










### <a name="Config.Event">func</a> (\*Config) [Event](/src/target/udp_config.go?s=2610:2702#L69)
``` go
func (cfg *Config) Event(context interface{}, event string, format string, a ...interface{})
```
Event fires events back to the user for important events.




### <a name="Config.Validate">func</a> (\*Config) [Validate](/src/target/udp_config.go?s=2144:2179#L44)
``` go
func (cfg *Config) Validate() error
```
Validate checks the configuration to required items.




## <a name="ConnHandler">type</a> [ConnHandler](/src/target/handlers.go?s=158:310#L2)
``` go
type ConnHandler interface {
    // Bind is called to set the reader and writer.
    Bind(context interface{}, listener *net.UDPConn) (io.Reader, io.Writer)
}
```
ConnHandler is implemented by the user to bind the listener
to a reader and writer for processing.










## <a name="OptEvent">type</a> [OptEvent](/src/target/udp_config.go?s=817:921#L12)
``` go
type OptEvent struct {
    Event func(context interface{}, event string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.










## <a name="OptIntPool">type</a> [OptIntPool](/src/target/udp_config.go?s=420:760#L4)
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










## <a name="OptUserPool">type</a> [OptUserPool](/src/target/udp_config.go?s=148:312#L1)
``` go
type OptUserPool struct {
    RecvPool *pool.Pool // User provided work pool for the receive work.
    SendPool *pool.Pool // User provided work pool for the send work.
}
```
OptUserPool declares fields for the user to pass their own
work pools for configuration.










## <a name="ReqHandler">type</a> [ReqHandler](/src/target/handlers.go?s=503:963#L11)
``` go
type ReqHandler interface {
    // Read is provided the user-defined reader and must return the data read
    // off the wire and the length. Returning io.EOF or a non temporary error
    // will show down the listener.
    Read(context interface{}, reader io.Reader) (*net.UDPAddr, []byte, int, error)

    // Process is used to handle the processing of the request. This method
    // is called on a routine from a pool of routines.
    Process(context interface{}, r *Request)
}
```
ReqHandler is implemented by the user to implement the processing
of request messages from the client.










## <a name="Request">type</a> [Request](/src/target/handlers.go?s=1015:1136#L23)
``` go
type Request struct {
    UDP     *UDP
    UDPAddr *net.UDPAddr
    IsIPv6  bool
    ReadAt  time.Time
    Data    []byte
    Length  int
}
```
Request is the message received by the client.










### <a name="Request.Work">func</a> (\*Request) [Work](/src/target/handlers.go?s=1421:1472#L39)
``` go
func (r *Request) Work(context interface{}, id int)
```
Work implements the worker inteface for processing messages. This is called
from a routine in the work pool.




## <a name="RespHandler">type</a> [RespHandler](/src/target/handlers.go?s=1711:1869#L47)
``` go
type RespHandler interface {
    // Write is provided the user-defined writer and the data to write.
    Write(context interface{}, r *Response, writer io.Writer)
}
```
RespHandler is implemented by the user to implement the processing
of the response messages to the client.










## <a name="Response">type</a> [Response](/src/target/handlers.go?s=1917:2059#L53)
``` go
type Response struct {
    UDPAddr  *net.UDPAddr
    Data     []byte
    Length   int
    Complete func(r *Response)
    // contains filtered or unexported fields
}
```
Response is message to send to the client.










### <a name="Response.Work">func</a> (\*Response) [Work](/src/target/handlers.go?s=2195:2247#L65)
``` go
func (r *Response) Work(context interface{}, id int)
```
Work implements the worker interface for sending messages. Called by
AsyncSend via the d.client.Do(context, &resp) method call.




## <a name="UDP">type</a> [UDP](/src/target/udp.go?s=751:1048#L21)
``` go
type UDP struct {
    Config
    Name string
    // contains filtered or unexported fields
}
```
UDP manages message to a specific ip address and port.







### <a name="New">func</a> [New](/src/target/udp.go?s=1099:1167#L44)
``` go
func New(context interface{}, name string, cfg Config) (*UDP, error)
```
New creates a new manager to service clients.





### <a name="UDP.Addr">func</a> (\*UDP) [Addr](/src/target/udp.go?s=6607:6636#L286)
``` go
func (d *UDP) Addr() net.Addr
```
Addr returns the local listening network address.




### <a name="UDP.Do">func</a> (\*UDP) [Do](/src/target/udp.go?s=6078:6134#L264)
``` go
func (d *UDP) Do(context interface{}, r *Response) error
```
Do will post the request to be sent by the client worker pool.




### <a name="UDP.Start">func</a> (\*UDP) [Start](/src/target/udp.go?s=2765:2811#L118)
``` go
func (d *UDP) Start(context interface{}) error
```
Start begins to accept data.




### <a name="UDP.StatsRecv">func</a> (\*UDP) [StatsRecv](/src/target/udp.go?s=6360:6395#L276)
``` go
func (d *UDP) StatsRecv() pool.Stat
```
StatsRecv returns the current snapshot of the recv pool stats.




### <a name="UDP.StatsSend">func</a> (\*UDP) [StatsSend](/src/target/udp.go?s=6490:6525#L281)
``` go
func (d *UDP) StatsSend() pool.Stat
```
StatsSend returns the current snapshot of the send pool stats.




### <a name="UDP.Stop">func</a> (\*UDP) [Stop](/src/target/udp.go?s=5365:5410#L230)
``` go
func (d *UDP) Stop(context interface{}) error
```
Stop shuts down the manager and closes all connections.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
